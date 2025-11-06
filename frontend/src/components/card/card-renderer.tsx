import { JSX } from 'react';
import { match } from 'ts-pattern';
import {
  Card as CardType,
  FrontBackCard as FrontBackCardType,
  FillBlanksCard as FillBlanksCardType,
  MultipleChoiceCard as MultipleChoiceCardType,
  OrderedCard as OrderedCardType,
} from '@/types/card';
import { CardComponentProps, CardRendererProps } from './types';
import { FrontBackCard, FrontBackCardThumbnail } from './widgets/front-back-card';
import { FillBlanksCard, FillBlanksCardThumbnail } from './widgets/fill-blanks-card';
import { MultipleChoiceCard, MultipleChoiceCardThumbnail } from './widgets/multiple-choice-card';
import { OrderedCard, OrderedCardThumbnail } from './widgets/ordered-card';

import { CardThumbnail } from './card-thumbnail';
import { cn } from '@/lib/utils';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

/**
 * Renders an interactive card.
 */
export function RenderCard({ card, className }: CardComponentProps<CardType>): JSX.Element {
  const cardComponent = match(card)
    .with({ type: 'front_back' }, () => <FrontBackCard card={card as FrontBackCardType} className={cn(
      className,
      // Counteract padding for rulers
      "[&>hr]:-mx-10 [&>hr]:w-auto"
    )} />)
    .with({ type: 'blanks' }, () => <FillBlanksCard card={card as FillBlanksCardType} className={className} />)
    .with({ type: 'multiple_choice' }, () => (
      <MultipleChoiceCard card={card as MultipleChoiceCardType} className={className} />
    ))
    .with({ type: 'ordered' }, () => <OrderedCard card={card as OrderedCardType} className={className} />)
    .exhaustive();

  return <div className="flex flex-1 flex-col items-center justify-center px-4 sm:px-6 lg:px-8">
    <Card className="w-full max-w-sm -translate-y-30 transform px-10 py-5 text-xl sm:max-w-md md:max-w-lg lg:max-w-xl xl:max-w-2xl">
      {cardComponent}

    </Card>
  </div>
}

/**
 * Renders a non-interactive card thumbnail.
 */
export function RenderCardThumbnail({ card, className, deckId }: CardRendererProps<CardType>): JSX.Element {
  const cardComponent = match(card)
    .with({ type: 'front_back' }, () => (
      <FrontBackCardThumbnail card={card as FrontBackCardType} className={cn(
          className,
          // TODO: Counteract padding for rulers
        )}
      />
    ))
    .with({ type: 'blanks' }, () => <FillBlanksCardThumbnail card={card as FillBlanksCardType} className={className} />)
    .with({ type: 'multiple_choice' }, () => (
      <MultipleChoiceCardThumbnail card={card as MultipleChoiceCardType} className={className} />
    ))
    .with({ type: 'ordered' }, () => <OrderedCardThumbnail card={card as OrderedCardType} className={className} />)
    .exhaustive();

  const tags = [card.type];

  return (
    <CardThumbnail card={card} deckId={deckId} tags={tags}>
      {cardComponent}
    </CardThumbnail>
  );
}
