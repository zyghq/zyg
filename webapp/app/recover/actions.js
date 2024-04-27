"use server";

import { createClient } from "@/utils/supabase/server";

export async function recover(values) {
  const supabase = createClient();
  const { email } = values;
  try {
    const { data, error } = await supabase.auth.resetPasswordForEmail(email, {
      redirectTo: "http://localhost:3000/auth/callback/?next=/auth/reset/", // make sure to add this to Redirect Urls.
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
    console.log("*** response after recover ***");
    console.log(data);
    console.log("*** response after recover ***");
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
