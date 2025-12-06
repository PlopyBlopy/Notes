import styles from "./tag-select.module.css";
import { Tag } from "@/features/tag";
import { useStore } from "@/shared/hook/store";

interface TagsSelectProps {
  value: number[];
  onChange: (value: number[]) => void;
  placeholder?: string;
}

export const TagsSelect = ({ value = [], onChange, placeholder = "Теги еще не созданы" }: TagsSelectProps) => {
  const { tagArr } = useStore();

  const handleToggle = (option: number) => {
    if (value.includes(option)) {
      // Если тег уже выбран - удаляем
      onChange(value.filter((item) => item !== option));
    } else {
      // Если тег не выбран - добавляем
      onChange([...value, option]);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.tagsRow}>
        {tagArr
          .filter((tag) => value.includes(tag.id))
          .map((tag, index) => (
            <div key={`selected-${index}`} className={styles.selectedTag} onClick={() => handleToggle(tag.id)} title="Кликните чтобы удалить">
              <Tag tag={tag} />
            </div>
          ))}
        {/* Затем доступные теги */}
        {tagArr
          .filter((tag) => !value.includes(tag.id))
          .map((tag, index) => (
            <div key={`available-${index}`} className={styles.availableTag} onClick={() => handleToggle(tag.id)} title="Кликните чтобы добавить">
              <Tag tag={tag} />
            </div>
          ))}

        {/* Плейсхолдер, когда ничего не выбрано и нет тегов */}
        {value.length === 0 && tagArr.length === 0 && <span className={styles.placeholder}>{placeholder}</span>}
      </div>
    </div>
  );
};
