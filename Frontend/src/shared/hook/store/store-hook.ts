import { store } from "@/entities/store";
import { type Card, type ColorInfo, type NotesFilter, type TagInfo, type ThemeInfo } from "@/shared/api";
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
    updCards: (filter: NotesFilter, cursor: number) => {
      store.then((storeInstance) => {
        storeInstance.UpdateCards(filter, cursor);
      });
    },
    updNoteCompleted: (id: number, completed: boolean) => {
      store.then((storeInstance) => {
        storeInstance.UpdateNoteCompleted(id, completed);
      });
    },
    delNote: (id: number) => {
      store.then((storeInstance) => {
        storeInstance.DeleteNote(id);
      });
    },
  };

  return { cards, tags, tagArr, tagColors, cardColors, cardColorArr, themes, themeArr, ...actions };
}
