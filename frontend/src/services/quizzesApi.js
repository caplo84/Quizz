import apiClient from './apiClient.js';

export const quizzesApi = {
  getBySlug: async (slug) => {
    const response = await apiClient.get(`/quizzes/${slug}`);
    return response.data || response;
  },
  getQuestions: async (slug) => {
    const response = await apiClient.get(`/quizzes/${slug}/questions`);
    return response.data || response;
  },
  createAttempt: async (slug, data = {}) => {
    const response = await apiClient.post(`/quizzes/${slug}/attempts`, data);
    return response.data || response;
  },
  getAttempt: async (slug, attemptId) => {
    const response = await apiClient.get(`/quizzes/${slug}/attempts/${attemptId}`);
    return response.data || response;
  },
  submitAttempt: async (slug, attemptId, data) => {
    const response = await apiClient.put(`/quizzes/${slug}/attempts/${attemptId}`, data);
    return response.data || response;
  },
};