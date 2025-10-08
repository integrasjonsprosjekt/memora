'use client';

import { JSX, ReactNode } from 'react';
import { Card } from '@/components/ui/card';
import styles from './card.module.css';
import { Badge } from '@/components/ui/badge';
import React from 'react';
import { useRouter } from 'next/navigation';

/**
 * CardWrapper component - Base wrapper used for all cards in the application.
 * Provides consistent styling and structure for card components.
 *
 * @param {Object} props - Component props
 * @param {ReactNode} props.children - Child elements to be rendered inside the card
 * @param {string} props.className - Optional additional CSS classes to apply to the card
 */
export default function CardWrapper({
  children,
  id,
  className,
  tags,
}: {
  children: ReactNode;
  id: string;
  className?: string;
  tags?: string[];
}): JSX.Element {
  const router = useRouter();

  return (
    <Card
      className={`${styles.card} flex h-fit max-h-[250px] min-h-[125px] w-full cursor-pointer flex-col gap-0 rounded-2xl p-2 ${className ?? ''}`}
      onClick={() => router.push(`/cards/${id}`)}
    >
      <div className="flex-1 overflow-y-auto">{children}</div>
      <div className="mt-auto pt-2">
        {tags?.map((tag, index) => (
          <Badge key={index} variant="outline">
            {tag}
          </Badge>
        ))}
      </div>
    </Card>
  );
}
