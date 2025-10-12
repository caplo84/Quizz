import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { topicsApi, quizzesApi } from "../../services/api.js";

export const fetchTopics = createAsyncThunk(
  "quiz/fetchTopics",
  async () => {
    return await topicsApi.getAll();
  }
);

export const fetchQuizzesByTopic = createAsyncThunk(
  "quiz/fetchQuizzesByTopic",
  async (topic) => {
    return await topicsApi.getQuizzes(topic);
  }
);

export const fetchQuizBySlug = createAsyncThunk(
  "quiz/fetchQuizBySlug",
  async (slug) => {
    return await quizzesApi.getBySlug(slug);
  }
);

export const fetchRandomQuestions = createAsyncThunk(
  "quiz/fetchRandomQuestions",
  async ({ topicSlug, limit = 10, excludeIds = [] }) => {
    const excludeParam = excludeIds.length > 0 ? excludeIds.join(',') : '';
    const response = await fetch(
      `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/topics/${topicSlug}/questions/random?limit=${limit}&exclude=${excludeParam}`
    );
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.message || `Failed to fetch random questions: ${response.status}`);
    }
    const data = await response.json();
    
    // Ensure we have questions
    if (!data.data || data.data.length === 0) {
      throw new Error(`No questions available for ${topicSlug} topic`);
    }
    
    return data;
  }
);

export const fetchQuestionsWithAnswers = createAsyncThunk(
  "quiz/fetchQuestionsWithAnswers",
  async ({ topicSlug, questionIds }) => {
    // We'll fetch specific questions by their IDs with answers
    // For now, since we don't have a specific endpoint, we'll use the same random endpoint
    // but this should ideally be a different endpoint that accepts question IDs
    const response = await fetch(
      `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/topics/${topicSlug}/questions/random?limit=10&include_answers=true`
    );
    if (!response.ok) {
      throw new Error(`Failed to fetch questions with answers: ${response.status}`);
    }
    const data = await response.json();
    return data.data;
  }
);

const initialState = {
  index: 0,
  score: 0,
  correctAnswer: "",
  chosenAnswer: "",
  questions: [],
  // Random quiz specific state
  isRandomQuiz: false,
  randomTopic: null,
  usedQuestionIds: [],
  totalBatches: 0,
  currentBatch: 1,
  // Answer tracking for review
  userAnswers: [], // Array to store user answers for review
};

const quizSlice = createSlice({
  name: "quiz",
  initialState,
  reducers: {
    setIndex(state) {
      state.index++;
    },
    jumpToQuestion(state, action) {
      state.index = action.payload;
    },
    setScore(state) {
      state.score++;
    },
    setAnswer(state, action) {
      state.correctAnswer = action.payload;
    },
    setChosenAnswer(state, action) {
      state.chosenAnswer = action.payload;
    },
    setQuestions(state, action) {
      state.questions = action.payload;
      // Reset user answers when new questions are loaded
      state.userAnswers = [];
    },
    // Add user answer for review
    addUserAnswer(state, action) {
      const { questionIndex, userAnswer, correctAnswer, isCorrect } = action.payload;
      state.userAnswers[questionIndex] = {
        userAnswer,
        correctAnswer,
        isCorrect,
        question: state.questions[questionIndex]?.question,
        questionId: state.questions[questionIndex]?.id
      };
    },
    // Random quiz specific actions
    startRandomQuiz(state, action) {
      state.isRandomQuiz = true;
      state.randomTopic = action.payload.topic;
      state.usedQuestionIds = [];
      state.currentBatch = 1;
      state.totalBatches = 1;
      state.index = 0;
      state.score = 0;
      state.correctAnswer = "";
      state.chosenAnswer = "";
    },
    addUsedQuestionIds(state, action) {
      state.usedQuestionIds = [...state.usedQuestionIds, ...action.payload];
    },
    startNewBatch(state) {
      state.currentBatch++;
      state.totalBatches++;
      state.index = 0;
      state.correctAnswer = "";
      state.chosenAnswer = "";
    },

    resetQuiz(state) {
      state.index = initialState.index;
      state.score = initialState.score;
      state.correctAnswer = initialState.correctAnswer;
      state.chosenAnswer = initialState.chosenAnswer;
      state.questions = initialState.questions;
      state.isRandomQuiz = initialState.isRandomQuiz;
      state.randomTopic = initialState.randomTopic;
      state.usedQuestionIds = initialState.usedQuestionIds;
      state.totalBatches = initialState.totalBatches;
      state.currentBatch = initialState.currentBatch;
      state.userAnswers = initialState.userAnswers;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(fetchRandomQuestions.fulfilled, (state, action) => {
      state.questions = action.payload.data;
      // Add the new question IDs to used questions
      const newQuestionIds = action.payload.data.map(q => q.id);
      state.usedQuestionIds = [...state.usedQuestionIds, ...newQuestionIds];
    });
    builder.addCase(fetchQuestionsWithAnswers.fulfilled, (state, action) => {
      // Update correct answers in userAnswers
      const questionsWithAnswers = action.payload;
      state.userAnswers.forEach((userAnswer, index) => {
        if (userAnswer.questionId) {
          const matchingQuestion = questionsWithAnswers.find(q => q.id === userAnswer.questionId);
          if (matchingQuestion) {
            const correctChoice = matchingQuestion.choices?.find(choice => choice.is_correct);
            state.userAnswers[index].correctAnswer = correctChoice?.choice_text || '';
          }
        }
      });
    });
  },
});

export const {
  setIndex,
  jumpToQuestion,
  setScore,
  setAnswer,
  setQuestions,
  setChosenAnswer,
  resetQuiz,
  startRandomQuiz,
  addUsedQuestionIds,
  startNewBatch,
  addUserAnswer,

} = quizSlice.actions;

export default quizSlice.reducer;
