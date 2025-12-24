import { useState } from "react";
import { api } from "./api";
import { ServerCog } from "lucide-react";

interface SetupProps {
  onSetupComplete: () => void;
}

export default function Setup({ onSetupComplete }: SetupProps) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      // Call the First-Run Setup endpoint
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
    <div className="min-h-screen flex items-center justify-center bg-gray-950 px-4">
      <div className="max-w-md w-full bg-gray-900 rounded-xl shadow-2xl p-8 border border-gray-800">
        <div className="flex justify-center mb-6">
          <div className="p-3 bg-green-600 rounded-full bg-opacity-10">
            <ServerCog className="w-8 h-8 text-green-500" />
          </div>
        </div>

        <h2 className="text-2xl font-bold text-center text-white mb-2">
          System Setup
        </h2>
        <p className="text-gray-400 text-center mb-8">
          Create your admin account to get started.
        </p>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-2">
              Choose Username
            </label>
            <input
              type="text"
              required
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-3 bg-gray-950 border border-gray-800 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent outline-none text-white transition-all"
              placeholder="e.g. admin"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-2">
              Choose Password
            </label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-3 bg-gray-950 border border-gray-800 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent outline-none text-white transition-all"
              placeholder="Make it strong!"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-green-600 hover:bg-green-700 text-white font-semibold rounded-lg transition-colors"
          >
            {loading ? "Initializing System..." : "Complete Setup"}
          </button>
        </form>
      </div>
    </div>
  );
}
