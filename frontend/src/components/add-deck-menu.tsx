'use client';
import { useForm } from "react-hook-form";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "./ui/dialog";
import { Form, FormControl, FormField, FormItem, FormLabel } from "./ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import { deckSchema } from "@/lib/deckSchema";
import { z } from "zod";
import { useState } from "react";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { createDeck } from "@/app/api";
import { useRouter } from "next/navigation";
import { EmailInput } from "./email-input";

export function AddDeckMenu() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const form = useForm<z.infer<typeof deckSchema>>({
    resolver: zodResolver(deckSchema),
    defaultValues: {},
  });

  const onSubmit = form.handleSubmit(async (values) => {
    setLoading(true);
    try {
      const res = await createDeck(values.title, values.shared_emails ?? []);
      if (res.success) {
        router.push(`/decks/${values.title}`);
      } else {
        alert(res.message);
      }
    } catch (err) {
      console.error(err);
      alert("Something went wrong creating the deck");
    } finally {
      setLoading(false);
    }
  });
  return (
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
          {loading ? "Loading..." : "Add deck"}
        </Button>
      </form>
    </Form>
  );
}
