import { NoteList } from "@/widgets/note-list/note-list";
import styles from "./main-page.module.css";
import { ThemeRow } from "@/widgets/theme-row";
import { TagsSelect } from "@/widgets/tag-select";
import { SearchBar } from "@/widgets/search-bar";
import { useStore } from "@/shared/hook/store";

export const MainPage = () => {
  const { filter, updNoteFilter } = useStore();

  const handleSearchChange = (search: string) => {
    updNoteFilter({ search: search });
  };

  const handleThemeChange = (selectedTheme: number) => {
    updNoteFilter({ themeId: selectedTheme });
  };

  const handleTagsChange = (selectedTags: number[]) => {
    updNoteFilter({ tagIds: selectedTags });
  };

  if (filter === undefined) {
    return null;
  }

  return (
    <div className={styles.container}>
      <ThemeRow value={filter.themeId} onChange={handleThemeChange} />
      <SearchBar value={filter.search} onSearch={handleSearchChange} placeholder="Поиск заметок..." delay={300} />
      <TagsSelect value={filter.tagIds} onChange={handleTagsChange} placeholder="Теги еще не созданы" />
      <NoteList />
    </div>
  );
};
