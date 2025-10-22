'use client';

import { JSX, useState } from 'react';
import styles from '../card.module.css';
import { FrontBackCard as FrontBackCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { ClientMarkdownRenderer } from '@/components/markdown/client-markdown-renderer';

export function FrontBackCard({ card, className }: CardComponentProps<FrontBackCardType>): JSX.Element {
  const [isFlipped, setIsFlipped] = useState(false);

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

export function FrontBackCardThumbnail({ card, className }: CardComponentProps<FrontBackCardType>): JSX.Element {
  return (
    <div className={className}>
      <p>{card.front}</p>

      <hr className="border-border tap-highlight-transparent my-2 w-full border-t border-dashed" />

      <p>{card.back.length > 100 ? card.back.substring(0, 100) + '...' : card.back}</p>
    </div>
  );
}
