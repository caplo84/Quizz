
import { ChevronLeft } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const QuizHeader = ({ topicName, currentQuestion, totalQuestions, answeredCount = 0, darkMode }) => {
  const navigate = useNavigate();
  const progress = totalQuestions > 0 ? Math.round((currentQuestion / totalQuestions) * 100) : 0;

  const handleExit = () => {
    navigate('/');
  };

  return (
    <header className={`quiz-header ${darkMode ? 'bg-dark-navy border-navy' : 'bg-white border-light-grey'}`}>
      <div className="quiz-header-inner">
        <div className="quiz-header-top">
          <div className="quiz-header-left">
            <button 
              className={`exit-button ${darkMode ? 'hover:bg-navy' : 'hover:bg-light-grey'}`}
              onClick={handleExit}
              aria-label="Exit quiz"
            >
              <ChevronLeft size={18} className={darkMode ? 'text-light-bluish' : 'text-dark-navy'} />
            </button>
            <div className="topic-info">
              <h1 className={`quiz-topic-title ${darkMode ? 'text-white' : 'text-dark-navy'}`}>
                {topicName}
              </h1>
              <span className={`quiz-topic-subtitle ${darkMode ? 'text-light-grey' : 'text-grey-navy'}`}>
                Question {currentQuestion} of {totalQuestions}
              </span>
            </div>
          </div>

          <div className="quiz-header-right">
            <span className={`quiz-stat ${darkMode ? 'text-light-grey' : 'text-grey-navy'}`}>
              {progress}% complete
            </span>
            <span className={`quiz-stat ${darkMode ? 'text-light-grey' : 'text-grey-navy'}`}>
              {answeredCount}/{totalQuestions} answered
            </span>
          </div>
        </div>

        <div className={`quiz-header-progress-track ${darkMode ? 'bg-navy' : 'bg-light-grey'}`}>
          <div
            className="quiz-header-progress-fill"
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>
    </header>
  );
};

export default QuizHeader;
