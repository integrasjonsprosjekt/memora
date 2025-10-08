import { ReactNode } from "react";
import { plainToInstance } from "class-transformer";
import "reflect-metadata";

export type CardType = "front_back" | "blanks" | "multiple_choice" | "ordered";

type CardDescription = {
  title: ReactNode;
  body: ReactNode;
};

export abstract class BaseCard {
  constructor(
    public id: string,
    public type: CardType,
  ) {}

  // Abstract method to render card component
  // TOOD: Add render method
  // abstract render(): JSX.Element;
  // Abstract method to display card information
  abstract display(): CardDescription;
}

export class FrontBackCardType extends BaseCard {
  constructor(
    id: string,
    public front: string,
    public back: string,
  ) {
    super(id, "front_back");
  }

  display(): CardDescription {
    return {
      title: this.front,
      body: this.back.length > 100 ? this.back.substring(0, 100) + "..." : this.back,
    };
  }
}

export class FillBlanksCardType extends BaseCard {
  constructor(
    id: string,
    public question: string,
    public answers: string[],
  ) {
    super(id, "blanks");
  }

  display(): CardDescription {
    return {
      title: this.question,
      body: this.answers.join(", "),
    };
  }
}

export class MultipleChoiceCardType extends BaseCard {
  constructor(
    id: string,
    public question: string,
    public options: { [option: string]: boolean },
  ) {
    super(id, "multiple_choice");
  }

  display(): CardDescription {
    return {
      title: this.question,
      body: Object.keys(this.options).join(", "),
    };
  }
}

export class OrderedCardType extends BaseCard {
  constructor(
    id: string,
    public question: string,
    public options: string[],
  ) {
    super(id, "ordered");
  }

  display(): CardDescription {
    return {
      title: this.question,
      body: this.options.join(", "),
    };
  }
}

// Union type for all cards
export type Card = FrontBackCardType | FillBlanksCardType | MultipleChoiceCardType | OrderedCardType;

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
