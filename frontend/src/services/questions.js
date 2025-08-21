import { api, getQuiz as getQuizLegacy } from './api.js';

// Export the legacy function for backward compatibility
export async function getQuiz() {
  return getQuizLegacy();
}

// Export the new API functions
export const {
  getTopics,
  getQuizzes,
  getQuiz: getQuizBySlug,
  getQuizQuestions,
  createAttempt,
  submitAttempt,
  getAttemptResults,
  healthCheck,
} = api;

export default api;
