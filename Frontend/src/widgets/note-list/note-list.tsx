import { NoteCard } from "@/features/note-card";
import styles from "./note-list.module.css";
import { useStore } from "@/shared/hook/store";

type Props = {
  onUpdateCards: () => void;
};

export const NoteList = ({ onUpdateCards }: Props) => {
  const { cards, updNoteCompleted, delNote } = useStore();

  const handleComplete = async (id: number, completed: boolean) => {
    await updNoteCompleted(id, completed);
    await onUpdateCards();
  };

  const handleEdit = (id: number) => {
    console.log(`edit: ${id}`);
  };

  const handleDelete = async (id: number) => {
    await delNote(id);
    await onUpdateCards();
  };

  return (
    <div className={styles.container}>
      {cards?.map((c, i) => (
        <NoteCard key={`card-${i}`} card={c} onComplete={handleComplete} onEdit={handleEdit} onDelete={handleDelete} />
      ))}
    </div>
  );
};
