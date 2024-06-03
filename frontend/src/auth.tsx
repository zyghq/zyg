import React from "react";
import {
  Session,
  User,
  createClient,
  // AuthResponse,
  SupabaseClient,
} from "@supabase/supabase-js";

export interface AuthContext {
  client: SupabaseClient;
  // signInWithPassword: (
  //   email: string,
  //   password: string
  // ) => Promise<AuthResponse>;
  // signUp: (email: string, password: string) => Promise<AuthResponse>;
  session: Session | null;
  user: User | null;
  isLoading: boolean;
}

const supaClient = createClient(
  import.meta.env.VITE_SUPABASE_URL,
  import.meta.env.VITE_SUPABASE_ANON_KEY
);

export const AuthContext = React.createContext<AuthContext | null>(null);

export const AuthProvider: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  const [client] = React.useState(() => supaClient);
  const [session, setSession] = React.useState<Session | null>(null);
  const [user, setUser] = React.useState<User | null>(null);
  const [isLoading, setIsLoading] = React.useState(true);

  // const signInWithPassword = React.useCallback(
  //   async (email: string, password: string) => {
  //     const { data, error } = await supaClient.auth.signInWithPassword({
  //       email,
  //       password,
  //     });
  //     if (error) {
  //       console.error("error signing with password", error);
  //       return { error, data };
  //     }
  //     setSession(data.session);
  //     setUser(data.session?.user || null);
  //     return { error, data };
  //   },
  //   []
  // );

  // const signUp = React.useCallback(async (email: string, password: string) => {
  //   const { data, error } = await supaClient.auth.signUp({
  //     email,
  //     password,
  //   });
  //   if (error) {
  //     console.error("error signing up", error);
  //     return { data, error };
  //   }
  //   setSession(data.session);
  //   setUser(data.session?.user || null);
  //   return { error, data };
  // }, []);

  React.useEffect(() => {
    const { data: listener } = supaClient.auth.onAuthStateChange(
      (_event, session) => {
        setSession(session);
        setUser(session?.user || null);
        setIsLoading(false);
      }
    );

    const setData = async () => {
      const {
        data: { session },
        error,
      } = await supaClient.auth.getSession();
      if (error) {
        throw error;
      }

      setSession(session);
      setUser(session?.user || null);
      setIsLoading(false);
    };

    setData();

    return () => {
      listener?.subscription.unsubscribe();
    };
  }, []);

  const value = {
    client,
    session,
    user,
    isLoading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  const context = React.useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within a AuthProvider");
  }
  return context;
}
