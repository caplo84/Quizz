import { useEffect, useMemo, useState } from "react";
import { Link, useLoaderData, useNavigate } from "react-router-dom";
import { useDispatch, useSelector } from "react-redux";
import { BarChart2, LogIn, LogOut, Search, Shuffle, Trophy, User, Zap } from "lucide-react";
import { selectIcon, selectQuiz, setDarkMode, setQuizzes } from "./homeSlice";
import { startRandomQuiz } from "../quiz/quizSlice";
import TopicCardModern from "./TopicCardModern";
import { clearCurrentUser, getCurrentUser } from "../../utils/attemptHistory";

export default function HomeModern() {
  const data = useLoaderData();
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { darkMode } = useSelector((state) => state.home);

  const [search, setSearch] = useState("");
  const [currentUser, setCurrentUser] = useState(null);

  useEffect(() => {
    dispatch(setQuizzes(Array.isArray(data) ? data : []));
    setCurrentUser(getCurrentUser());
  }, [data, dispatch]);

  useEffect(() => {
    document.documentElement.classList.toggle("dark", darkMode);
  }, [darkMode]);

  const topics = useMemo(() => {
    const source = Array.isArray(data) ? data : [];
    return source.filter(
      (topic) =>
        Number(topic.quizCount || 0) > 0 &&
        topic.title?.toLowerCase().includes(search.toLowerCase()),
    );
  }, [data, search]);

  const handleSelect = (topic) => {
    dispatch(selectQuiz(topic.title));
    dispatch(selectIcon(topic.icon));
    navigate(`/${topic.slug}`);
  };

  const handleRandom = () => {
    if (!topics.length) return;
    const random = topics[Math.floor(Math.random() * topics.length)];

    dispatch(selectQuiz(random.title));
    dispatch(selectIcon(random.icon));
    dispatch(
      startRandomQuiz({
        topic: {
          id: random.id,
          name: random.title,
          slug: random.slug,
          icon: random.icon,
        },
      }),
    );
    navigate(`/${random.slug}?random=true`);
  };

  const handleLogout = () => {
    clearCurrentUser();
    setCurrentUser(null);
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors text-base">
      <div className="relative overflow-hidden bg-gradient-to-br from-indigo-600 via-violet-600 to-purple-700 text-white">
        <div
          className="absolute inset-0 opacity-20"
          style={{
            backgroundImage:
              "radial-gradient(circle at 20% 50%, white 1px, transparent 1px), radial-gradient(circle at 80% 80%, white 1px, transparent 1px)",
            backgroundSize: "40px 40px",
          }}
        />
        <div className="relative max-w-7xl mx-auto px-6 py-16 sm:py-24">
          <div className="flex justify-between items-start mb-8 gap-4">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <div className="bg-white/20 backdrop-blur-sm p-2 rounded-xl">
                  <Zap className="w-6 h-6" />
                </div>
                <span className="font-bold text-xl">Quizz</span>
              </div>
              <h1 className="text-5xl sm:text-6xl font-black mb-4 leading-tight tracking-tight">
                New Frontend Integrated
                <br />
                <span className="text-white/80">Modern UI on Current Stack</span>
              </h1>
              <p className="text-white/80 text-xl">
                {topics.length} topics ready · Existing API contracts preserved
              </p>
            </div>

            <div className="flex items-center gap-2 flex-wrap justify-end">
              <Link
                to="/leaderboard"
                className="flex items-center gap-1.5 bg-white/20 backdrop-blur-sm rounded-xl px-3 py-2 text-sm font-medium hover:bg-white/30 transition-colors"
              >
                <Trophy className="w-3.5 h-3.5" />
                <span className="hidden sm:inline">Leaderboard</span>
              </Link>
              <Link
                to="/dashboard"
                className="flex items-center gap-1.5 bg-white/20 backdrop-blur-sm rounded-xl px-3 py-2 text-sm font-medium hover:bg-white/30 transition-colors"
              >
                <BarChart2 className="w-3.5 h-3.5" />
                <span className="hidden sm:inline">My Progress</span>
              </Link>

              <button
                type="button"
                onClick={() => dispatch(setDarkMode())}
                className="bg-white/20 backdrop-blur-sm rounded-xl px-3 py-2 text-sm font-medium hover:bg-white/30 transition-colors"
              >
                {darkMode ? "☀️" : "🌙"}
              </button>

              {currentUser ? (
                <>
                  <div className="flex items-center gap-1.5 bg-white/20 backdrop-blur-sm rounded-xl px-3 py-2 text-sm font-medium">
                    <User className="w-3.5 h-3.5" />
                    <span className="max-w-[140px] truncate">{currentUser.name}</span>
                  </div>
                  <button
                    type="button"
                    onClick={handleLogout}
                    className="bg-white/20 backdrop-blur-sm rounded-xl p-2 hover:bg-white/30 transition-colors"
                    title="Logout"
                  >
                    <LogOut className="w-4 h-4" />
                  </button>
                </>
              ) : (
                <Link
                  to="/login"
                  className="flex items-center gap-1.5 bg-white text-indigo-700 rounded-xl px-3 py-2 text-sm font-semibold hover:bg-white/90 transition-colors"
                >
                  <LogIn className="w-3.5 h-3.5" /> Sign In
                </Link>
              )}
            </div>
          </div>

          <div className="flex gap-3 max-w-2xl mt-2">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-white/60" />
              <input
                type="text"
                placeholder="Search topics..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full pl-10 pr-4 py-3.5 rounded-xl bg-white/20 backdrop-blur-sm border border-white/30 text-white placeholder:text-white/60 focus:outline-none focus:ring-2 focus:ring-white/50"
              />
            </div>
            <button
              type="button"
              onClick={handleRandom}
              className="flex items-center gap-2 px-5 py-3.5 rounded-xl bg-white text-indigo-700 font-semibold hover:bg-white/90 transition-all shadow-lg"
            >
              <Shuffle className="w-4 h-4" /> Random
            </button>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-6 py-12">
        {topics.length === 0 ? (
          <div className="text-center py-20 text-gray-500 dark:text-gray-300">
            <p className="text-lg font-semibold">No topics found</p>
            <p className="text-sm mt-2">Try a different search term.</p>
          </div>
        ) : (
          <>
            <div className="flex items-center justify-between mb-6">
              <h2 className="font-bold text-xl text-gray-900 dark:text-white">
                {search ? `Results for "${search}"` : "All Topics"}
                <span className="ml-2 text-sm font-normal text-gray-500 dark:text-gray-300">
                  ({topics.length})
                </span>
              </h2>
              <Link
                to="/legacy"
                className="text-xs text-gray-600 dark:text-gray-300 underline-offset-2 hover:underline"
              >
                Open legacy home
              </Link>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-7">
              {topics.map((topic) => (
                <TopicCardModern key={topic.slug} topic={topic} onClick={handleSelect} />
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
