import { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { Plus as PlusIcon, Search as MagnifyingGlassIcon } from 'lucide-react';
import adminApi from '../../services/adminApi';
import TopicList from './components/TopicList';
import TopicEditor from './components/TopicEditor';

const TopicManagement = () => {
  const location = useLocation();
  const [topics, setTopics] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [deleteConfirm, setDeleteConfirm] = useState(null);
  const [editingTopic, setEditingTopic] = useState(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadTopics();
  }, []);

  useEffect(() => {
    if (location.pathname.endsWith('/topics/new')) {
      setShowCreateModal(true);
    }
  }, [location.pathname]);

  const loadTopics = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await adminApi.getAllTopics();
      setTopics(data || []);
    } catch (error) {
      console.error('Failed to load topics:', error);
      setError('Failed to load topics. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    try {
      await adminApi.deleteTopic(id);
      setTopics(topics.filter((t) => t.id !== id));
      setDeleteConfirm(null);
    } catch (error) {
      console.error('Failed to delete topic:', error);
      setError('Failed to delete topic. It may have associated quizzes.');
    }
  };

  const handleSaveTopic = async (payload) => {
    if (!payload.name || !payload.slug) {
      setError('Topic name and slug are required.');
      return;
    }

    try {
      setSaving(true);
      setError('');
      if (editingTopic) {
        await adminApi.updateTopic(editingTopic.id, payload);
      } else {
        await adminApi.createTopic(payload);
      }

      setShowCreateModal(false);
      setEditingTopic(null);
      await loadTopics();
    } catch (e) {
      setError(e.message || 'Failed to save topic.');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Topic Management</h1>
          <p className="text-gray-600 mt-1">Manage quiz topics and categories</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          <PlusIcon className="w-5 h-5" />
          Create Topic
        </button>
      </div>

      {/* Search Bar */}
      <div className="bg-white rounded-lg shadow-md p-4">
        <div className="relative">
          <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
          <input
            type="text"
            placeholder="Search topics by name or slug..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
      </div>

      {error ? <div className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-lg p-3">{error}</div> : null}

      <TopicList
        topics={topics}
        searchTerm={searchTerm}
        onEdit={setEditingTopic}
        onDeleteRequest={setDeleteConfirm}
      />

      {/* Create/Edit Modal */}
      {(showCreateModal || editingTopic) && (
        <TopicEditor
          topic={editingTopic}
          loading={saving}
          onSave={handleSaveTopic}
          onClose={() => {
            setShowCreateModal(false);
            setEditingTopic(null);
          }}
        />
      )}

      {/* Delete Confirmation Modal */}
      {deleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 className="text-xl font-bold text-gray-900 mb-4">Confirm Delete</h3>
            <p className="text-gray-600 mb-6">
              Are you sure you want to delete this topic? This may affect associated quizzes.
            </p>
            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setDeleteConfirm(null)}
                className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={() => handleDelete(deleteConfirm)}
                className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default TopicManagement;
