import apiClient from './apiClient.js';

export const topicsApi = {
  getAll: () => apiClient.get('/topics'),
  getById: (id) => apiClient.get(`/topics/${id}`),
  getQuizzes: (topic) => apiClient.get(`/topics/${topic}/quizzes`),
};