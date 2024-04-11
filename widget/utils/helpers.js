export const isAuthenticated = async (cookies) => {
  if (!cookies) {
    return false;
  }
  const token = cookies.get("__zygtoken");
  if (!token) {
    return false;
  }
  const { value = "" } = token;
  try {
    const resp = await fetch(`${process.env.ZYG_API_URL}/-/me/`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${value}`,
      },
    });
    if (!resp.ok) {
      return false;
    }
    const data = await resp.json();
    console.log("authenticated:", data);
    return true;
  } catch (err) {
    console.error("error when authenticating customer:", err);
    return false;
  }
};
