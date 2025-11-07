export default function Layout({ children }: { children: React.ReactNode }) {
  return <div className="m-auto flex h-lvh w-lvw flex-col items-center justify-center">{children}</div>;
}
