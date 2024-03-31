// import { cookies } from "next/headers";
// import { createClient } from "@/utils/supabase/actions";
// import { NextResponse } from "next/server";

// export async function POST(request) {
//   const cookieStore = cookies();
//   const supabase = createClient(cookieStore);
//   const json = await request.json();
//   const user = await supabase.auth.getUser();
//   console.log(user);
//   console.log(json);
//   return NextResponse.json(json, { status: 200 });
// }
