import { JSX, Fragment } from 'react';
import { FillBlanksCard as FillBlanksCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { FillBlanksCardInteractive } from './fill-blanks-card-interactive';
import { MarkdownRenderer } from '@/components/markdown';

export function FillBlanksCard({
  card,
  className,
  onAnswerChange,
}: CardComponentProps<FillBlanksCardType>): JSX.Element {
  return <FillBlanksCardInteractive card={card} className={className} onAnswerChange={onAnswerChange} />;
}

export function FillBlanksCardThumbnail({ card, className }: CardComponentProps<FillBlanksCardType>): JSX.Element {
  const parts = card.question.split('{}');

  return (
    <div className={className}>
      {parts.map((part, index) => {
        const hasTrailingSpace = part.endsWith(' ');
        const nextPart = parts[index + 1];
        const hasLeadingSpace = nextPart?.startsWith(' ');

        return (
          <Fragment key={index}>
            {/* Render part */}
            <MarkdownRenderer inline>{part}</MarkdownRenderer>
            {/* Check if there is a corresponding answer */}
            {index < parts.length - 1 && index < card.answers.length && (
              <>
                {hasTrailingSpace && ' '}
                <span className="bg-accent rounded-sm border border-dashed border-[var(--border)] px-1 py-[0.02rem]">
                  <MarkdownRenderer inline>{card.answers[index]}</MarkdownRenderer>
                </span>
                {hasLeadingSpace && ' '}
              </>
            )}
          </Fragment>
        );
      })}
    </div>
  );
}
