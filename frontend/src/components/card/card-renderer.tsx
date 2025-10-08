import { JSX } from 'react';
import { match } from 'ts-pattern';
import {
  Card,
  CardComponentProps,
  FrontBackCard as FrontBackCardType,
  FillBlanksCard as FillBlanksCardType,
  MultipleChoiceCard as MultipleChoiceCardType,
  OrderedCard as OrderedCardType,
} from './types';
import CardWrapper from './card-wrapper';
import { FrontBackCard, FrontBackCardThumbnail } from './widgets/front-back-card';
import { FillBlanksCard, FillBlanksCardThumbnail } from './widgets/fill-blanks-card';
import { MultipleChoiceCard, MultipleChoiceCardThumbnail } from './widgets/multiple-choice-card';
import { OrderedCard, OrderedCardThumbnail } from './widgets/ordered-card';

/**
 * Renders an interactive card.
 */
export function RenderCard({ card, className }: CardComponentProps<Card>): JSX.Element {
  const cardComponent = match(card)
    .with({ type: 'front_back' }, () => (
      <FrontBackCard card={card as FrontBackCardType} className={className} />
    ))
    .with({ type: 'blanks' }, () => (
      <FillBlanksCard card={card as FillBlanksCardType} className={className} />
    ))
    .with({ type: 'multiple_choice' }, () => (
      <MultipleChoiceCard card={card as MultipleChoiceCardType} className={className} />
    ))
    .with({ type: 'ordered' }, () => (
      <OrderedCard card={card as OrderedCardType} className={className} />
    ))
    .exhaustive();

  // TODO: Should the card boilerplate be moved from CardPage to here?
  //       Or should the boilerplate be removed from RenderCardThumbnail and moved to DeckPage?
  return cardComponent;
}

/**
 * Renders a non-interactive card thumbnail.
 */
export function RenderCardThumbnail({ card, className }: CardComponentProps<Card>): JSX.Element {
  const cardComponent = match(card)
    .with({ type: 'front_back' }, () => (
      <FrontBackCardThumbnail card={card as FrontBackCardType} className={className} />
    ))
    .with({ type: 'blanks' }, () => (
      <FillBlanksCardThumbnail card={card as FillBlanksCardType} className={className} />
    ))
    .with({ type: 'multiple_choice' }, () => (
      <MultipleChoiceCardThumbnail card={card as MultipleChoiceCardType} className={className} />
    ))
    .with({ type: 'ordered' }, () => (
      <OrderedCardThumbnail card={card as OrderedCardType} className={className} />
    ))
    .exhaustive();

  return (
    <CardWrapper tags={[card.type]} id={card.id}>
      {cardComponent}
    </CardWrapper>
  );
}
