import apiClient from './apiClient.js';

export const topicsApi = {
  getAll: async () => {
    const response = await apiClient.get('/topics');
    return response.data || response;
  },
  getById: async (id) => {
    const response = await apiClient.get(`/topics/${id}`);
    return response.data || response;
  },
  getQuizzes: async (topicSlug) => {
    // Use the backend endpoint that accepts topic slugs directly
    const response = await apiClient.get(`/topics/${topicSlug}/quizzes`);
    const quizzes = response.data || response;
    if (!Array.isArray(quizzes)) return [];

    // Hide unpublished/inactive quizzes from public flows
    return quizzes.filter((quiz) => quiz?.is_active !== false);
  },
};