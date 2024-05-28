"use server";

import { revalidatePath } from "next/cache";
import { createClient } from "@/utils/supabase/server";

export async function signup(values) {
  const supabase = createClient();
  try {
    const { data, error } = await supabase.auth.signUp({ ...values });
    //
    // catch error from supabase
    if (error) {
      const { name, message } = error;
      return {
        ok: null,
        data: null,
        error: {
          name,
          message,
        },
      };
    }
    console.log("*** response after signup ***");
    console.log(data);
    console.log("*** response after signup ***");
    const { user } = data;
    const { email } = user;
    revalidatePath("/login/", "layout");
    return {
      ok: true,
      data: {
        email: email,
      },
      error: null,
    };
  } catch (err) {
    //
    // catch all errors - notify webmaster
    console.error(err);
    return {
      ok: null,
      error: {
        name: "oops",
        message: "something went wrong!",
      },
    };
  }
}
