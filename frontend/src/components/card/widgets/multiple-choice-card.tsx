import { JSX } from 'react';
import { MultipleChoiceCard as MultipleChoiceCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { MultipleChoiceCardInteractive } from './multiple-choice-card-interactive';

export function MultipleChoiceCard({ card, className }: CardComponentProps<MultipleChoiceCardType>): JSX.Element {
  return <MultipleChoiceCardInteractive card={card} className={className} />;
}

export function MultipleChoiceCardThumbnail({
  card,
  className,
}: CardComponentProps<MultipleChoiceCardType>): JSX.Element {
  const keys = Object.keys(card.options);
  const correctAnswers = keys.filter((key) => card.options[key]);
  const isMultipleChoice = correctAnswers.length > 1;

  return (
    <div className={className}>
      {card.question && <p className="pb-2 font-bold">{card.question}</p>}
      {isMultipleChoice ? (
        <div className="flex flex-col space-y-3">
          {keys.map((key) => (
            <div className="flex items-center space-x-2" key={key}>
              <Checkbox id={key} checked={card.options[key]} className="pointer-events-none" />
              <Label htmlFor={key} className="pointer-events-none">
                {key}
              </Label>
            </div>
          ))}
        </div>
      ) : (
        <RadioGroup defaultValue={correctAnswers[0]}>
          {keys.map((key) => (
            <div className="flex items-center space-x-2" key={key}>
              <RadioGroupItem value={key} id={key} className="pointer-events-none" />
              <Label htmlFor={key} className="pointer-events-none">
                {key}
              </Label>
            </div>
          ))}
        </RadioGroup>
      )}
    </div>
  );
}
