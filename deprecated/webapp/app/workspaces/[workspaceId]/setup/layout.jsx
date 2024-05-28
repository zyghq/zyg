import { redirect } from "next/navigation";
import { createClient } from "@/utils/supabase/server";
import { isAuthenticated } from "@/utils/supabase/helpers";

export default async function WorkspaceSetupLayout({ children }) {
  const supabase = createClient();

  if (!(await isAuthenticated(supabase))) {
    return redirect("/login/");
  }

  return <div className="flex flex-col justify-center p-4">{children}</div>;
}
