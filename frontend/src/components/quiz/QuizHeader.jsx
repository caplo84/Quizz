
import { X, Clock } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const QuizHeader = ({ topicName, currentQuestion, totalQuestions, timeRemaining, darkMode }) => {
  const navigate = useNavigate();

  const handleExit = () => {
    if (window.confirm('Are you sure you want to exit? Your progress will be lost.')) {
      navigate('/');
    }
  };

  return (
    <header className={`quiz-header ${darkMode ? 'bg-dark-navy border-navy' : 'bg-white border-light-grey'}`}>
      <div className="quiz-header-left">
        <button 
          className={`exit-button ${darkMode ? 'hover:bg-navy' : 'hover:bg-light-grey'}`}
          onClick={handleExit}
          aria-label="Exit quiz"
        >
          <X size={24} className={darkMode ? 'text-light-bluish' : 'text-dark-navy'} />
        </button>
        <div className="topic-info">
          <h1 className={`text-2xl font-bold ${darkMode ? 'text-white' : 'text-dark-navy'}`}>
            {topicName}
          </h1>
          <span className={`text-sm ${darkMode ? 'text-light-grey' : 'text-grey-navy'}`}>
            Question {currentQuestion} of {totalQuestions}
          </span>
        </div>
      </div>
      
      <div className="quiz-header-right">
        {timeRemaining !== null && timeRemaining !== undefined && (
          <div className="timer">
            <Clock size={20} />
            <span className="font-semibold">
              {Math.floor(timeRemaining / 60)}:{(timeRemaining % 60).toString().padStart(2, '0')}
            </span>
          </div>
        )}
      </div>
    </header>
  );
};

export default QuizHeader;
