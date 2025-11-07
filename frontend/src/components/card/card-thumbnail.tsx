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
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';

export function CardThumbnail({
  card,
  className,
  deckId,
  tags,
  children,
}: CardRendererProps<CardType> & { tags?: string[]; children: JSX.Element }): JSX.Element {
  const [open, setOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const router = useRouter();

  async function handleDelete() {
    const res = await deleteCard(deckId, card.id);
    if (res.success) {
      toast.success('Card deleted');
      router.refresh();
    } else {
      console.error(res.message);
      toast.error('Failed to delete card');
    }
    setDeleteDialogOpen(false);
  }

  return (
    <>
      <Card
        className={`flex h-fit max-h-[250px] min-h-[125px] w-full cursor-pointer flex-col gap-0 rounded-2xl p-2 ${className ?? ''}`}
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
            <ContextMenuItem variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
              <Trash2 />
              Delete
            </ContextMenuItem>
          </ContextMenuContent>
        </ContextMenu>
      </Card>
      <EditCardMenu open={open} onOpenChange={setOpen} deckId={deckId} card={card} />
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
            <AlertDialogAction onClick={handleDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
