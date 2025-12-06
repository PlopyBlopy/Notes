import { NoteList } from "@/widgets/note-list/note-list";
import styles from "./main-page.module.css";
import { ThemeRow } from "@/widgets/theme-row";
import { useEffect, useState } from "react";
import { TagsSelect } from "@/widgets/tag-select";
import { SearchBar } from "@/widgets/search-bar";
import { type NotesFilter } from "@/shared/api";
import { useStore } from "@/shared/hook/store";

export const MainPage = () => {
  const { updCards } = useStore();

  const [filter, setFilter] = useState<NotesFilter>({
    completed: false,
    search: "",
    limit: 20,
    themeId: 0,
    tagIds: [],
  });
  const [cursor, setCursor] = useState<number>(0);

  useEffect(() => {
    const loadNotes = async () => {
      updCards(filter, cursor);
    };

    loadNotes();
  }, [filter]);

  const handleSearchChange = (search: string) => {
    setFilter((prev) => ({ ...prev, search: search }));
    setCursor(0);
  };

  const handleThemeChange = (selectedTheme: number) => {
    setFilter((prev) => ({ ...prev, themeId: selectedTheme }));
    setCursor(0);
  };

  const handleTagsChange = (selectedTags: number[]) => {
    setFilter((prev) => ({ ...prev, tagIds: selectedTags }));
    setCursor(0);
  };

  return (
    <div className={styles.container}>
      <ThemeRow value={filter.themeId} onChange={handleThemeChange} />
      <SearchBar value={filter.search} onSearch={handleSearchChange} placeholder="Поиск заметок..." delay={300} />
      <TagsSelect value={filter.tagIds} onChange={handleTagsChange} placeholder="Теги еще не созданы" />
      <NoteList />
    </div>
  );
};
