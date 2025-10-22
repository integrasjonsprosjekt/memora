import { JSX } from 'react';
import { FillBlanksCard as FillBlanksCardType } from '@/types/card';
import { CardComponentProps } from '../types';

export function FillBlanksCard({ card, className }: CardComponentProps<FillBlanksCardType>): JSX.Element {
  return <div className={className}>{card.question}</div>;
}

export function FillBlanksCardThumbnail({ card, className }: CardComponentProps<FillBlanksCardType>): JSX.Element {
  const parts = card.question.split('{}');
  return (
    <div className={className}>
      {parts.map((part, index) => (
        <span key={index}>
          {/* Render part */}
          {part}
          {/* Check if there is a corresponding answer */}
          {index < parts.length - 1 && index < card.answers.length && (
            <span className="bg-accent rounded-sm border border-dashed border-[var(--border)] px-1 py-[0.02rem]">
              {card.answers[index]}
            </span>
          )}
        </span>
      ))}
    </div>
  );
}
