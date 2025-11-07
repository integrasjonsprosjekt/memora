'use client';

import { JSX, useState, useEffect } from 'react';
import { FrontBackCard as FrontBackCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { ClientMarkdownRenderer } from '@/components/markdown/client-markdown-renderer';
import { cn } from '@/lib/utils';

export function FrontBackCardInteractive({ card, className, onAnswerChange }: CardComponentProps<FrontBackCardType>): JSX.Element {
  const [isFlipped, setIsFlipped] = useState(false);

  useEffect(() => {
    // Notify parent when card is flipped (considered "answered")
    if (onAnswerChange) {
      onAnswerChange(isFlipped);
    }
  }, [isFlipped, onAnswerChange]);

  return (
    <div
      className={cn(
        className,
        "cursor-pointer py-5"
      )}
      onClick={() => setIsFlipped(!isFlipped)}
    >
      <div>
        <ClientMarkdownRenderer>{card.front}</ClientMarkdownRenderer>
      </div>

      <hr
        className={cn(
          "border-border tap-highlight-transparent w-full border-t border-dashed transition-all duration-300",
          isFlipped ? "opacity-100 my-10" : "opacity-0 my-0"
        )}
      />

      <div
        className={cn(
          "transition-all duration-300 origin-top",
          isFlipped ? "opacity-100 scale-y-100" : "opacity-0 scale-y-0 h-0"
        )}
      >
        <ClientMarkdownRenderer>{card.back}</ClientMarkdownRenderer>
      </div>
    </div>
  );
}
