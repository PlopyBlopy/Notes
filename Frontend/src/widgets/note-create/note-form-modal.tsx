import { useState } from "react";
import styles from "./note-form-modal.module.css";
import type { NoteForm } from "./note-form-modal.model";
import { ButtonIcon } from "@/shared/components/button-icon";
import { Icons } from "@/shared/assets/icons";
import { PrimaryButtonSubmit } from "@/shared/components/primary-button-submit";
import { ColorPicker } from "@/features/color-picker";
import { DropdownTheme } from "@/features/dropdown-theme";
import { TagsSelect } from "../tag-select";
import { useStore } from "@/shared/hook/store";

type Props = {
  label: string;
  submitLabel: string;
  onClose: () => void;
  onSubmit: (form: NoteForm) => void;
  initForm?: NoteForm;
};

export const NoteFormModal = ({ label, submitLabel, onClose, onSubmit, initForm }: Props) => {
  const [form, setForm] = useState<NoteForm>({
    title: "",
    description: "",
    themeId: 0,
    tagIds: [],
    noteColorId: 0,
    ...initForm,
  });

  const { themeArr, cardColorArr } = useStore();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmitForm = (e: React.FormEvent) => {
    e.preventDefault();

    setForm((prev) => ({ ...prev, tagIds: prev.tagIds.sort() }));

    onSubmit(form);
  };

  return (
    <div>
      <div className={styles.header}>
        {label}
        <ButtonIcon onClick={onClose} IconComponent={Icons.elements.close} label={"close"} variant={"danger"} />
      </div>
      <form className={styles.middle} onSubmit={handleSubmitForm}>
        <div>
          <label>Название</label>
          <input
            type="text"
            id="title"
            name="title"
            required={true}
            value={form.title}
            onChange={handleChange}
            placeholder="Название"
            autoComplete="off"
            className={styles.formInput}
          />
        </div>

        <div>
          <label>Описание</label>
          <textarea
            id="description"
            name="description"
            required={false}
            value={form.description}
            onChange={handleChange}
            placeholder="Описание..."
            className={styles.formInput}
            rows={4}
          />
        </div>
        <div>
          <label>Темы</label>
          <DropdownTheme
            options={themeArr}
            value={themeArr[form.themeId]}
            onChange={(themeId) => setForm((prev) => ({ ...prev, themeId: themeId }))}
            placeholder="Тема"
          />
        </div>
        <div>
          <label>Теги</label>
          <TagsSelect value={form.tagIds} onChange={(tags) => setForm((prev) => ({ ...prev, tagIds: tags }))} />
        </div>
        <div>
          <ColorPicker
            options={cardColorArr}
            value={form.noteColorId}
            onColorSelectId={(color) => setForm((prev) => ({ ...prev, noteColorId: color }))}
            placeholder="Цвет темы"
          />
        </div>

        <div className={styles.footer}>
          <PrimaryButtonSubmit text={submitLabel} />
        </div>
      </form>
    </div>
  );
};
