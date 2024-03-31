import { Header } from "@/components/header";

export const metadata = {
  title: "All Threads - Zyg AI",
};

export default function DashboardLayout({ children }) {
  return (
    <div vaul-drawer-wrapper="">
      <div className="flex flex-col">
        <Header />
        <div className="flex flex-1 flex-col">{children}</div>
      </div>
    </div>
  );
}
