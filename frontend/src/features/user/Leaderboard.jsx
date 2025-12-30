import { useMemo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Link, useNavigate } from "react-router-dom";
import { ChevronLeft, Clock3, Filter, Moon, Sun, Target, Trophy } from "lucide-react";
import { getLeaderboard } from "../../utils/attemptHistory";
import { setDarkMode } from "../home/homeSlice";

function formatTime(seconds) {
  if (!seconds) return "—";
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}m ${secs}s`;
}

export default function Leaderboard() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { darkMode } = useSelector((state) => state.home);
  const [selectedTopic, setSelectedTopic] = useState("all");

  const allAttempts = getLeaderboard("all");
  const topics = useMemo(() => {
    const unique = new Set(allAttempts.map((a) => a.topicSlug).filter(Boolean));
    return Array.from(unique);
  }, [allAttempts]);

  const ranked = useMemo(() => getLeaderboard(selectedTopic), [selectedTopic]);

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <header className="sticky top-0 z-40 bg-white/85 dark:bg-gray-800/85 backdrop-blur-md border-b border-gray-200 dark:border-gray-700">
        <div className="max-w-3xl mx-auto px-4 py-3 flex items-center gap-3">
          <button onClick={() => navigate("/")} className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
            <ChevronLeft className="w-5 h-5" />
          </button>
          <Trophy className="w-5 h-5 text-yellow-500" />
          <h1 className="font-bold flex-1 text-gray-900 dark:text-white">Leaderboard</h1>
          <button
            type="button"
            onClick={() => dispatch(setDarkMode())}
            className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            aria-label="Toggle dark mode"
            title="Toggle dark mode"
          >
            {darkMode ? <Sun className="w-4 h-4 text-yellow-300" /> : <Moon className="w-4 h-4 text-gray-700" />}
          </button>
          <Link to="/" className="text-xs text-gray-500 dark:text-gray-300">Home</Link>
        </div>
      </header>

      <div className="max-w-3xl mx-auto px-4 py-8">
        <div className="flex items-center gap-2 mb-6 overflow-x-auto pb-1">
          <Filter className="w-4 h-4 text-gray-500 dark:text-gray-300 shrink-0" />
          <button
            onClick={() => setSelectedTopic("all")}
            className={`px-3 py-1.5 rounded-full text-sm font-medium whitespace-nowrap transition-colors ${
              selectedTopic === "all"
                ? "bg-indigo-600 text-white"
                : "bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-200"
            }`}
          >
            All Topics
          </button>
          {topics.map((topic) => (
            <button
              key={topic}
              onClick={() => setSelectedTopic(topic)}
              className={`px-3 py-1.5 rounded-full text-sm font-medium whitespace-nowrap transition-colors ${
                selectedTopic === topic
                  ? "bg-indigo-600 text-white"
                  : "bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-200"
              }`}
            >
              {topic}
            </button>
          ))}
        </div>

        {ranked.length === 0 ? (
          <div className="text-center py-20 text-gray-500 dark:text-gray-300">
            <Trophy className="w-12 h-12 mx-auto mb-3 opacity-25" />
            <p>No attempts yet. Complete a quiz to appear here.</p>
          </div>
        ) : (
          <div className="space-y-2">
            {ranked.map((attempt, idx) => (
              <div key={attempt.id} className="flex items-center gap-4 p-4 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
                <div className="w-8 text-center font-bold text-lg text-gray-600 dark:text-gray-200">#{idx + 1}</div>
                <div className="flex-1 min-w-0">
                  <p className="font-semibold text-sm truncate text-gray-900 dark:text-white">{attempt.userName || "Anonymous"}</p>
                  <p className="text-xs text-gray-500 dark:text-gray-300">{attempt.topicName || "Unknown topic"}</p>
                </div>
                <div className="flex items-center gap-4 text-sm shrink-0">
                  <div className="flex items-center gap-1 text-gray-500 dark:text-gray-300">
                    <Clock3 className="w-3.5 h-3.5" />
                    {formatTime(attempt.durationSeconds)}
                  </div>
                  <div className="flex items-center gap-1 text-gray-500 dark:text-gray-300">
                    <Target className="w-3.5 h-3.5" />
                    {attempt.score}/{attempt.totalQuestions}
                  </div>
                  <div
                    className={`font-bold text-base w-14 text-right ${
                      attempt.percentage >= 80
                        ? "text-green-500"
                        : attempt.percentage >= 60
                          ? "text-yellow-500"
                          : "text-red-500"
                    }`}
                  >
                    {attempt.percentage}%
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
