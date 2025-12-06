import styles from "./marked-word.module.css";

export interface Style {
  color?: string;
  backgroundColor?: string;
}

interface Props extends Style {
  text?: string;
}

export const MarkedWord = ({ text = "MarkedWord", color, backgroundColor }: Props) => {
  const style: React.CSSProperties = {
    color: color,
    backgroundColor: backgroundColor,
  };

  return (
    <div style={style} className={styles.word}>
      {text}
    </div>
  );
};
