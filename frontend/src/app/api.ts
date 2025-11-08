import { fetchApi } from '@/lib/api/config';
import { CardType } from '@/types/card';
import { User } from 'firebase/auth';

interface UpdateDeckPayload {
  title?: string;
  addedEmails?: string[];
  removedEmails?: string[];
}

export async function createCard(user: User, id: string, type: CardType, data: Record<string, unknown>) {
  try {
    const body = { type, ...data };

    const result = await fetchApi(`decks/${id}/cards`, {
      user,
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });

    return { success: true, data: result };
  } catch (error) {
    console.error('Failed to create card:', error);
    return { success: false, message: 'Failed to create card' };
  }
}

export async function updateCard(user: User, deckId: string, cardId: string, data: Record<string, unknown>) {
  try {
    const body = { ...data };

    const result = await fetchApi(`decks/${deckId}/cards/${cardId}`, {
      user,
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });

    return { success: true, data: result };
  } catch (error) {
    console.error('Failed to update card:', error);
    return { success: false, message: 'Failed to update card' };
  }
}

export async function deleteCard(user: User, deckId: string, cardId: string) {
  try {
    await fetchApi(`decks/${deckId}/cards/${cardId}`, {
      user,
      method: 'DELETE',
    });

    return { success: true };
  } catch (error) {
    console.error('Failed to delete card:', error);
    return { success: false, message: 'Failed to delete card' };
  }
}

// ---------- Decks ----------

export async function createDeck(user: User, owner_id: string, title: string, shared_emails?: string[]) {
  try {
    const body = {
      owner_id,
      title,
      shared_emails,
    };

    const result = await fetchApi('decks', {
      user,
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });

    return { success: true, data: result };
  } catch (error) {
    console.error('Failed to create deck:', error);
    return { success: false, message: 'Failed to create deck' };
  }
}

export async function updateDeck(user: User, deckId: string, payload: UpdateDeckPayload) {
  try {
    // Update title
    const resultTitle = await fetchApi(`decks/${deckId}`, {
      user,
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: payload.title }),
    });

    // Update emails if needed
    if (payload.addedEmails?.length) {
      await fetchApi(`decks/${deckId}/emails`, {
        user,
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          opp: 'add',
          shared_emails: payload.addedEmails,
        }),
      });
    }

    if (payload.removedEmails?.length) {
      await fetchApi(`decks/${deckId}/emails`, {
        user,
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          opp: 'remove',
          shared_emails: payload.removedEmails,
        }),
      });
    }

    return { success: true, data: resultTitle };
  } catch (error) {
    console.error('Failed to update deck:', error);
    return { success: false, message: 'Failed to update deck' };
  }
}

export async function deleteDeck(user: User, deckId: string) {
  try {
    await fetchApi(`decks/${deckId}`, {
      user,
      method: 'DELETE',
    });

    return { success: true };
  } catch (error) {
    console.error('Failed to delete deck:', error);
    return { success: false, message: 'Failed to delete deck' };
  }
}
