export default async function DeckDashboardPage({ params }: { params: Promise<{ deckId: string }> }) {
  const { deckId } = await params;

  return (
    <div>
      <h1>Viewing deck {deckId}&apos;s dashboard</h1>
    </div>
  );
}
