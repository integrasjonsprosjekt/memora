'use client';
import { useForm } from 'react-hook-form';
import { Form, FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form';
import { zodResolver } from '@hookform/resolvers/zod';
import { deckSchema } from '@/lib/deckSchema';
import { z } from 'zod';
import { useState } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { createDeck } from '@/app/api';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/auth';
import { EmailInput } from './email-input';
import { Dialog, DialogContent, DialogHeader } from '@/components/ui/dialog';
import { DialogTitle } from '@radix-ui/react-dialog';
import { toast } from 'sonner';
import { Spinner } from '@/components/ui/spinner';

interface AddDeckMenuProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function AddDeckMenu({ open, onOpenChange }: AddDeckMenuProps) {
  const router = useRouter();
  const { user } = useAuth();
  const [loading, setLoading] = useState(false);

  const form = useForm<z.infer<typeof deckSchema>>({
    resolver: zodResolver(deckSchema),
    defaultValues: {},
  });

  const onSubmit = form.handleSubmit(async (values) => {
    if (!user) return;

    setLoading(true);
    try {
      const res = await createDeck(user, user.uid, values.title, values.shared_emails ?? []);
      if (res.success) {
        form.reset();
        onOpenChange(false);
        router.push(`/decks/${res.data.id}`);
        toast.success('Deck created', {
          duration: 1500,
        });
      } else {
        console.error(res.message);
        toast.error('Failed to create deck');
      }
    } catch (err) {
      console.error(err);
      toast.error('Failed to create deck');
    } finally {
      setLoading(false);
    }
  });
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add new Deck</DialogTitle>
        </DialogHeader>
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
              {loading ? <Spinner /> : 'Add deck'}
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
