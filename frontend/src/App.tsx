/* eslint-disable react-hooks/set-state-in-effect */
import { useEffect, useState } from "react";
import { api } from "./Utils/api";
import Login from "./Auth/Login";
import Setup from "./Auth/Setup";
import Dashboard from "./Dashboard/Dashboard";
import { Loader2 } from "lucide-react";
import { useTheme } from "./Context/ThemeContext";

type AppState = "loading" | "setup" | "login" | "dashboard";

function App() {
  const [view, setView] = useState<AppState>("loading");
  const { accentStyles } = useTheme();

  const checkSystemStatus = async () => {
    try {
      const statusRes = await api.get("/system/status");
      if (statusRes.data.status === "setup_required") {
        setView("setup");
        return;
      }
      try {
        await api.get("/me");
        setView("dashboard");
      } catch (e) {
        console.error("Not logged in?", e);
        setView("login");
      }
    } catch (error) {
      console.error("Backend offline?", error);
      setView("login");
    }
  };

  useEffect(() => {
    checkSystemStatus();
  }, []);

  if (view === "loading") {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-950 flex flex-col items-center justify-center transition-colors duration-300">
        <Loader2 className={`w-12 h-12 animate-spin ${accentStyles.text}`} />
        <p className="mt-4 text-gray-500 text-sm font-medium animate-pulse">
          Initializing System...
        </p>
      </div>
    );
  }

  return (
    <div className="fade-in">
      {view === "setup" && (
        <Setup onSetupComplete={() => setView("dashboard")} />
      )}
      {view === "login" && <Login onLogin={() => setView("dashboard")} />}
      {view === "dashboard" && <Dashboard />}
    </div>
  );
}

export default App;
