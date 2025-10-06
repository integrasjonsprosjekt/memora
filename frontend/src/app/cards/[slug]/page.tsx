export default async function CardPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;

  return (
    <div>
      <h1>Viewing card {slug}</h1>
    </div>
  );
}
