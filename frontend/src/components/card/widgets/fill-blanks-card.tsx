import { FillBlanksCardType } from '../types';

export default function FillBlanksCard({
  card,
  className,
}: {
  card: FillBlanksCardType;
  className?: string;
}) {
  return <div className={className}>{card.question}</div>;
}
