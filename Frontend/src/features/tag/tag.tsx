import { type TagInfo } from "@/shared/api";
import { MarkedWord } from "@/shared/components/marked-word";
import { useStore } from "@/shared/hook/store";

type Props = {
  tag: TagInfo | undefined;
};

export const Tag = ({ tag }: Props) => {
  if (!tag) {
    console.warn("Tag component received undefined tag");
    return null;
  }

  const { tagColors } = useStore();

  return <MarkedWord text={tag.title} color={"var(--tag-text-color)"} backgroundColor={tagColors.get(tag.colorId)?.variable} />;
};
