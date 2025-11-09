'use client';

import { JSX, useState, useEffect, useRef } from 'react';
import { MultipleChoiceCard as MultipleChoiceCardType } from '@/types/card';
import { CardComponentProps } from '../types';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';

export function MultipleChoiceCardInteractive({
  card,
  className,
  onAnswerChange,
}: CardComponentProps<MultipleChoiceCardType>): JSX.Element {
  const keys = Object.keys(card.options);
  const correctAnswers = keys.filter((key) => card.options[key]);
  const isMultipleChoice = correctAnswers.length > 1;

  const [selectedOptions, setSelectedOptions] = useState<Record<string, boolean>>(
    keys.reduce((acc, key) => ({ ...acc, [key]: false }), {})
  );
  const [selectedRadio, setSelectedRadio] = useState<string>('');

  // Store the latest callback in a ref to avoid it being a dependency
  const onAnswerChangeRef = useRef(onAnswerChange);

  useEffect(() => {
    onAnswerChangeRef.current = onAnswerChange;
  }, [onAnswerChange]);

  const handleCheckboxChange = (key: string, checked: boolean) => {
    setSelectedOptions((prev) => ({ ...prev, [key]: checked }));
  };

  useEffect(() => {
    // Notify parent of checkbox selection changes
    if (onAnswerChangeRef.current && isMultipleChoice) {
      onAnswerChangeRef.current(selectedOptions);
    }
  }, [selectedOptions, isMultipleChoice]);

  useEffect(() => {
    // Notify parent of radio selection changes
    if (onAnswerChangeRef.current && selectedRadio) {
      onAnswerChangeRef.current(selectedRadio);
    }
  }, [selectedRadio]);

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
