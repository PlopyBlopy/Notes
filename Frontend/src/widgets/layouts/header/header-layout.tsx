import { ThemeToggle } from "@/widgets/theme-toggle";
import styles from "./header-layout.module.css";
import { Icon } from "@/shared/components/icon";
import { Icons } from "@/shared/assets/icons";
import { PrimaryButtonIcon } from "@/shared/components/primary-button-icon";
import { Modal } from "@/shared/components/modal";
import { NoteCreate } from "@/widgets/note-create/note-create";
import { useState } from "react";
import { type CreateNote } from "@/shared/api";
import { useStore } from "@/shared/hook/store";

export const HeaderLayout = () => {
  const [isModalOpen, setModalOpen] = useState(false);
  const { postNote } = useStore();

  const createNote = async (note: CreateNote) => {
    try {
      await postNote(note);
      handleClose();
    } catch (err) {
      console.error("Ошибка при создании заметки:", err);
    }
  };

  const handleOpen = () => {
    setModalOpen(true);
  };
  const handleClose = () => {
    setModalOpen(false);
  };
  return (
    <div className={styles.container}>
      <div className={styles.items}>
        <Icon IconComponent={Icons.elements.note} />
        <PrimaryButtonIcon onClick={handleOpen} text="Новая заметка" IconComponent={Icons.elements.plus} />
        <ThemeToggle />
        <Modal isOpen={isModalOpen} onClose={handleClose}>
          <NoteCreate onClose={handleClose} onSubmit={createNote} />
        </Modal>
      </div>
    </div>
  );
};
