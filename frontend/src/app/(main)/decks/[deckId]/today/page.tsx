export default async function DeckTodayPage({ params }: { params: Promise<{ deckId: string }> }) {
  const { deckId } = await params;

  return (
    <div>
      <h1>Viewing deck {deckId}&apos;s daily overview</h1>
    </div>
  );
}
