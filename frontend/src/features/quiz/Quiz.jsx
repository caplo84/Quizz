import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useParams, useSearchParams } from "react-router-dom";
import { setQuestions, resetQuiz, fetchRandomQuestions } from "./quizSlice";
import { topicsApi, quizzesApi } from "../../services/api.js";
import QuizPageNew from "./QuizPageNew";

function Quiz() {
  const { type } = useParams(); // Get the topic name from URL
  const [searchParams] = useSearchParams();
  const isRandomQuiz = searchParams.get('random') === 'true';
  const { questions, index, isRandomQuiz: storeIsRandomQuiz, randomTopic, usedQuestionIds } = useSelector((state) => state.quiz);
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
        

        
        if (isRandomQuiz && storeIsRandomQuiz && randomTopic) {
          // For random quiz, only reset question-specific state
          try {
            const result = await dispatch(fetchRandomQuestions({
              topicSlug: randomTopic.slug,
              limit: 10,
              excludeIds: usedQuestionIds
            })).unwrap();
            
            if (!result.data || result.data.length === 0) {
              setError(`No questions available for ${randomTopic.name} topic yet. Check back later!`);
              setLoading(false);
              return;
            }
            
            // Transform questions to match expected format
            const transformedQuestions = result.data.map((q, index) => {
              try {
                // Find the correct answer from choices
                const correctChoice = q.choices?.find(choice => choice.is_correct);
                const answer = correctChoice ? correctChoice.choice_text || correctChoice.text || correctChoice.option : null;
                
                return {
                  id: q.id,
                  question: q.question_text || q.text || q.question,
                  options: q.choices?.map(choice => choice.choice_text || choice.text || choice.option) || q.options || [],
                  answer: answer,
                  // Content separation fields
                  question_image_url: q.question_image_url,
                  question_image_alt: q.question_image_alt,
                  question_code: q.question_code,
                  question_code_language: q.question_code_language,
                  choice_codes: q.choices?.map(choice => choice.choice_code) || [],
                  choice_code_languages: q.choices?.map(choice => choice.choice_code_language) || []
                };
              } catch (error) {
                console.error(`Error transforming question ${index}:`, error, q);
                throw error;
              }
            });
            
            dispatch(setQuestions(transformedQuestions));
            
            // Set current quiz info for random topic
            setCurrentQuiz({
              title: `${randomTopic.name} Quiz`,
              slug: randomTopic.slug,
              isRandom: true,
              topic: randomTopic
            });
            
          } catch (randomError) {
            console.error('Failed to load random questions:', randomError);
            setError(randomError.message || 'Failed to load random questions');
          }
        } else {
          // Reset quiz state for regular topic selection
          dispatch(resetQuiz());
          
          // For regular topic selection, ALSO use random questions API
          
          // First get topic info to get topic data
          try {
            const topicQuizzes = await topicsApi.getQuizzes(type);
            
            if (!topicQuizzes || topicQuizzes.length === 0) {
              setError(`No quizzes available for ${type} topic yet. Check back later!`);
              setLoading(false);
              return;
            }
            
            // Now load 10 random questions using the same API
            const result = await dispatch(fetchRandomQuestions({
              topicSlug: type,
              limit: 10,
              excludeIds: [] // No exclusions for first batch of regular topic
            })).unwrap();
            
            if (!result.data || result.data.length === 0) {
              setError(`No questions available for ${type} topic yet. Check back later!`);
              setLoading(false);
              return;
            }
            
            
            // Transform questions to match expected format
            const transformedQuestions = result.data.map(q => ({
              id: q.id,
              question: q.question_text || q.text || q.question,
              options: q.choices?.map(choice => choice.choice_text || choice.text || choice.option) || q.options || [],
              answer: q.choices?.find(choice => choice.is_correct)?.choice_text || q.answer,
              // Content separation fields
              question_image_url: q.question_image_url,
              question_image_alt: q.question_image_alt,
              question_code: q.question_code,
              question_code_language: q.question_code_language,
              choice_codes: q.choices?.map(choice => choice.choice_code) || [],
              choice_code_languages: q.choices?.map(choice => choice.choice_code_language) || []
            }));
            
            dispatch(setQuestions(transformedQuestions));
            
            // Set current quiz info for regular topic
            setCurrentQuiz({
              title: `${topicQuizzes[0]?.title || type} Quiz`,
              slug: type,
              isRandom: false,
              topic: {
                id: topicQuizzes[0]?.topic_id,
                name: topicQuizzes[0]?.title || type,
                slug: type
              }
            });
            
          } catch (error) {
            console.error('Failed to load questions for topic:', error);
            setError(error.message || 'Failed to load questions');
          }
        }
        
      } catch (error) {
        console.error('Failed to load quiz:', error);
        setError(error.message);
      } finally {
        setLoading(false);
      }
    }
    
    if (type) {
      loadQuizData();
    }
  }, [type, dispatch, isRandomQuiz, storeIsRandomQuiz, randomTopic, usedQuestionIds]);

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
        <h2 className="text-xl font-semibold text-gray-700 mb-2">No Questions Available</h2>
        <p className="text-gray-500 max-w-md">{error || 'There are currently no questions for this topic. Please try another topic or check back later.'}</p>
        <button
          onClick={() => window.location.href = '/'}
          className="mt-6 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          Back to Home
        </button>
      </div>
    );
  }

  return (
    <div>
      {currentQuestion && currentQuiz && (
        <QuizPageNew 
          question={currentQuestion} 
          quiz={currentQuiz}
        />
      )}
    </div>
  );
}

export default Quiz;
