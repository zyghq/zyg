import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
} from "react";

// const NODE_ENV = process.env.NODE_ENV || "development";

export const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

export function useAuthProvider() {
  const [authUser, setAuthUser] = useState(null);
  const [isAuthLoading, setIsAuthLoading] = useState(true);

  console.log("*************** useAuthProvider *****************");
  console.log("authUser", authUser);
  console.log("isAuthLoading", isAuthLoading);
  console.log("*************** useAuthProvider *****************");

  const postMessage = useCallback((data) => {
    let postable;
    if (typeof data === "object") {
      postable = JSON.stringify(data);
    }
    postable = data;
    window.parent.postMessage(postable, "*");
  }, []);

  const logout = useCallback(async () => {
    try {
      // safer way is to delete from the server.
      const resp = await fetch("/api/auth/", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (resp.ok) {
        const data = await resp.json();
        const { ok } = data;
        if (ok) {
          setAuthUser(null);
          postMessage("auth:logout");
        }
        return;
      }
      console.log("unsure about logout.");
    } catch (err) {
      console.log("error when logging out with error", err);
    }
  }, [postMessage]);

  const authenticate = useCallback(
    async (token) => {
      try {
        const resp = await fetch("/api/auth/", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ token }),
        });
        if (resp.ok) {
          const data = await resp.json();
          setAuthUser({ ...data });
          postMessage("auth:authenticated");
        }
        const { status } = resp;
        if (status === 401) {
          await logout();
          postMessage("auth:logout");
        }
      } catch (err) {
        console.log("authenticate error:", err);
        postMessage("auth:error");
      } finally {
        setIsAuthLoading(false);
      }
    },
    [postMessage, logout]
  );

  const me = useCallback(async () => {
    try {
      const resp = await fetch("/api/auth/", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (resp.ok) {
        const data = await resp.json();
        setAuthUser({ ...data });
        return;
      }
    } catch (err) {
      console.log("me error:", err);
    } finally {
      setIsAuthLoading(false);
    }
  }, []);

  // onload
  useEffect(() => {
    const unsubscribe = async () => {
      await me();
      // if (NODE_ENV === "development") {
      //   await me();
      // }
    };
    unsubscribe();
  }, [me]);

  // on window post message
  useEffect(() => {
    const onMessageHandler = (e) => {
      if (typeof e.data === "object") {
        return;
      }

      if (e.data === "auth:logout") return;
      if (e.data === "auth:authenticated") return;
      if (e.data === "auth:error") return;

      try {
        const data = JSON.parse(e.data);
        const { event = "", payload = {} } = data;
        if (event === "authenticate") {
          // message to authenticate
          const authToken = payload?.authToken || "";
          console.log("authToken", authToken);
          authenticate(authToken); // not worried about awaiting
        }
      } catch (err) {
        console.log(e.data);
        console.log(typeof e.data);
        console.error("evt message parse err:", err);
      }
    };
    window.addEventListener("message", onMessageHandler);
    return () => {
      window.removeEventListener("message", onMessageHandler);
    };
  }, [authenticate]);

  return {
    authUser,
    isAuthLoading,
    logout,
  };
}
