import { Header } from "@/components/settings/header";
import Sidebar from "@/components/settings/sidebar";

export const metadata = {
  title: "Settings - Zyg AI",
};

export default function SettingsLayout({ children }) {
  return (
    <div vaul-drawer-wrapper="">
      <div className="flex flex-col">
        <Header />
        <div className="flex flex-col">
          <div className="flex">
            <div className="hidden min-w-80 flex-col border-r lg:flex">
              <Sidebar className="h-[calc(100dvh-8rem)]" />
            </div>
            <main className="flex-1">{children}</main>
          </div>
        </div>
      </div>
    </div>
  );
}
