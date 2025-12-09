import { NoteCard } from "@/features/note-card";
import styles from "./note-list.module.css";
import { useStore } from "@/shared/hook/store";

export const NoteList = () => {
  const { cards, updNoteCompleted, delNote } = useStore();

  const handleComplete = async (id: number, completed: boolean) => {
    await updNoteCompleted(id, completed);
  };

  const handleEdit = (id: number) => {
    console.log(`edit: ${id}`);
  };

  const handleDelete = async (id: number) => {
    await delNote(id);
  };

  return (
    <div className={styles.container}>
      {cards?.map((c) => (
        <NoteCard key={`card-${c.note.id}`} card={c} onComplete={handleComplete} onEdit={handleEdit} onDelete={handleDelete} />
      ))}
    </div>
  );
};
