'use client';

import { RenderCardThumbnail } from '@/components/card';
import { getApiEndpoint } from '@/config/api';
import { Button } from '@/components/ui/button';
import { SquarePen, Trash2 } from 'lucide-react';
import { useState, useEffect } from 'react';
import { EditCardMenu } from '@/components/edit-card-menu';
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
import { deleteCard } from '@/app/api';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { Card as CardType } from '@/types/card';

export default function CardPage({ params }: { params: Promise<{ deckId: string; cardId: string }> }) {
  const [deckId, setDeckId] = useState<string>('');
  const [cardId, setCardId] = useState<string>('');
  const [card, setCard] = useState<CardType | null>(null);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const router = useRouter();

  useEffect(() => {
    params.then(({ deckId, cardId }) => {
      setDeckId(deckId);
      setCardId(cardId);

      fetch(getApiEndpoint(`/v1/decks/${deckId}/cards/${cardId}`), {
        cache: 'no-store',
      })
        .then((res) => res.json())
        .then((data) => setCard(data));
    });
  }, [params]);

  async function handleDelete() {
    const res = await deleteCard(deckId, cardId);
    if (res.success) {
      toast.success('Card deleted', { icon: <Trash2 size={16} /> });
      router.push(`/decks/${deckId}`);
    } else {
      console.error(res.message);
      toast.error('Failed to delete card');
    }
    setDeleteDialogOpen(false);
  }

  if (!card) {
    return null;
  }

  return (
    <>
      <div className="flex flex-1 flex-col items-center justify-between px-4 pb-8 sm:px-6 lg:px-8">
        <div className="flex w-full flex-1 items-center justify-center">
          <RenderCardThumbnail
            key={card.id}
            card={card}
            deckId={deckId}
            clickable={false}
            className="max-w-sm -translate-y-30 transform px-10 py-5 text-xl sm:max-w-md md:max-w-lg lg:max-w-xl xl:max-w-2xl"
          />
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => setEditOpen(true)}>
            <SquarePen />
            Edit
          </Button>
          <Button variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
            <Trash2 />
            Delete
          </Button>
        </div>
      </div>
      <EditCardMenu open={editOpen} onOpenChange={setEditOpen} deckId={deckId} card={card} />
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
