package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/PlopyBlopy/notebot/config"
	service "github.com/PlopyBlopy/notebot/internal/adapters/note_service"
	"github.com/PlopyBlopy/notebot/internal/controller/middleware"
	"github.com/PlopyBlopy/notebot/internal/controller/router"
	addcardcolor "github.com/PlopyBlopy/notebot/internal/metadata/add_card_color"
	addtag "github.com/PlopyBlopy/notebot/internal/metadata/add_tag"
	addtagcolor "github.com/PlopyBlopy/notebot/internal/metadata/add_tag_color"
	addtheme "github.com/PlopyBlopy/notebot/internal/metadata/add_theme"
	getcardcolor "github.com/PlopyBlopy/notebot/internal/metadata/get_card_color"
	gettagcolor "github.com/PlopyBlopy/notebot/internal/metadata/get_tag_colors"
	gettags "github.com/PlopyBlopy/notebot/internal/metadata/get_tags"
	getthemes "github.com/PlopyBlopy/notebot/internal/metadata/get_themes"
	addnote "github.com/PlopyBlopy/notebot/internal/note/add_note"
	deletenote "github.com/PlopyBlopy/notebot/internal/note/delete_note"
	getfilterednotecards "github.com/PlopyBlopy/notebot/internal/note/get_filtered_note_cards"
	"github.com/PlopyBlopy/notebot/pkg/httpserver"
	"github.com/PlopyBlopy/notebot/pkg/logger"
	"github.com/PlopyBlopy/notebot/pkg/note"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

func main() {
	// launch
	logger.NewLogger()
	log.Info().Msg("bot launch")

	c, err := config.InitConfig()
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	if err := App(context.Background(), c); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize application:")
	}
}

func App(ctx context.Context, c config.Config) error {
	/* dependencies */
	metadataManager, err := note.NewMetadataManager(&c.Metadata)
	if err != nil {
		return err
	}

	indexManager, err := note.NewIndexManager(metadataManager)
	if err != nil {
		return err
	}

	noteIndexManager, err := note.NewNoteIndexManager(metadataManager, indexManager)
	if err != nil {
		return err
	}

	noteManager, err := note.NewNoteManager(metadataManager, noteIndexManager, indexManager)
	if err != nil {
		return err
	}

	err = indexManager.Scan()
	if err != nil {
		return err
	}

	metadataService, err := service.NewMetadataService(metadataManager)
	if err != nil {
		return err
	}

	noteService, err := service.NewNoteService(noteManager)
	if err != nil {
		return err
	}

	/* usecases */
	addnoteUcecase := addnote.NewUsecase(noteService)
	addtagUsecase := addtag.NewUsecase(metadataService)
	addthemeUsecase := addtheme.NewUsecase(metadataService)
	addtagcolorUsecase := addtagcolor.NewUsecase(metadataService)
	addcardcolorUsecase := addcardcolor.NewUsecase(metadataService)

	getfilterednotecardsUcecase := getfilterednotecards.NewUsecase(noteService)
	getcardcolorUcecase := getcardcolor.NewUsecase(metadataService)
	gettagsUcecase := gettags.NewUsecase(metadataService)
	gettagcolorUcecase := gettagcolor.NewUsecase(metadataService)
	getthemesUcecase := getthemes.NewUsecase(metadataService)

	deletenoteUcecase := deletenote.NewUsecase(noteService)

	/* http handlers */
	addNoteHandler := addnote.NewHttpHandler(addnoteUcecase)
	addtagHandler := addtag.NewHttpHandler(addtagUsecase)
	addthemeHandler := addtheme.NewHttpHandler(addthemeUsecase)
	addtagcolorHandler := addtagcolor.NewHttpHandler(addtagcolorUsecase)
	addcardcolorHandler := addcardcolor.NewHttpHandler(addcardcolorUsecase)

	getfilterednotecardsHandler := getfilterednotecards.NewHttpHandler(getfilterednotecardsUcecase)
	getcardcolorHandler := getcardcolor.NewHttpHandler(getcardcolorUcecase)
	gettagsHandler := gettags.NewHttpHandler(gettagsUcecase)
	gettagcolorHandler := gettagcolor.NewHttpHandler(gettagcolorUcecase)
	getthemesHandler := getthemes.NewHttpHandler(getthemesUcecase)

	deletenoteHandler := deletenote.NewHttpHandler(deletenoteUcecase)

	/* http router  */
	sm := router.NewServeMux()

	sm.AddHandler("POST /note", addNoteHandler)
	sm.AddHandler("POST /note/tag", addtagHandler)
	sm.AddHandler("POST /note/theme", addthemeHandler)
	sm.AddHandler("POST /note/tag/color", addtagcolorHandler)
	sm.AddHandler("POST /note/card/color", addcardcolorHandler)

	sm.AddHandler("GET /note/card/filtered", getfilterednotecardsHandler)
	sm.AddHandler("GET /note/card/color", getcardcolorHandler)
	sm.AddHandler("GET /note/tag", gettagsHandler)
	sm.AddHandler("GET /note/tag/color", gettagcolorHandler)
	sm.AddHandler("GET /note/theme", getthemesHandler)

	sm.AddHandler("DELETE /note", deletenoteHandler)

	httpRouter := sm.InitRouter(1)

	/* middleware */
	httpWrapRouter := middleware.NewMiddlewareBuilder().GlobalExceptionHandler(c).PanicMiddleware().Build(httpRouter)

	/* cors */
	corsOptions := cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"POST", "GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}

	corsRouter := cors.New(corsOptions).Handler(httpWrapRouter)

	/* chans */
	var mainErr error
	errChan := make(chan error, 1)
	done := make(chan struct{})
	var wg sync.WaitGroup

	/* http server */
	httpServer, err := httpserver.NewHttpServer(corsRouter, c.HttpServer)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func(c httpserver.HttpServerConfig) {
		defer close(errChan)
		defer wg.Done()

		log.Info().Msg("http server started")
		log.Info().Msgf("server listening on %s:%s", c.Host, c.Port)

		serverErr := make(chan error, 1)
		go func() {
			serverErr <- httpServer.Start()
		}()

		for {
			select {
			case err := <-serverErr:
				if err != nil {
					errChan <- err
				}
				return
			case _, ok := <-errChan:
				if ok {

				}
			case <-done:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := httpServer.Shutdown(ctx); err != nil {
					errChan <- fmt.Errorf("graceful shutdown failed: %w", err)
					log.Error().Msg("httpServer shutdown error")
				} else {
					log.Info().Msg("httpServer shutdown")
				}
				return
			}
		}

	}(c.HttpServer)

	/* Signal for Shutdown or Close */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	shouldExit := false
	for !shouldExit {
		select {
		case <-quit:
			log.Info().Msg("A termination signal is received, and the server stops...")
			close(done)
			if err := <-errChan; err != nil {
				mainErr = fmt.Errorf("server stopped with error: %w", err)
			}
			log.Info().Msg("stopped")
			log.Info().Msg("wait shutdown")
			shouldExit = true
		case err, ok := <-errChan:
			if err != nil {
				mainErr = fmt.Errorf("crashed: %w", err)
			}
			if ok {
				close(done)
				shouldExit = true
			}
		}
	}

	// Close or Shutdown

	wg.Wait()

	log.Info().Msg("app closed")

	return mainErr
}

func Shutdown(f func() error) func() error {
	return f
}
