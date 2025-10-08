"use server";

import { CardType } from "@/types/cards";

const apiUrl = "http://localhost:8080/api/v1";

export async function createCard(type: CardType, data: Record<string, unknown>) {
  const body = { type, ...data };

  const response = await fetch(`${apiUrl}/cards`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  const result = await response.json();

  if (!response.ok) {
    // Return a failure object instead of throwing
    return { success: false, message: result.error || "Failed to create card" };
  }

  // Success
  return { success: true, data: result };
}
