import { Icons } from "@/shared/assets/icons";
import styles from "./primary-button-icon.module.css";

export interface Style {
  variant?: "light" | "greyDarker" | "accent" | "danger"; // Вместо color
}

interface Props extends Style {
  onClick: () => void;
  IconComponent?: React.ComponentType<{ className?: string }>;
  label: string;
}

export const ButtonIcon = ({ onClick, IconComponent = Icons.default, label, variant = "greyDarker" }: Props) => {
  const variantClass = styles[variant] || styles.greyDarker;

  return (
    <button title={label} className={`${styles.button} ${variantClass}`} onClick={onClick}>
      <IconComponent className={styles.icon} />
    </button>
  );
};
