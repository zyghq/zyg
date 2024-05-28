"use server";

/**
 * SA - server actions to reset password as per supabase docs
 * here: https://supabase.com/docs/guides/auth/auth-password-reset
 *
 *
 * The password reset form is protected route hence, the user should be already
 * authenticated before accessing this route. We are using callback route to authenticate first.
 */

import { createClient } from "@/utils/supabase/server";

export async function reset(values) {
  const supabase = createClient();
  const { password } = values;
  try {
    const { data, error } = await supabase.auth.updateUser({
      password,
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
    console.log("*** response after reset ***");
    console.log(data);
    console.log("*** response after reset ***");
    return {
      ok: true,
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
