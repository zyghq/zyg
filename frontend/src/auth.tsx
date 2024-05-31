import React from "react";
import { SupabaseClient, Session, User } from "@supabase/supabase-js"; // Import the type if using TypeScript

export type AuthContextType = {
  client: SupabaseClient;
  session: Session | null;
  user: User | null;
  isLoading: boolean;
};

export const createAuthContext = (supaClient: SupabaseClient) => {
  const Context = React.createContext<AuthContextType | null>(null);

  const Provider: React.FC<{ children?: React.ReactNode }> = ({ children }) => {
    const [client] = React.useState(() => supaClient);

    const [isLoading, setIsLoading] = React.useState<boolean>(true);
    const [session, setSession] = React.useState<Session | null>(null);
    const [user, setUser] = React.useState<User | null>(null);

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
      isLoading,
      session,
      user,
    };

    return <Context.Provider value={value}>{children}</Context.Provider>;
  };

  return {
    Context,
    Provider,
    useContext: () => React.useContext(Context),
    client: supaClient,
  };
};
