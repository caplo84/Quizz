// API base URL configuration
const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

// Generic API request function
async function apiRequest(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error(`API request failed for ${endpoint}:`, error);
    throw error;
  }
}

// API functions for the quiz application
export const api = {
  // Get all topics
  getTopics: () => apiRequest('/topics'),

  // Get all quizzes (optionally filtered by topic)
  getQuizzes: (topicId = null) => {
    const query = topicId ? `?topic_id=${topicId}` : '';
    return apiRequest(`/quizzes${query}`);
  },

  // Get a specific quiz by slug
  getQuiz: (slug) => apiRequest(`/quizzes/${slug}`),

  // Get questions for a specific quiz
  getQuizQuestions: (slug) => apiRequest(`/quizzes/${slug}/questions`),

  // Create a new quiz attempt
  createAttempt: (slug) => apiRequest(`/quizzes/${slug}/attempts`, {
    method: 'POST',
  }),

  // Submit answers for a quiz attempt
  submitAttempt: (slug, attemptId, answers) => apiRequest(`/quizzes/${slug}/attempts/${attemptId}`, {
    method: 'PUT',
    body: JSON.stringify({ answers }),
  }),

  // Get attempt results
  getAttemptResults: (slug, attemptId) => apiRequest(`/quizzes/${slug}/attempts/${attemptId}`),

  // Health check
  healthCheck: () => apiRequest('/health', { 
    headers: { 'Content-Type': 'application/json' } 
  }),
};

// Legacy function for backward compatibility
export async function getQuiz() {
  try {
    const topics = await api.getTopics();
    const quizzes = await api.getQuizzes();
    
    // Transform the data to match the expected format
    return {
      topics: topics.data || topics,
      quizzes: quizzes.data || quizzes,
    };
  } catch (error) {
    console.error('Failed to get quiz data:', error);
    // Fallback to local data if API is not available
    try {
      const res = await fetch("/data.json");
      if (!res.ok) throw Error("Failed in getting quiz");
      const data = await res.json();
      return data?.quizzes;
    } catch (fallbackError) {
      throw new Error("Both API and local data failed to load");
    }
  }
}

export default api;
