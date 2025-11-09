'use client';

import { JSX, useState, useEffect } from 'react';
import { FillBlanksCard as FillBlanksCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { Input } from '@/components/ui/input';

export function FillBlanksCardInteractive({
  card,
  className,
  onAnswerChange,
}: CardComponentProps<FillBlanksCardType>): JSX.Element {
  const parts = card.question.split('{}');
  const [userAnswers, setUserAnswers] = useState<string[]>(new Array(card.answers.length).fill(''));

  useEffect(() => {
    if (onAnswerChange) {
      onAnswerChange(userAnswers);
    }
  }, [userAnswers, onAnswerChange]);

  const handleInputChange = (index: number, value: string) => {
    const newAnswers = [...userAnswers];
    newAnswers[index] = value;
    setUserAnswers(newAnswers);
  };

  return (
    <div className={className}>
      {parts.map((part, index) => (
        <span key={index} className="inline-flex items-baseline">
          {/* Render text part */}
          <span>{part}</span>
          {/* Render input field if there's a corresponding answer */}
          {index < parts.length - 1 && index < card.answers.length && (
            <Input
              type="text"
              value={userAnswers[index]}
              onChange={(e) => handleInputChange(index, e.target.value)}
              className="mx-1 inline-block h-7 w-32 px-2 text-sm"
              placeholder="..."
            />
          )}
        </span>
      ))}
    </div>
  );
}
