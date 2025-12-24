import { useState } from "react";
import { api } from "./api";
import { ServerCog } from "lucide-react";
import { useTheme } from "./ThemeContext";

interface SetupProps {
  onSetupComplete: () => void;
}

export default function Setup({ onSetupComplete }: SetupProps) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const { accentStyles } = useTheme();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await api.post("/setup", { username, password });
      onSetupComplete();
    } catch (err) {
      console.error("Setup failed:", err);
      alert("Setup failed. Check console.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-950 px-4 transition-colors duration-300">
      <div className="max-w-md w-full bg-white dark:bg-gray-900 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-800 p-8 md:p-10">
        <div className="flex justify-center mb-6">
          <div
            className={`p-4 rounded-full bg-opacity-10 ${accentStyles.lightBg}`}
          >
            <ServerCog className={`w-8 h-8 ${accentStyles.text}`} />
          </div>
        </div>

        <h2 className="text-2xl font-bold text-center text-gray-900 dark:text-white mb-2">
          System Setup
        </h2>
        <p className="text-gray-500 text-center mb-8">
          Configure your admin credentials to initialize the platform.
        </p>

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
              Admin Username
            </label>
            <input
              type="text"
              required
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className={`w-full px-4 py-3 bg-gray-50 dark:bg-gray-950 border border-gray-200 dark:border-gray-800 rounded-xl focus:outline-none focus:ring-2 focus:border-transparent transition-all text-gray-900 dark:text-white ${accentStyles.ring}`}
              placeholder="e.g. admin"
            />
          </div>
          <div>
            <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
              Master Password
            </label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={`w-full px-4 py-3 bg-gray-50 dark:bg-gray-950 border border-gray-200 dark:border-gray-800 rounded-xl focus:outline-none focus:ring-2 focus:border-transparent transition-all text-gray-900 dark:text-white ${accentStyles.ring}`}
              placeholder="Create a strong password"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className={`w-full py-3.5 font-bold text-white rounded-xl shadow-lg transition-all ${accentStyles.bg} ${accentStyles.bgHover}`}
          >
            {loading ? "Initializing..." : "Complete Setup"}
          </button>
        </form>
      </div>
    </div>
  );
}