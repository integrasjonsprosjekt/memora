'use client';

import { useEffect, useMemo, useState } from 'react';
import { useForm, UseFormReturn } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Form } from './ui/form';

import CardTypeSelector from './card-type-selector';
import { FrontBackField } from './front-back-field';
import { BlankFields } from './blank-field';
import { MultipleChoiceFields } from './multiple-choice-field';

import { CardInput, cardInputSchemas, CardPayload, cardPayloadSchemas } from '@/lib/cardSchemas';
import { CardType } from '@/types/card';
import { createCard } from '@/app/api';
import { Button } from './ui/button';
import { OrderedFields } from './ordered-field';
import { useRouter } from 'next/navigation';

interface AddCardMenuProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  deckId: string;
}

export function AddCardMenu({ open, onOpenChange, deckId }: AddCardMenuProps) {
  const [cardType, setCardType] = useState<CardType>('front_back');
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  // Pick the schema for the current card type dynamically
  const inputSchema = useMemo(() => cardInputSchemas[cardType], [cardType]);
  const payloadSchema = useMemo(() => cardPayloadSchemas[cardType], [cardType]);

  const form = useForm<CardInput>({
    resolver: zodResolver(inputSchema),
    defaultValues: {},
  });

  useEffect(() => {
    form.reset({ type: cardType } as CardInput);
  }, [cardType, form]);

  const onSubmit = form.handleSubmit(async (values) => {
    setLoading(true);
    try {
      const payload: CardPayload = payloadSchema.parse(values);
      const res = await createCard(deckId, cardType, payload);

      if (res.success) {
        form.reset({ type: cardType } as CardInput);
        onOpenChange(false);
        router.refresh();
        alert('Card created successfully!');
      } else {
        alert(`Failed to create card: ${res.message}`);
      }
    } catch (err) {
      console.error(err);
      alert('Something went wrong creating the card');
    } finally {
      setLoading(false);
    }
  });

  const renderFields = () => {
    switch (cardType) {
      case 'front_back':
        return <FrontBackField form={form as UseFormReturn<CardInput & { type: 'front_back' }>} />;
      case 'blanks':
        return <BlankFields form={form as UseFormReturn<CardInput & { type: 'blanks' }>} />;
      case 'multiple_choice':
        return <MultipleChoiceFields form={form as UseFormReturn<CardInput & { type: 'multiple_choice' }>} />;
      case 'ordered':
        return <OrderedFields form={form as UseFormReturn<CardInput & { type: 'ordered' }>} />;
      default:
        return null;
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Card</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={onSubmit} className="flex flex-col gap-4 py-4">
            <CardTypeSelector value={cardType} onChange={setCardType} />
            {renderFields()}
            <Button type="submit" disabled={loading} className="mt-4">
              {/*Ensures that the button is disabled while loading*/}
              {loading ? 'Loading...' : 'Add card'}
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
