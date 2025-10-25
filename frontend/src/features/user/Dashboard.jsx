import { Link, useNavigate } from "react-router-dom";
import { Award, BarChart3, BookOpen, ChevronLeft, Clock3, Target } from "lucide-react";
import { getAttempts, getCurrentUser } from "../../utils/attemptHistory";

function formatDuration(seconds) {
  if (!seconds) return "0m";
  const hours = Math.floor(seconds / 3600);
  const mins = Math.floor((seconds % 3600) / 60);
  return hours > 0 ? `${hours}h ${mins}m` : `${mins}m`;
}

export default function Dashboard() {
  const navigate = useNavigate();
  const user = getCurrentUser();
  const attempts = getAttempts();

  const myAttempts = user
    ? attempts.filter((item) => item.userName === user.name)
    : attempts;

  const totalAttempts = myAttempts.length;
  const avgScore = totalAttempts
    ? Math.round(myAttempts.reduce((sum, item) => sum + (item.percentage || 0), 0) / totalAttempts)
    : 0;
  const bestScore = totalAttempts ? Math.max(...myAttempts.map((item) => item.percentage || 0)) : 0;
  const totalTime = myAttempts.reduce((sum, item) => sum + (item.durationSeconds || 0), 0);

  const stats = [
    { label: "Quizzes Taken", value: totalAttempts, icon: BookOpen, color: "text-indigo-600" },
    { label: "Avg Score", value: `${avgScore}%`, icon: Target, color: "text-green-500" },
    { label: "Best Score", value: `${bestScore}%`, icon: Award, color: "text-yellow-500" },
    { label: "Time Spent", value: formatDuration(totalTime), icon: Clock3, color: "text-purple-500" },
  ];

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <header className="sticky top-0 z-40 bg-white/85 dark:bg-gray-800/85 backdrop-blur-md border-b border-gray-200 dark:border-gray-700">
        <div className="max-w-4xl mx-auto px-4 py-3 flex items-center gap-3">
          <button onClick={() => navigate("/")} className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
            <ChevronLeft className="w-5 h-5" />
          </button>
          <BarChart3 className="w-5 h-5 text-indigo-600" />
          <h1 className="font-bold flex-1 text-gray-900 dark:text-white">My Progress</h1>
          <Link to="/leaderboard" className="text-xs text-gray-500 dark:text-gray-300">Leaderboard</Link>
        </div>
      </header>

      <div className="max-w-4xl mx-auto px-4 py-8 space-y-6">
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
          {stats.map(({ label, value, icon: Icon, color }) => (
            <div key={label} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl p-4 text-center">
              <Icon className={`w-6 h-6 mx-auto mb-2 ${color}`} />
              <div className="font-black text-2xl text-gray-900 dark:text-white">{value}</div>
              <div className="text-xs text-gray-500 dark:text-gray-300 mt-0.5">{label}</div>
            </div>
          ))}
        </div>

        {totalAttempts === 0 ? (
          <div className="text-center py-16 text-gray-500 dark:text-gray-300">
            <p className="font-semibold">No quiz attempts yet</p>
            <p className="text-sm mt-1">Complete a quiz to populate your dashboard.</p>
            <button
              type="button"
              onClick={() => navigate("/")}
              className="mt-4 px-5 py-2 bg-indigo-600 text-white rounded-xl text-sm font-semibold"
            >
              Start a Quiz
            </button>
          </div>
        ) : (
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl p-6">
            <h2 className="font-bold mb-4 text-gray-900 dark:text-white">Recent Attempts</h2>
            <div className="space-y-2">
              {myAttempts.slice(0, 10).map((attempt) => (
                <div key={attempt.id} className="flex items-center gap-3 p-3 rounded-xl bg-gray-100 dark:bg-gray-700/50">
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm truncate text-gray-900 dark:text-white">{attempt.topicName}</p>
                    <p className="text-xs text-gray-500 dark:text-gray-300">{new Date(attempt.createdAt).toLocaleString()}</p>
                  </div>
                  <div className="text-xs text-gray-500 dark:text-gray-300">{attempt.score}/{attempt.totalQuestions}</div>
                  <div className={`font-bold text-sm w-12 text-right ${attempt.percentage >= 80 ? "text-green-500" : attempt.percentage >= 60 ? "text-yellow-500" : "text-red-500"}`}>
                    {attempt.percentage}%
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
