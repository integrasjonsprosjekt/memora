'use client';

import { Breadcrumb, BreadcrumbItem, BreadcrumbList } from '@/components/ui/breadcrumb';
import { Deck } from '@/types/deck';
import { DeckBreadcrumb } from './deck-breadcrumb';
import { useAuth } from '@/context/auth';
import { fetchApi } from '@/lib/api/config';
import { useEffect, useState, use } from 'react';
import { Skeleton } from '@/components/ui/skeleton';

export default function Layout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ deckId: string }>;
}>) {
  const { deckId } = use(params);
  const { user } = useAuth();
  const [deckTitle, setDeckTitle] = useState<string>(deckId);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchDeck() {
      if (!user) {
        setLoading(false);
        return;
      }

      try {
        const deck = await fetchApi<Deck>(`decks/${deckId}`, { user });
        setDeckTitle(deck.title || deckId);
      } catch (error) {
        console.error('Failed to fetch deck for breadcrumb:', error);
        setDeckTitle(deckId);
      } finally {
        setLoading(false);
      }
    }

    fetchDeck();
  }, [user, deckId]);

  return (
    <Breadcrumb>
      <BreadcrumbList>
        {loading ? (
          <BreadcrumbItem>
            <Skeleton className="h-4 w-24" />
          </BreadcrumbItem>
        ) : (
          <DeckBreadcrumb deckId={deckId} deckTitle={deckTitle} />
        )}
        {children}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
