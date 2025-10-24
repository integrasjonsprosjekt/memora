'use client';

import {
  ContextMenu,
  ContextMenuTrigger,
  ContextMenuContent,
  ContextMenuSeparator,
  ContextMenuItem,
} from '@/components/ui/context-menu';
import { SquarePen, Trash2 } from 'lucide-react';
import { JSX, useState } from 'react';
import { Card as CardType } from '@/types/card';
import { CardRendererProps } from './types';
import styles from './card.module.css';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import Link from 'next/link';
import { EditCardMenu } from '@/components/edit-card-menu';

interface CardButtonProps {
  deckId: string;
  cardId: string;
  cardType: CardType;
  initialData: Partial<CardType>;
}

export function CardThumbnail({
  card,
  className,
  deckId,
  tags,
  children,
}: CardRendererProps<CardType> & { tags?: string[]; children: JSX.Element }): JSX.Element {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Card
        className={`${styles.card} flex h-fit max-h-[250px] min-h-[125px] w-full cursor-pointer flex-col gap-0 rounded-2xl p-2 ${className ?? ''}`}
      >
        <ContextMenu>
          <ContextMenuTrigger asChild className="flex flex-1 flex-col">
            <Link href={`/decks/${deckId}/cards/${card.id}`}>
              <div className="flex-1 overflow-y-auto">{children}</div>
              <div className="mt-auto pt-2">
                {tags?.map((tag, index) => (
                  <Badge key={index} variant="outline">
                    {tag}
                  </Badge>
                ))}
              </div>
            </Link>
          </ContextMenuTrigger>
          <ContextMenuContent>
            <ContextMenuItem onClick={() => setOpen(true)}>
              <SquarePen />
              Edit
            </ContextMenuItem>
            <ContextMenuSeparator />
            <ContextMenuItem variant="destructive">
              <Trash2 />
              Delete
            </ContextMenuItem>
          </ContextMenuContent>
        </ContextMenu>
      </Card>
      <EditCardMenu open={open} onOpenChange={setOpen} deckId={deckId ? deckId : ''} card={card} />
    </>
  );
}
