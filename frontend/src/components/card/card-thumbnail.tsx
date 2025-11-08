'use client';

import {
  ContextMenu,
  ContextMenuTrigger,
  ContextMenuContent,
  ContextMenuSeparator,
  ContextMenuItem,
} from '@/components/ui/context-menu';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { SquarePen, Trash2 } from 'lucide-react';
import { JSX, useState } from 'react';
import { Card as CardType } from '@/types/card';
import { CardRendererProps } from './types';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import Link from 'next/link';
import { EditCardMenu } from '@/components/edit-card-menu';
import { deleteCard } from '@/app/api';
import { toast } from 'sonner';
import { cn } from '@/lib/utils';
import { useAuth } from '@/context/auth';

export function CardThumbnail({
  card,
  className,
  deckId,
  tags,
  children,
  clickable = true,
  onSuccess,
}: CardRendererProps<CardType> & {
  tags?: string[];
  clickable?: boolean;
  children: JSX.Element;
  onSuccess?: () => void;
}): JSX.Element {
  const [open, setOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const { user } = useAuth();

  async function handleDelete() {
    if (!user) {
      toast.error('You must be signed in to delete cards');
      return;
    }

    const res = await deleteCard(user, deckId, card.id);
    if (res.success) {
      toast.success('Card deleted', { icon: <Trash2 size={16} /> });
      onSuccess?.();
    } else {
      console.error(res.message);
      toast.error('Failed to delete card');
    }
    setDeleteDialogOpen(false);
  }

  const cardContent = (
    <>
      <div className="flex-1 overflow-y-auto">{children}</div>
      <div className="mt-auto pt-2">
        {tags?.map((tag, index) => (
          <Badge key={index} variant="outline">
            {tag}
          </Badge>
        ))}
      </div>
    </>
  );
  const cardContentClassName = 'flex flex-1 flex-col';

  return (
    <>
      <Card
        className={cn(
          'flex h-fit min-h-[125px] w-full flex-col gap-0 rounded-2xl p-2',
          clickable ? 'cursor-pointer' : '',
          className
        )}
      >
        <ContextMenu>
          <ContextMenuTrigger asChild>
            {clickable ? (
              <Link href={`/decks/${deckId}/cards/${card.id}`} className={cardContentClassName}>
                {cardContent}
              </Link>
            ) : (
              <div className={cardContentClassName}>{cardContent}</div>
            )}
          </ContextMenuTrigger>
          <ContextMenuContent>
            <ContextMenuItem onClick={() => setOpen(true)}>
              <SquarePen />
              Edit
            </ContextMenuItem>
            <ContextMenuSeparator />
            <ContextMenuItem variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
              <Trash2 />
              Delete
            </ContextMenuItem>
          </ContextMenuContent>
        </ContextMenu>
      </Card>
      <EditCardMenu open={open} onOpenChange={setOpen} deckId={deckId} card={card} onSuccess={onSuccess} />
      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the card.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
