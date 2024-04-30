// import { Header } from "@/components/header";

export const metadata = {
  title: "All Threads - Zyg AI",
};

export default function ThreadItemLayout({ children }) {
  return (
    <div vaul-drawer-wrapper="">
      <div className="flex min-h-screen flex-col">
        <div className="flex flex-1 flex-col">{children}</div>
      </div>
    </div>
  );
}
