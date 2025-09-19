import { useDispatch, useSelector } from "react-redux";
import { selectIcon, selectQuiz } from "./homeSlice";
import { useNavigate } from "react-router-dom";

function HomeItem({ img, text, slug, source = 'manual', quizCount = 0 }) {
  const { darkMode } = useSelector((state) => state.home);
  const dispatch = useDispatch();
  const navigate = useNavigate();

  function setQuiz() {
    // Only navigate if there are quizzes available for now, just allow all navigation
    dispatch(selectQuiz(text));
    dispatch(selectIcon(img));
    navigate(`/${slug}`); // Use slug instead of text for navigation
  }

  const bgColors = {
    HTML: "#FFF1E9",
    CSS: "#E0FDEF",
    JavaScript: "#EBF0FF",
    Accessibility: "#F6E7FF",
  };

  // Use specific colors for original topics, generic for external sources
  const isExternalSource = source !== 'manual';
  const backgroundColor = isExternalSource 
    ? (darkMode ? "#374151" : "#F3F4F6") 
    : bgColors[text];

  const bgStyle = {
    backgroundColor: backgroundColor,
    borderRadius: "0.4rem",
  };

  return (
    <li
      className={`mobile:gap-6 mobile:rounded-[1.2rem] flex cursor-pointer items-center gap-12 rounded-[2.4rem] p-8 shadow-sm transition-all duration-300 hover:shadow-md hover:scale-105 ${
        darkMode ? "bg-navy" : "bg-white"
      }`}
      onClick={setQuiz}
    >
      <div className="relative">
        <img
          src={img}
          alt={text}
          className="h-[4rem] w-[4rem] p-2"
          style={bgStyle}
        />
        {isExternalSource && (
          <div 
            className="absolute -top-1 -right-1 w-5 h-5 bg-gray-800 rounded-full flex items-center justify-center"
            title={`Source: ${source}`}
          >
            {source === 'github' ? (
              <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 0C4.477 0 0 4.484 0 10.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0110 4.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.203 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.942.359.31.678.921.678 1.856 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0020 10.017C20 4.484 15.522 0 10 0z" clipRule="evenodd" />
              </svg>
            ) : source === 'quizapi' ? (
              <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z" clipRule="evenodd" />
              </svg>
            ) : source === 'ai-generated' ? (
              <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M11.3 1.046A1 1 0 0112 2v5h4a1 1 0 01.82 1.573l-7 10A1 1 0 018 18v-5H4a1 1 0 01-.82-1.573l7-10a1 1 0 011.12-.38z" clipRule="evenodd" />
              </svg>
            ) : (
              <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
            )}
          </div>
        )}
      </div>
      <div className="flex-1">
        <p
          className={`mobile:text-[1.8rem] text-[2.8rem] font-medium leading-[100%] transition-all duration-300 ${
            darkMode ? "text-white" : "text-dark-navy"
          }`}
        >
          {text}
        </p>
        {isExternalSource && (
          <p className={`text-sm mt-1 ${
            darkMode ? "text-gray-400" : "text-gray-600"
          }`}>
            Source: {source} • {quizCount > 0 ? `${quizCount} quizzes` : 'No quizzes yet'}
          </p>
        )}
      </div>
    </li>
  );
}

export default HomeItem;
