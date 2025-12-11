import { store } from "@/entities/store";
import { type Card, type ColorInfo, type CreateNote, type NotesFilter, type TagInfo, type ThemeInfo, type UpdateNote } from "@/shared/api";
import { useEffect, useState } from "react";

export function useStore() {
  const [cards, setCards] = useState<Card[]>([]);
  const [tags, setTags] = useState<Map<number, TagInfo>>(new Map());
  const [tagArr, setTagArr] = useState<TagInfo[]>([]);
  const [themes, setThemes] = useState<Map<number, ThemeInfo>>(new Map());
  const [themeArr, setThemeArr] = useState<ThemeInfo[]>([]);
  const [tagColors, setTagColors] = useState<Map<number, ColorInfo>>(new Map());
  const [cardColors, setCardColors] = useState<Map<number, ColorInfo>>(new Map());
  const [cardColorArr, setCardColorArr] = useState<ColorInfo[]>([]);

  const [filter, setFilter] = useState<NotesFilter>();

  useEffect(() => {
    let unsubscribe: (() => void) | null = null;
    let mounted = true;

    store.then((storeInstance) => {
      if (!mounted) return;

      setCards(storeInstance.GetCards());
      setTags(storeInstance.GetTags());
      setTagArr(storeInstance.GetTagArr());
      setThemes(storeInstance.GetThemes());
      setThemeArr(storeInstance.GetThemeArr());
      setTagColors(storeInstance.GetTagColors());
      setCardColors(storeInstance.GetCardColors());
      setCardColorArr(storeInstance.GetCardColorArr());
      setFilter(storeInstance.GetFilter());

      unsubscribe = storeInstance.Subscribe(() => {
        if (!mounted) return;
        setCards(storeInstance.GetCards());
        setTags(storeInstance.GetTags());
        setTagArr(storeInstance.GetTagArr());
        setThemes(storeInstance.GetThemes());
        setThemeArr(storeInstance.GetThemeArr());
        setTagColors(storeInstance.GetTagColors());
        setCardColors(storeInstance.GetCardColors());
        setCardColorArr(storeInstance.GetCardColorArr());
        setFilter(storeInstance.GetFilter());
      });
    });

    return () => {
      mounted = false;
      if (unsubscribe) {
        unsubscribe();
      }
    };
  }, []);

  const actions = {
    postNote: (note: CreateNote) => {
      store.then((storeInstance) => {
        storeInstance.PostNote(note);
      });
    },
    updCards: () => {
      store.then((storeInstance) => {
        storeInstance.UpdateCards();
      });
    },
    updNote: (note: UpdateNote) => {
      store.then((storeInstance) => {
        storeInstance.UpdateNote(note);
      });
    },
    updNoteCompleted: (id: number, completed: boolean) => {
      store.then((storeInstance) => {
        storeInstance.UpdateNoteCompleted(id, completed);
      });
    },
    updNoteFilter: (filter: Partial<NotesFilter>) => {
      store.then((storeInstance) => {
        storeInstance.UpdateFilter(filter);
      });
    },
    resetCursor: () => {
      store.then((storeInstance) => {
        storeInstance.ResetCursor();
      });
    },
    delNote: (id: number) => {
      store.then((storeInstance) => {
        storeInstance.DeleteNote(id);
      });
    },
  };

  return { cards, tags, tagArr, tagColors, cardColors, cardColorArr, themes, themeArr, filter, ...actions };
}
