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
    // Map topic slug to topic ID (based on backend data)
    const topicIdMap = {
      'html': 1,
      'css': 2,
      'javascript': 3,
      'accessibility': 4
    };
    
    const topicId = topicIdMap[topicSlug.toLowerCase()];
    if (!topicId) {
      throw new Error(`Unknown topic: ${topicSlug}`);
    }
    
    const response = await apiClient.get(`/topics/${topicSlug}/quizzes?topic_id=${topicId}`);
    return response.data || response;
  },
};