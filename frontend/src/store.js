import { configureStore } from "@reduxjs/toolkit";
import homeReducer from "./features/home/homeSlice";
import quizReducer from "./features/quiz/quizSlice";
import authReducer from "./features/admin/authSlice";

const store = configureStore({
  reducer: {
    home: homeReducer,
    quiz: quizReducer,
    auth: authReducer,
  },
});

export default store;
