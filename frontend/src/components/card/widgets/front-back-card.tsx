import { JSX } from 'react';
import { FrontBackCard as FrontBackCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { FrontBackCardInteractive } from './front-back-card-interactive';

export function FrontBackCard({ card, className }: CardComponentProps<FrontBackCardType>): JSX.Element {
  return <FrontBackCardInteractive card={card} className={className} />;
}

export function FrontBackCardThumbnail({ card, className }: CardComponentProps<FrontBackCardType>): JSX.Element {
  return (
    <div className={className}>
      <p>{card.front}</p>

      <hr className="border-border tap-highlight-transparent my-2 w-full border-t border-dashed" />

      <p>{card.back.length > 100 ? card.back.substring(0, 100) + '...' : card.back}</p>
    </div>
  );
}
