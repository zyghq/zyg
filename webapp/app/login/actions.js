"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { createClient } from "@/utils/supabase/actions";

async function getOrCreateZygAccount(token) {
  console.log("getOrCreateZygAccount token ->", token);
  try {
    const response = await fetch(`${process.env.ZYG_API_URL}/accounts/auth/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({}),
    });
    if (!response.ok) {
      const { status, statusText } = response;
      return [
        new Error(
          `error creating Zyg auth account with withs status: ${status} and statusText: ${statusText}`
        ),
        null,
      ];
    }
    const data = await response.json();
    return [null, data];
  } catch (err) {
    return [err, null];
  }
}

export async function login(values) {
  const cookieStore = cookies();
  const supabase = createClient(cookieStore);
  try {
    const { data, error } = await supabase.auth.signInWithPassword({
      ...values,
    });
    //
    // catch error from supabase
    if (error) {
      console.log(`error: ${error}`); // notify webmaster
      const { name, message } = error;
      return {
        ok: null,
        error: {
          name: name,
          message: message,
        },
      };
    }
    const { session } = data;
    const { access_token } = session;
    const [accountErr, zygAccount] = await getOrCreateZygAccount(access_token);
    if (accountErr) {
      throw accountErr;
    }
    console.log("*** response after login ***");
    console.log(zygAccount);
    console.log("*** response after login ***");
    const { email } = zygAccount;
    revalidatePath("/");
    return {
      ok: true,
      error: null,
      data: {
        email: email,
      },
    };
  } catch (err) {
    //
    // catch all errors - notify webmaster
    console.error(err);
    return {
      ok: null,
      data: null,
      error: {
        name: "oops",
        message: "something went wrong!",
      },
    };
  }
}
