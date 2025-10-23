'use client'

import { deckSchema } from "@/lib/deckSchema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import z from "zod";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "./ui/dialog";
import { updateDeck } from "@/app/api";
import { Form, FormControl, FormField, FormItem, FormLabel } from "./ui/form";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { EmailInput } from "./email-input";

interface EditDeckMenuProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  deckId: string;
  initialData: Record<string, unknown>;
};

export function EditDeckMenu({ open, onOpenChange, deckId, initialData }: EditDeckMenuProps) {
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const form = useForm<z.infer<typeof deckSchema>>({
    resolver: zodResolver(deckSchema),
    defaultValues: initialData,
  });

  const onSubmit = form.handleSubmit(async (values) => {
    setLoading(true);
    try {
      const prevEmails = (initialData.shared_emails ?? []) as string[];
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
        alert("Deck updated successfully! values: " + JSON.stringify(values));
        alert(JSON.stringify(payload));
        router.refresh();
      } else {
        alert(`Failed to update deck: ${res.message}`);
      }
    } catch (err) {
      console.error(err);
      alert("Something went wrong updating the deck");
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
                        value={field.value ?? ""}
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
                {loading ? "Loading..." : "Edit card"}
              </Button>
            </form>
          </Form>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
