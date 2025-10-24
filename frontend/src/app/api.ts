'use server';

import { getApiEndpoint } from '@/config/api';
import { CardType } from '@/types/card';

interface UpdateDeckPayload {
  title?: string;
  addedEmails?: string[];
  removedEmails?: string[];
}

export async function createCard(id: string, type: CardType, data: Record<string, unknown>) {
  const body = { type, ...data };

  const response = await fetch(getApiEndpoint(`/v1/decks/${id}/cards`), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });

  const result = await response.json();

  if (!response.ok) {
    return { success: false, message: result.error || 'Failed to create card' };
  }

  return { success: true, data: result };
}

export async function updateCard(deckId: string, cardId: string, data: Record<string, unknown>) {
  const body = { ...data };

  const response = await fetch(getApiEndpoint(`/v1/decks/${deckId}/cards/${cardId}`), {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });

  const result = await response.json();

  if (!response.ok) {
    return { success: false, message: result.error || 'Failed to update card' };
  }

  return { success: true, data: result };
}

export async function deleteCard(deckId: string, cardId: string) {
  const response = await fetch(getApiEndpoint(`/v1/decks/${deckId}/cards/${cardId}`), {
    method: 'DELETE',
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
    return { success: false, message: result?.error || 'Failed to delete card' };
  }

  return { success: true, data: result };
}

// ---------- Decks ----------

export async function createDeck(title: string, shared_emails?: string[]) {
  const body = {
    owner_id: 'zb8FOCiyYCfXmqpPpxaU7Kjm5p92', // TODO: Remove hardcoded user ID
    title,
    shared_emails,
  };

  const response = await fetch(getApiEndpoint('/v1/decks'), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });

  const result = await response.json();

  if (!response.ok) {
    return { success: false, message: result.error || 'Failed to create deck' };
  }

  return { success: true, data: result };
}

export async function updateDeck(deckId: string, payload: UpdateDeckPayload) {
  const responseTitle = await fetch(getApiEndpoint(`/v1/decks/${deckId}`), {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title: payload.title }),
  });

  const resultTitle = await responseTitle.json();

  if (!responseTitle.ok) {
    return { success: false, message: resultTitle.error || 'Failed to update deck title' };
  }

  const addBody = {
    opp: 'add',
    shared_emails: payload.addedEmails,
  };
  const removeBody = {
    opp: 'remove',
    shared_emails: payload.removedEmails,
  };

  const emailResponses: Response[] = [];

  if (payload.addedEmails?.length) {
    emailResponses.push(
      await fetch(getApiEndpoint(`/v1/decks/${deckId}/emails`), {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(addBody),
      })
    );
  }
  if (payload.removedEmails?.length) {
    emailResponses.push(
      await fetch(getApiEndpoint(`/v1/decks/${deckId}/emails`), {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(removeBody),
      })
    );
  }

  let resultEmails;
  for (const res of emailResponses) {
    resultEmails = await res.json().catch(() => []);
    if (!res.ok) {
      return { success: false, message: resultEmails.error || 'Failed to update deck emails' };
    }
  }

  return { success: true, data: resultTitle + resultEmails };
}

export async function deleteDeck(deckId: string) {
  const response = await fetch(getApiEndpoint(`/v1/decks/${deckId}`), {
    method: 'DELETE',
  });

  if (!response.ok) {
    return { success: false, message: 'Failed to delete deck' };
  }

  return { success: true };
}
