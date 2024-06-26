import { redirect } from "next/navigation";
import { createClient } from "@/utils/supabase/server";

export const metadata = {
  title: "Reset Password | Zyg AI",
};

export default async function Layout({ children }) {
  const supabase = createClient();
  const { data, error } = await supabase.auth.getUser();
  if (error || !data?.user) {
    redirect("/login/");
  }

  return (
    <div className="flex min-h-screen flex-col justify-center p-4">
      {children}
    </div>
  );
}
