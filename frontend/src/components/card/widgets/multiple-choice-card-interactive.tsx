'use client';

import { JSX, useState } from 'react';
import { MultipleChoiceCard as MultipleChoiceCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';

export function MultipleChoiceCardInteractive({
  card,
  className,
}: CardComponentProps<MultipleChoiceCardType>): JSX.Element {
  const keys = Object.keys(card.options);
  const correctAnswers = keys.filter((key) => card.options[key]);
  const isMultipleChoice = correctAnswers.length > 1;

  const [selectedOptions, setSelectedOptions] = useState<Record<string, boolean>>(
    keys.reduce((acc, key) => ({ ...acc, [key]: false }), {})
  );
  const [selectedRadio, setSelectedRadio] = useState<string>('');

  const handleCheckboxChange = (key: string, checked: boolean) => {
    setSelectedOptions((prev) => ({ ...prev, [key]: checked }));
  };

  return (
    <div className={className}>
      {card.question && <p className="pb-2 font-bold">{card.question}</p>}
      {isMultipleChoice ? (
        <div className="flex flex-col space-y-3">
          {keys.map((key) => (
            <div className="flex items-center space-x-2" key={key}>
              <Checkbox
                id={key}
                checked={selectedOptions[key]}
                onCheckedChange={(checked) => handleCheckboxChange(key, checked as boolean)}
              />
              <Label htmlFor={key} className="cursor-pointer">
                {key}
              </Label>
            </div>
          ))}
        </div>
      ) : (
        <RadioGroup value={selectedRadio} onValueChange={setSelectedRadio}>
          {keys.map((key) => (
            <div className="flex items-center space-x-2" key={key}>
              <RadioGroupItem value={key} id={key} />
              <Label htmlFor={key} className="cursor-pointer">
                {key}
              </Label>
            </div>
          ))}
        </RadioGroup>
      )}
    </div>
  );
}
