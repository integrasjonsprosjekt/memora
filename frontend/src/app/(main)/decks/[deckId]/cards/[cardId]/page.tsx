'use client';

import { RenderCardThumbnail } from '@/components/card';
import { Button } from '@/components/ui/button';
import { SquarePen, Trash2 } from 'lucide-react';
import { useState, useEffect, useCallback } from 'react';
import { use } from 'react';
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
import { useAuth } from '@/context/auth';
import { fetchApi } from '@/lib/api/config';
import { Skeleton } from '@/components/ui/skeleton';

export default function CardPage({ params }: { params: Promise<{ deckId: string; cardId: string }> }) {
  const { deckId, cardId } = use(params);
  const { user } = useAuth();
  const [card, setCard] = useState<CardType | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [refreshKey, setRefreshKey] = useState(0);
  const router = useRouter();

  const fetchCard = useCallback(
    async (showLoading = false) => {
      if (!user) return;

      try {
        if (showLoading) {
          setLoading(true);
        }
        setError(null);

        const cardData = await fetchApi<CardType>(`decks/${deckId}/cards/${cardId}`, { user });
        setCard(cardData);
        setRefreshKey((prev) => prev + 1);
      } catch (error) {
        console.error(error);
        setError('Failed to load card');
        toast.error('Failed to load card');
      } finally {
        if (showLoading) {
          setLoading(false);
        }
      }
    },
    [user, deckId, cardId]
  );

  useEffect(() => {
    fetchCard(true);
  }, [fetchCard]);

  async function handleDelete() {
    if (!user) return;

    const res = await deleteCard(user, deckId, cardId);
    if (res.success) {
      toast.success('Card deleted', { icon: <Trash2 size={16} /> });
      router.push(`/decks/${deckId}`);
    } else {
      console.error(res.message);
      toast.error('Failed to delete card');
    }
    setDeleteDialogOpen(false);
  }

  if (loading) {
    return (
      <div className="flex h-full flex-col px-4 pb-8 sm:px-6 lg:px-8">
        <div className="flex min-h-0 flex-1 items-center justify-center overflow-auto pt-8 pb-24">
          <Skeleton className="h-[400px] w-full max-w-sm rounded-2xl sm:max-w-md md:max-w-lg lg:max-w-xl xl:max-w-2xl" />
        </div>
        <div className="flex justify-center gap-2 pt-8">
          <Skeleton className="h-10 w-20 rounded-md" />
          <Skeleton className="h-10 w-24 rounded-md" />
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="p-8">
        <h1 className="text-3xl font-bold">Authentication Required</h1>
        <p className="mt-4 text-[var(--muted-foreground)]">Please sign in to view this card.</p>
      </div>
    );
  }

  if (error || !card) {
    return (
      <div className="p-8">
        <h1 className="text-3xl font-bold">Error loading card</h1>
        <p className="mt-4 text-[var(--muted-foreground)]">{error || 'There was a problem loading this card.'}</p>
      </div>
    );
  }

  return (
    <>
      <div className="flex h-full flex-col px-4 pb-8 sm:px-6 lg:px-8">
        <div className="flex min-h-0 flex-1 items-center justify-center overflow-auto pt-8 pb-24">
          <RenderCardThumbnail
            key={`${card.id}-${refreshKey}`}
            card={card}
            deckId={deckId}
            clickable={false}
            className="max-w-sm px-10 py-5 text-xl sm:max-w-md md:max-w-lg lg:max-w-xl xl:max-w-2xl"
          />
        </div>
        <div className="flex justify-center gap-2 pt-8">
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
      <EditCardMenu open={editOpen} onOpenChange={setEditOpen} deckId={deckId} card={card} onSuccess={fetchCard} />
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
