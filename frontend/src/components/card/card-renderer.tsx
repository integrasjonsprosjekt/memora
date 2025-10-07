import { Card, createCardFromJson } from './types';
import { JSX } from 'react';
import CardWrapper from './card-wrapper';

// TODO: Rename these
export function CardRenderer({ card, className }: { card: Card; className?: string }): JSX.Element {
  const cardInstance = createCardFromJson(card);

  return cardInstance.render(className);
}

export function CardDisplay({ card }: { card: Card }): JSX.Element {
  const cardInstance = createCardFromJson(card);

  return (
    <CardWrapper tags={[card.type]} id={card.id}>
      {cardInstance.display()}
    </CardWrapper>
  );
}
