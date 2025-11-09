'use client';

import { RenderCard } from '@/components/card';
import { use, useState, useEffect, useCallback } from 'react';
import { useAuth } from '@/context/auth';
import { Card as CardType } from '@/types/card';
import { fetchApi } from '@/lib/api/config';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { toast } from 'sonner';
import { ArrowRight, ArrowLeft, SearchCheck, PartyPopper } from 'lucide-react';
import { useRouter } from 'next/navigation';

interface DueCardsResponse {
  cards: CardType[];
  next_cursor: string;
  has_more: boolean;
}

export default function Page({ params }: { params: Promise<{ deckId: string }> }) {
  const { deckId } = use(params);
  const { user } = useAuth();
  const router = useRouter();

  const [cards, setCards] = useState<CardType[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);
  const [isChecked, setIsChecked] = useState(false);
  const [cursor, setCursor] = useState<string>('');
  const [hasMore, setHasMore] = useState(false);
  const [fetchingMore, setFetchingMore] = useState(false);

  const fetchDueCards = useCallback(
    async (cursorParam: string = '', isInitial: boolean = false) => {
      if (!user) return;

      try {
        if (isInitial) {
          setLoading(true);
        } else {
          setFetchingMore(true);
        }
        setError(null);

        const limit = 10;
        const url = `decks/${deckId}/cards/due?limit=${limit}${cursorParam ? `&cursor=${cursorParam}` : ''}`;
        const response = await fetchApi<DueCardsResponse>(url, { user });

        if (isInitial) {
          setCards(response.cards);
        } else {
          setCards((prev) => [...prev, ...response.cards]);
        }
        setCursor(response.next_cursor);
        setHasMore(response.has_more);
      } catch (error) {
        console.error(error);
        setError('Failed to load due cards');
        toast.error('Failed to load due cards');
      } finally {
        if (isInitial) {
          setLoading(false);
        } else {
          setFetchingMore(false);
        }
      }
    },
    [user, deckId]
  );

  useEffect(() => {
    fetchDueCards('', true);
  }, [fetchDueCards]);

  useEffect(() => {
    // Reset checked state when moving to a new card
    setIsChecked(false);
    setRefreshKey((prev) => prev + 1);

    // If we're approaching the end and there are more cards, fetch the next batch
    if (cards.length > 0 && currentIndex >= cards.length - 3 && hasMore && !fetchingMore) {
      fetchDueCards(cursor, false);
    }
  }, [currentIndex, cards.length, hasMore, cursor, fetchingMore, fetchDueCards]);

  const handleCheckOrNext = () => {
    if (!isChecked) {
      setIsChecked(true);
    } else if (currentIndex < cards.length - 1) {
      setCurrentIndex((prev) => prev + 1);
    }
  };

  const handlePrevious = () => {
    if (currentIndex > 0) {
      setCurrentIndex((prev) => prev - 1);
    }
  };

  if (loading) {
    return (
      <div className="flex flex-1 flex-col items-center justify-between px-4 pb-8 sm:px-6 lg:px-8">
        <div className="flex w-full flex-1 items-center justify-center">
          <Skeleton className="h-[400px] w-full max-w-sm -translate-y-30 rounded-2xl sm:max-w-md md:max-w-lg lg:max-w-xl xl:max-w-2xl" />
        </div>
        <div className="flex flex-col items-center gap-4">
          <Skeleton className="h-6 w-24 rounded-md" />
          <Skeleton className="h-10 w-28 rounded-md" />
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="p-8">
        <h1 className="text-3xl font-bold">Authentication Required</h1>
        <p className="mt-4 text-[var(--muted-foreground)]">Please sign in to study this deck.</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-1 flex-col items-center justify-center px-4 pb-8 sm:px-6 lg:px-8">
        <div className="p-8 text-center">
          <h1 className="text-3xl font-bold">Error loading cards</h1>
          <p className="mt-4 text-[var(--muted-foreground)]">{error}</p>
          <Button variant="outline" onClick={() => fetchDueCards('', true)} className="mt-4">
            Try Again
          </Button>
        </div>
      </div>
    );
  }

  if (cards.length === 0) {
    return (
      <div className="p-8">
        <h1 className="text-3xl font-bold">No Cards Due</h1>
        <p className="mt-4 text-[var(--muted-foreground)]">There are no cards due for review in this deck.</p>
      </div>
    );
  }

  const currentCard = cards[currentIndex];
  const isLastCard = currentIndex >= cards.length - 1 && !hasMore;
  const isFirstCard = currentIndex === 0;
  const showCompletion = isLastCard && isChecked;

  return (
    <div className="flex flex-1 flex-col items-center justify-between px-4 pb-8 sm:px-6 lg:px-8">
      <div className="flex w-full flex-1 items-center justify-center">
        <RenderCard
          key={`${currentCard.id}-${refreshKey}`}
          card={currentCard}
          className="max-w-sm -translate-y-30 transform px-10 py-5 text-xl sm:max-w-md md:max-w-lg lg:max-w-xl xl:max-w-2xl"
        />
      </div>
      <div className="flex flex-col items-center gap-4">
        <p className="text-sm text-[var(--muted-foreground)]">
          Card {currentIndex + 1} of {cards.length}
          {hasMore && '+'}
        </p>
        {!showCompletion ? (
          <div className="flex gap-2">
            <Button onClick={handlePrevious} size="lg" variant="outline" disabled={isFirstCard}>
              <ArrowLeft className="mr-2" size={20} />
              Previous
            </Button>
            <Button onClick={handleCheckOrNext} size="lg">
              {isChecked && !isLastCard ? (
                <>
                  Next
                  <ArrowRight className="ml-2" size={20} />
                </>
              ) : (
                <>
                  Check
                  <SearchCheck className="ml-2" size={20} />
                </>
              )}
            </Button>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-between text-center">
            <p className="text-lg font-semibold">Study session complete!</p>
            <p className="mt-2 text-[var(--muted-foreground)]">You&apos;ve reviewed all cards in this session.</p>
            <div className="mt-4 flex flex-row items-center gap-2">
              {!isFirstCard && (
                <Button onClick={handlePrevious} size="lg" variant="outline">
                  <ArrowLeft className="mr-2" size={20} />
                  Previous
                </Button>
              )}
              {isChecked && isLastCard && (
                <Button size="lg" onClick={() => router.push(`/decks/${deckId}`)}>
                  Done
                  <PartyPopper className="ml-2" size={20} />
                </Button>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
