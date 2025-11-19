'use client';

import { RenderCard, RenderCardThumbnail } from '@/components/card';
import { use, useState, useEffect, useCallback, useRef } from 'react';
import { useAuth } from '@/context/auth';
import { Card as CardType } from '@/types/card';
import { fetchApi } from '@/lib/api/config';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { toast } from 'sonner';
import { ArrowLeft, SearchCheck, PartyPopper, Check, X } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { validateAnswer } from '@/lib/card-validation';
import { UserAnswer } from '@/components/card/types';
import { cn } from '@/lib/utils';

type Rating = 'again' | 'hard' | 'good' | 'easy';

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
  const [isCorrect, setIsCorrect] = useState<boolean | null>(null);
  const [userAnswer, setUserAnswer] = useState<UserAnswer | null>(null);
  const [cursor, setCursor] = useState<string>('');
  const [hasMore, setHasMore] = useState(false);
  const [fetchingMore, setFetchingMore] = useState(false);
  const [flipTrigger, setFlipTrigger] = useState(0);
  const [isRating, setIsRating] = useState(false);
  const [hasRated, setHasRated] = useState(false);
  const checkButtonRef = useRef<(() => void) | null>(null);
  const doneButtonRef = useRef<(() => void) | null>(null);

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

        // Validate that cards is an array
        const cardsArray = Array.isArray(response.cards) ? response.cards : [];

        if (isInitial) {
          setCards(cardsArray);
        } else {
          setCards((prev) => [...prev, ...cardsArray]);
        }
        setCursor(response.next_cursor || '');
        setHasMore(response.has_more || false);
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
    // Reset state when moving to a new card
    setIsChecked(false);
    setIsCorrect(null);
    setUserAnswer(null);
    setFlipTrigger(0);
    setHasRated(false);
    setRefreshKey((prev) => prev + 1);

    // If we're approaching the end and there are more cards, fetch the next batch
    if (cards.length > 0 && currentIndex >= cards.length - 3 && hasMore && !fetchingMore) {
      fetchDueCards(cursor, false);
    }
  }, [currentIndex, cards.length, hasMore, cursor, fetchingMore, fetchDueCards]);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Enter') {
        // Don't submit if user is typing in an input or textarea
        const activeElement = document.activeElement;
        if (activeElement instanceof HTMLTextAreaElement) {
          return;
        }

        if (!isChecked && !isRating && checkButtonRef.current) {
          checkButtonRef.current();
        } else if (doneButtonRef.current) {
          doneButtonRef.current();
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [isChecked, isRating]);

  const handleAnswerChange = (answer: UserAnswer) => {
    setUserAnswer(answer);
  };

  const handleRating = async (rating: Rating) => {
    if (!user || isRating) return;

    setIsRating(true);
    try {
      await fetchApi(`decks/${deckId}/cards/${currentCard.id}/progress`, {
        method: 'PUT',
        user,
        body: JSON.stringify({ rating }),
      });

      setHasRated(true);

      // Move to next card if not the last one
      if (currentIndex < cards.length - 1) {
        setCurrentIndex((prev) => prev + 1);
      }
    } catch (error) {
      console.error('Failed to update card progress:', error);
      toast.error('Failed to save rating');
    } finally {
      setIsRating(false);
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
          <Button variant="outline" onClick={() => fetchDueCards('', true)} className="mt-4 cursor-pointer">
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
  const showCompletion = isLastCard && isChecked && hasRated;

  const handleDone = () => {
    router.push(`/decks/${deckId}`);
  };

  const handleCheckOrNext = () => {
    if (!isChecked) {
      // For front-back cards, trigger the flip before validation
      if (currentCard.type === 'front_back') {
        setFlipTrigger((prev) => prev + 1);
      }

      // Validate the answer
      const correct = validateAnswer(currentCard, userAnswer);
      setIsCorrect(correct);
      setIsChecked(true);

      // Show feedback toast
      if (correct) {
        toast.success('Correct!');
      } else {
        toast.error('Incorrect');
      }
    } else if (currentIndex < cards.length - 1) {
      setCurrentIndex((prev) => prev + 1);
    }
  };

  // Store the functions in refs so the event listener can access them
  checkButtonRef.current = handleCheckOrNext;
  doneButtonRef.current = showCompletion ? handleDone : null;

  return (
    <div className="flex flex-1 flex-col items-center justify-between px-4 pb-8 sm:px-6 lg:px-8">
      <div className="flex w-full flex-1 items-center justify-center">
        <div className="grid w-full max-w-sm -translate-y-30 transform sm:max-w-md md:max-w-lg lg:max-w-xl xl:max-w-2xl">
          {isChecked && isCorrect === false ? (
            <RenderCardThumbnail
              card={currentCard}
              deckId={deckId}
              clickable={false}
              className="col-start-1 row-start-1 px-10 py-5 text-xl"
            />
          ) : (
            <RenderCard
              key={`${currentCard.id}-${refreshKey}`}
              card={currentCard}
              className="col-start-1 row-start-1 px-10 py-5 text-xl"
              onAnswerChange={handleAnswerChange}
              flipTrigger={flipTrigger}
            />
          )}
        </div>
      </div>
      <div className="flex flex-col items-center gap-4">
        <p className="text-sm text-[var(--muted-foreground)]">
          Card {currentIndex + 1} of {cards.length}
          {hasMore && '+'}
        </p>
        {!showCompletion ? (
          <>
            {!isChecked ? (
              <div className="flex gap-2">
                <Button onClick={handlePrevious} size="lg" variant="outline" disabled={isFirstCard}>
                  <ArrowLeft className="mr-2" size={20} />
                  Previous
                </Button>
                <Button onClick={handleCheckOrNext} size="lg">
                  {isChecked && isCorrect !== null ? (
                    <>{isCorrect ? <Check size={20} /> : <X size={20} />}</>
                  ) : (
                    <>
                      Check
                      <SearchCheck className="ml-2" size={20} />
                    </>
                  )}
                </Button>
              </div>
            ) : (
              <div className="flex flex-col gap-3">
                <p
                  className={cn(
                    'text-center text-sm font-semibold',
                    isCorrect === true && 'text-green-500',
                    isCorrect === false && 'text-red-500'
                  )}
                >
                  {currentCard.type !== 'front_back' ? (isCorrect ? 'Correct!' : 'Incorrect') : ''}
                </p>
                <div className="flex gap-2">
                  <Button
                    onClick={() => handleRating('again')}
                    variant="outline"
                    size="lg"
                    disabled={isRating}
                    className="!cursor-pointer !border-red-500 !bg-red-500/10 !text-red-500 hover:!bg-red-500/20"
                  >
                    Again
                  </Button>
                  <Button
                    onClick={() => handleRating('hard')}
                    variant="outline"
                    size="lg"
                    disabled={isRating}
                    className="!cursor-pointer !border-orange-500 !bg-orange-500/10 !text-orange-500 hover:!bg-orange-500/20"
                  >
                    Hard
                  </Button>
                  <Button
                    onClick={() => handleRating('good')}
                    variant="outline"
                    size="lg"
                    disabled={isRating}
                    className="!cursor-pointer !border-blue-500 !bg-blue-500/10 !text-blue-500 hover:!bg-blue-500/20"
                  >
                    Good
                  </Button>
                  <Button
                    onClick={() => handleRating('easy')}
                    variant="outline"
                    size="lg"
                    disabled={isRating}
                    className="!cursor-pointer !border-green-500 !bg-green-500/10 !text-green-500 hover:!bg-green-500/20"
                  >
                    Easy
                  </Button>
                </div>
              </div>
            )}
          </>
        ) : (
          <>
            <div className="text-center">
              <p className="text-lg font-semibold">Study session complete!</p>
              <p className="mt-2 text-[var(--muted-foreground)]">You&apos;ve reviewed all cards in this session.</p>
            </div>
            <div className="flex flex-row items-center gap-2">
              {!isFirstCard && (
                <Button onClick={handlePrevious} size="lg" variant="outline" className="cursor-pointer">
                  <ArrowLeft className="mr-2" size={20} />
                  Previous
                </Button>
              )}
              <Button size="lg" onClick={handleDone}>
                Done
                <PartyPopper className="ml-2" size={20} />
              </Button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
