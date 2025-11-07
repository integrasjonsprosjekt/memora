'use client';

import { JSX, useState, useEffect, useCallback } from 'react';
import { match } from 'ts-pattern';
import { Check, X } from 'lucide-react';
import {
  Card as CardType,
  FrontBackCard as FrontBackCardType,
  FillBlanksCard as FillBlanksCardType,
  MultipleChoiceCard as MultipleChoiceCardType,
  OrderedCard as OrderedCardType,
} from '@/types/card';
import { CardComponentProps, CardRendererProps, UserAnswer } from './types';
import { FrontBackCard, FrontBackCardThumbnail } from './widgets/front-back-card';
import { FillBlanksCard, FillBlanksCardThumbnail } from './widgets/fill-blanks-card';
import { MultipleChoiceCard, MultipleChoiceCardThumbnail } from './widgets/multiple-choice-card';
import { OrderedCard, OrderedCardThumbnail } from './widgets/ordered-card';

import { CardThumbnail } from './card-thumbnail';
import { cn } from '@/lib/utils';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

/**
 * Renders an interactive card.
 */
export function RenderCard({ card, className }: CardComponentProps<CardType>): JSX.Element {
  const [isVerified, setIsVerified] = useState(false);
  const [userAnswer, setUserAnswer] = useState<UserAnswer | null>(null);
  const [feedbackState, setFeedbackState] = useState<'correct' | 'incorrect' | null>(null);

  const verifyAnswer = useCallback((): boolean => {
    return match(card)
      .with({ type: 'front_back' }, () => {
        // Front-back cards are always considered correct (user reveals the answer)
        return true;
      })
      .with({ type: 'blanks' }, () => {
        const fillBlanksCard = card as FillBlanksCardType;
        if (!userAnswer || !Array.isArray(userAnswer)) return false;

        // Check if all blanks are filled correctly
        return fillBlanksCard.answers.every((correctAnswer, index) => {
          const userAns = userAnswer[index]?.trim().toLowerCase();
          const correctAns = correctAnswer.trim().toLowerCase();
          return userAns === correctAns;
        });
      })
      .with({ type: 'multiple_choice' }, () => {
        const multipleChoiceCard = card as MultipleChoiceCardType;
        if (!userAnswer) return false;

        const correctAnswers = Object.keys(multipleChoiceCard.options).filter(
          (key) => multipleChoiceCard.options[key]
        );

        // Check if it's single or multiple choice
        if (correctAnswers.length === 1) {
          // Single choice - userAnswer should be a string
          return typeof userAnswer === 'string' && userAnswer === correctAnswers[0];
        } else {
          // Multiple choice - userAnswer should be an object with boolean values
          if (typeof userAnswer === 'string' || typeof userAnswer === 'boolean' || Array.isArray(userAnswer)) {
            return false;
          }
          return Object.keys(multipleChoiceCard.options).every(
            (key) => userAnswer[key] === multipleChoiceCard.options[key]
          );
        }
      })
      .with({ type: 'ordered' }, () => {
        const orderedCard = card as OrderedCardType;
        if (!userAnswer || !Array.isArray(userAnswer)) return false;

        // Check if the order matches exactly
        return orderedCard.options.every((item, index) => item === userAnswer[index]);
      })
      .exhaustive();
  }, [card, userAnswer]);

  const handleButtonClick = useCallback(() => {
    if (feedbackState !== null && !isVerified) return; // Prevent multiple clicks during feedback

    if (isVerified) {
      // Handle "Next" action
      console.log('Moving to next card');
      // Reset for next card
      setIsVerified(false);
      setUserAnswer(null);
      setFeedbackState(null);
      // TODO: Call parent callback or navigate to next card
    } else {
      // Verify the current answer
      const isCorrect = verifyAnswer();
      console.log('Answer is correct:', isCorrect);

      if (isCorrect) {
        setFeedbackState('correct');
        setIsVerified(true);
        // Remove feedback after 1 second
        setTimeout(() => {
          setFeedbackState(null);
        }, 1000);
      } else {
        setFeedbackState('incorrect');
        console.log('Set feedback to incorrect');
        // Remove feedback after 1 second
        setTimeout(() => {
          setFeedbackState(null);
        }, 1000);
      }
    }
  }, [isVerified, feedbackState, verifyAnswer]);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Enter') {
        handleButtonClick();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [handleButtonClick]);

  const cardComponent = match(card)
    .with({ type: 'front_back' }, () => (
      <FrontBackCard
        card={card as FrontBackCardType}
        className={cn(
          className,
          // Counteract padding for rulers
          "[&>hr]:-mx-10 [&>hr]:w-auto"
        )}
        onAnswerChange={setUserAnswer}
      />
    ))
    .with({ type: 'blanks' }, () => (
      <FillBlanksCard
        card={card as FillBlanksCardType}
        className={className}
        onAnswerChange={setUserAnswer}
      />
    ))
    .with({ type: 'multiple_choice' }, () => (
      <MultipleChoiceCard
        card={card as MultipleChoiceCardType}
        className={className}
        onAnswerChange={setUserAnswer}
      />
    ))
    .with({ type: 'ordered' }, () => (
      <OrderedCard
        card={card as OrderedCardType}
        className={className}
        onAnswerChange={setUserAnswer}
      />
    ))
    .exhaustive();

  return (
    <div className="flex flex-1 flex-col items-center justify-between px-4 sm:px-6 lg:px-8 pb-8">
      <div className="flex flex-1 items-center justify-center w-full">
        <Card className="w-full max-w-sm -translate-y-30 transform px-10 py-5 text-xl sm:max-w-md md:max-w-lg lg:max-w-xl xl:max-w-2xl">
          {cardComponent}
        </Card>
      </div>

      <div className="flex justify-center w-full">
        <Button
          onClick={handleButtonClick}
          className={cn(
            "min-w-[120px] transition-all duration-300",
            feedbackState === 'correct' && "!bg-green-500/20 !border-green-500 !border-1 !text-green-500",
            feedbackState === 'incorrect' && "!bg-red-500/20 !border-red-500 !border-1 !text-red-500"
          )}
          variant={isVerified ? "default" : "outline"}
        >
          <span
            key={`${feedbackState}-${isVerified}`}
            className="inline-flex items-center justify-center animate-in fade-in duration-200"
          >
            {feedbackState === 'correct' ? (
              <Check size={5} />
            ) : feedbackState === 'incorrect' ? (
              <X size={5} />
            ) : isVerified ? (
              'Next'
            ) : (
              'Check'
            )}
          </span>
        </Button>
      </div>
    </div>
  );
}

/**
 * Renders a non-interactive card thumbnail.
 */
export function RenderCardThumbnail({ card, className, deckId }: CardRendererProps<CardType>): JSX.Element {
  const cardComponent = match(card)
    .with({ type: 'front_back' }, () => (
      <FrontBackCardThumbnail card={card as FrontBackCardType} className={cn(
          className,
          // TODO: Counteract padding for rulers
        )}
      />
    ))
    .with({ type: 'blanks' }, () => <FillBlanksCardThumbnail card={card as FillBlanksCardType} className={className} />)
    .with({ type: 'multiple_choice' }, () => (
      <MultipleChoiceCardThumbnail card={card as MultipleChoiceCardType} className={className} />
    ))
    .with({ type: 'ordered' }, () => <OrderedCardThumbnail card={card as OrderedCardType} className={className} />)
    .exhaustive();

  const tags = [card.type];

  return (
    <CardThumbnail card={card} deckId={deckId} tags={tags}>
      {cardComponent}
    </CardThumbnail>
  );
}
