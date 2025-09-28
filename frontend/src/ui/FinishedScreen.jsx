import { useDispatch, useSelector } from "react-redux";
import { useNavigate, useSearchParams } from "react-router-dom";
import { resetQuiz, fetchRandomQuestions, startNewBatch } from "../features/quiz/quizSlice";
import { resetMode } from "../features/home/homeSlice";
import { useEffect, useState } from "react";

function FinishedScreen() {
  const { name, icon } = useSelector((state) => state.home);
  const { score, questions, isRandomQuiz, randomTopic, usedQuestionIds, userAnswers } = useSelector((state) => state.quiz);
  const { darkMode } = useSelector((state) => state.home);
  const [searchParams] = useSearchParams();
  const isRandomFinish = searchParams.get('random') === 'true';
  const [showReview, setShowReview] = useState(false);

  const dispatch = useDispatch();
  const navigate = useNavigate();

  // Debug the random quiz state
  console.log('🎯 FinishedScreen Debug:');
  console.log('  isRandomFinish:', isRandomFinish);
  console.log('  isRandomQuiz:', isRandomQuiz);
  console.log('  randomTopic:', randomTopic);
  console.log('  Combined condition:', isRandomFinish && isRandomQuiz);

  const bgColors = {
    HTML: "#FFF1E9",
    CSS: "#E0FDEF",
    JavaScript: "#EBF0FF",
    Accessibility: "#F6E7FF",
  };

  const bgStyle = {
    backgroundColor: bgColors[name],
  };

  function playAgain() {
    if (isRandomFinish && randomTopic) {
      // For random quiz, start new batch with same topic
      dispatch(startNewBatch());
      navigate(`/${randomTopic.slug}?random=true`);
    } else if (randomTopic) {
      // If we have topic info from random quiz but no URL param, use it
      dispatch(startNewBatch());
      navigate(`/${randomTopic.slug}?random=true`);
    } else {
      // For regular quiz, reset and go to same topic selection
      dispatch(resetQuiz());
      dispatch(resetMode());
      navigate("/");
    }
  }

  function backToHome() {
    dispatch(resetQuiz());
    dispatch(resetMode());
    navigate("/");
  }

  useEffect(() => {
    if (!name || !icon) {
      dispatch(resetQuiz());
      dispatch(resetMode());
      navigate("/");
    }
  }, [icon, name, dispatch, navigate]);

  return (
    <div className="desktop:grid-cols-1 desktop:gap-24 mobile:gap-16 grid grid-cols-2">
      {!showReview ? (
        <>
          {/* Score Screen */}
          <div>
            <h2
              className={`mobile:text-[4rem] text-[6.4rem] font-light leading-[100%] transition-all duration-300 ${
                darkMode ? "text-white" : "text-dark-navy"
              }`}
            >
              Quiz completed <br />
              <span className="font-medium">You scored...</span>
            </h2>
          </div>
          <div>
            <div
              className={`mobile:rounded-[1.2rem] flex flex-col items-center gap-16 rounded-[2.4rem] p-20 transition-all duration-300 ${
                darkMode ? "bg-navy" : "bg-white"
              }`}
            >
              <div className="mobile:gap-8 flex items-center gap-12">
                <img
                  src={icon}
                  alt={name}
                  className="h-[4rem] w-[4rem] p-2"
                  style={bgStyle}
                />
                <p
                  className={`mobile:text-[1.8rem] text-[2.8rem] font-medium leading-[100%] transition-all duration-300 ${
                    darkMode ? "text-white" : "text-dark-navy"
                  }`}
                >
                  {name}
                </p>
              </div>
              <div className="flex flex-col items-center gap-6">
                <p
                  className={`mobile:text-[8.8rem] text-[14.4rem] font-medium leading-[100%] transition-all duration-300 ${
                    darkMode ? "text-white" : "text-dark-navy"
                  }`}
                >
                  {score}
                </p>
                <p
                  className={`mobile:text-[1.8rem] text-[2.4rem] leading-[150%] transition-all duration-300 ${
                    darkMode ? "text-light-bluish" : "text-dark-navy"
                  }`}
                >
                  out of {questions.length}
                </p>
              </div>
            </div>
            
            {/* Review button */}
            <button
              className={`mobile:text-[1.6rem] mobile:p-5 mobile:rounded-[1rem] w-full mt-4 rounded-[2rem] p-[2rem] text-[2rem] font-medium leading-[100%] transition-all duration-300 border-2 ${
                darkMode 
                  ? 'border-gray-400 text-gray-400 hover:bg-gray-400 hover:text-dark-navy' 
                  : 'border-gray-600 text-gray-600 hover:bg-gray-600 hover:text-white'
              }`}
              onClick={() => setShowReview(true)}
            >
              📝 Review Answers
            </button>
            
            {/* Quiz completion buttons - same for all quizzes */}
            <div className="space-y-4">
              <button
                className="mobile:text-[1.8rem] mobile:p-7 mobile:rounded-[1.2rem] col-start-2 mt-12 w-full rounded-[2.4rem] bg-purple p-[3.2rem] text-[2.8rem] font-medium leading-[100%] text-white transition-all duration-300 hover:bg-purple/80"
                onClick={playAgain}
              >
                🎲 Play Again {randomTopic ? `(${randomTopic.name})` : `(${name})`}
              </button>
              <button
                className={`mobile:text-[1.8rem] mobile:p-7 mobile:rounded-[1.2rem] w-full rounded-[2.4rem] p-[3.2rem] text-[2.8rem] font-medium leading-[100%] transition-all duration-300 border-2 ${
                  darkMode 
                    ? 'border-white text-white hover:bg-white hover:text-dark-navy' 
                    : 'border-dark-navy text-dark-navy hover:bg-dark-navy hover:text-white'
                }`}
                onClick={backToHome}
              >
                🏠 Back to Home
              </button>
            </div>
          </div>
        </>
      ) : (
        <>
          {/* Review Screen */}
          <div className="col-span-2">
            <div className="flex items-center justify-between mb-8">
              <h2
                className={`text-[4rem] font-light leading-[100%] transition-all duration-300 ${
                  darkMode ? "text-white" : "text-dark-navy"
                }`}
              >
                Answer Review
              </h2>
              <button
                className={`text-[1.6rem] px-6 py-3 rounded-xl font-medium transition-all duration-300 ${
                  darkMode 
                    ? 'bg-gray-700 text-white hover:bg-gray-600' 
                    : 'bg-gray-200 text-dark-navy hover:bg-gray-300'
                }`}
                onClick={() => setShowReview(false)}
              >
                ← Back to Score
              </button>
            </div>
            
            <div className="space-y-6 max-h-[60vh] overflow-y-auto">
              {userAnswers.map((answer, index) => (
                <div
                  key={index}
                  className={`p-6 rounded-xl border-2 transition-all duration-300 ${
                    answer.isCorrect
                      ? darkMode 
                        ? 'border-green-400 bg-green-900/20' 
                        : 'border-green-500 bg-green-50'
                      : darkMode
                        ? 'border-red-400 bg-red-900/20'
                        : 'border-red-500 bg-red-50'
                  }`}
                >
                  <div className="flex items-start gap-4">
                    <span className={`text-2xl ${answer.isCorrect ? '✅' : '❌'}`}>
                      {answer.isCorrect ? '✅' : '❌'}
                    </span>
                    <div className="flex-1">
                      <h3 className={`text-[1.8rem] font-medium mb-3 ${
                        darkMode ? 'text-white' : 'text-dark-navy'
                      }`}>
                        Question {index + 1}: {answer.question}
                      </h3>
                      
                      <div className="space-y-2">
                        <p className={`text-[1.4rem] ${
                          darkMode ? 'text-gray-300' : 'text-gray-700'
                        }`}>
                          <span className="font-medium">Your answer:</span> {answer.userAnswer}
                        </p>
                        
                        {!answer.isCorrect && (
                          <p className={`text-[1.4rem] ${
                            darkMode ? 'text-green-300' : 'text-green-700'
                          }`}>
                            <span className="font-medium">Correct answer:</span> {answer.correctAnswer}
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
}

export default FinishedScreen;
