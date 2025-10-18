
import { Check, Flag } from 'lucide-react';

const QuizSidebar = ({ 
  totalQuestions, 
  currentQuestion, 
  answeredQuestions = [], 
  flaggedQuestions = [],
  onQuestionSelect,
  darkMode 
}) => {
  return (
    <aside className={`quiz-sidebar-compact ${darkMode ? 'bg-dark-navy' : 'bg-white'}`}>
      <h3 className={`sidebar-title ${darkMode ? 'text-white' : 'text-dark-navy'}`}>
        Questions
      </h3>
      <div className="question-navigator-grid">
        {Array.from({ length: totalQuestions }, (_, i) => {
          const isAnswered = answeredQuestions.includes(i);
          const isFlagged = flaggedQuestions.includes(i);
          const isCurrent = i === currentQuestion;
          
          return (
            <button
              key={i}
              className={`question-nav-btn ${isCurrent ? 'current' : ''} ${isAnswered ? 'answered' : ''} ${darkMode ? 'dark-mode' : ''}`}
              onClick={() => onQuestionSelect(i)}
              aria-label={`Go to question ${i + 1}`}
            >
              <span className="question-num">{i + 1}</span>
              {isAnswered && <Check size={12} className="check-badge" />}
              {isFlagged && <Flag size={10} className="flag-badge" fill="currentColor" />}
            </button>
          );
        })}
      </div>
    </aside>
  );
};

export default QuizSidebar;
