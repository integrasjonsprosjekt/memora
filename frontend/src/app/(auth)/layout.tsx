export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="w-lvh h-lvh m-auto flex flex-col items-center justify-center">
      {children}
    </div>
  );
}
