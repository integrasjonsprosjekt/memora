'use client';

import { JSX, useState } from 'react';
import { FrontBackCard as FrontBackCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { ClientMarkdownRenderer } from '@/components/markdown/client-markdown-renderer';
import { cn } from '@/lib/utils';

export function FrontBackCardInteractive({ card, className }: CardComponentProps<FrontBackCardType>): JSX.Element {
  const [isFlipped, setIsFlipped] = useState(false);

  return (
    <div className={cn(className, 'cursor-pointer py-5')} onClick={() => setIsFlipped(!isFlipped)}>
      <div>
        <ClientMarkdownRenderer>{card.front}</ClientMarkdownRenderer>
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
        <ClientMarkdownRenderer>{card.back}</ClientMarkdownRenderer>
      </div>
    </div>
  );
}
