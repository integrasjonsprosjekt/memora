'use client';

import { RenderCardThumbnail } from '@/components/card/card-renderer';
import { Card } from '@/types/card';
import { Deck, deckDefaults } from '@/types/deck';
import { AddCardButton } from '@/components/add-card-button';
import { DeckLayout } from '@/components/deck-layout';
import { withDefaults } from '@/lib/utils';
import { toast } from 'sonner';
import { useAuth } from '@/context/auth';
import { fetchApi } from '@/lib/api/config';
import { useEffect, useState, useCallback } from 'react';
import { use } from 'react';
import { Skeleton } from '@/components/ui/skeleton';

export default function DeckPage({ params }: { params: Promise<{ deckId: string }> }) {
  const { deckId } = use(params);
  const { user } = useAuth();
  const [deck, setDeck] = useState<Deck | null>(null);
  const [cards, setCards] = useState<Card[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [hasError, setHasError] = useState(false);

  const fetchDeck = useCallback(
    async (showLoading = false) => {
      if (!user) return;

      try {
        if (showLoading) {
          setLoading(true);
        }
        setHasError(false);

        const [deckData, cardsResponse] = await Promise.all([
          fetchApi<Deck>(`decks/${deckId}`, { user }),
          fetchApi<{ cards: Card[] }>(`decks/${deckId}/cards`, { user }),
        ]);

        setDeck(withDefaults(deckData, deckDefaults));
        setCards(cardsResponse.cards || []);
      } catch (error) {
        console.error(error);
        setHasError(true);
        toast.error('Failed to load deck');
      } finally {
        if (showLoading) {
          setLoading(false);
        }
      }
    },
    [user, deckId]
  );

  useEffect(() => {
    fetchDeck(true);
  }, [fetchDeck]);

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-4 sm:px-6 lg:px-8">
        <header className="mb-8 lg:mb-12">
          <Skeleton className="h-9 w-64" />
          <Skeleton className="mt-2 h-6 w-32" />
        </header>
        <DeckLayout>
          <Skeleton className="h-[250px] w-full rounded-2xl" />
          <Skeleton className="h-[250px] w-full rounded-2xl" />
          <Skeleton className="h-[250px] w-full rounded-2xl" />
        </DeckLayout>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="p-8">
        <h1 className="text-3xl font-bold">Authentication Required</h1>
        <p className="mt-4 text-[var(--muted-foreground)]">Please sign in to view this deck.</p>
      </div>
    );
  }

  if (hasError || !deck) {
    return (
      <div className="p-8">
        <h1 className="text-3xl font-bold">Error loading deck</h1>
        <p className="mt-4 text-[var(--muted-foreground)]">There was a problem loading this deck.</p>
      </div>
    );
  }

  const cardsList = cards || [];

  return (
    <div className="container mx-auto px-4 py-4 sm:px-6 lg:px-8">
      <header className="mb-8 lg:mb-12">
        <h1 className="text-2xl font-bold sm:text-3xl">{deck.title}</h1>
        <p className="mt-1 text-lg text-[var(--muted-foreground)]">
          {cardsList.length} card{cardsList.length !== 1 ? 's' : ''}
        </p>
      </header>

      <DeckLayout>
        <AddCardButton deckId={deckId} onSuccess={fetchDeck} />
        {cardsList.map((card) => (
          <RenderCardThumbnail
            key={card.id}
            card={card}
            deckId={deckId}
            className="max-h-[250px]"
            onSuccess={fetchDeck}
          />
        ))}
      </DeckLayout>
    </div>
  );
}
