'use client';

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

/**
 * Renders an interactive card.
 */
export function RenderCard({
  card,
  className,
  onAnswerChange,
  flipTrigger,
}: CardComponentProps<CardType> & { flipTrigger?: number }): JSX.Element {
  const cardComponent = match(card)
    .with({ type: 'front_back' }, () => (
      <FrontBackCard
        card={card as FrontBackCardType}
        className={cn(
          // Counteract padding for rulers
          '[&>hr]:-mx-10 [&>hr]:w-auto'
        )}
        onAnswerChange={onAnswerChange}
        flipTrigger={flipTrigger}
      />
    ))
    .with({ type: 'blanks' }, () => (
      <FillBlanksCard card={card as FillBlanksCardType} onAnswerChange={onAnswerChange} />
    ))
    .with({ type: 'multiple_choice' }, () => (
      <MultipleChoiceCard card={card as MultipleChoiceCardType} onAnswerChange={onAnswerChange} />
    ))
    .with({ type: 'ordered' }, () => <OrderedCard card={card as OrderedCardType} onAnswerChange={onAnswerChange} />)
    .exhaustive();

  return <Card className={cn('w-full px-10 py-5', className)}>{cardComponent}</Card>;
}

/**
 * Renders a non-interactive card thumbnail.
 */
export function RenderCardThumbnail({
  card,
  className,
  deckId,
  clickable = true,
  onSuccess,
}: CardRendererProps<CardType> & { clickable?: boolean; onSuccess?: () => void }): JSX.Element {
  const cardComponent = match(card)
    .with({ type: 'front_back' }, () => (
      // TODO: Counteract padding for rulers
      <FrontBackCardThumbnail card={card as FrontBackCardType} />
    ))
    .with({ type: 'blanks' }, () => <FillBlanksCardThumbnail card={card as FillBlanksCardType} />)
    .with({ type: 'multiple_choice' }, () => <MultipleChoiceCardThumbnail card={card as MultipleChoiceCardType} />)
    .with({ type: 'ordered' }, () => <OrderedCardThumbnail card={card as OrderedCardType} />)
    .exhaustive();

  const tags = [card.type];

  return (
    <CardThumbnail
      card={card}
      deckId={deckId}
      tags={tags}
      className={className}
      clickable={clickable}
      onSuccess={onSuccess}
    >
      <div className="overflow-hidden">{cardComponent}</div>
    </CardThumbnail>
  );
}
