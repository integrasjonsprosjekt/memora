export default async function DeckTodayPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;

  return (
    <div>
      <h1>Viewing deck {slug}&apos;s daily overview</h1>
    </div>
  );
}
