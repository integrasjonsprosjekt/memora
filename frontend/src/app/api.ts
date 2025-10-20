"use server";

import { CardType } from "@/types/cards";

interface UpdateDeckPayload {
  title?: string;
  addedEmails?: string[];
  removedEmails?: string[];
}

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

  // 204 means "No Content" — don't try to parse JSON
  if (response.status === 204) {
    return { success: true };
  }

  let result;
  try {
    result = await response.json();
  } catch {
    result = null;
  }

  if (!response.ok) {
    return { success: false, message: result?.error || "Failed to delete card" };
  }

  return { success: true, data: result };
}

// ---------- Decks ----------

export async function createDeck(title: string, shared_emails?: string[]) {
  const body = {
    owner_id: "zb8FOCiyYCfXmqpPpxaU7Kjm5p92", // TODO: Remove hardcoded user ID
    title,
    shared_emails
  };

  const response = await fetch(`${process.env.API_URI}/v1/decks`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  const result = await response.json();

  if (!response.ok) {
    return { success: false, message: result.error || "Failed to create deck" };
  }

  return { success: true, data: result };
}

export async function updateDeck(deck_id: string, payload: UpdateDeckPayload) {

  const responseTitle = await fetch(`${process.env.API_URI}/v1/decks/${deck_id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({title: payload.title}),
  });

  const resultTitle = await responseTitle.json();

  if (!responseTitle.ok) {
    return { success: false, message: resultTitle.error || "Failed to update deck title" };
  }

  const addBody = {
    opp: "add",
    shared_emails: payload.addedEmails
  }
  const removeBody = {
    opp: "remove",
    shared_emails: payload.removedEmails
  }

  const emailResponses: Response[] = []

  if (payload.addedEmails?.length) {
    emailResponses.push(await fetch(`${process.env.API_URI}/v1/decks/${deck_id}/emails`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(addBody),
    }));
  }
  if (payload.removedEmails?.length) {
    emailResponses.push(await fetch(`${process.env.API_URI}/v1/decks/${deck_id}/emails`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(removeBody),
    }));
  }

  let resultEmails;
  for (const res of emailResponses) {
    resultEmails = await res.json().catch(() => []);
    if (!res.ok) {
      return { success: false, message: resultEmails.error || "Failed to update deck emails" };
    }
  }

  return { success: true, data: resultTitle + resultEmails };
}

export async function deleteDeck(deck_id: string) {
  const response = await fetch(`${process.env.API_URI}/v1/decks/${deck_id}`, {
    method: "DELETE",
  });

  if (!response.ok) {
    return { success: false, message: "Failed to delete deck" };
  }

  return { success: true };
}
