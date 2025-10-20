"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { EditCardMenu } from "./edit-card-menu";
import { CardType } from "@/types/cards";

interface CardButtonProps {
  deckId: string;
  cardId: string;
  cardType: CardType;
  initialData: Record<string, unknown>;
}

export function EditCardButton({ deckId, cardId, cardType, initialData }: CardButtonProps) {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button
        className="mx-20"
        onClick={() => setOpen(true)}
      >
        Edit card
      </Button>
      <EditCardMenu open={open} onOpenChange={setOpen} deckId={deckId} cardId={cardId} cardType={cardType} initialData={initialData} />
    </>
  );
}
