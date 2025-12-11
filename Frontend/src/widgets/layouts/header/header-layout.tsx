import type { CreateNote } from "@/shared/api";
import { Icons } from "@/shared/assets/icons";
import { Icon } from "@/shared/components/icon";
import { Modal } from "@/shared/components/modal";
import { PrimaryButtonIcon } from "@/shared/components/primary-button-icon";
import { useStore } from "@/shared/hook/store";
import type { NoteForm } from "@/widgets/note-create";
import { NoteFormModal } from "@/widgets/note-create/note-form-modal";
import { ThemeToggle } from "@/widgets/theme-toggle";
import { useState } from "react";
import styles from "./header-layout.module.css";

export const HeaderLayout = () => {
  const [isModalOpen, setModalOpen] = useState(false);
  const { postNote } = useStore();

  const handleOpen = () => {
    setModalOpen(true);
  };
  const handleClose = () => {
    setModalOpen(false);
  };

  const handlePostNote = async (form: NoteForm) => {
    handleClose();

    const note: CreateNote = {
      ...form,
    };

    await postNote(note);
  };

  return (
    <div className={styles.container}>
      <div className={styles.items}>
        <Icon IconComponent={Icons.elements.note} />
        <PrimaryButtonIcon onClick={handleOpen} text={"Добавить заметку"} IconComponent={Icons.elements.plus} />
        <ThemeToggle />
        <Modal isOpen={isModalOpen} onClose={handleClose}>
          <NoteFormModal label={"Новая заметка"} submitLabel="Создать" onClose={handleClose} onSubmit={(val) => handlePostNote(val)} />
        </Modal>
      </div>
    </div>
  );
};
