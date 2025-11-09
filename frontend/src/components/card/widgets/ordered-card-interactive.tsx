'use client';

import { JSX, useState, useEffect, useRef } from 'react';
import { OrderedCard as OrderedCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { GripVertical } from 'lucide-react';
import { Skeleton } from '@/components/ui/skeleton';

export function OrderedCardInteractive({
  card,
  className,
  onAnswerChange,
}: CardComponentProps<OrderedCardType>): JSX.Element {
  const [items, setItems] = useState<string[] | null>(null);
  const [draggedItem, setDraggedItem] = useState<string | null>(null);
  const hasShuffled = useRef(false);

  useEffect(() => {
    // Shuffle on mount to avoid hydration mismatch
    if (!hasShuffled.current) {
      hasShuffled.current = true;
      const shuffled = [...card.options];
      for (let i = shuffled.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
      }
      setItems(shuffled);
    }
  }, [card.options]);

  useEffect(() => {
    if (onAnswerChange && items) {
      onAnswerChange(items);
    }
  }, [items, onAnswerChange]);

  useEffect(() => {
    const handleMouseUp = () => {
      if (draggedItem !== null) {
        setDraggedItem(null);
      }
    };

    document.addEventListener('mouseup', handleMouseUp);
    return () => {
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [draggedItem]);

  const handleDragStart = (item: string) => {
    setDraggedItem(item);
  };

  const handleDragOver = (e: React.DragEvent, index: number) => {
    e.preventDefault();
    if (draggedItem === null || items === null) return;

    const draggedIndex = items.indexOf(draggedItem);
    if (draggedIndex === index) return;

    const newItems = [...items];
    newItems.splice(draggedIndex, 1);
    newItems.splice(index, 0, draggedItem);

    setItems(newItems);
  };

  const handleDragEnd = () => {
    setDraggedItem(null);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDraggedItem(null);
  };

  if (!items) {
    return (
      <div className="space-y-2">
        <Skeleton className="h-4 w-[250px]" />
        <Skeleton className="h-4 w-[120px]" />
        <Skeleton className="h-4 w-[150px]" />
      </div>
    );
  }

  return (
    <div className={className}>
      {card.question && <p className="pb-2 font-bold">{card.question}</p>}
      <div className="space-y-2" onDrop={handleDrop}>
        {items.map((item, index) => (
          <div key={`container-${index}`} className="relative flex items-center gap-3">
            {/* Static number outside */}
            <div className="text-muted-foreground w-4 text-sm font-semibold">{index + 1}</div>

            {/* Container with background */}
            <div className="relative flex-1">
              {/* Static background */}
              <div className="border-green bg-accent/40 pointer-events-none absolute inset-0 rounded-xl border border-dashed"></div>

              {/* Draggable item */}
              <div
                draggable
                onDragStart={() => handleDragStart(item)}
                onDragOver={(e) => handleDragOver(e, index)}
                onDragEnd={handleDragEnd}
                className={`hover:bg-muted/50 relative flex cursor-move items-center gap-2 rounded p-2 transition-colors select-none ${
                  draggedItem === item ? 'opacity-50' : ''
                }`}
              >
                <GripVertical size={15} className="text-muted-foreground" />
                <span>{item}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
