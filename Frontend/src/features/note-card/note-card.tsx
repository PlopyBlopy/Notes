import { useState, useRef, useEffect } from "react";
import { type Card } from "@/shared/api";
import styles from "./note-card.module.css";
import { ButtonIcon } from "@/shared/components/button-icon";
import { Icons } from "@/shared/assets/icons";
import { MarkedWord } from "@/shared/components/marked-word";
import { Checkbox } from "@/shared/components/checkbox/checkbox";
import { Tag } from "../tag";
import { useStore } from "@/shared/hook/store";

interface Style {
  backgroundColor?: string;
}

interface PropNoteContainer extends Style {
  card: Card;
}

export const NoteCard = ({ card }: PropNoteContainer) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [hasLongText, setHasLongText] = useState(false);
  const descriptionRef = useRef<HTMLDivElement>(null);
  const { tags, themes, cardColors } = useStore();

  useEffect(() => {
    if (descriptionRef.current) {
      const element = descriptionRef.current;
      const lineHeight = parseFloat(getComputedStyle(element).lineHeight);
      const maxHeight = lineHeight * 3;
      const isLongText = element.scrollHeight > maxHeight;

      setHasLongText(isLongText);
    }
  }, [card.note.description]);

  if (!card) {
    console.warn("noteMetadata component received undefined noteMetadata");
    return null;
  }

  const toggleDescription = () => {
    if (hasLongText) {
      setIsExpanded(!isExpanded);
    }
  };

  const handleEdit = () => {};
  const handleDelete = () => {};

  const style: React.CSSProperties = {
    backgroundColor: cardColors.get(card.noteColorId)?.variable,
  };

  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString("ru-RU", {
        day: "numeric",
        month: "short",
        year: "numeric",
      });
    } catch {
      return dateString;
    }
  };

  return (
    <div style={style} className={styles.container}>
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <Checkbox />
          <div className={styles.title}>{card.note.title}</div>
        </div>
        <MarkedWord text={themes.get(card.themeId)?.title} color="var(--text-color-primary)" backgroundColor="var(--color-light)" />
      </div>

      <div className={styles.middle}>
        <div className={styles.descriptionWrapper}>
          <div
            ref={descriptionRef}
            className={`${styles.description} ${
              isExpanded ? styles.descriptionExpanded : styles.descriptionCollapsed
            } ${hasLongText ? styles.canExpand : ""}`}
            onClick={toggleDescription}
            title={hasLongText ? (isExpanded ? "Свернуть описание" : "Раскрыть описание") : undefined}
            style={{
              // Динамически вычисляем max-height для анимации
              maxHeight: isExpanded ? `${descriptionRef.current?.scrollHeight}px` : hasLongText ? "4.5em" : "none",
            }}
          >
            {card.note.description}
          </div>
        </div>
        <div className={styles.tags}>
          {card.tagIds.map((tagId, index) => (
            <Tag key={index} tag={tags.get(tagId)} />
          ))}
        </div>
      </div>

      <div className={styles.footer}>
        <div>{formatDate(card.createdAt)}</div>
        <div className={styles.footerLeft}>
          <ButtonIcon onClick={handleEdit} IconComponent={Icons.elements.edit} label="edit" variant="greyDark" />
          <ButtonIcon onClick={handleDelete} IconComponent={Icons.elements.delete} label="delete" variant="danger" />
        </div>
      </div>
    </div>
  );
};
