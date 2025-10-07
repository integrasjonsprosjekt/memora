'use client';

import { useState, ReactNode } from 'react';
import styles from '../card.module.css';

interface FlipWrapperProps {
  children: ReactNode;
  className?: string;
}

export default function FrontBackCard({ children, className }: FlipWrapperProps) {
  const [isFlipped, setIsFlipped] = useState(false);

  return (
    <div
      className={`${styles.front_back} ${isFlipped ? styles.front_back_flipped : ''} ${className}`}
      onClick={() => setIsFlipped(!isFlipped)}
    >
      {children}
    </div>
  );
}
