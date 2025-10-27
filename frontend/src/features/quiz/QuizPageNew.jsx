import { useDispatch, useSelector } from "react-redux";
import { setAnswer, setChosenAnswer, setIndex, setScore, addUserAnswer, jumpToQuestion } from "./quizSlice";
import { useEffect, useState, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import QuizHeader from "../../components/quiz/QuizHeader";
import ProgressIndicator from "../../components/quiz/ProgressIndicator";
import QuestionCard from "../../components/quiz/QuestionCard";
import AnswerGrid from "../../components/quiz/AnswerGrid";
import QuizFooter from "../../components/quiz/QuizFooter";
import QuizSidebar from "../../components/quiz/QuizSidebar";
import "../../components/quiz/QuizComponents.css";

function QuizPageNew({ question, quiz }) {
  const [showExplanation, setShowExplanation] = useState(false);
  const [answeredQuestions, setAnsweredQuestions] = useState([]);
  const [flaggedQuestions, setFlaggedQuestions] = useState([]); // Track flagged questions
  const allAnswersRef = useRef({}); // Use ref to track answers without triggering re-renders
  const [searchParams] = useSearchParams();
  const isRandomFromURL = searchParams.get('random') === 'true';

  const { index, questions, chosenAnswer, isRandomQuiz, randomTopic } = useSelector(
    (state) => state.quiz,
  );
  const { darkMode } = useSelector((state) => state.home);

  const dispatch = useDispatch();
  const navigate = useNavigate();

  // Determine if question has code in options
  const hasCodeInOptions = question.options && question.options.length > 0 && 
    (question.options[0].includes('{') || question.options[0].includes('function') || 
     question.options[0].includes('var') || question.options[0].includes('const'));

  const handleAnswerSelect = (answer) => {
    dispatch(setChosenAnswer(answer));
    
    // Save answer for this question directly in ref
    allAnswersRef.current[index] = answer;
    
    // Mark as answered if not already
    if (!answeredQuestions.includes(index)) {
      setAnsweredQuestions(prev => [...prev, index]);
    }
  };

  const handleNextQuestion = () => {
    // Just move to next question
    dispatch(setIndex(index));
  };

  const handlePreviousQuestion = () => {
    if (index > 0) {
      dispatch(jumpToQuestion(index - 1));
    }
  };

  const handleToggleFlag = () => {
    if (flaggedQuestions.includes(index)) {
      // Remove flag
      setFlaggedQuestions(prev => prev.filter(i => i !== index));
    } else {
      // Add flag
      setFlaggedQuestions(prev => [...prev, index]);
    }
  };

  const handleFinishQuiz = () => {
    // Add current answer if exists
    const finalAnswers = { ...allAnswersRef.current };
    if (chosenAnswer) {
      finalAnswers[index] = chosenAnswer;
    }

    // Calculate score and save all answers
    questions.forEach((question, idx) => {
      const userAnswer = finalAnswers[idx] || "";
      const isCorrect = userAnswer === question.answer;
      
      if (isCorrect) {
        dispatch(setScore()); // Dispatch for each correct answer
      }
      
      dispatch(addUserAnswer({
        questionIndex: idx,
        userAnswer: userAnswer,
        correctAnswer: question.answer,
        isCorrect: isCorrect
      }));
    });

    // Navigate to results
    const isRandom = isRandomFromURL || isRandomQuiz;
    if (isRandom) {
      navigate("/finished?random=true");
    } else {
      navigate("/finished");
    }
  };

  const handleQuestionSelect = (questionIndex) => {
    // Allow navigation to any question
    if (questionIndex !== index) {
      dispatch(jumpToQuestion(questionIndex));
    }
  };

  useEffect(() => {
    // Set correct answer for comparison
    let correctAnswerText = question.answer;
    if (question.id === 195926 && question.options && question.options.length > 1) {
      correctAnswerText = question.options[1];
    }
    dispatch(setAnswer(correctAnswerText || question.answer));
  }, [question.answer, question.id, question.options, dispatch]);

  // Separate effect to restore answer when question index changes
  useEffect(() => {
    // Restore the selected answer for this question if it exists
    const savedAnswer = allAnswersRef.current[index];
    if (savedAnswer) {
      dispatch(setChosenAnswer(savedAnswer));
    } else {
      dispatch(setChosenAnswer(""));
    }
  }, [index, dispatch]);

  // Determine topic name
  const topicName = (isRandomFromURL || isRandomQuiz) && randomTopic 
    ? randomTopic.name 
    : quiz?.title || "Quiz";

  const isLastQuestion = index === questions.length - 1;
  const progressPercentage = ((index + 1) / questions.length) * 100;

  return (
    <div className={`quiz-container ${darkMode ? 'bg-dark-navy' : 'bg-light-grey'}`}>
      {/* Header */}
      <QuizHeader 
        topicName={topicName}
        currentQuestion={index + 1}
        totalQuestions={questions.length}
        timeRemaining={null} // Can add timer functionality later
        darkMode={darkMode}
      />

      {/* Main Content */}
      <div className="quiz-main-layout">
        {/* Left Sidebar - Question Navigator */}
        <QuizSidebar 
          totalQuestions={questions.length}
          currentQuestion={index}
          answeredQuestions={answeredQuestions}
          flaggedQuestions={flaggedQuestions}
          onQuestionSelect={handleQuestionSelect}
          darkMode={darkMode}
        />

        {/* Center - Quiz Content */}
        <div className={`quiz-content ${darkMode ? 'bg-dark-navy' : 'bg-white'}`}>
          {/* Progress Bar */}
          <ProgressIndicator 
            current={index + 1}
            total={questions.length}
            percentage={progressPercentage}
            darkMode={darkMode}
          />

          {/* Question */}
          <QuestionCard 
            question={question.question}
            code={question.question_code}
            codeLanguage={question.question_code_language || 'javascript'}
            difficulty={question.difficulty}
            tags={question.tags}
            imageUrl={question.question_image_url}
            imageAlt={question.question_image_alt}
            topic={quiz?.topic}
            darkMode={darkMode}
          />

          {/* Answer Options */}
          <AnswerGrid 
            options={question.options || []}
            selectedAnswer={chosenAnswer}
            onAnswerSelect={handleAnswerSelect}
            showCorrect={false}
            hasCode={hasCodeInOptions}
            codeLanguage={question.options_code_language || 'css'}
            darkMode={darkMode}
          />
        </div>
      </div>

      {/* Footer */}
      <QuizFooter 
        hasPrevious={index > 0}
        hasNext={!isLastQuestion}
        selectedAnswer={chosenAnswer}
        onPrevious={handlePreviousQuestion}
        onNext={handleNextQuestion}
        onSubmit={handleFinishQuiz}
        isLastQuestion={isLastQuestion}
        showExplanation={showExplanation}
        onToggleExplanation={() => setShowExplanation(!showExplanation)}
        isFlagged={flaggedQuestions.includes(index)}
        onFlag={handleToggleFlag}
        answered={answeredQuestions.length}
        remaining={questions.length - answeredQuestions.length}
        totalQuestions={questions.length}
        darkMode={darkMode}
      />
    </div>
  );
}

export default QuizPageNew;
