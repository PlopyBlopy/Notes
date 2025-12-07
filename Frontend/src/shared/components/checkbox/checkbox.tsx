import styles from "./checkbox.module.css";

type Props = {
  isCompleted: boolean;
  onComplete: () => void;
};

export const Checkbox = ({ isCompleted, onComplete }: Props) => {
  return (
    <input
      className={styles.container}
      type="checkbox"
      aria-label="chekbox"
      disabled={isCompleted}
      onChange={() => {
        onComplete();
      }}
    />
  );
};
