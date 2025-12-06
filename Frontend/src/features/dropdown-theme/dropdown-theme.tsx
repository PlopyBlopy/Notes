import { useState } from "react";
import styles from "./dropdown-theme.module.css";
import type { ThemeInfo } from "@/shared/api";

type Props = {
  options: ThemeInfo[];
  value: ThemeInfo;
  onChange: (value: number) => void;
  placeholder?: string;
};

export const DropdownTheme = ({ options, value, onChange, placeholder = "Выберите опцию" }: Props) => {
  const [isOpen, setIsOpen] = useState<boolean>(false);

  const handleSelect = (option: number) => {
    onChange(option);
    setIsOpen(false);
  };

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation();
    onChange(0);
    setIsOpen(false);
  };

  return (
    <div className={styles.container}>
      <div className={styles.header} onClick={() => setIsOpen(!isOpen)}>
        <span className={value ? styles.selected : styles.placeholder}>{value?.title || placeholder}</span>
        <div className={styles.controls}>
          {value && (
            <span className={styles.clear} onClick={handleClear}>
              ×
            </span>
          )}
          <span className={styles.arrow}>{isOpen ? "▲" : "▼"}</span>
        </div>
      </div>

      {isOpen && (
        <div className={styles.list}>
          {options.map((option) => (
            <div
              key={option.id}
              className={`${styles.item} ${value.id === option.id ? styles.itemSelected : ""}`}
              onClick={() => handleSelect(option.id)}
            >
              {option?.title}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
