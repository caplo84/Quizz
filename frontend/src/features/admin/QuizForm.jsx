import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { ArrowLeft as ArrowLeftIcon } from 'lucide-react';
import adminApi from '../../services/adminApi';
import QuizEditor from './components/QuizEditor';

const QuizForm = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const isEditMode = Boolean(id);

  const [loading, setLoading] = useState(false);
  const [submitError, setSubmitError] = useState('');
  const [validationErrors, setValidationErrors] = useState([]);
  const [topics, setTopics] = useState([]);
  const [formData, setFormData] = useState({
    title: '',
    slug: '',
    topic: '',
    description: '',
    questions: [],
  });

  useEffect(() => {
    loadTopics();
    if (isEditMode) {
      loadQuiz();
    }
  }, [id]);

  const loadTopics = async () => {
    try {
      const data = await adminApi.getAllTopics();
      setTopics(data || []);
    } catch (error) {
      console.error('Failed to load topics:', error);
    }
  };

  const loadQuiz = async () => {
    try {
      setLoading(true);
      const data = await adminApi.getQuizById(id);
      setFormData((prev) => ({
        ...prev,
        ...data,
        topic: data.topicId || data.topic_id || data.topic || '',
        questions: data.questions || [],
      }));
    } catch (error) {
      console.error('Failed to load quiz:', error);
      alert('Failed to load quiz. Redirecting to quiz list.');
      navigate('/admin/quizzes');
    } finally {
      setLoading(false);
    }
  };

  const generateSlug = (title) => {
    return title
      .toLowerCase()
      .replace(/[^a-z0-9\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
      .trim();
  };

  const handleTitleChange = (e) => {
    setSubmitError('');
    const newTitle = e.target.value;
    setFormData({
      ...formData,
      title: newTitle,
      slug: generateSlug(newTitle),
    });
  };

  const handleChange = (e) => {
    setSubmitError('');
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const validateForm = () => {
    const errors = [];

    if (!formData.title?.trim()) {
      errors.push('Title is required.');
    }

    if (!formData.slug?.trim()) {
      errors.push('Slug is required.');
    }

    if (!formData.topic) {
      errors.push('Please select a topic.');
    }

    if (!Array.isArray(formData.questions) || formData.questions.length === 0) {
      errors.push('Add at least one question.');
    }

    (formData.questions || []).forEach((question, qIndex) => {
      const questionNumber = qIndex + 1;
      if (!question.question?.trim()) {
        errors.push(`Question ${questionNumber}: question text is required.`);
      }

      const choices = Array.isArray(question.choices) ? question.choices : [];
      if (choices.length === 0 || choices.some((choice) => !String(choice || '').trim())) {
        errors.push(`Question ${questionNumber}: all choices must be filled.`);
      }

      if (!Array.isArray(question.correctAnswers) || question.correctAnswers.length === 0) {
        errors.push(`Question ${questionNumber}: mark at least one correct answer.`);
      }
    });

    return errors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    setSubmitError('');
    const errors = validateForm();
    setValidationErrors(errors);

    if (errors.length > 0) {
      return;
    }

    try {
      setLoading(true);
      if (isEditMode) {
        await adminApi.updateQuiz(id, formData);
      } else {
        await adminApi.createQuiz(formData);
      }
      navigate('/admin/quizzes');
    } catch (error) {
      console.error('Failed to save quiz:', error);
      setSubmitError('Failed to save quiz. Please check the data and try again.');
    } finally {
      setLoading(false);
    }
  };

  if (loading && isEditMode) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <button
          onClick={() => navigate('/admin/quizzes')}
          className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
        >
          <ArrowLeftIcon className="w-6 h-6" />
        </button>
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            {isEditMode ? 'Edit Quiz' : 'Create New Quiz'}
          </h1>
          <p className="text-gray-600 mt-1">
            {isEditMode ? 'Update quiz details and questions' : 'Add a new quiz to your collection'}
          </p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {validationErrors.length > 0 && (
          <div className="bg-amber-50 border border-amber-200 rounded-lg p-4">
            <h3 className="text-sm font-semibold text-amber-800 mb-2">Please fix the following:</h3>
            <ul className="list-disc list-inside text-sm text-amber-700 space-y-1">
              {validationErrors.map((item) => (
                <li key={item}>{item}</li>
              ))}
            </ul>
          </div>
        )}

        {submitError && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-sm text-red-700">
            {submitError}
          </div>
        )}

        {/* Basic Information */}
        <div className="bg-white rounded-lg shadow-md p-6 space-y-4">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Basic Information</h2>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Title <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              name="title"
              value={formData.title}
              onChange={handleTitleChange}
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="e.g., JavaScript Basics"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Slug <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              name="slug"
              value={formData.slug}
              onChange={handleChange}
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-gray-50"
              placeholder="javascript-basics"
            />
            <p className="text-xs text-gray-500 mt-1">Auto-generated from title, but can be edited</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Topic
            </label>
            <select
              name="topic"
              value={formData.topic}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Select a topic</option>
              {topics.map((topic) => (
                <option key={topic.id} value={topic.id}>
                  {topic.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Description
            </label>
            <textarea
              name="description"
              value={formData.description}
              onChange={handleChange}
              rows="3"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Brief description of the quiz..."
            />
          </div>
        </div>

        <QuizEditor
          questions={formData.questions}
          onQuestionsChange={(questions) => setFormData((prev) => ({ ...prev, questions }))}
        />

        {/* Form Actions */}
        <div className="flex gap-3 justify-end">
          <button
            type="button"
            onClick={() => navigate('/admin/quizzes')}
            className="px-6 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={loading}
            className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {loading ? 'Saving...' : (isEditMode ? 'Update Quiz' : 'Create Quiz')}
          </button>
        </div>
      </form>
    </div>
  );
};

export default QuizForm;
