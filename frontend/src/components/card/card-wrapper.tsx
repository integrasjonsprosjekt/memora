import { JSX, ReactNode } from 'react';
import { Card } from '@/components/ui/card';
import styles from './card.module.css';
import { Badge } from '@/components/ui/badge';

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
  className,
  tags,
}: {
  children: ReactNode;
  className?: string;
  tags?: string[];
}): JSX.Element {
  return (
    <Card className={`${styles.card} rounded-2xl ${className ?? ''}`}>
      {children}
      {tags?.map((tag, index) => (
        <Badge key={index} variant="outline">
          {tag}
        </Badge>
      ))}
    </Card>
  );
}
