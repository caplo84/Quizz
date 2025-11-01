import apiClient from './apiClient';
import config from '../config/api';

/**
 * Admin API Service
 * Handles all admin-related API operations
 */

const adminApi = {
  normalizeApiResponse: (response) => response?.data ?? response,

  toTopicPayload: (topicData) => ({
    name: topicData.name,
    slug: topicData.slug,
    description: topicData.description || '',
    icon_url: topicData.icon_url || topicData.icon || '',
  }),

  toQuizPayload: (quizData) => {
    const questions = (quizData.questions || []).map((question, qIndex) => ({
      question_text: question.question,
      question_type: 'multiple_choice',
      points: 1,
      order_index: qIndex + 1,
      choices: (question.choices || []).map((choiceText, cIndex) => ({
        choice_text: choiceText,
        is_correct: (question.correctAnswers || []).includes(cIndex),
        order_index: cIndex + 1,
      })),
    }));

    return {
      title: quizData.title,
      slug: quizData.slug,
      description: quizData.description || '',
      topic_id: Number(quizData.topic),
      difficulty_level: quizData.difficulty_level || 'medium',
      time_limit_minutes: Number(quizData.time_limit_minutes || 30),
      total_questions: questions.length,
      is_active: quizData.is_active !== undefined ? Boolean(quizData.is_active) : true,
      questions,
    };
  },

  normalizeTopic: (topic) => ({
    ...topic,
    description: topic.description || '',
    icon: topic.icon || topic.icon_url || '',
  }),

  normalizeQuiz: (quiz, topic = null) => ({
    ...quiz,
    topic: quiz.topic?.name || quiz.topic_name || topic?.name || quiz.topic || '',
    topicId: quiz.topic_id || topic?.id || quiz.topicId,
    topicSlug: quiz.topic?.slug || quiz.topic_slug || topic?.slug || quiz.topicSlug || '',
    totalQuestions: quiz.total_questions || quiz.totalQuestions || quiz.questions?.length || 0,
    isActive: quiz.is_active !== undefined ? Boolean(quiz.is_active) : true,
  }),

  // ==================== Quiz Operations ====================
  
  /**
   * Get all quizzes
   */
  getAllQuizzes: async () => {
    try {
      const topicsRaw = await apiClient.get('/topics');
      const topics = (adminApi.normalizeApiResponse(topicsRaw) || []).map(adminApi.normalizeTopic);

      const quizzesByTopic = await Promise.all(
        topics.map(async (topic) => {
          const quizzesRaw = await apiClient.get(`/topics/${topic.slug}/quizzes`);
          const quizzes = adminApi.normalizeApiResponse(quizzesRaw) || [];
          return quizzes.map((quiz) => adminApi.normalizeQuiz(quiz, topic));
        }),
      );

      return quizzesByTopic.flat();
    } catch (error) {
      console.error('Failed to fetch quizzes:', error);
      throw error;
    }
  },

  /**
   * Get single quiz by ID
   */
  getQuizById: async (id) => {
    try {
      const quizRaw = await apiClient.get(`/admin/quizzes/${id}`);
      const quiz = adminApi.normalizeApiResponse(quizRaw);

      if (!quiz) {
        throw new Error('Quiz not found');
      }

      const mappedQuestions = (quiz.questions || []).map((question) => {
        const sortedChoices = [...(question.choices || [])].sort(
          (a, b) => (a.order_index || 0) - (b.order_index || 0),
        );

        const correctAnswers = sortedChoices
          .map((choice, index) => (choice.is_correct ? index : -1))
          .filter((index) => index !== -1);

        return {
          id: question.id,
          question: question.question_text || '',
          choices: sortedChoices.map((choice) => choice.choice_text || ''),
          correctAnswers,
        };
      });

      return {
        id: quiz.id,
        title: quiz.title || '',
        slug: quiz.slug || '',
        description: quiz.description || '',
        topic: quiz.topic_id || quiz.topic?.id || '',
        difficulty_level: quiz.difficulty_level || 'medium',
        time_limit_minutes: quiz.time_limit_minutes || 30,
        is_active: quiz.is_active !== undefined ? Boolean(quiz.is_active) : true,
        questions: mappedQuestions,
      };
    } catch (error) {
      console.error(`Failed to fetch quiz ${id}:`, error);
      throw error;
    }
  },

  /**
   * Create new quiz
   */
  createQuiz: async (quizData) => {
    try {
      return await apiClient.post('/admin/quizzes', adminApi.toQuizPayload(quizData));
    } catch (error) {
      console.error('Failed to create quiz:', error);
      throw error;
    }
  },

  /**
   * Update existing quiz
   */
  updateQuiz: async (id, quizData) => {
    try {
      return await apiClient.put(`/admin/quizzes/${id}`, adminApi.toQuizPayload(quizData));
    } catch (error) {
      console.error(`Failed to update quiz ${id}:`, error);
      throw error;
    }
  },

  /**
   * Publish/Unpublish quiz by toggling is_active
   */
  setQuizPublicationStatus: async (id, shouldPublish) => {
    try {
      const existingQuiz = await adminApi.getQuizById(id);
      const payload = {
        ...existingQuiz,
        is_active: Boolean(shouldPublish),
      };

      return await adminApi.updateQuiz(id, payload);
    } catch (error) {
      console.error(`Failed to set publication status for quiz ${id}:`, error);
      throw error;
    }
  },

  /**
   * Delete quiz
   */
  deleteQuiz: async (id) => {
    try {
      return await apiClient.delete(`/admin/quizzes/${id}`);
    } catch (error) {
      console.error(`Failed to delete quiz ${id}:`, error);
      throw error;
    }
  },

  // ==================== Topic Operations ====================

  /**
   * Get all topics
   */
  getAllTopics: async () => {
    try {
      const response = await apiClient.get('/topics');
      return (adminApi.normalizeApiResponse(response) || []).map(adminApi.normalizeTopic);
    } catch (error) {
      console.error('Failed to fetch topics:', error);
      throw error;
    }
  },

  /**
   * Get single topic by ID
   */
  getTopicById: async (id) => {
    try {
      const topics = await adminApi.getAllTopics();
      const topic = topics.find((item) => String(item.id) === String(id));
      if (!topic) {
        throw new Error('Topic not found');
      }
      return topic;
    } catch (error) {
      console.error(`Failed to fetch topic ${id}:`, error);
      throw error;
    }
  },

  /**
   * Create new topic
   */
  createTopic: async (topicData) => {
    try {
      return await apiClient.post('/admin/topics', adminApi.toTopicPayload(topicData));
    } catch (error) {
      console.error('Failed to create topic:', error);
      throw error;
    }
  },

  /**
   * Update existing topic
   */
  updateTopic: async (id, topicData) => {
    try {
      return await apiClient.put(`/admin/topics/${id}`, adminApi.toTopicPayload(topicData));
    } catch (error) {
      console.error(`Failed to update topic ${id}:`, error);
      throw error;
    }
  },

  /**
   * Delete topic
   */
  deleteTopic: async (id) => {
    try {
      return await apiClient.delete(`/admin/topics/${id}`);
    } catch (error) {
      console.error(`Failed to delete topic ${id}:`, error);
      throw error;
    }
  },

  // ==================== Sync Operations ====================

  /**
   * Trigger GitHub sync
   */
  triggerGitHubSync: async () => {
    try {
      const response = await fetch(`${config.BASE_URL}/api/admin/sync/github`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });

      if (!response.ok) {
        throw new Error(`Failed to trigger sync: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to trigger GitHub sync:', error);
      throw error;
    }
  },

  /**
   * Get sync status
   */
  getSyncStatus: async () => {
    try {
      const response = await fetch(`${config.BASE_URL}/api/admin/sync/github/status`, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch sync status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to fetch sync status:', error);
      throw error;
    }
  },

  triggerSync: async () => adminApi.triggerGitHubSync(),

  downloadAllTopicImages: async () => {
    const response = await fetch(`${config.BASE_URL}/api/admin/download-all-topic-images`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    });

    if (!response.ok) {
      throw new Error(`Failed to download topic images: ${response.status}`);
    }

    return response.json();
  },

  correctQuestions: async (payload) => {
    const response = await fetch(`${config.BASE_URL}/api/admin/questions/correct`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload || {}),
    });

    if (!response.ok) {
      throw new Error(`Failed to correct questions: ${response.status}`);
    }

    return response.json();
  },

  reviewQuizBeforePublish: async (quizSlug) => {
    if (!quizSlug) {
      throw new Error('Quiz slug is required for pre-publish review.');
    }

    return adminApi.correctQuestions({
      quiz_slug: quizSlug,
      dry_run: true,
      review_only: true,
      batch_size: 100,
      confidence_threshold: 0.85,
    });
  },

  // ==================== AI Settings ====================

  getAISettings: async () => {
    const response = await fetch(`${config.BASE_URL}/api/admin/ai/settings`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch AI settings: ${response.status}`);
    }

    return response.json();
  },

  updateAISettings: async (payload) => {
    const response = await fetch(`${config.BASE_URL}/api/admin/ai/settings`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload || {}),
    });

    if (!response.ok) {
      throw new Error(`Failed to update AI settings: ${response.status}`);
    }

    return response.json();
  },

  // ==================== Statistics ====================

  /**
   * Get dashboard statistics
   */
  getStatistics: async () => {
    try {
      const [quizzes, topics] = await Promise.all([
        adminApi.getAllQuizzes(),
        adminApi.getAllTopics(),
      ]);

      const totalQuestions = quizzes.reduce(
        (sum, quiz) => sum + (quiz.totalQuestions || quiz.questions?.length || 0),
        0,
      );

      return {
        totalQuizzes: quizzes.length,
        totalTopics: topics.length,
        totalQuestions,
        lastSyncDate: null,
      };
    } catch (error) {
      console.error('Failed to fetch statistics:', error);
      return {
        totalQuizzes: 0,
        totalTopics: 0,
        totalQuestions: 0,
        lastSyncDate: null,
      };
    }
  },
};

export default adminApi;
