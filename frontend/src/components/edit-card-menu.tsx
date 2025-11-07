'use client';

import { useEffect, useState } from 'react';
import { useForm, UseFormReturn } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Form } from '@/components/ui/form';

import { FrontBackField } from './front-back-field';
import { BlankFields } from './blank-field';
import { MultipleChoiceFields } from './multiple-choice-field';

import { CardInput, cardInputSchemas, CardPayload, cardPayloadSchemas } from '@/lib/cardSchemas';
import { Card } from '@/types/card';
import { updateCard } from '@/app/api';
import { Button } from '@/components/ui/button';
import { OrderedFields } from './ordered-field';
import { useRouter } from 'next/navigation';
import { match } from 'ts-pattern';
import normalizeCardData from '@/lib/normalizeCardData';
import { toast } from 'sonner';
import { Spinner } from '@/components/ui/spinner';

interface EditCardMenuProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  card: Card;
  deckId: string;
}

export function EditCardMenu({ open, onOpenChange, card, deckId }: EditCardMenuProps) {
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const inputSchema = cardInputSchemas[card.type];
  const payloadSchema = cardPayloadSchemas[card.type];

  const form = useForm<CardInput>({
    resolver: zodResolver(inputSchema),
    defaultValues: card as CardInput,
  });

  useEffect(() => {
    if (card) {
      form.reset({ type: card.type, ...normalizeCardData(card) } as CardInput);
    }
  }, [card, form]);

  const onSubmit = form.handleSubmit(async (values) => {
    setLoading(true);
    try {
      const payload: CardPayload = payloadSchema.parse(values);

      const res = await updateCard(deckId, card.id, payload);

      if (res.success) {
        form.reset({});
        onOpenChange(false);
        toast.success('Card updated');
        router.refresh();
      } else {
        console.error(res.message);
        toast.error('Failed to update card');
      }
    } catch (err) {
      console.error(err);
      toast.error('Failed to update card');
    } finally {
      setLoading(false);
    }
  });

  const renderFields = match(form.getValues())
    .with({ type: 'front_back' }, () => (
      <FrontBackField form={form as UseFormReturn<CardInput & { type: 'front_back' }>} />
    ))
    .with({ type: 'blanks' }, () => <BlankFields form={form as UseFormReturn<CardInput & { type: 'blanks' }>} />)
    .with({ type: 'multiple_choice' }, () => (
      <MultipleChoiceFields form={form as UseFormReturn<CardInput & { type: 'multiple_choice' }>} />
    ))
    .with({ type: 'ordered' }, () => <OrderedFields form={form as UseFormReturn<CardInput & { type: 'ordered' }>} />)
    .exhaustive();

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Card</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={onSubmit} className="flex flex-col gap-4 py-4">
            {renderFields}
            <Button type="submit" disabled={loading} className="mt-4">
              {/*Ensures that the button is disabled while loading*/}
              {loading ? <Spinner /> : 'Edit card'}
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
