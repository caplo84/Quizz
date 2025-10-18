
import { ChevronLeft, ChevronRight, Send, Flag, CheckCircle, Clock } from 'lucide-react';

const QuizFooter = ({ 
  hasPrevious, 
  hasNext, 
  selectedAnswer, 
  onPrevious, 
  onNext, 
  onSubmit,
  isLastQuestion,
  showExplanation,
  isFlagged,
  onFlag,
  answered,
  remaining,
  totalQuestions,
  darkMode 
}) => {
  const allAnswered = answered === totalQuestions;
  
  return (
    <footer className={`quiz-footer-new ${darkMode ? 'bg-dark-navy border-navy' : 'bg-white border-light-grey'}`}>
      {/* Left: Flag + Previous Buttons */}
      <div className="footer-left">
        <button 
          className={`btn-flag-compact ${isFlagged ? 'flagged' : ''} ${darkMode ? 'dark-mode' : ''}`}
          onClick={onFlag}
          title={isFlagged ? "Remove flag" : "Flag for review"}
        >
          <Flag size={22} fill={isFlagged ? "currentColor" : "none"} />
          <span>Flag</span>
        </button>
        
        {hasPrevious && (
          <button 
            className={`btn-secondary-nav ${darkMode ? 'dark-mode' : ''}`}
            onClick={onPrevious}
          >
            <ChevronLeft size={22} />
            <span>Previous</span>
          </button>
        )}
      </div>
      
      {/* Center: Stats */}
      <div className="footer-stats">
        <div className="stat-item">
          <CheckCircle size={22} className="stat-icon answered" />
          <span className={`stat-text ${darkMode ? 'text-light-bluish' : 'text-grey-navy'}`}>
            <strong>{answered}</strong> Answered
          </span>
        </div>
        <div className="stat-item">
          <Clock size={22} className="stat-icon remaining" />
          <span className={`stat-text ${darkMode ? 'text-light-bluish' : 'text-grey-navy'}`}>
            <strong>{remaining}</strong> Remaining
          </span>
        </div>
      </div>
      
      {/* Right: Next + Submit Buttons */}
      <div className="footer-right">
        {hasNext && (
          <button 
            className="btn-primary-gradient"
            onClick={onNext}
          >
            <span>Next</span>
            <ChevronRight size={22} />
          </button>
        )}
        
        <button 
          className={`btn-submit ${allAnswered ? 'enabled' : 'disabled'}`}
          onClick={onSubmit}
          disabled={!allAnswered}
          title={allAnswered ? "Submit your quiz answers" : `${remaining} question${remaining !== 1 ? 's' : ''} remaining`}
        >
          <Send size={22} />
          <span>Submit</span>
        </button>
      </div>
    </footer>
  );
};

export default QuizFooter;
