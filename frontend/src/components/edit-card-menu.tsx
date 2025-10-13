"use client";

import { useEffect, useMemo, useState } from "react";
import { useForm, UseFormReturn } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Form } from "./ui/form";

import { FrontBackField } from "./front-back-field";
import { BlankFields } from "./blank-field";
import { MultipleChoiceFields } from "./multiple-choice-field";

import { cardSchemas } from "@/lib/cardSchemas";
import buildCardPayload from "@/lib/cardBuildPayload";
import { CardType } from "@/types/cards";
import { updateCard } from "@/app/api";
import { Button } from "./ui/button";
import { OrderedFields } from "./ordered-field";
import { useRouter } from "next/navigation";

interface EditCardMenuProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  deckId: string;
  cardId: string;
  cardType: CardType;
  initialData: Record<string, any>;
}

type CardForms = {
  front_back: z.infer<typeof cardSchemas.front_back>;
  blanks: z.infer<typeof cardSchemas.blanks>;
  multiple_choice: z.infer<typeof cardSchemas.multiple_choice>;
  ordered: z.infer<typeof cardSchemas.ordered>;
};

export function EditCardMenu({ open, onOpenChange, deckId, cardId, cardType, initialData }: EditCardMenuProps) {
  const [loading, setLoading] = useState(false);

  // Pick the schema for the current card type dynamically
  const currentSchema = useMemo(() => cardSchemas[cardType], [cardType]);

  const form = useForm<CardForms[typeof cardType]>({
    resolver: zodResolver(currentSchema),
    defaultValues: initialData,
  });

  useEffect(() => {
    if (initialData) {
      form.reset(initialData);
    }
  }, [initialData, form]);

  const router = useRouter();

  const onSubmit = form.handleSubmit(async (values) => {
    setLoading(true);
    try {
      const basePayload = buildCardPayload(cardType, values);

      let payload: Record<string, unknown>;

      switch (cardType) {
        case "front_back":
          payload = { ...basePayload, id: cardId, type: cardType };
          break;
        case "blanks":
          payload = { ...basePayload, id: cardId, type: cardType };
          break;
        case "multiple_choice":
          payload = { ...basePayload, id: cardId, type: cardType };
          break;
        case "ordered":
          payload = {...basePayload, id: cardId, type: cardType };
          break;
        default:
          throw new Error("Unknown card type");
      }

      const res = await updateCard(deckId, cardId, payload);

      if (res.success) {
        form.reset({});
        onOpenChange(false);
        alert("Card updated successfully!");
        router.refresh();
      } else {
        alert(`Failed to update card: ${res.message}`);
      }
    } catch (err) {
      console.error(err);
      alert("Something went wrong updating the card");
    } finally {
      setLoading(false);
    }
  });

  const renderFields = () => {
    switch (cardType) {
      case "front_back":
        return <FrontBackField form={form as UseFormReturn<CardForms["front_back"]>} />;
      case "blanks":
        return <BlankFields form={form as UseFormReturn<CardForms["blanks"]>} />;
      case "multiple_choice":
        return <MultipleChoiceFields form={form as UseFormReturn<CardForms["multiple_choice"]>} />;
      case "ordered":
        return <OrderedFields form={form as UseFormReturn<CardForms["ordered"]>} />;
      default:
        return null;
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Card</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={onSubmit} className="flex flex-col gap-4 py-4">
            {renderFields()}
            <Button type="submit" disabled={loading} className="mt-4">
              {/*Ensures that the button is disabled while loading*/}
              {loading ? "Loading..." : "Edit card"}
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
