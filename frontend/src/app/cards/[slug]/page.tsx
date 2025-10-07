export default async function CardPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;

  const card = await fetch(`${process.env.API_URI}/v1/cards/${slug}`).then((res) => res.json());

  return (
    <div>
      <h1>Viewing card {slug}</h1>
    </div>
  );
}
