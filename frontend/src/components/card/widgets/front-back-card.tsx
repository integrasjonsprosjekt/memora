import { JSX } from 'react';
import { FrontBackCard as FrontBackCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { FrontBackCardInteractive } from './front-back-card-interactive';
import { MarkdownRenderer } from '@/components/markdown';

export function FrontBackCard({
  card,
  className,
  onAnswerChange,
  flipTrigger,
}: CardComponentProps<FrontBackCardType> & { flipTrigger?: number }): JSX.Element {
  return (
    <FrontBackCardInteractive
      card={card}
      className={className}
      onAnswerChange={onAnswerChange}
      flipTrigger={flipTrigger}
    />
  );
}

export function FrontBackCardThumbnail({ card, className }: CardComponentProps<FrontBackCardType>): JSX.Element {
  return (
    <div className={className}>
      <div>
        <MarkdownRenderer>{card.front}</MarkdownRenderer>
      </div>

      <hr className="border-border tap-highlight-transparent my-2 w-full border-t border-dashed" />

      {/*<p>{card.back.length > 100 ? card.back.substring(0, 100) + '...' : card.back}</p>*/}
      <div>
        <MarkdownRenderer>{card.back}</MarkdownRenderer>
      </div>
    </div>
  );
}
