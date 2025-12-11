export type FilteredNotes = {
  cards: Card[];
  cursor: number;
};

export type Card = {
  note: Note;
  completed: boolean;
  themeId: number;
  tagIds: number[];
  noteColorId: number;
  createdAt: string;
};

export type Note = {
  id: number;
  title: string;
  description: string;
};

export type CreateNote = {
  title: string;
  description: string;
  themeId: number;
  tagIds: number[];
  noteColorId: number;
};

export type UpdateNote = {
  id: number;
  title: string;
  description: string;
  themeId: number;
  tagIds: number[];
  noteColorId: number;
};

export type ThemeInfo = {
  id: number;
  title: string;
};

export type TagInfo = {
  id: number;
  title: string;
  colorId: number;
};

export type ColorInfo = {
  id: number;
  name: string;
  variable: string;
};

export type NotesFilter = {
  completed: boolean;
  search: string;
  limit: number;
  themeId: number;
  tagIds: number[];
};
