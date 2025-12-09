import type { ColorInfo, CreateNote, NotesFilter, TagInfo, ThemeInfo, FilteredNotes } from "./api.model";

const apiurl = "http://localhost:8080/api/v1";

export const postNote = async (note: CreateNote) => {
  const response = await fetch(`${apiurl}/note`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(note),
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
};

// export const postTag = async (tag: TagInfo) => {};

// export const postCompleteNote = async (id: number) => {};

export const getNotes = async (filter: NotesFilter, cursor: number): Promise<FilteredNotes> => {
  const params = new URLSearchParams({
    completed: `${filter.completed}`,
    search: `${filter.search}`,
    limit: `${filter.limit}`,
    cursor: `${cursor}`,
    themeId: `${filter.themeId}`,
    tagIds: `${filter.tagIds}`,
  });

  const response = await fetch(`${apiurl}/note/card/filtered?${params}`);
  const data: FilteredNotes = await response.json();

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return data || [];
};

export const getTags = async (): Promise<TagInfo[]> => {
  const response = await fetch(`${apiurl}/note/tag`);
  const data = await response.json();
  return data;
};

export const getTagColors = async (): Promise<ColorInfo[]> => {
  const response = await fetch(`${apiurl}/note/tag/color`);
  const data = await response.json();
  return data;
};

export const getThemes = async (): Promise<ThemeInfo[]> => {
  const response = await fetch(`${apiurl}/note/theme`);
  const data = await response.json();
  return data;
};

export const getCardColors = async (): Promise<ColorInfo[]> => {
  const response = await fetch(`${apiurl}/note/card/color`);
  const data = await response.json();
  return data;
};

// export const patchNote = async (note: UpdateNote) => {};

export const patchNoteCompleted = async (id: number, completed: boolean) => {
  const params = new URLSearchParams({
    id: `${id}`,
    completed: `${completed}`,
  });

  const response = await fetch(`${apiurl}/note?${params}`, {
    method: "PATCH",
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
};

export const deleteNote = async (id: number) => {
  const params = new URLSearchParams({
    id: `${id}`,
  });

  const response = await fetch(`${apiurl}/note?${params}`, {
    method: "DELETE",
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
};
