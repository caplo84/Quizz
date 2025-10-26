
import { Check } from 'lucide-react';

const QuizSidebar = ({ 
  totalQuestions, 
  currentQuestion, 
  answeredQuestions = [], 
  onQuestionSelect,
  darkMode,
  canSubmit = false,
  onSubmit,
}) => {
  return (
    <aside className={`quiz-sidebar-compact ${darkMode ? 'bg-dark-navy' : 'bg-white'}`}>
      <h3 className={`sidebar-title ${darkMode ? 'text-white' : 'text-dark-navy'}`}>
        Questions
      </h3>
      <div className="question-navigator-grid">
        {Array.from({ length: totalQuestions }, (_, i) => {
          const isAnswered = answeredQuestions.includes(i);
          const isCurrent = i === currentQuestion;
          
          return (
            <button
              key={i}
              className={`question-nav-btn ${isCurrent ? 'current' : ''} ${isAnswered ? 'answered' : ''} ${darkMode ? 'dark-mode' : ''}`}
              onClick={() => onQuestionSelect(i)}
              aria-label={`Go to question ${i + 1}`}
            >
              <span className="question-num">Q{i + 1}</span>
              {isAnswered && <Check size={12} className="check-badge" />}
            </button>
          );
        })}
      </div>

      {onSubmit && (
        <div className="mt-3 pt-3 border-t border-slate-200 dark:border-slate-700">
          <button
            type="button"
            onClick={onSubmit}
            disabled={!canSubmit}
            className="w-full px-3 py-2 rounded-xl bg-violet-600 text-white text-sm font-semibold hover:bg-violet-700 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Submit Quiz
          </button>
        </div>
      )}
    </aside>
  );
};

export default QuizSidebar;
