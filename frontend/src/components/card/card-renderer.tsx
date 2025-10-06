import {
  Card,
  FrontBackCardType as FrontBackCardType,
  FillBlanksCardType as FillBlanksCardType,
  createCardFromJson,
} from './types';
import { FrontBackCard } from './impl/front-back-card';
import { FillBlanksCard } from './impl/blanks-card';
import { JSX } from 'react';
import CardWrapper from './card-wrapper';

// TODO: This should use the render method of the card
export function CardRenderer({ card }: { card: Card }): JSX.Element {
  switch (card.type) {
    case 'front_back': {
      return <FrontBackCard card={card as FrontBackCardType} />;
    }
    case 'blanks': {
      return <FillBlanksCard card={card as FillBlanksCardType} />;
    }
    case 'multiple_choice':
      return <div>Work in progress</div>;
    case 'ordered':
      return <div>Work in progress</div>;
    default:
      // TODO: Handle unknown card types
      return <div>Unknown card type</div>;
  }
}

export function CardDisplay({ card }: { card: Card }): JSX.Element {
  const cardInstance = createCardFromJson(card);
  const cardDisplay = cardInstance.display();

  return (
    <CardWrapper tags={[card.type]}>
      {cardDisplay.title}
      {cardDisplay.body}
    </CardWrapper>
  );
}
