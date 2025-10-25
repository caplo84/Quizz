// API base URL configuration
const EXTERNAL_API_BASE_URL = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1`;
const SAME_ORIGIN_API_BASE_URL = '/api/v1';

// Generic API request function
async function apiRequest(endpoint, options = {}) {
  const requestConfig = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  const baseCandidates = [EXTERNAL_API_BASE_URL, SAME_ORIGIN_API_BASE_URL];
  let lastError = null;

  for (const base of baseCandidates) {
    const url = `${base}${endpoint}`;

    try {
      const response = await fetch(url, requestConfig);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      const result = await response.json();

      // Backend returns data wrapped in {success: true, data: [...]}
      if (result.success && result.data) {
        return result.data;
      }

      return result;
    } catch (error) {
      lastError = error;

      // Retry with next base only for network-level failures
      if (error instanceof TypeError) {
        continue;
      }

      throw error;
    }
  }

  console.error(`💥 API request failed for ${endpoint}:`, lastError);
  throw lastError;
}

// API functions for the quiz application
export const api = {
  // Get all topics
  getTopics: () => apiRequest('/topics'),

  // Get all quizzes (optionally filtered by topic)
  getQuizzes: (topicId = null) => {
    if (topicId) {
      // Use the topics/:topic/quizzes endpoint for specific topics
      return apiRequest(`/topics/${topicId}/quizzes`);
    } else {
      // For getting all quizzes, we'll need to fetch all topics first
      // This is handled in the getQuiz() legacy function
      return apiRequest('/topics'); // This will be processed differently
    }
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
    
    // Handle different possible response structures
    let topicsArray = topics;
    if (topics && topics.data && Array.isArray(topics.data)) {
      topicsArray = topics.data;
    } else if (!Array.isArray(topics)) {
      topicsArray = [];
    }
    
    // Transform topics into the quiz list format expected by the frontend
    // Each topic represents a quiz category that users can select
    const quizItems = [];
    
    if (topicsArray && topicsArray.length > 0) {
      for (const topic of topicsArray) {
        
        // Determine source - if has icon_url, it's original manual data, otherwise GitHub
        const source = topic.icon_url ? 'manual' : 'github';
        
        // For external topics, use a generic programming icon, for original topics use specific icons
        const iconUrl = topic.icon_url || '/icon-js.svg'; // fallback to JS icon for external topics
        
        // Check if this topic has quizzes by trying to fetch them
        let quizCount = 0;
        try {
          const topicQuizzesResponse = await apiRequest(`/topics/${topic.slug}/quizzes`);
          const topicQuizzes = Array.isArray(topicQuizzesResponse)
            ? topicQuizzesResponse
            : Array.isArray(topicQuizzesResponse?.data)
              ? topicQuizzesResponse.data
              : [];

          if (Array.isArray(topicQuizzes) && topicQuizzes.length > 0) {
            const activeQuizzes = topicQuizzes.filter((quiz) => quiz?.is_active !== false);
            quizCount = activeQuizzes.length;
          }
        } catch (error) {
          // Keep existing UX resilient if a topic quiz endpoint intermittently fails.
          if (source === 'manual') {
            quizCount = 1;
          }
        }

        if (quizCount > 0) {
          quizItems.push({
            title: topic.name,
            icon: iconUrl,
            slug: topic.slug || topic.name.toLowerCase().replace(/\s+/g, '-'),
            id: topic.id,
            description: topic.description || `Test your knowledge of ${topic.name}`,
            source: source,
            quizCount: quizCount,
          });
        }
      }
    }
    
    // Return array of quiz items for the home page
    return quizItems;
  } catch (error) {
    console.error('Failed to get quiz data from backend:', error);
    throw new Error(`Failed to load quiz data: ${error.message}`);
  }
}

export { topicsApi } from './topicsApi.js';
export { quizzesApi } from './quizzesApi.js';
export { healthApi } from './healthApi.js';
export { default as apiClient } from './apiClient.js';
