/* eslint-disable react-refresh/only-export-components */
import { createContext, useContext, useEffect, useState } from "react";

// --- Types ---
export type Theme = "dark" | "light" | "system";
export type AccentColor =
  | "blue"
  | "violet"
  | "emerald"
  | "rose"
  | "amber"
  | "cyan";

interface ThemeProviderProps {
  children: React.ReactNode;
  defaultTheme?: Theme;
  defaultAccent?: AccentColor;
  storageKey?: string;
}

interface ThemeProviderState {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  accent: AccentColor;
  setAccent: (accent: AccentColor) => void;
  // Helper to get dynamic classes based on current accent
  accentStyles: {
    text: string;
    bg: string;
    bgHover: string;
    border: string;
    ring: string;
    lightBg: string; // Faint background for badges
  };
}

// --- Style Maps ---
// We map accents to specific Tailwind classes so the compiler picks them up.
const COLOR_MAP: Record<AccentColor, ThemeProviderState["accentStyles"]> = {
  blue: {
    text: "text-blue-600 dark:text-blue-400",
    bg: "bg-blue-600",
    bgHover: "hover:bg-blue-700",
    border: "border-blue-600",
    ring: "focus:ring-blue-500",
    lightBg: "bg-blue-50 dark:bg-blue-500/10",
  },
  violet: {
    text: "text-violet-600 dark:text-violet-400",
    bg: "bg-violet-600",
    bgHover: "hover:bg-violet-700",
    border: "border-violet-600",
    ring: "focus:ring-violet-500",
    lightBg: "bg-violet-50 dark:bg-violet-500/10",
  },
  emerald: {
    text: "text-emerald-600 dark:text-emerald-400",
    bg: "bg-emerald-600",
    bgHover: "hover:bg-emerald-700",
    border: "border-emerald-600",
    ring: "focus:ring-emerald-500",
    lightBg: "bg-emerald-50 dark:bg-emerald-500/10",
  },
  rose: {
    text: "text-rose-600 dark:text-rose-400",
    bg: "bg-rose-600",
    bgHover: "hover:bg-rose-700",
    border: "border-rose-600",
    ring: "focus:ring-rose-500",
    lightBg: "bg-rose-50 dark:bg-rose-500/10",
  },
  amber: {
    text: "text-amber-600 dark:text-amber-400",
    bg: "bg-amber-600",
    bgHover: "hover:bg-amber-700",
    border: "border-amber-600",
    ring: "focus:ring-amber-500",
    lightBg: "bg-amber-50 dark:bg-amber-500/10",
  },
  cyan: {
    text: "text-cyan-600 dark:text-cyan-400",
    bg: "bg-cyan-600",
    bgHover: "hover:bg-cyan-700",
    border: "border-cyan-600",
    ring: "focus:ring-cyan-500",
    lightBg: "bg-cyan-50 dark:bg-cyan-500/10",
  },
};

const initialState: ThemeProviderState = {
  theme: "system",
  setTheme: () => null,
  accent: "blue",
  setAccent: () => null,
  accentStyles: COLOR_MAP.blue,
};

const ThemeContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider({
  children,
  defaultTheme = "system",
  defaultAccent = "blue",
  storageKey = "gofiles-prefs",
}: ThemeProviderProps) {
  // 1. Initialize State
  const [theme, setThemeState] = useState<Theme>(
    () => (localStorage.getItem(`${storageKey}-theme`) as Theme) || defaultTheme
  );

  const [accent, setAccentState] = useState<AccentColor>(
    () =>
      (localStorage.getItem(`${storageKey}-accent`) as AccentColor) ||
      defaultAccent
  );

  // 2. Effect: Handle Dark/Light Mode
  useEffect(() => {
    const root = window.document.documentElement;
    root.classList.remove("light", "dark");

    if (theme === "system") {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
        .matches
        ? "dark"
        : "light";
      root.classList.add(systemTheme);
    } else {
      root.classList.add(theme);
    }
  }, [theme]);

  // 3. Setters with persistence
  const setTheme = (t: Theme) => {
    localStorage.setItem(`${storageKey}-theme`, t);
    setThemeState(t);
  };

  const setAccent = (a: AccentColor) => {
    localStorage.setItem(`${storageKey}-accent`, a);
    setAccentState(a);
  };

  const value = {
    theme,
    setTheme,
    accent,
    setAccent,
    accentStyles: COLOR_MAP[accent],
  };

  return (
    <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
  );
}

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (context === undefined)
    throw new Error("useTheme must be used within a ThemeProvider");
  return context;
};