"use client";

import { useEffect, useMemo, useState } from "react";
import { useForm, UseFormReturn } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Form } from "./ui/form";

import CardTypeSelector from "./card-type-selector";
import { FrontBackField } from "./front-back-field";
import { BlankFields } from "./blank-field";
import { MultipleChoiceFields } from "./multiple-choice-field";

import { cardSchemas } from "@/lib/cardSchemas";
import buildCardPayload from "@/lib/cardBuildPayload";
import { CardType } from "@/types/cards";
import { createCard } from "@/app/api";
import { Button } from "./ui/button";
import { OrderedFields } from "./ordered-field";

interface AddCardMenuProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

type CardForms = {
  front_back: z.infer<typeof cardSchemas.front_back>;
  blanks: z.infer<typeof cardSchemas.blanks>;
  multiple_choice: z.infer<typeof cardSchemas.multiple_choice>;
  ordered: z.infer<typeof cardSchemas.ordered>;
};

export function AddCardMenu({ open, onOpenChange }: AddCardMenuProps) {
  const [cardType, setCardType] = useState<CardType>("front_back");
  const [loading, setLoading] = useState(false);

  // Pick the schema for the current card type dynamically
  const currentSchema = useMemo(() => cardSchemas[cardType], [cardType]);

  const form = useForm<CardForms[typeof cardType]>({
    resolver: zodResolver(currentSchema),
    defaultValues: {},
  });

  useEffect(() => {
    form.reset({});
  }, [cardType, form]);

  const onSubmit = form.handleSubmit(async (values) => {
    setLoading(true);
    try {
      const payload = buildCardPayload(cardType, values);
      const res = await createCard(cardType, payload);

      if (res.success) {
        form.reset({});
        onOpenChange(false);
        alert("Card created successfully!");
      } else {
        alert(`Failed to create card: ${res.message}`);
      }
    } catch (err) {
      console.error(err);
      alert("Something went wrong creating the card");
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
          <DialogTitle>Create Card</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={onSubmit} className="flex flex-col gap-4 py-4">
            <CardTypeSelector value={cardType} onChange={setCardType} />
            {renderFields()}
            <Button type="submit" disabled={loading} className="mt-4">
              {/*Ensures that the button is disabled while loading*/}
              {loading ? "Loading..." : "Add card"}
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
