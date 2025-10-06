export default async function DeckDashboardPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;

  return (
    <div>
      <h1>Viewing deck {slug}&apos;s dashboard</h1>
    </div>
  );
}
