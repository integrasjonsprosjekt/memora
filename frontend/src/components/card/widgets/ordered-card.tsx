import { JSX } from 'react';
import { OrderedCard as OrderedCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { OrderedCardInteractive } from './ordered-card-interactive';

export function OrderedCard({ card, className, onAnswerChange }: CardComponentProps<OrderedCardType>): JSX.Element {
  return <OrderedCardInteractive card={card} className={className} onAnswerChange={onAnswerChange} />;
}

export function OrderedCardThumbnail({ card, className }: CardComponentProps<OrderedCardType>): JSX.Element {
  return (
    <div className={className}>
      {card.question && <p className="pb-2 font-bold">{card.question}</p>}
      <ol className="marker:text-muted-foreground/50 list-decimal pl-5 marker:text-xs">
        {card.options.map((key) => (
          <li key={key}>{key}</li>
        ))}
      </ol>
    </div>
  );
}
