/**
 * Checks if the user is authenticated.
 * @param {Object} supabase - The Supabase client object.
 * @returns {Promise<boolean>} - A promise that resolves to a boolean indicating whether the user is authenticated or not.
 */
export const isAuthenticated = async (supabase) => {
  const { data, error } = await supabase.auth.getUser();
  if (error || !data?.user) {
    return false;
  }
  return true;
};

/**
 * Retrieves the authentication token from Supabase.
 * @param {Object} - The Supabase client object.
 * @returns {Promise<string>} - The authentication token.
 * @deprecated This function is deprecated. Please use the `getSession` function instead.
 */
export const getAuthToken = async (supabase) => {
  const { data, error } = await supabase.auth.getSession();
  if (error) return "";
  const accessToken = data?.session?.access_token || "";
  return accessToken;
};

/**
 * Retrieves the session from Supabase.
 * A thin wrapper on top of `supabase.auth.getSession`.
 * @param {Object} - The Supabase client object.
 * @returns {Promise<{token: String, error: Object}>} - The session string and error object.
 */
export const getSession = async (supabase) => {
  const { data, error } = await supabase.auth.getSession();
  if (error) {
    console.error("error fetching supabase auth session", error);
    return { token: null, error };
  }
  const accessToken = data?.session?.access_token || "";
  return { token: accessToken, error: null };
};
