import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useParams } from "react-router-dom";
import { setQuestions, resetQuiz } from "./quizSlice";
import { topicsApi, quizzesApi } from "../../services/api.js";
import QuizPage from "./QuizPage";

function Quiz() {
  const { type } = useParams(); // Get the topic name from URL
  const { questions, index } = useSelector((state) => state.quiz);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [currentQuiz, setCurrentQuiz] = useState(null);

  const dispatch = useDispatch();

  const currentQuestion = questions[index];

  useEffect(() => {
    async function loadQuizData() {
      try {
        setLoading(true);
        setError(null);
        
        console.log('🔄 Loading quiz for topic:', type);
        console.log('🔍 Type value:', type, 'Type of:', typeof type);
        console.log('🔗 API Base URL:', `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1`);
        
        // Reset quiz state
        dispatch(resetQuiz());
        
        // First get quizzes for this topic
        const topicQuizzes = await topicsApi.getQuizzes(type);
        console.log('📋 Topic quizzes:', topicQuizzes);
        console.log('📋 Topic quizzes type:', typeof topicQuizzes, 'Is Array:', Array.isArray(topicQuizzes));
        
        if (!topicQuizzes || topicQuizzes.length === 0) {
          setError(`No quizzes available for ${type} topic yet. Check back later!`);
          setLoading(false);
          return;
        }
        
        // Take the first quiz for this topic (or you could let user select)
        const quiz = topicQuizzes[0];
        
        if (!quiz || !quiz.slug) {
          setError(`Invalid quiz data for ${type} topic. Please try again later.`);
          setLoading(false);
          return;
        }
        
        setCurrentQuiz(quiz);
        
        // Now get questions for this specific quiz
        const questions = await quizzesApi.getQuestions(quiz.slug);
        console.log('❓ Quiz questions:', questions);
        console.log('❓ First question structure:', questions[0]);
        console.log('❓ First question choices:', questions[0]?.choices);
        
        if (!questions || questions.length === 0) {
          setError(`No questions available for this ${quiz.title || type} quiz yet. Check back later!`);
          setLoading(false);
          return;
        }
        
        // Transform questions to match expected format
        const transformedQuestions = questions.map(q => ({
          question: q.question_text || q.text || q.question,
          options: q.choices?.map(choice => choice.choice_text || choice.text || choice.option) || q.options || [],
          answer: q.choices?.find(choice => choice.is_correct)?.choice_text || q.answer
        }));
        
        dispatch(setQuestions(transformedQuestions));
        
      } catch (error) {
        console.error('💥 Failed to load quiz:', error);
        setError(error.message);
      } finally {
        setLoading(false);
      }
    }
    
    if (type) {
      loadQuizData();
    }
  }, [type, dispatch]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-lg">Loading quiz...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-64 text-center">
        <div className="mb-4">
          <svg className="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.172 16.172a4 4 0 015.656 0M9 12h6m-6-4h6m2 5.291A7.962 7.962 0 0112 20a7.962 7.962 0 01-5-1.709M15 3H9a2 2 0 00-2 2v1.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 000 1.414l6.414 6.414A1 1 0 019 21.586V20a2 2 0 002-2h6a2 2 0 002-2V5a2 2 0 00-2-2z" />
          </svg>
        </div>
        <h2 className="text-xl font-semibold text-gray-700 mb-2">Quiz Coming Soon!</h2>
        <p className="text-gray-500 max-w-md">{error}</p>
        <button 
          onClick={() => window.history.back()} 
          className="mt-6 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          Go Back
        </button>
      </div>
    );
  }

  return (
    <div>
      {currentQuestion && currentQuiz && (
        <QuizPage 
          key={currentQuestion.question} 
          question={currentQuestion} 
          quiz={currentQuiz}
        />
      )}
    </div>
  );
}

export default Quiz;
