import { Sidebar } from "@/components/sidebar";
import Title from "@/components/title";
import ThreadsTab from "@/components/threads-tab";

// TODO: pass the API server data from page/layout to child components.
// example we can pass the threads data to the ThreadTabs component.
export default function DashboardPage() {
  return (
    <div className="grid lg:grid-cols-5">
      <Sidebar className="hidden border-r lg:block" />
      <main className="col-span-3 lg:col-span-4">
        <div className="container">
          <Title title="Threads" />
          <ThreadsTab />
        </div>
      </main>
    </div>
  );
}
