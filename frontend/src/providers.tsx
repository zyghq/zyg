import React from "react";
import { StoreApi } from "zustand";

import { buildStore, buildAccountStore } from "@/db/store";

const createZustandContext = <TInitial, TStore extends StoreApi<any>>(
  getStore: (initial: TInitial) => TStore
) => {
  const Context = React.createContext(null as any as TStore);

  const Provider = (props: {
    children?: React.ReactNode;
    initialValue: TInitial;
  }) => {
    const [store] = React.useState(() => getStore(props.initialValue));

    return <Context.Provider value={store}>{props.children}</Context.Provider>;
  };

  return {
    useContext: () => React.useContext(Context),
    Context,
    Provider,
  };
};

const workspaceStore = createZustandContext(buildStore);
export const WorkspaceStoreContext = workspaceStore.Context;
export const WorkspaceStoreProvider = workspaceStore.Provider;

 
export function useWorkspaceStore() {
  const context = React.useContext(WorkspaceStoreContext);
  if (!context) {
    throw new Error(
      "useWorkspaceStore must be used within a WorkspaceStoreProvider"
    );
  }
  return context;
}

const accountStore = createZustandContext(buildAccountStore);
export const AccoutStoreContext = accountStore.Context;
export const AccountStoreProvider = accountStore.Provider;

 
export function useAccountStore() {
  const context = React.useContext(AccoutStoreContext);
  if (!context) {
    throw new Error(
      "useAccountStore must be used within an AccountStoreProvider"
    );
  }
  return context;
}

// theme provider from shadcn/ui
type Theme = "dark" | "light" | "system";

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
};

type ThemeProviderState = {
  theme: Theme;
  setTheme: (theme: Theme) => void;
};

const initialState: ThemeProviderState = {
  theme: "system",
  setTheme: () => null,
};

export const ThemeProviderContext =
  React.createContext<ThemeProviderState>(initialState);

export function ThemeProvider({
  children,
  defaultTheme = "system",
  storageKey = "vite-ui-theme",
  ...props
}: ThemeProviderProps) {
  const [theme, setTheme] = React.useState<Theme>(
    () => (localStorage.getItem(storageKey) as Theme) || defaultTheme
  );

  React.useEffect(() => {
    const root = window.document.documentElement;

    root.classList.remove("light", "dark");

    if (theme === "system") {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
        .matches
        ? "dark"
        : "light";

      root.classList.add(systemTheme);
      return;
    }

    root.classList.add(theme);
  }, [theme]);

  const value = {
    theme,
    setTheme: (theme: Theme) => {
      localStorage.setItem(storageKey, theme);
      setTheme(theme);
    },
  };

  return (
    <ThemeProviderContext.Provider {...props} value={value}>
      {children}
    </ThemeProviderContext.Provider>
  );
}
