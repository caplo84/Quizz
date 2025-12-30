import { useDispatch, useSelector } from "react-redux";
import { setAnswer, setChosenAnswer, setScore, addUserAnswer, jumpToQuestion } from "./quizSlice";
import { useEffect, useState, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { ChevronLeft, ArrowLeft, ArrowRight, CheckSquare, Moon, Sun } from "lucide-react";
import CodeBlock from "../../ui/CodeBlock";
import AnswerInput from "../../components/quiz/AnswerInput";
import { setDarkMode } from "../home/homeSlice";

const hasAnsweredValue = (value) => {
  if (Array.isArray(value)) {
    return value.length > 0;
  }

  if (typeof value === "number") {
    return !Number.isNaN(value);
  }

  if (typeof value === "string") {
    return value.trim().length > 0;
  }

  return value !== undefined && value !== null;
};

const normalizeText = (value) => String(value || "").trim().toLowerCase();

function QuizPageNew({ question, quiz }) {
  const [answeredQuestions, setAnsweredQuestions] = useState([]);
  const allAnswersRef = useRef({}); // Use ref to track answers without triggering re-renders
  const [searchParams] = useSearchParams();
  const isRandomFromURL = searchParams.get('random') === 'true';

  const { index, questions, chosenAnswer, isRandomQuiz, randomTopic } = useSelector(
    (state) => state.quiz,
  );
  const { darkMode } = useSelector((state) => state.home);

  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleAnswerSelect = (answer) => {
    dispatch(setChosenAnswer(answer));
    
    // Save answer for this question directly in ref
    allAnswersRef.current[index] = answer;

    // Mark / unmark as answered based on current value
    setAnsweredQuestions((prev) => {
      const alreadyAnswered = prev.includes(index);
      const isAnsweredNow = hasAnsweredValue(answer);

      if (isAnsweredNow && !alreadyAnswered) {
        return [...prev, index];
      }
      if (!isAnsweredNow && alreadyAnswered) {
        return prev.filter((questionIndex) => questionIndex !== index);
      }
      return prev;
    });
  };

  const handleFinishQuiz = () => {
    // Add current answer if exists
    const finalAnswers = { ...allAnswersRef.current };
    if (hasAnsweredValue(chosenAnswer)) {
      finalAnswers[index] = chosenAnswer;
    }

    // Calculate score and save all answers by question type
    questions.forEach((question, idx) => {
      const questionType = question.question_type || "single";
      const userAnswer = finalAnswers[idx];

      let isCorrect = false;
      let userAnswerDisplay = "";
      let correctAnswerDisplay = "";

      if (questionType === "multiple") {
        const selectedIndices = Array.isArray(userAnswer) ? userAnswer : [];
        const correctIndices = Array.isArray(question.answer) ? question.answer : [];

        isCorrect =
          correctIndices.length === selectedIndices.length &&
          correctIndices.every((correctIndex) => selectedIndices.includes(correctIndex));

        userAnswerDisplay = selectedIndices
          .map((answerIndex) => question.options?.[answerIndex]?.text || "")
          .filter(Boolean)
          .join(", ");
        correctAnswerDisplay = correctIndices
          .map((answerIndex) => question.options?.[answerIndex]?.text || "")
          .filter(Boolean)
          .join(", ");
      } else if (questionType === "text_input") {
        const userText = typeof userAnswer === "string" ? userAnswer : "";
        const correctText = question.answer || question.correct_text_answer || "";
        isCorrect = normalizeText(userText) === normalizeText(correctText);
        userAnswerDisplay = userText;
        correctAnswerDisplay = String(correctText);
      } else {
        const selectedIndex = typeof userAnswer === "number" ? userAnswer : null;
        const correctIndex = typeof question.answer === "number" ? question.answer : null;
        isCorrect = selectedIndex !== null && correctIndex !== null && selectedIndex === correctIndex;
        userAnswerDisplay = selectedIndex !== null ? question.options?.[selectedIndex]?.text || "" : "";
        correctAnswerDisplay = correctIndex !== null ? question.options?.[correctIndex]?.text || "" : "";
      }
      
      if (isCorrect) {
        dispatch(setScore()); // Dispatch for each correct answer
      }
      
      dispatch(addUserAnswer({
        questionIndex: idx,
        userAnswer: userAnswerDisplay,
        correctAnswer: correctAnswerDisplay,
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

  const handlePreviousQuestion = () => {
    if (index > 0) {
      dispatch(jumpToQuestion(index - 1));
    }
  };

  useEffect(() => {
    // Set correct answer for comparison
    dispatch(setAnswer(question.answer));
  }, [question.answer, question.id, question.options, dispatch]);

  // Separate effect to restore answer when question index changes
  useEffect(() => {
    // Restore the selected answer for this question if it exists
    const savedAnswer = allAnswersRef.current[index];
    if (hasAnsweredValue(savedAnswer)) {
      dispatch(setChosenAnswer(savedAnswer));
    } else {
      dispatch(setChosenAnswer(""));
    }
  }, [index, dispatch]);

  useEffect(() => {
    window.scrollTo({ top: 0, behavior: "smooth" });
  }, [index]);

  // Determine topic name
  const topicName = (isRandomFromURL || isRandomQuiz) && randomTopic 
    ? randomTopic.name 
    : quiz?.title || "Quiz";

  const answeredCount = answeredQuestions.length;
  const progress = questions.length > 0 ? Math.round(((index + 1) / questions.length) * 100) : 0;

  return (
    <div className={`min-h-screen ${darkMode ? "bg-slate-900" : "bg-slate-100"}`}>
      <header className={`sticky top-0 z-30 border-b ${darkMode ? "bg-slate-900/95 border-slate-700" : "bg-white/95 border-slate-200"} backdrop-blur`}>
        <div className="max-w-6xl mx-auto px-4 md:px-6 py-3">
          <div className="flex items-center justify-between gap-3">
            <div className="flex items-center gap-2 min-w-0">
              <button
                type="button"
                onClick={() => navigate("/")}
                className={`p-1.5 rounded-lg ${darkMode ? "hover:bg-slate-800" : "hover:bg-slate-100"}`}
              >
                <ChevronLeft size={18} className={darkMode ? "text-slate-200" : "text-slate-700"} />
              </button>
              <div className="min-w-0">
                <h1 className={`font-bold text-base truncate ${darkMode ? "text-white" : "text-slate-900"}`}>{topicName}</h1>
                <p className={`text-xs ${darkMode ? "text-slate-300" : "text-slate-500"}`}>Question {index + 1} of {questions.length}</p>
              </div>
            </div>
            <div className={`text-xs sm:text-sm ${darkMode ? "text-slate-300" : "text-slate-600"}`}>
              <span>{progress}% complete</span>
              <span className="mx-2">·</span>
              <span>{answeredCount}/{questions.length} answered</span>
              <button
                type="button"
                onClick={() => dispatch(setDarkMode())}
                className={`inline-flex items-center justify-center ml-2 p-1 rounded ${darkMode ? "hover:bg-slate-800" : "hover:bg-slate-100"}`}
                aria-label="Toggle dark mode"
                title="Toggle dark mode"
              >
                {darkMode ? <Sun size={13} className="text-yellow-300" /> : <Moon size={13} className="text-slate-700" />}
              </button>
            </div>
          </div>
          <div className={`mt-2 h-1.5 rounded-full ${darkMode ? "bg-slate-700" : "bg-slate-200"}`}>
            <div className="h-full rounded-full bg-violet-600 transition-all" style={{ width: `${progress}%` }} />
          </div>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-4 md:px-6 py-5 grid grid-cols-1 lg:grid-cols-[190px_1fr] gap-4">
        <aside className={`rounded-xl border p-3 h-fit ${darkMode ? "bg-slate-800 border-slate-700" : "bg-white border-slate-200"}`}>
          <p className={`text-xs font-semibold mb-2 ${darkMode ? "text-slate-200" : "text-slate-700"}`}>Questions</p>
          <div className="space-y-1.5">
            {Array.from({ length: questions.length }, (_, i) => {
              const isAnswered = answeredQuestions.includes(i);
              const isCurrent = i === index;

              return (
                <button
                  key={i}
                  type="button"
                  onClick={() => handleQuestionSelect(i)}
                  className={`w-full text-left px-2.5 py-1.5 rounded-lg text-sm border transition-colors ${
                    isCurrent
                      ? "bg-violet-600 text-white border-violet-600"
                      : isAnswered
                        ? darkMode
                          ? "bg-emerald-900/40 text-emerald-200 border-emerald-700"
                          : "bg-emerald-100 text-emerald-800 border-emerald-300"
                        : darkMode
                          ? "bg-slate-800 text-slate-300 border-slate-700 hover:bg-slate-700"
                          : "bg-white text-slate-600 border-slate-200 hover:bg-slate-50"
                  }`}
                >
                  Q{i + 1}
                </button>
              );
            })}
          </div>
          <button
            type="button"
            onClick={handleFinishQuiz}
            disabled={answeredCount !== questions.length}
            className="mt-3 w-full px-3 py-2 rounded-lg bg-violet-600 text-white text-sm font-semibold hover:bg-violet-700 disabled:opacity-40 disabled:cursor-not-allowed"
          >
            <span className="inline-flex items-center gap-1.5"><CheckSquare size={14} /> Submit Quiz</span>
          </button>
          <p className={`mt-2 text-[11px] text-center ${darkMode ? "text-slate-400" : "text-slate-500"}`}>
            Answer all {questions.length} questions
          </p>
        </aside>

        <section className={`rounded-xl border p-4 sm:p-5 ${darkMode ? "bg-slate-800 border-slate-700" : "bg-white border-slate-200"}`}>
          <div className="flex items-start gap-3 mb-4">
            <span className="w-7 h-7 rounded-lg bg-violet-100 text-violet-700 flex items-center justify-center text-sm font-bold shrink-0">
              {index + 1}
            </span>
            <h2 className={`text-lg sm:text-xl font-semibold leading-snug ${darkMode ? "text-white" : "text-slate-900"}`}>
              {question.question}
            </h2>
          </div>

          {question.question_code && (
            <div className="mb-4">
              <CodeBlock code={question.question_code} language={question.question_code_language || "javascript"} />
            </div>
          )}

          {question.question_image_url && (
            <img
              src={question.question_image_url}
              alt={question.question_image_alt || "Question image"}
              className="w-full max-h-80 object-contain rounded-lg border border-slate-200 mb-4"
            />
          )}

          <AnswerInput
            question={question}
            answer={chosenAnswer}
            onAnswer={handleAnswerSelect}
            darkMode={darkMode}
          />

          <div className="mt-4 flex items-center justify-between">
            <button
              type="button"
              onClick={handlePreviousQuestion}
              disabled={index === 0}
              className={`inline-flex items-center gap-1.5 px-3 py-2 rounded-lg border text-sm font-medium transition-colors disabled:opacity-40 disabled:cursor-not-allowed ${
                darkMode ? "border-slate-600 text-slate-200 hover:bg-slate-700" : "border-slate-300 text-slate-700 hover:bg-slate-50"
              }`}
            >
              <ArrowLeft size={14} /> Previous
            </button>

            {index < questions.length - 1 ? (
              <button
                type="button"
                onClick={() => handleQuestionSelect(index + 1)}
                className="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-violet-600 text-white text-sm font-semibold hover:bg-violet-700"
              >
                Next <ArrowRight size={14} />
              </button>
            ) : (
              <button
                type="button"
                onClick={handleFinishQuiz}
                disabled={answeredCount !== questions.length}
                className="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-violet-600 text-white text-sm font-semibold hover:bg-violet-700 disabled:opacity-40 disabled:cursor-not-allowed"
              >
                <CheckSquare size={14} /> Finish Quiz
              </button>
            )}
          </div>
        </section>
      </main>
    </div>
  );
}

export default QuizPageNew;
