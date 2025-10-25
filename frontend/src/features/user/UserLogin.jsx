import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { LogIn, User, Zap } from "lucide-react";
import { setCurrentUser } from "../../utils/attemptHistory";

export default function UserLogin() {
  const [name, setName] = useState("");
  const navigate = useNavigate();

  const handleSubmit = (event) => {
    event.preventDefault();
    if (!name.trim()) return;

    setCurrentUser({
      name: name.trim(),
      email: `${name.trim().toLowerCase().replace(/\s+/g, ".")}@local.quiz`,
    });

    navigate("/", { replace: true });
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center p-4">
      <div className="w-full max-w-sm text-center">
        <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-indigo-600 to-purple-600 flex items-center justify-center mx-auto mb-6 shadow-xl">
          <Zap className="w-8 h-8 text-white" />
        </div>
        <h1 className="text-3xl font-black mb-2 text-gray-900 dark:text-white">Welcome to Quizz</h1>
        <p className="text-gray-500 dark:text-gray-300 mb-8">
          Sign in locally to track your results on this device.
        </p>

        <form onSubmit={handleSubmit} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl p-6 space-y-4">
          <label className="text-left block text-sm text-gray-600 dark:text-gray-300">Display name</label>
          <div className="relative">
            <User className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Your name"
              className="w-full pl-10 pr-3 py-3 border border-gray-200 dark:border-gray-600 rounded-xl bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
              required
            />
          </div>
          <button
            type="submit"
            className="w-full py-3 bg-indigo-600 text-white rounded-xl font-semibold hover:bg-indigo-500 transition-colors flex items-center justify-center gap-2"
          >
            <LogIn className="w-4 h-4" /> Continue
          </button>
        </form>
      </div>
    </div>
  );
}
