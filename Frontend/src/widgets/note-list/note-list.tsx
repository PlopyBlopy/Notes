import { NoteCard } from "@/features/note-card";
import styles from "./note-list.module.css";
import { useStore } from "@/shared/hook/store";

type Props = {
  onUpdateCards: () => void;
};

export const NoteList = ({ onUpdateCards }: Props) => {
  const { cards, delNote } = useStore();

  const handleComplete = (id: number) => {
    console.log(`completed: ${id}`);
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
