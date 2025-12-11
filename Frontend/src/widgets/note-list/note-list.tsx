import { NoteCard } from "@/features/note-card";
import styles from "./note-list.module.css";
import { useStore } from "@/shared/hook/store";
import type { Card, UpdateNote } from "@/shared/api";
import type { NoteForm } from "../note-create";
import { useEffect, useState } from "react";
import { Modal } from "@/shared/components/modal";
import { NoteFormModal } from "../note-create/note-form-modal";

export const NoteList = () => {
  const [updatedCard, setUpdatedCard] = useState<Card | undefined>(undefined);
  const [form, setForm] = useState<NoteForm | undefined>(undefined);
  const [isModalOpen, setModalOpen] = useState(false);
  const { cards, updNote, updNoteCompleted, delNote } = useStore();

  useEffect(() => {
    const handleEdit = () => {
      if (updatedCard === undefined) {
        return;
      }
      const form: NoteForm = {
        title: updatedCard.note.title,
        description: updatedCard.note.description,
        ...updatedCard,
      };

      setForm(form);

      handleOpen();
    };
    handleEdit();
  }, [updatedCard]);

  const handleComplete = async (id: number, completed: boolean) => {
    await updNoteCompleted(id, completed);
  };

  const handleDelete = async (id: number) => {
    await delNote(id);
  };

  const handleSubmit = async (form: NoteForm) => {
    handleClose();

    if (updatedCard === undefined) {
      return;
    }

    const note: UpdateNote = {
      id: updatedCard.note.id,
      ...form,
    };

    updNote(note);
  };

  const handleOpen = () => {
    setModalOpen(true);
  };
  const handleClose = () => {
    setModalOpen(false);
    setUpdatedCard(undefined);
  };

  return (
    <div>
      <div>
        <Modal isOpen={isModalOpen} onClose={handleClose}>
          <NoteFormModal
            label={"Изменить заметку"}
            submitLabel="Изменить"
            onClose={handleClose}
            onSubmit={(val) => handleSubmit(val)}
            initForm={form}
          />
        </Modal>
      </div>
      <div className={styles.container}>
        {cards?.map((c) => (
          <NoteCard key={`card-${c.note.id}`} card={c} onComplete={handleComplete} onEdit={(card) => setUpdatedCard(card)} onDelete={handleDelete} />
        ))}
      </div>
    </div>
  );
};
