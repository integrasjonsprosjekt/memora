'use client';

import { JSX, useState, useEffect } from 'react';
import { FrontBackCard as FrontBackCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { ClientMarkdownRenderer } from '@/components/markdown/client-markdown-renderer';

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
      className={`${styles.front_back} ${isFlipped ? styles.front_back_flipped : ''} ${className}`}
      onClick={() => setIsFlipped(!isFlipped)}
    >
      <div className={styles.front}>
        <ClientMarkdownRenderer>{card.front}</ClientMarkdownRenderer>
      </div>

      <hr className="border-border tap-highlight-transparent my-5 w-full border-t border-dashed" />

      <div className={styles.back}>
        <ClientMarkdownRenderer>{card.back}</ClientMarkdownRenderer>
      </div>
    </div>
  );
}
