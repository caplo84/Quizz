import { Code, Globe, Cpu, Terminal, BookOpen, Layers, Zap, Shield, Box } from "lucide-react";

const iconMap = {
  code: Code,
  globe: Globe,
  cpu: Cpu,
  terminal: Terminal,
  book: BookOpen,
  layers: Layers,
  zap: Zap,
  shield: Shield,
  box: Box,
};

const colorMap = {
  "#8B5CF6": "from-violet-500 to-purple-600",
  "#06B6D4": "from-cyan-500 to-teal-600",
  "#F59E0B": "from-amber-500 to-orange-600",
  "#EF4444": "from-red-500 to-pink-600",
  "#10B981": "from-emerald-500 to-green-600",
  "#3B82F6": "from-blue-500 to-indigo-600",
  "#EC4899": "from-pink-500 to-rose-600",
  "#F97316": "from-orange-500 to-amber-600",
};

function inferIcon(slug = "") {
  if (slug.includes("aws") || slug.includes("cloud")) return "globe";
  if (slug.includes("react") || slug.includes("frontend")) return "layers";
  if (slug.includes("security") || slug.includes("cyber")) return "shield";
  if (slug.includes("bash") || slug.includes("linux")) return "terminal";
  if (slug.includes("go") || slug.includes("c") || slug.includes("java")) return "cpu";
  return "code";
}

export default function TopicCardModern({ topic, onClick }) {
  const Icon = iconMap[inferIcon(topic.slug)] || Code;
  const gradient = colorMap[topic.color] || "from-violet-500 to-purple-600";

  return (
    <button
      type="button"
      onClick={() => onClick(topic)}
      className="w-full cursor-pointer text-left group transition-transform hover:-translate-y-1"
    >
      <div className="relative overflow-hidden rounded-2xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 p-6 transition-shadow hover:shadow-xl">
        <div
          className={`absolute top-0 right-0 w-32 h-32 bg-gradient-to-br ${gradient} opacity-10 rounded-full transform translate-x-8 -translate-y-8 group-hover:scale-150 transition-transform duration-500`}
        />
        <div
          className={`w-12 h-12 rounded-xl bg-gradient-to-br ${gradient} flex items-center justify-center mb-4 shadow-lg`}
        >
          <Icon className="w-6 h-6 text-white" />
        </div>
        <h3 className="font-bold text-lg mb-1 text-gray-900 dark:text-white">{topic.title}</h3>
        {topic.description && (
          <p className="text-sm text-gray-500 dark:text-gray-300 line-clamp-2 mb-3">{topic.description}</p>
        )}
        <div className="flex items-center justify-between">
          <span className="text-xs font-medium text-gray-500 dark:text-gray-300">
            {topic.quizCount || 0} quizzes
          </span>
          <span className={`text-xs font-semibold px-2.5 py-1 rounded-full bg-gradient-to-r ${gradient} text-white`}>
            Start →
          </span>
        </div>
      </div>
    </button>
  );
}
