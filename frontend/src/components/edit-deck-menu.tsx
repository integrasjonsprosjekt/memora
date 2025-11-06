'use client';

import { deckSchema } from '@/lib/deckSchema';
import { zodResolver } from '@hookform/resolvers/zod';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import z from 'zod';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from './ui/dialog';
import { updateDeck } from '@/app/api';
import { Form, FormControl, FormField, FormItem, FormLabel } from './ui/form';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { EmailInput } from './email-input';
import { getApiEndpoint } from '@/config/api';
import { Deck } from '@/types/deck';

interface EditDeckMenuProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  deckId: string;
}

interface DeckApiResponse {
  Title: string;
  SharedEmails: string[];
}

function normalizeDeckData(apiData: DeckApiResponse): Partial<Deck> {
  return {
    title: apiData.Title ?? '',
    shared_emails: apiData.SharedEmails ?? [],
  };
}

export function EditDeckMenu({ open, onOpenChange, deckId }: EditDeckMenuProps) {
  const [loading, setLoading] = useState(false);
  const [deck, setDeck] = useState<Partial<Deck> | null>(null);
  const router = useRouter();

  const form = useForm<z.infer<typeof deckSchema>>({
    resolver: zodResolver(deckSchema),
    defaultValues: {},
  });

  useEffect(() => {
    if (!open) return;
    const fetchDeck = async () => {
      try {
        const res = await fetch(getApiEndpoint(`/v1/decks/${deckId}`), { cache: 'no-store' });
        const data: DeckApiResponse = await res.json();
        const normalizedData = normalizeDeckData(data);
        setDeck(normalizedData);
        form.reset({
          title: normalizedData.title ?? '',
          shared_emails: normalizedData.shared_emails ?? [],
        });
      } catch (error) {
        console.error('Error fetching deck:', error);
        setDeck(null);
      }
    };

    fetchDeck();
  }, [deckId, open, form]);

  const onSubmit = form.handleSubmit(async (values) => {
    setLoading(true);
    try {
      const prevEmails = (deck?.shared_emails ?? []) as string[];
      const newEmails = (values.shared_emails ?? []) as string[];

      const addedEmails = newEmails.filter((email) => !prevEmails.includes(email));
      const removedEmails = prevEmails.filter((email) => !newEmails.includes(email));

      const payload = {
        title: values.title,
        addedEmails,
        removedEmails,
      };

      const res = await updateDeck(deckId, payload);

      if (res.success) {
        //form.reset({});
        onOpenChange(false);
        router.refresh();
      } else {
        alert(`Failed to update deck: ${res.message}`);
      }
    } catch (err) {
      console.error(err);
      alert('Something went wrong updating the deck');
    } finally {
      setLoading(false);
    }
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Deck</DialogTitle>
          <Form {...form}>
            <form onSubmit={onSubmit} className="flex flex-col gap-4 py-4">
              <FormField
                control={form.control}
                name="title"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Title</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Deck title"
                        value={field.value ?? ''}
                        onChange={field.onChange}
                        onBlur={field.onBlur}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="shared_emails"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Shared Emails</FormLabel>
                    <FormControl>
                      <EmailInput value={field.value ?? []} onChange={field.onChange} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <Button type="submit" disabled={loading} className="mt-4">
                {/*Ensures that the button is disabled while loading*/}
                {loading ? 'Loading...' : 'Edit deck'}
              </Button>
            </form>
          </Form>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
