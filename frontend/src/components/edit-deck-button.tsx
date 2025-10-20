'use client'

import { useState } from "react";
import { Button } from "./ui/button";
import { EditDeckMenu } from "./edit-deck-menu";

interface EditDeckButtonProps {
  deckId: string;
  initialData: Record<string, unknown>;
};

export function EditDeckButton({deckId, initialData}: EditDeckButtonProps) {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button onClick={() => setOpen(true)}>
        Edit Deck
      </Button>
      <EditDeckMenu open={open} onOpenChange={setOpen} deckId={deckId} initialData={initialData} />
    </>
  );
}
