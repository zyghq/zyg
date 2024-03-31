/**
 * God bless this code
 * This callback route made sure we are able to reset password
 * This is the callback route for the reset password form.
 *
 * How this works:
 * In short - you get a token from the supabase in the email reset link,
 * - that link is checked against if its valid or not, then
 * - this route is called with `code` query params that is used to authenticate the user
 * - now there are lot of other stuff going on here for the security stuff, but this the short version of it.
 *
 * More: https://supabase.com/docs/guides/auth/auth-helpers/nextjs#install-nextjs-auth-helpers-library
 *
 * */

import { cookies } from "next/headers";
import { createClient } from "@/utils/supabase/actions";
import { NextResponse } from "next/server";

export async function GET(request) {
  const requestUrl = new URL(request.url);
  const code = requestUrl.searchParams.get("code");
  const next = requestUrl.searchParams.get("next") || `/auth/reset/`;

  if (code) {
    const cookieStore = cookies();
    const supabase = createClient(cookieStore);
    await supabase.auth.exchangeCodeForSession(code);
  }

  const url = `${requestUrl.origin}${next}`;
  // redirect to as specified by `next` otherwise redirect auth reset
  return NextResponse.redirect(url);
}
