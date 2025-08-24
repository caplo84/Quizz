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

const initialState = {
  index: 0,
  score: 0,
  correctAnswer: "",
  chosenAnswer: "",
  questions: [],
};

const quizSlice = createSlice({
  name: "quiz",
  initialState,
  reducers: {
    setIndex(state) {
      state.index++;
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
    },
    resetQuiz(state) {
      state.index = initialState.index;
      state.score = initialState.score;
      state.correctAnswer = initialState.correctAnswer;
      state.chosenAnswer = initialState.chosenAnswer;
      state.questions = initialState.questions;
    },
  },
});

export const {
  setIndex,
  setScore,
  setAnswer,
  setQuestions,
  setChosenAnswer,
  resetQuiz,
} = quizSlice.actions;

export default quizSlice.reducer;
