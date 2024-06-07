// import React from "react";
// import { Session, createClient, SupabaseClient } from "@supabase/supabase-js";

// export interface AuthContext {
//   client: SupabaseClient;
//   session: Session | null;
//   isLoading: boolean;
//   isAuthenticated: boolean;
// }

// const supaClient = createClient(
//   import.meta.env.VITE_SUPABASE_URL,
//   import.meta.env.VITE_SUPABASE_ANON_KEY
// );

// export const AuthContext = React.createContext<AuthContext | null>(null);

// export const AuthProvider: React.FC<{ children?: React.ReactNode }> = ({
//   children,
// }) => {
//   const [client] = React.useState(() => supaClient);
//   const [session, setSession] = React.useState<Session | null>(null);
//   const [isLoading, setIsLoading] = React.useState(true);

//   const isAuthenticated = React.useMemo(() => !!session, [session]);

//   React.useEffect(() => {
//     const { data: listener } = supaClient.auth.onAuthStateChange(
//       (_event, session) => {
//         setSession(session);
//         setIsLoading(false);
//       }
//     );

//     const setData = async () => {
//       const {
//         data: { session },
//         error,
//       } = await supaClient.auth.getSession();
//       if (error) {
//         throw error;
//       }

//       setSession(session);
//       setIsLoading(false);
//     };

//     setData();

//     return () => {
//       listener?.subscription.unsubscribe();
//     };
//   }, []);

//   // console.log("*** in auth context start ****");
//   // console.log("isAuthenticated", isAuthenticated);
//   // console.log("session", session);
//   // console.log("isLoading", isLoading);
//   // console.log("*** in auth context end ****");

//   const value = {
//     client,
//     session,
//     isLoading,
//     isAuthenticated,
//   };

//   return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
// };

// // eslint-disable-next-line react-refresh/only-export-components
// export function useAuth() {
//   const context = React.useContext(AuthContext);
//   if (!context) {
//     throw new Error("useAuth must be used within a AuthProvider");
//   }
//   return context;
// }
