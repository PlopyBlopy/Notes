import { type TagInfo } from "@/shared/api";
import { MarkedWord } from "@/shared/components/marked-word";
import { useStore } from "@/shared/hook/store";

type Props = {
  tag: TagInfo | undefined;
};

export const Tag = ({ tag }: Props) => {
  if (!tag) {
    return null;
  }

  const { tagColors } = useStore();

  return <MarkedWord text={tag.title} color={"var(--text-color-primary)"} backgroundColor={tagColors.get(tag.colorId)?.variable} />;
};
