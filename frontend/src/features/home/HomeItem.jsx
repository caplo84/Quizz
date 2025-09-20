import { useDispatch, useSelector } from "react-redux";
import { selectIcon, selectQuiz } from "./homeSlice";
import { useNavigate } from "react-router-dom";

function HomeItem({ img, text, slug, source = 'manual', quizCount = 0 }) {
  const { darkMode } = useSelector((state) => state.home);
  const dispatch = useDispatch();
  const navigate = useNavigate();

  function setQuiz() {
    dispatch(selectQuiz(text));
    dispatch(selectIcon(img));
    navigate(`/${slug}`);
  }

  // Define colors for different technologies to match the image
  const getTopicColor = (topicName) => {
    const colors = {
      'HTML': '#ff6b47',
      'CSS': '#4285f4', 
      'JavaScript': '#ffd600',
      'Accessibility': '#a729f5',
      'React': '#61dafb',
      'Angular': '#dd0031',
      'Python': '#3776ab',
      'Java': '#ed8b00',
      'C++': '#9c88ff',
      'Node.js': '#8bc34a',
      'PHP': '#777bb4',
      'Go Programming': '#00add8',
      'Kotlin': '#7f52ff',
      'Swift': '#fa7343',
      'MongoDB': '#47a248',
      'MySQL': '#4479a1',
      'Android Development': '#3ddc84',
      'AWS Lambda': '#ff9900',
      'Machine Learning': '#ff6f00',
      'Git Version Control': '#f05032',
      'Django': '#092e20',
      'C Programming': '#a8b9cc',
      'C#': '#239120',
      'jQuery': '#0769ad',
      'Linux': '#fcc624',
      'Apache Maven': '#c71a36',
      'Microsoft Azure': '#0078d4',
      'MATLAB': '#e16737',
      'R Programming': '#276dc3',
      'Objective-C': '#438eff',
      'Bash Scripting': '#4eaa25',
      'JSON': '#000000',
      '.NET Framework': '#512bd4',
      'Frontend Development': '#61dafb',
      'Google Cloud Platform': '#4285f4',
    };
    
    return colors[topicName] || '#6b7280';
  };

  const topicColor = getTopicColor(text);

  return (
    <div
      className={`group relative rounded-2xl p-6 cursor-pointer transition-all duration-500 hover:scale-105 hover:-translate-y-2 shadow-xl transform ${
        darkMode 
          ? 'bg-gray-800/90 backdrop-blur-lg hover:shadow-2xl hover:shadow-purple-500/20 border border-gray-700' 
          : 'bg-white/95 backdrop-blur-lg hover:shadow-2xl hover:shadow-slate-500/10 border border-slate-200/50'
      }`}
      onClick={setQuiz}
    >
      {/* Hover glow effect */}
      <div className={`absolute inset-0 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500 ${
        darkMode 
          ? 'bg-gradient-to-br from-purple-400/20 to-indigo-600/20' 
          : 'bg-gradient-to-br from-slate-400/10 to-blue-400/10'
      }`}></div>
      
      {/* Content */}
      <div className="relative z-10 text-center">
        {/* Icon Container */}
        <div className="flex justify-center mb-4">
          <div 
            className="w-14 h-14 rounded-xl flex items-center justify-center group-hover:rotate-12 transition-transform duration-500 shadow-lg"
            style={{ backgroundColor: topicColor }}
        >
          {img ? (
            <img
              src={img}
              alt={text}
              className="w-8 h-8 object-contain filter brightness-0 invert"
            />
          ) : (
            <div className="text-white font-bold text-sm">
              {text.substring(0, 2).toUpperCase()}
            </div>
          )}
        </div>
      </div>
      
      {/* Topic Name */}
      <div className="text-center">
        <h3 className={`font-semibold text-base leading-tight group-hover:scale-105 transition-transform duration-300 ${
          darkMode ? 'text-white group-hover:text-purple-400' : 'text-slate-700 group-hover:text-blue-600'
        }`}>
          {text}
        </h3>
        
        {/* Quiz Count for External Sources */}
        {source !== 'manual' && quizCount > 0 && (
          <p className={`text-xs mt-1 transition-colors duration-300 ${
            darkMode ? 'text-gray-400 group-hover:text-gray-300' : 'text-slate-500 group-hover:text-slate-600'
          }`}>
            {quizCount} quiz{quizCount > 1 ? 'es' : ''}
          </p>
        )}
      </div>
      
      {/* Source Indicator (small badge) */}
      {source !== 'manual' && (
        <div className={`absolute top-2 right-2 w-5 h-5 rounded-full flex items-center justify-center opacity-70 ${
          darkMode ? 'bg-gray-600' : 'bg-gray-600'
        }`}>
          {source === 'github' ? (
            <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 0C4.477 0 0 4.484 0 10.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0110 4.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.203 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.942.359.31.678.921.678 1.856 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0020 10.017C20 4.484 15.522 0 10 0z" clipRule="evenodd" />
            </svg>
          ) : (
            <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z" clipRule="evenodd" />
            </svg>
          )}
        </div>
      )}

      {/* Animated dot */}
      <div className={`absolute top-2 left-2 w-2 h-2 rounded-full opacity-0 group-hover:opacity-100 animate-ping transition-opacity duration-300 ${
        darkMode ? 'bg-purple-400' : 'bg-blue-500'
      }`}></div>
      
      {/* Bottom accent line */}
      <div className={`absolute bottom-0 left-0 right-0 h-1 rounded-b-2xl transform scale-x-0 group-hover:scale-x-100 transition-transform duration-500 origin-left ${
        darkMode ? 'bg-gradient-to-r from-purple-500 to-indigo-600' : 'bg-gradient-to-r from-blue-500 to-slate-600'
      }`}></div>
      </div>
    </div>
  );
}

export default HomeItem;
