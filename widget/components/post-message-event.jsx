"use client";

import * as React from "react";
import { useAuth } from "@/lib/auth";

export default function PostMessageEvent() {
  const auth = useAuth();
  const { authUser, isAuthLoading } = auth;

  console.log("authUser", authUser);
  console.log("isAuthLoading", isAuthLoading);

  // const [authUser, setAuthUser] = React.useState(null);

  // const logout = React.useCallback(async () => {
  //   try {
  //     // fastest way is to delete from the client side if possible.
  //     // for this to work httpOnly: false - at the server side setting the cookie.
  //     Cookies.remove("__zygtoken");
  //     // safer way is to delete from the server side.
  //     const resp = await fetch("/api/auth/", {
  //       method: "DELETE",
  //       headers: {
  //         "Content-Type": "application/json",
  //       },
  //     });
  //     if (resp.ok) {
  //       const data = await resp.json();
  //       console.log("log out", data);
  //     }
  //   } catch (err) {
  //     console.log("logout error:", err);
  //   }
  // }, []);

  // const authenticate = React.useCallback(
  //   async (token) => {
  //     try {
  //       const resp = await fetch("/api/auth/", {
  //         method: "POST",
  //         headers: {
  //           "Content-Type": "application/json",
  //         },
  //         body: JSON.stringify({ token }),
  //       });
  //       if (resp.ok) {
  //         const data = await resp.json();
  //         const { customerId } = data;
  //         console.log(`authenticated customer: ${customerId}`);
  //       }
  //       const { status } = resp;
  //       if (status === 401) {
  //         logout();
  //       }
  //     } catch (err) {
  //       console.log("authentication error:", err);
  //     }
  //   },
  //   [logout]
  // );

  // React.useEffect(() => {
  //   const onMessageHandler = (e) => {
  //     if (typeof e.data === "object") {
  //       return;
  //     }
  //     try {
  //       const data = JSON.parse(e.data);
  //       const { event = "", payload = {} } = data;
  //       if (event === "authenticate") {
  //         const { accessToken } = payload;
  //         authenticate(accessToken);
  //       }
  //     } catch (err) {
  //       console.log(e.data);
  //       console.log(typeof e.data);
  //       console.error("evt message parse err:", err);
  //     }
  //   };

  //   window.addEventListener("message", onMessageHandler);
  //   return () => {
  //     window.removeEventListener("message", onMessageHandler);
  //   };
  // }, [authenticate]);

  return <div></div>;
}
