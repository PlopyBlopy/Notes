import {
  getCardColors,
  getNotes,
  getTagColors,
  getTags,
  getThemes,
  type Card,
  type ColorInfo,
  type FilteredNotes,
  type NotesFilter,
  type TagInfo,
  type ThemeInfo,
} from "@/shared/api";

class Store {
  private _cards: Card[] = [];
  private _tags: Map<number, TagInfo> = new Map();
  private _tagArr: TagInfo[] = [];
  private _tagColors: Map<number, ColorInfo> = new Map();
  private _cardColors: Map<number, ColorInfo> = new Map();
  private _cardColorArr: ColorInfo[] = [];
  private _themes: ThemeInfo[] = [];

  private listeners = new Set<() => void>();

  static async Init(): Promise<Store> {
    const store = new Store();
    await Promise.all([store.initCards(), store.initTags(), store.initTagColors(), store.initCardColors(), store.initThemes()]);
    store.notify();
    return store;
  }

  private initCards = async () => {
    const filter: NotesFilter = {
      completed: false,
      search: "",
      limit: 20,
      themeId: 0,
      tagIds: [],
    };
    const data: FilteredNotes = await getNotes(filter, 0);

    this._cards = data.cards;
  };

  private initTags = async () => {
    const tags: TagInfo[] = await getTags();

    tags.forEach((tag) => {
      this._tags.set(tag.id, tag);
    });

    this._tagArr = tags;
  };

  private initTagColors = async () => {
    const colors: ColorInfo[] = await getTagColors();

    colors.forEach((color) => {
      this._tagColors.set(color.id, color);
    });
  };

  private initCardColors = async () => {
    const colors: ColorInfo[] = await getCardColors();

    colors.forEach((color) => {
      this._cardColors.set(color.id, color);
    });

    this._cardColorArr = colors;
  };

  private initThemes = async () => {
    const themes: ThemeInfo[] = await getThemes();
    this._themes = themes;
  };

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

  public GetThemes(): ThemeInfo[] {
    return this._themes;
  }

  public async UpdateCards(filter: NotesFilter, cursor: number) {
    const data: FilteredNotes = await getNotes(filter, cursor);

    this._cards = data.cards;

    this.notify();
  }

  Subscribe(callback: () => void): () => void {
    this.listeners.add(callback);
    return () => this.listeners.delete(callback);
  }

  private notify(): void {
    this.listeners.forEach((cb) => cb());
  }
}

export const store = Store.Init();
