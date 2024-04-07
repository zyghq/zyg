import WidgetCloseButton from "@/components/widget-close-btn";
import GoHomeButton from "@/components/gohome-btn";

export function Header({ label }) {
  return (
    <header className="flex items-center justify-between p-4">
      <div className="text-xl">{label}</div>
      <WidgetCloseButton />
    </header>
  );
}

export function ThreadHeader() {
  return (
    <div className="flex items-center justify-start py-4 border-b px-1">
      <GoHomeButton />
      <div>
        <div className="flex flex-col">
          <div className="font-semibold">Zyg Team</div>
          <div className="text-xs text-muted-foreground">
            Ask us anything, share your feedback.
          </div>
        </div>
      </div>
      <div className="ml-auto">
        <WidgetCloseButton />
      </div>
    </div>
  );
}
