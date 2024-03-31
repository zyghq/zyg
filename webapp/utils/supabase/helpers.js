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
 */
export const getAuthToken = async (supabase) => {
  const { data, error } = await supabase.auth.getSession();
  if (error) return "";
  const accessToken = data?.session?.access_token || "";
  return accessToken;
};
