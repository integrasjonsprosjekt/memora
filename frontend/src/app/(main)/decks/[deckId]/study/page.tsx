'use client';

import { RenderCard } from '@/components/card';
import { use, useState, useEffect, useCallback, useMemo } from 'react';
import { useAuth } from '@/context/auth';
import { Card as CardType } from '@/types/card';
import { fetchApi } from '@/lib/api/config';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { toast } from 'sonner';
import { ArrowRight, ArrowLeft, SearchCheck } from 'lucide-react';

export default function Page({ params }: { params: Promise<{ deckId: string }> }) {
  const cardIds = useMemo(() => ['hG20zJsdM6MUIXko1otD', '8h3NSUhG6coQHp49yAPC', 'Cu62mrzNAmkEPGJm1ECk'], []);

  const { deckId } = use(params);
  const { user } = useAuth();

  const [currentIndex, setCurrentIndex] = useState(0);
  const [card, setCard] = useState<CardType | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);
  const [isChecked, setIsChecked] = useState(false);

  const fetchCard = useCallback(
    async (cardId: string) => {
      if (!user) return;

      try {
        setLoading(true);
        setError(null);

        const cardData = await fetchApi<CardType>(`decks/${deckId}/cards/${cardId}`, { user });
        setCard(cardData);
        setRefreshKey((prev) => prev + 1);
      } catch (error) {
        console.error(error);
        setError('Failed to load card');
        toast.error('Failed to load card');
      } finally {
        setLoading(false);
      }
    },
    [user, deckId]
  );

  useEffect(() => {
    if (cardIds.length > 0 && currentIndex < cardIds.length) {
      fetchCard(cardIds[currentIndex]);
      setIsChecked(false);
    }
  }, [currentIndex, fetchCard, cardIds]);

  const handleCheckOrNext = () => {
    if (!isChecked) {
      setIsChecked(true);
    } else if (currentIndex < cardIds.length - 1) {
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

  if (cardIds.length === 0) {
    return (
      <div className="p-8">
        <h1 className="text-3xl font-bold">No Cards Available</h1>
        <p className="mt-4 text-[var(--muted-foreground)]">This deck has no cards to study.</p>
      </div>
    );
  }

  if (error || !card) {
    return (
      <div className="flex flex-1 flex-col items-center justify-center px-4 pb-8 sm:px-6 lg:px-8">
        <div className="p-8 text-center">
          <h1 className="text-3xl font-bold">Error loading card</h1>
          <p className="mt-4 text-[var(--muted-foreground)]">{error || 'There was a problem loading this card.'}</p>
          <Button variant="outline" onClick={() => fetchCard(cardIds[currentIndex])} className="mt-4">
            Try Again
          </Button>
        </div>
      </div>
    );
  }

  const isLastCard = currentIndex >= cardIds.length - 1;
  const isFirstCard = currentIndex === 0;
  const showCompletion = isLastCard && isChecked;

  return (
    <div className="flex flex-1 flex-col items-center justify-between px-4 pb-8 sm:px-6 lg:px-8">
      <div className="flex w-full flex-1 items-center justify-center">
        <RenderCard
          key={`${card.id}-${refreshKey}`}
          card={card}
          className="max-w-sm -translate-y-30 transform px-10 py-5 text-xl sm:max-w-md md:max-w-lg lg:max-w-xl xl:max-w-2xl"
        />
      </div>
      <div className="flex flex-col items-center gap-4">
        <p className="text-sm text-[var(--muted-foreground)]">
          Card {currentIndex + 1} of {cardIds.length}
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
          <div className="text-center">
            <p className="text-lg font-semibold">Study session complete!</p>
            <p className="mt-2 text-[var(--muted-foreground)]">You&apos;ve reviewed all cards in this session.</p>
            {!isFirstCard && (
              <Button onClick={handlePrevious} size="lg" variant="outline" className="mt-4">
                <ArrowLeft className="mr-2" size={20} />
                Previous
              </Button>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
