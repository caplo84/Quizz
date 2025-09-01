// API base URL configuration
const API_BASE_URL = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1`;

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
    console.log('🌐 Making API request to:', url);
    const response = await fetch(url, config);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
    }
    
    const result = await response.json();
    console.log('📡 API response:', result);
    
    // Backend returns data wrapped in {success: true, data: [...]}
    if (result.success && result.data) {
      return result.data;
    }
    
    // Fallback if structure is different
    return result;
  } catch (error) {
    console.error(`💥 API request failed for ${endpoint}:`, error);
    throw error;
  }
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
    console.log('🔄 Fetching data from backend API...');
    console.log('🔗 API_BASE_URL:', API_BASE_URL);
    
    const topics = await api.getTopics();
    console.log('📂 Raw topics response:', topics);
    console.log('📂 Topics type:', typeof topics, 'Array?', Array.isArray(topics));
    
    // Handle different possible response structures
    let topicsArray = topics;
    if (topics && topics.data && Array.isArray(topics.data)) {
      topicsArray = topics.data;
    } else if (!Array.isArray(topics)) {
      console.warn('⚠️ Topics is not an array:', topics);
      topicsArray = [];
    }
    
    // Transform topics into the quiz list format expected by the frontend
    // Each topic represents a quiz category that users can select
    const quizItems = [];
    
    if (topicsArray && topicsArray.length > 0) {
      for (const topic of topicsArray) {
        console.log('🏷️ Processing topic:', topic);
        // Each topic becomes a selectable quiz item on the home page
        quizItems.push({
          title: topic.name,  // e.g., "HTML", "CSS", "JavaScript"
          icon: topic.icon_url || `/icon-${topic.name.toLowerCase()}.svg`,  // Use backend icon_url or fallback
          slug: topic.slug || topic.name.toLowerCase(),
          id: topic.id,
          description: topic.description || `Test your knowledge of ${topic.name}`
        });
      }
    } else {
      console.warn('⚠️ No topics found or topics array is empty');
    }
    
    console.log('🎯 Final quiz items for home page:', quizItems);
    
    // Return array of quiz items for the home page
    return quizItems;
  } catch (error) {
    console.error('💥 Failed to get quiz data from backend:', error);
    console.error('💥 Error details:', error.stack);
    throw new Error(`Failed to load quiz data: ${error.message}`);
  }
}

export { topicsApi } from './topicsApi.js';
export { quizzesApi } from './quizzesApi.js';
export { healthApi } from './healthApi.js';
export { default as apiClient } from './apiClient.js';
