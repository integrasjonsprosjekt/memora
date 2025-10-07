import { MultipleChoiceCardType } from '../types';

export default function MultipleChoiceCard({
  card,
  className,
}: {
  card: MultipleChoiceCardType;
  className?: string;
}) {
  return <div className={className}>{card.question}</div>;
}
