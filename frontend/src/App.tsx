/* eslint-disable react-hooks/set-state-in-effect */
import { useEffect, useState } from "react";
import { api } from "./api";
import Login from "./Login";
import Setup from "./Setup";
import { Loader2 } from "lucide-react";
import Dashboard from "./Dashboard";

type AppState = "loading" | "setup" | "login" | "dashboard";

function App() {
  const [view, setView] = useState<AppState>("loading");

  const checkSystemStatus = async () => {
    try {
      // 1. Check if system is configured (First Run?)
      const statusRes = await api.get("/system/status");

      if (statusRes.data.status === "setup_required") {
        setView("setup");
        return;
      }

      // 2. If configured, check if we are logged in
      try {
        await api.get("/me");
        setView("dashboard"); // Cookie is valid!
      } catch (e) {
        console.error("Not logged in?", e);
        setView("login"); // Cookie invalid/missing
      }
    } catch (error) {
      console.error("Backend offline?", error);
      setView("login"); // Fallback
    }
  };

  useEffect(() => {
    checkSystemStatus();
  }, []);

  if (view === "loading") {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center text-white">
        <Loader2 className="w-10 h-10 animate-spin text-blue-500" />
      </div>
    );
  }

  if (view === "setup") {
    return <Setup onSetupComplete={() => setView("dashboard")} />;
  }

  if (view === "login") {
    return <Login onLogin={() => setView("dashboard")} />;
  }

  if (view === "dashboard") {
    // Pass a way to logout if needed, or just let the Dashboard handle its own refresh
    return <Dashboard />;
  }
}

export default App;
