'use client';

import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { CardType } from '@/types/card';

interface Props {
  value: CardType;
  onChange: (value: CardType) => void;
}

export default function CardTypeSelector({ value, onChange }: Props) {
  return (
    <div className="flex flex-col">
      <label className="mb-2 text-sm font-medium">Card type</label>
      <Select value={value} onValueChange={(v) => onChange(v as CardType)}>
        <SelectTrigger>
          <SelectValue placeholder="Select a card type" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="front_back">Front / Back</SelectItem>
          <SelectItem value="blanks">Fill in the Blanks</SelectItem>
          <SelectItem value="multiple_choice">Multiple Choice</SelectItem>
          <SelectItem value="ordered">Ordered</SelectItem>
        </SelectContent>
      </Select>
    </div>
  );
}
