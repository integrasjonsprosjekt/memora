import { JSX } from 'react';
import { match } from 'ts-pattern';
import {
  Card as CardType,
  FrontBackCard as FrontBackCardType,
  FillBlanksCard as FillBlanksCardType,
  MultipleChoiceCard as MultipleChoiceCardType,
  OrderedCard as OrderedCardType,
} from '@/types/card';
import { CardComponentProps } from './types';
import { FrontBackCard, FrontBackCardThumbnail } from './widgets/front-back-card';
import { FillBlanksCard, FillBlanksCardThumbnail } from './widgets/fill-blanks-card';
import { MultipleChoiceCard, MultipleChoiceCardThumbnail } from './widgets/multiple-choice-card';
import { OrderedCard, OrderedCardThumbnail } from './widgets/ordered-card';
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from '@/components/ui/context-menu';
import { SquarePen, Trash2 } from 'lucide-react';
import { Card } from '@/components/ui/card';
import styles from './card.module.css';
import { Badge } from '@/components/ui/badge';
import Link from 'next/link';

/**
 * Renders an interactive card.
 */
export function RenderCard({ card, className }: CardComponentProps<CardType>): JSX.Element {
  const cardComponent = match(card)
    .with({ type: 'front_back' }, () => <FrontBackCard card={card as FrontBackCardType} className={className} />)
    .with({ type: 'blanks' }, () => <FillBlanksCard card={card as FillBlanksCardType} className={className} />)
    .with({ type: 'multiple_choice' }, () => (
      <MultipleChoiceCard card={card as MultipleChoiceCardType} className={className} />
    ))
    .with({ type: 'ordered' }, () => <OrderedCard card={card as OrderedCardType} className={className} />)
    .exhaustive();

  // TODO: Should the card boilerplate be moved from CardPage to here?
  //       Or should the boilerplate be removed from RenderCardThumbnail and moved to DeckPage?
  return cardComponent;
}

/**
 * Renders a non-interactive card thumbnail.
 */
export function RenderCardThumbnail({ card, className, deckId }: CardComponentProps<CardType>): JSX.Element {
  const cardComponent = match(card)
    .with({ type: 'front_back' }, () => (
      <FrontBackCardThumbnail card={card as FrontBackCardType} className={className} />
    ))
    .with({ type: 'blanks' }, () => <FillBlanksCardThumbnail card={card as FillBlanksCardType} className={className} />)
    .with({ type: 'multiple_choice' }, () => (
      <MultipleChoiceCardThumbnail card={card as MultipleChoiceCardType} className={className} />
    ))
    .with({ type: 'ordered' }, () => <OrderedCardThumbnail card={card as OrderedCardType} className={className} />)
    .exhaustive();

  const tags = [card.type];

  return (
    <Card
      className={`${styles.card} flex h-fit max-h-[250px] min-h-[125px] w-full cursor-pointer flex-col gap-0 rounded-2xl p-2 ${className ?? ''}`}
    >
      <ContextMenu>
        <ContextMenuTrigger asChild className="flex flex-1 flex-col">
          <Link href={`/decks/${deckId}/cards/${card.id}`}>
            <div className="flex-1 overflow-y-auto">{cardComponent}</div>
            <div className="mt-auto pt-2">
              {tags.map((tag, index) => (
                <Badge key={index} variant="outline">
                  {tag}
                </Badge>
              ))}
            </div>
          </Link>
        </ContextMenuTrigger>
        <ContextMenuContent>
          <ContextMenuItem>
            <SquarePen />
            Edit
          </ContextMenuItem>
          <ContextMenuSeparator />
          <ContextMenuItem variant="destructive">
            <Trash2 />
            Delete
          </ContextMenuItem>
        </ContextMenuContent>
      </ContextMenu>
    </Card>
  );
}
