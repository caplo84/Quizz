import apiClient from './apiClient.js';

export const quizzesApi = {
  getBySlug: (slug) => apiClient.get(`/quizzes/${slug}`),
  getQuestions: (slug) => apiClient.get(`/quizzes/${slug}/questions`),
  createAttempt: (slug, data) => apiClient.post(`/quizzes/${slug}/attempts`, data),
  getAttempt: (slug, attemptId) => apiClient.get(`/quizzes/${slug}/attempts/${attemptId}`),
  submitAttempt: (slug, attemptId, data) => apiClient.put(`/quizzes/${slug}/attempts/${attemptId}`, data),
};