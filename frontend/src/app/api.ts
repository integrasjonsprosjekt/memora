"use server";

import { CardType } from "@/types/cards";

export async function createCard(id: string, type: CardType, data: Record<string, unknown>) {
  const body = { type, ...data };

  const response = await fetch(`${process.env.API_URI}/v1/decks/${id}/cards`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  const result = await response.json();

  if (!response.ok) {
    return { success: false, message: result.error || "Failed to create card" };
  }

  return { success: true, data: result };
}

export async function updateCard(deck_id: string, card_id: string, data: Record<string, unknown>) {
  const body = { ...data };

  const response = await fetch(`${process.env.API_URI}/v1/decks/${deck_id}/cards/${card_id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  const result = await response.json();

  if (!response.ok) {
    return { success: false, message: result.error || "Failed to update card" };
  }

  return { success: true, data: result };
}

export async function deleteCard(deck_id: string, card_id: string) {
  const response = await fetch(`${process.env.API_URI}/v1/decks/${deck_id}/cards/${card_id}`, {
    method: "DELETE",
  });

  const result = await response.json();

  if (!response.ok) {
    return { success: false, message: result.error || "Failed to delete card" };
  }

  return { success: true, data: result };
}
