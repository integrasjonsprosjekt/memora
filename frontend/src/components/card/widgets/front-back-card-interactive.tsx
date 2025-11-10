'use client';

import { JSX, useState, useEffect } from 'react';
import { FrontBackCard as FrontBackCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { MarkdownRenderer } from '@/components/markdown/markdown-renderer';
import { cn } from '@/lib/utils';

export function FrontBackCardInteractive({
  card,
  className,
  onAnswerChange,
  flipTrigger,
}: CardComponentProps<FrontBackCardType> & { flipTrigger?: number }): JSX.Element {
  const [isFlipped, setIsFlipped] = useState(false);

  useEffect(() => {
    if (flipTrigger && flipTrigger > 0) {
      setIsFlipped(true);
    }
  }, [flipTrigger]);

  useEffect(() => {
    if (onAnswerChange) {
      onAnswerChange(isFlipped);
    }
  }, [isFlipped, onAnswerChange]);

  return (
    <div className={cn(className, 'cursor-pointer py-5')} onClick={() => setIsFlipped(true)}>
      <div>
        <MarkdownRenderer>{card.front}</MarkdownRenderer>
      </div>

      <hr
        className={cn(
          'border-border tap-highlight-transparent w-full border-t border-dashed transition-all duration-300',
          isFlipped ? 'my-10 opacity-100' : 'my-0 opacity-0'
        )}
      />

      <div
        className={cn(
          'origin-top transition-all duration-300',
          isFlipped ? 'scale-y-100 opacity-100' : 'h-0 scale-y-0 opacity-0'
        )}
      >
        <MarkdownRenderer>{card.back}</MarkdownRenderer>
      </div>
    </div>
  );
}
