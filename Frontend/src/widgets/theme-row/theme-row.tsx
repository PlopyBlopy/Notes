import { useStore } from "@/shared/hook/store";
import styles from "./theme-row.module.css";

type Props = {
  value: number;
  onChange: (value: number) => void;
};

export const ThemeRow = ({ value, onChange }: Props) => {
  const { themeArr } = useStore();

  return (
    <div className={styles.container}>
      <div className={styles.themeContainer}>
        {themeArr.map((t, i) => (
          <button key={`theme-${i}`} className={value === t.id ? styles.selectButton : styles.button} onClick={() => onChange(t.id)}>
            {t.title}
          </button>
        ))}
      </div>
      <div className={styles.footer}></div>
    </div>
  );
};
