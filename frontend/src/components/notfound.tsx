import { Link } from "@tanstack/react-router";

export function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[100dvh] px-4 md:px-6 text-center">
      <div className="max-w-md space-y-4">
        <h1 className="text-5xl font-bold tracking-tighter">404</h1>
        <p className="text-muted-foreground text-lg">
          Oops, the page you're looking for doesn't exist.
        </p>
        <Link
          className="inline-flex h-10 items-center justify-center rounded-md bg-primary px-8 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50"
          to="/workspaces"
        >
          Go to Workspaces
        </Link>
      </div>
    </div>
  );
}
