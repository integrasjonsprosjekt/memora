import { OrderedCardType } from '../types';

export default function OrderedCard({
  card,
  className,
}: {
  card: OrderedCardType;
  className?: string;
}) {
  return <div className={className}>{card.question}</div>;
}
