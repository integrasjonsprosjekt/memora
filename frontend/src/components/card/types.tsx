import { JSX, ReactNode } from 'react';
import { plainToInstance } from 'class-transformer';
import 'reflect-metadata';
import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Checkbox } from '@/components/ui/checkbox';

export type CardType = 'front_back' | 'blanks' | 'multiple_choice' | 'ordered';

export abstract class BaseCard {
  constructor(
    public id: string,
    public type: CardType,
  ) {}

  // Abstract method to render card component
  // TOOD: Add render method
  // abstract render(): JSX.Element;
  // Abstract method to display card information
  abstract display(): JSX.Element;
}

export class FrontBackCardType extends BaseCard {
  constructor(
    id: string,
    public front: string,
    public back: string,
  ) {
    super(id, 'front_back');
  }

  display(): JSX.Element {
    return (
      <>
        <p>{this.front}</p>

        <hr className="border-border tap-highlight-transparent my-2 w-full border-t border-dashed" />

        <p>{this.back.length > 100 ? this.back.substring(0, 100) + '...' : this.back}</p>
      </>
    );
  }
}

export class FillBlanksCardType extends BaseCard {
  constructor(
    id: string,
    public question: string,
    public answers: string[],
  ) {
    super(id, 'blanks');
  }

  display(): JSX.Element {
    const parts = this.question.split('{}');
    return (
      <>
        <p>
          {parts.map((part, index) => (
            <span key={index}>
              {/* Render part */}
              {part}
              {/* Check if there is a corresponding answer */}
              {index < parts.length - 1 && index < this.answers.length && (
                <span className="bg-accent rounded-sm border border-dashed border-[var(--border)] px-1 py-[0.02rem]">
                  {this.answers[index]}
                </span>
              )}
            </span>
          ))}
        </p>
      </>
    );
  }
}

export class MultipleChoiceCardType extends BaseCard {
  constructor(
    id: string,
    public question: string,
    public options: { [option: string]: boolean },
  ) {
    super(id, 'multiple_choice');
  }

  display(): JSX.Element {
    const keys = Object.keys(this.options);
    const correctAnswers = keys.filter((key) => this.options[key]);
    const isMultipleChoice = correctAnswers.length > 1;

    return (
      <>
        {this.question && <p className="pb-2 font-bold">{this.question}</p>}
        {isMultipleChoice ? (
          <div className="flex flex-col space-y-3">
            {keys.map((key) => (
              <div className="flex items-center space-x-2" key={key}>
                <Checkbox id={key} checked={this.options[key]} className="pointer-events-none" />
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
      </>
    );
  }
}

export class OrderedCardType extends BaseCard {
  constructor(
    id: string,
    public question: string,
    public options: string[],
  ) {
    super(id, 'ordered');
  }

  display(): JSX.Element {
    return (
      <>
        {this.question && <p className="pb-2 font-bold">{this.question}</p>}
        <ol className="marker:text-muted-foreground/50 list-decimal pl-5 marker:text-xs">
          {this.options.map((key) => (
            <li key={key}>{key}</li>
          ))}
        </ol>
      </>
    );
  }
}

// Union type for all cards
export type Card =
  | FrontBackCardType
  | FillBlanksCardType
  | MultipleChoiceCardType
  | OrderedCardType;

// Factory method to convert JSON to appropriate card type
export function createCardFromJson(json: any): Card {
  const cardTypeMap = {
    front_back: FrontBackCardType,
    blanks: FillBlanksCardType,
    multiple_choice: MultipleChoiceCardType,
    ordered: OrderedCardType,
  };

  const CardClass = cardTypeMap[json.type as CardType];
  if (!CardClass) {
    throw new Error(`Unknown card type: ${json.type}`);
  }

  return plainToInstance(CardClass, json);
}
