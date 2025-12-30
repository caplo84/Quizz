import { createSlice } from "@reduxjs/toolkit";

const getInitialDarkMode = () => {
  if (typeof window === "undefined") return false;

  const savedTheme = window.localStorage.getItem("quizz-theme");
  if (savedTheme === "dark") return true;
  if (savedTheme === "light") return false;

  return window.matchMedia?.("(prefers-color-scheme: dark)")?.matches ?? false;
};

const initialState = {
  name: "",
  icon: "",
  quizzes: [],
  darkMode: getInitialDarkMode(),
};

const homeSlice = createSlice({
  name: "home",
  initialState,
  reducers: {
    selectQuiz(state, action) {
      state.name = action.payload;
    },
    selectIcon(state, action) {
      state.icon = action.payload;
    },
    setQuizzes(state, action) {
      state.quizzes = action.payload;
    },
    resetMode(state) {
      state.name = initialState.name;
      state.icon = initialState.icon;
      state.quizzes = initialState.quizzes;
    },
    setDarkMode(state, action) {
      if (typeof action.payload === "boolean") {
        state.darkMode = action.payload;
        return;
      }

      state.darkMode = !state.darkMode;
    },
  },
});

export const { selectQuiz, selectIcon, setQuizzes, resetMode, setDarkMode } =
  homeSlice.actions;
export default homeSlice.reducer;
