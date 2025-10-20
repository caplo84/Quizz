import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { 
  FileText as DocumentTextIcon, 
  Folder as FolderIcon, 
  HelpCircle as QuestionMarkCircleIcon,
  RefreshCw as ArrowPathIcon,
  Clock as ClockIcon,
  CheckCircle as CheckCircleIcon,
  AlertCircle as ExclamationCircleIcon
} from 'lucide-react';
import adminApi from '../../services/adminApi';

const AdminDashboard = () => {
  const [stats, setStats] = useState({
    totalQuizzes: 0,
    totalTopics: 0,
    totalQuestions: 0,
    lastSyncDate: null,
  });
  const [syncStatus, setSyncStatus] = useState(null);
  const [loading, setLoading] = useState(true);
  const [recentActivity] = useState([
    { action: 'Quiz Updated', target: 'JavaScript Basics', time: '8 minutes ago' },
    { action: 'Topic Created', target: 'Cloud Fundamentals', time: '1 hour ago' },
    { action: 'Sync Completed', target: 'GitHub Repository', time: '2 hours ago' },
  ]);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      setLoading(true);
      const [statisticsData, syncData] = await Promise.all([
        adminApi.getStatistics(),
        adminApi.getSyncStatus().catch(() => null),
      ]);
      
      setStats(statisticsData);
      setSyncStatus(syncData);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  const StatCard = ({ icon: Icon, label, value, color, link }) => (
    <Link 
      to={link}
      className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
    >
      <div className="flex items-center justify-between">
        <div>
          <p className="text-gray-500 text-sm font-medium mb-1">{label}</p>
          <p className="text-3xl font-bold text-gray-900">{value}</p>
        </div>
        <div className={`p-3 rounded-lg ${color}`}>
          <Icon className="w-8 h-8 text-white" />
        </div>
      </div>
    </Link>
  );

  const SyncStatusBadge = () => {
    if (!syncStatus) {
      return (
        <div className="flex items-center gap-2 text-gray-500">
          <ClockIcon className="w-5 h-5" />
          <span>No sync data</span>
        </div>
      );
    }

    if (syncStatus.status === 'running') {
      return (
        <div className="flex items-center gap-2 text-blue-600">
          <ArrowPathIcon className="w-5 h-5 animate-spin" />
          <span>Syncing...</span>
        </div>
      );
    }

    if (syncStatus.status === 'completed') {
      return (
        <div className="flex items-center gap-2 text-green-600">
          <CheckCircleIcon className="w-5 h-5" />
          <span>Last sync successful</span>
        </div>
      );
    }

    return (
      <div className="flex items-center gap-2 text-red-600">
        <ExclamationCircleIcon className="w-5 h-5" />
        <span>Last sync failed</span>
      </div>
    );
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Dashboard</h1>
        <p className="text-gray-600">Welcome to the Quiz Admin Panel</p>
      </div>

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard
          icon={DocumentTextIcon}
          label="Total Quizzes"
          value={stats.totalQuizzes}
          color="bg-blue-500"
          link="/admin/quizzes"
        />
        <StatCard
          icon={FolderIcon}
          label="Total Topics"
          value={stats.totalTopics}
          color="bg-green-500"
          link="/admin/topics"
        />
        <StatCard
          icon={QuestionMarkCircleIcon}
          label="Total Questions"
          value={stats.totalQuestions}
          color="bg-purple-500"
          link="/admin/quizzes"
        />
        <StatCard
          icon={ArrowPathIcon}
          label="Last Sync"
          value={stats.lastSyncDate ? new Date(stats.lastSyncDate).toLocaleDateString() : 'Never'}
          color="bg-orange-500"
          link="/admin/sync"
        />
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Sync Status */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-4">GitHub Sync Status</h2>
          <div className="space-y-4">
            <SyncStatusBadge />
            {syncStatus?.lastSyncDate && (
              <p className="text-sm text-gray-600">
                Last synced: {new Date(syncStatus.lastSyncDate).toLocaleString()}
              </p>
            )}
            <Link
              to="/admin/sync"
              className="inline-block px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Manage Sync
            </Link>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Quick Actions</h2>
          <div className="space-y-3">
            <Link
              to="/admin/quizzes/new"
              className="block w-full px-4 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors text-center font-medium"
            >
              + Create New Quiz
            </Link>
            <Link
              to="/admin/topics/new"
              className="block w-full px-4 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors text-center font-medium"
            >
              + Create New Topic
            </Link>
            <Link
              to="/admin/quizzes"
              className="block w-full px-4 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors text-center font-medium"
            >
              View All Quizzes
            </Link>
          </div>
        </div>
      </div>

      {/* Recent Activity (Placeholder) */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold text-gray-900 mb-4">Recent Activity</h2>
        <div className="space-y-3">
          {recentActivity.map((entry) => (
            <div
              key={`${entry.action}-${entry.time}`}
              className="flex items-center justify-between border-b border-gray-100 pb-3 last:border-b-0"
            >
              <div>
                <p className="text-sm font-medium text-gray-900">{entry.action}</p>
                <p className="text-xs text-gray-500">{entry.target}</p>
              </div>
              <span className="text-xs text-gray-500">{entry.time}</span>
            </div>
          ))}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-3">Analytics Snapshot</h2>
          <ul className="space-y-2 text-sm text-gray-700">
            <li>• Average questions per quiz: {stats.totalQuizzes ? Math.round(stats.totalQuestions / stats.totalQuizzes) : 0}</li>
            <li>• Topic coverage: {stats.totalTopics} active topics</li>
            <li>• Sync health: {syncStatus?.status || 'unknown'}</li>
          </ul>
        </div>

        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-3">Search & Filtering</h2>
          <p className="text-sm text-gray-700 mb-3">
            Advanced search and filtering are available in quiz and topic management pages.
          </p>
          <div className="flex gap-3">
            <Link to="/admin/quizzes" className="px-3 py-2 bg-gray-100 hover:bg-gray-200 rounded text-sm">
              Open Quiz Search
            </Link>
            <Link to="/admin/topics" className="px-3 py-2 bg-gray-100 hover:bg-gray-200 rounded text-sm">
              Open Topic Search
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AdminDashboard;
