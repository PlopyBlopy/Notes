import {
  deleteNote,
  getCardColors,
  getNotes,
  getTagColors,
  getTags,
  getThemes,
  patchNoteCompleted,
  postNote,
  putNote,
  type Card,
  type ColorInfo,
  type CreateNote,
  type FilteredNotes,
  type NotesFilter,
  type TagInfo,
  type ThemeInfo,
  type UpdateNote,
} from "@/shared/api";

class Store {
  private _cards: Card[] = [];
  private _tags: Map<number, TagInfo> = new Map();
  private _tagArr: TagInfo[] = [];
  private _themes: Map<number, ThemeInfo> = new Map();
  private _themesArr: ThemeInfo[] = [];
  private _tagColors: Map<number, ColorInfo> = new Map();
  private _cardColors: Map<number, ColorInfo> = new Map();
  private _cardColorArr: ColorInfo[] = [];

  private _filter: NotesFilter;
  private _cursor: number;

  private listeners = new Set<() => void>();

  constructor() {
    this._filter = {
      completed: false,
      search: "",
      limit: 20,
      themeId: 0,
      tagIds: [],
    };
    this._cursor = 0;
  }

  static async Init(): Promise<Store> {
    const store = new Store();
    await Promise.all([store.initTags(), store.initTagColors(), store.initCardColors(), store.initThemes()]);
    store.initCards();
    store.notify();
    return store;
  }

  Subscribe(callback: () => void): () => void {
    this.listeners.add(callback);
    return () => this.listeners.delete(callback);
  }

  private notify(): void {
    this.listeners.forEach((cb) => cb());
  }

  private async initCards() {
    const data: FilteredNotes = await getNotes(this._filter, this._cursor);

    this._cards = data.cards;
    this._cursor = data.cursor;
  }

  private initTags = async () => {
    const tags: TagInfo[] = await getTags();

    tags.forEach((tag) => {
      this._tags.set(tag.id, tag);
    });

    this._tagArr = tags;
  };

  private async initTagColors() {
    const colors: ColorInfo[] = await getTagColors();

    colors.forEach((color) => {
      this._tagColors.set(color.id, color);
    });
  }

  private async initCardColors() {
    const colors: ColorInfo[] = await getCardColors();

    colors.forEach((color) => {
      this._cardColors.set(color.id, color);
    });

    this._cardColorArr = colors;
  }

  private async initThemes() {
    const themes: ThemeInfo[] = await getThemes();

    themes.forEach((theme) => {
      this._themes.set(theme.id, theme);
    });

    this._themesArr = themes;
  }

  public async PostNote(note: CreateNote) {
    await postNote(note);
    await this.UpdateCards();
  }

  public GetCards(): Card[] {
    return this._cards;
  }

  public GetTags(): Map<number, TagInfo> {
    return this._tags;
  }
  public GetTagArr(): TagInfo[] {
    return this._tagArr;
  }

  public GetTagColors(): Map<number, ColorInfo> {
    return this._tagColors;
  }

  public GetCardColors(): Map<number, ColorInfo> {
    return this._cardColors;
  }

  public GetCardColorArr(): ColorInfo[] {
    return this._cardColorArr;
  }
  public GetThemes(): Map<number, ThemeInfo> {
    return this._themes;
  }

  public GetThemeArr(): ThemeInfo[] {
    return this._themesArr;
  }
  public GetFilter(): NotesFilter {
    return this._filter;
  }

  public async UpdateCards() {
    const data: FilteredNotes = await getNotes(this._filter, this._cursor);

    if (this._cursor === 0) {
      this._cards = data.cards;
    } else {
      this._cards = this._cards.concat(data.cards);
    }
    this._cursor = data.cursor;

    this.notify();
  }

  public async UpdateNote(note: UpdateNote) {
    await putNote(note);

    this.ResetCursor();
    this.UpdateCards();
  }

  public async UpdateNoteCompleted(id: number, completed: boolean) {
    await patchNoteCompleted(id, completed);

    this._cards = this._cards.filter((card) => card.note.id !== id);
    this._cursor = this._cursor - 1;

    this.notify();
  }

  public UpdateFilter(partialFilter: Partial<NotesFilter>) {
    this._filter = {
      ...this._filter,
      ...partialFilter,
    };

    this.ResetCursor();
    this.UpdateCards();
  }

  public ResetCursor() {
    this._cursor = 0;
  }

  public async DeleteNote(id: number) {
    await deleteNote(id);

    this._cards = this._cards.filter((card) => card.note.id !== id);
    this._cursor = this._cursor - 1;

    this.notify();
  }
}

export const store = Store.Init();
