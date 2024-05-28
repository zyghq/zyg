export default async function DashboardLayout({ children }) {
  return (
    <div vaul-drawer-wrapper="">
      <div className="flex flex-col">{children}</div>
    </div>
  );
}
