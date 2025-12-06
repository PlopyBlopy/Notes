import { NoteCard } from "@/features/note-card";
import styles from "./note-list.module.css";
import { useStore } from "@/shared/hook/store";

export const NoteList = () => {
  const { cards } = useStore();

  const handleComplete = () => {};
  const handleEdit = () => {};
  const handleDelete = () => {};

  return (
    <div className={styles.container}>
      {cards?.map((c, i) => (
        <NoteCard key={`card-${i}`} card={c} />
      ))}
    </div>
  );
};
