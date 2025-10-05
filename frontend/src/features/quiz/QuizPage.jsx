import { useDispatch, useSelector } from "react-redux";
import { setAnswer, setChosenAnswer, setIndex, setScore, addUserAnswer } from "./quizSlice";
import { useEffect, useState } from "react";
import QuizOptions from "./QuizOptions";
import { useNavigate, useSearchParams } from "react-router-dom";
import CodeBlock from "../../ui/CodeBlock";
import ImageDisplay from "../../ui/ImageDisplay";
import { buttonStyles } from "../../utils/commonStyles";

function QuizPage({ question, quiz }) {
  const [userAnswer, setUserAnswer] = useState("");
  const [isAnswered, setIsAnswered] = useState(false);
  const [isTransitioning, setIsTransitioning] = useState(false);
  const [searchParams] = useSearchParams();
  const isRandomFromURL = searchParams.get('random') === 'true';

  const { index, questions, chosenAnswer, correctAnswer, score, isRandomQuiz, randomTopic } = useSelector(
    (state) => state.quiz,
  );
  const { darkMode } = useSelector((state) => state.home);

  const dispatch = useDispatch();
  const navigate = useNavigate();

  function submitAnswer() {
    if (chosenAnswer === "") setIsAnswered(true);

    setUserAnswer(chosenAnswer);
  }

  function nextQuestion() {
    setIsTransitioning(true);
    
    setTimeout(() => {
      const isCorrect = userAnswer === correctAnswer;
      if (isCorrect) dispatch(setScore());

      // Track the user's answer for review
      dispatch(addUserAnswer({
        questionIndex: index,
        userAnswer: userAnswer,
        correctAnswer: correctAnswer,
        isCorrect: isCorrect
      }));

      dispatch(setIndex(index));
      dispatch(setChosenAnswer(""));
      setIsAnswered(false);
      setIsTransitioning(false);
    }, 300); // Small delay for smooth transition
  }

  function finishQuiz() {
    const isCorrect = userAnswer === correctAnswer;
    if (isCorrect) dispatch(setScore());

    // Track the final answer for review
    dispatch(addUserAnswer({
      questionIndex: index,
      userAnswer: userAnswer,
      correctAnswer: correctAnswer,
      isCorrect: isCorrect
    }));

    // Check if this is a random quiz (prioritize URL parameter)
    const isRandom = isRandomFromURL || isRandomQuiz;
    
    if (isRandom) {
      navigate("/finished?random=true");
    } else {
      navigate("/finished");
    }
  }

  const btnClass = buttonStyles.submit;

  useEffect(() => {
    // TEMPORARY FIX: Hard-code correct answer for Q36 (Question ID 195926)
    // Since API hides correct answers during quiz, we need this for Answer Review
    let correctAnswerText = question.answer;
    if (question.id === 195926 && question.options && question.options.length > 1) {
      // For the Android layout question, Choice B is correct
      correctAnswerText = question.options[1]; // Index 1 = Choice B
    }
    dispatch(setAnswer(correctAnswerText || question.answer));
  }, [question.answer, question.id, question.options, dispatch]);

  return (
    <div className={`grid grid-cols-2 gap-x-24 gap-y-12 desktop:grid-cols-1 desktop:gap-x-0 transition-all duration-300 ${
      isTransitioning ? 'opacity-50 scale-95' : 'opacity-100 scale-100'
    }`}>
      <div className="flex flex-col gap-11 desktop:pb-12 mobile:gap-5 mobile:pb-4">
        {/* Topic name for random quiz */}
        {(isRandomFromURL || isRandomQuiz) && randomTopic && (
          <div className={`text-[1.6rem] font-medium leading-[150%] transition-all duration-300 mobile:text-[1.2rem] ${
            darkMode ? "text-light-bluish" : "text-purple"
          }`}>
            Topic: {randomTopic.name}
          </div>
        )}
        
        <p
          className={`text-[2rem] italic leading-[150%] transition-all duration-300 mobile:text-[1.4rem] ${
            darkMode ? "text-light-bluish" : "text-dark-navy"
          }`}
        >
          Question {index + 1} of {questions.length}
        </p>
        <h2
          className={`text-[3.6rem] font-medium leading-[120%] transition-all duration-300 mobile:text-[2rem] ${
            darkMode ? "text-white" : "text-dark-navy"
          }`}
        >
          {question.question}
        </h2>
        
        {/* Question Code Block */}
        {question.question_code && (
          <div className="mt-6">
            <CodeBlock 
              code={question.question_code} 
              language={question.question_code_language}
              className="shadow-sm"
            />
          </div>
        )}
        
        {/* Question Image */}
        {question.question_image_url && (
          <div className="mt-6">
            <ImageDisplay 
              imageUrl={question.question_image_url} 
              altText={question.question_image_alt}
              topic="android"
              className="shadow-sm"
            />
          </div>
        )}
        <div className="mt-auto">
          <progress
            max={questions.length}
            className={`w-full ${
              darkMode
                ? "[&::-moz-progress-bar]:bg-navy [&::-webkit-progress-bar]:bg-navy"
                : "[&::-moz-progress-bar]:bg-white [&::-webkit-progress-bar]:bg-white"
            }  [&::-webkit-progress-bar]:rounded-full [&::-webkit-progress-bar]:p-[1.5px] [&::-webkit-progress-value]:rounded-full [&::-webkit-progress-value]:bg-purple`}
            value={score}
          />
        </div>
        {/* <div
          className={`relative mt-auto h-6 w-full rounded-full p-2 transition-all duration-300 before:absolute before:top-1/4 before:h-3 before:w-3/4 before:rounded-full before:bg-purple desktop:before:w-1/2 mobile:before:w-1/4 ${
            darkMode ? "bg-navy" : "bg-white"
          }`}
        ></div> */}
      </div>
      <div className="space-y-10 mobile:space-y-5">
        {(Array.isArray(question.options) ? question.options : []).map((item, index) => (
          <QuizOptions
            key={item}
            option={item}
            optionIndex={index}
            userAnswer={userAnswer}
            setIsAnswered={setIsAnswered}
            question={question}
          />
        ))}
      </div>
      <div className="col-start-2 desktop:col-start-1">
        {userAnswer === "" && (
          <button className={btnClass} onClick={submitAnswer}>
            Submit Answer
          </button>
        )}

        {userAnswer !== "" && index < questions.length - 1 ? (
          <button className={btnClass} onClick={nextQuestion}>
            Next Question
          </button>
        ) : (
          userAnswer !== "" &&
          index === questions.length - 1 && (
            <button className={btnClass} onClick={finishQuiz}>
              Submit Quiz
            </button>
          )
        )}

        <div
          className={`mt-12 flex items-center justify-center gap-8 ${
            isAnswered ? "opacity-100" : "opacity-0"
          }`}
        >
          <img src="/icon-incorrect.svg" alt="incorrect icon" />
          <p className="text-[2.4rem] leading-[150%] text-red mobile:text-[1.8rem]">
            Please select an answer
          </p>
        </div>
      </div>
    </div>
  );
}

export default QuizPage;
