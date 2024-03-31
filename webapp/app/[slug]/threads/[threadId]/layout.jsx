// import { Header } from "@/components/header";

export const metadata = {
  title: "All Threads - Zyg AI",
};

export default function ThreadItemLayout({ children }) {
  return (
    <div vaul-drawer-wrapper="">
      <div className="flex flex-col min-h-screen">
        <div className="flex flex-col flex-1">{children}</div>
      </div>
    </div>
  );
}
