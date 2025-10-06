import { FrontBackCardType as FrontBackCardType } from '../types';
import { MarkdownRenderer } from '@/components/markdown';
import CardWrapper from '../card-wrapper';

export function FrontBackCard({ card }: { card: FrontBackCardType }) {
  return (
    <CardWrapper>
      <div className="front">
        <MarkdownRenderer>{card.front}</MarkdownRenderer>
      </div>

      <hr className={`border-border tap-highlight-transparent w-full border-t border-dashed`} />

      <div className="back">
        <MarkdownRenderer>{card.back}</MarkdownRenderer>
      </div>
    </CardWrapper>
  );
}
