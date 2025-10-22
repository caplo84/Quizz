import React, { useState, useEffect } from 'react';
import { 
  RefreshCw as ArrowPathIcon, 
  Play as PlayIcon 
} from 'lucide-react';
import adminApi from '../../services/adminApi';
import SyncStatus from './components/SyncStatus';
import SyncHistory from './components/SyncHistory';

const SyncManagement = () => {
  const [syncStatus, setSyncStatus] = useState(null);
  const [loading, setLoading] = useState(true);
  const [syncing, setSyncing] = useState(false);
  const [syncHistory, setSyncHistory] = useState([]);

  useEffect(() => {
    loadSyncStatus();
    
    // Poll sync status every 5 seconds if syncing
    const interval = setInterval(() => {
      if (syncing) {
        loadSyncStatus();
      }
    }, 5000);

    return () => clearInterval(interval);
  }, [syncing]);

  const loadSyncStatus = async () => {
    try {
      const data = await adminApi.getSyncStatus();
      setSyncStatus(data);
      
      // Check if sync is running
      if (data?.status === 'running') {
        setSyncing(true);
      } else {
        setSyncing(false);

        if (data?.status) {
          setSyncHistory((prev) => {
            const latest = prev[0];
            if (latest && latest.status === 'running') {
              return [
                {
                  ...latest,
                  status: data.status,
                  message: data.message,
                  endTime: new Date().toISOString(),
                },
                ...prev.slice(1),
              ];
            }
            return prev;
          });
        }
      }
    } catch (error) {
      console.error('Failed to load sync status:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleTriggerSync = async () => {
    if (syncing) {
      alert('A sync operation is already in progress.');
      return;
    }

    const confirmed = window.confirm(
      'This will sync all quizzes from the GitHub repository. Continue?'
    );
    
    if (!confirmed) return;

    try {
      setSyncing(true);
      await adminApi.triggerGitHubSync();
      
      // Add to history
      setSyncHistory((prev) => [
        {
          id: Date.now(),
          status: 'running',
          startTime: new Date().toISOString(),
          message: 'Sync started by admin',
        },
        ...prev,
      ]);
      
      // Reload status after a delay
      setTimeout(() => {
        loadSyncStatus();
      }, 2000);
    } catch (error) {
      console.error('Failed to trigger sync:', error);
      setSyncHistory((prev) => [
        {
          id: Date.now(),
          status: 'failed',
          startTime: new Date().toISOString(),
          endTime: new Date().toISOString(),
          message: error.message || 'Failed to start sync',
        },
        ...prev,
      ]);
      setSyncing(false);
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
      <div>
        <h1 className="text-3xl font-bold text-gray-900">GitHub Sync Management</h1>
        <p className="text-gray-600 mt-1">Sync quizzes from the LinkedIn Skill Assessments repository</p>
      </div>

      {/* Current Status */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold text-gray-900 mb-4">Current Status</h2>
        
        <div className="space-y-4">
          <SyncStatus statusData={syncStatus} />
          
          {syncStatus?.lastSyncDate && (
            <div className="text-sm text-gray-600">
              <span className="font-medium">Last synced:</span>{' '}
              {new Date(syncStatus.lastSyncDate).toLocaleString()}
            </div>
          )}

          {syncStatus?.message && (
            <div className="p-4 bg-gray-50 rounded-lg text-sm text-gray-700">
              {syncStatus.message}
            </div>
          )}

          {syncStatus?.stats && (
            <div className="grid grid-cols-3 gap-4 pt-4 border-t border-gray-200">
              <div>
                <p className="text-sm text-gray-500">Quizzes Synced</p>
                <p className="text-2xl font-bold text-gray-900">{syncStatus.stats.quizzesSynced || 0}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Topics Added</p>
                <p className="text-2xl font-bold text-gray-900">{syncStatus.stats.topicsAdded || 0}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Errors</p>
                <p className="text-2xl font-bold text-gray-900">{syncStatus.stats.errors || 0}</p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Actions */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold text-gray-900 mb-4">Actions</h2>
        
        <button
          onClick={handleTriggerSync}
          disabled={syncing}
          className="flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
        >
          {syncing ? (
            <>
              <ArrowPathIcon className="w-5 h-5 animate-spin" />
              Syncing...
            </>
          ) : (
            <>
              <PlayIcon className="w-5 h-5" />
              Trigger GitHub Sync
            </>
          )}
        </button>

        <div className="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
          <p className="text-sm text-yellow-800">
            <strong>Note:</strong> This operation will fetch all quiz data from the LinkedIn Skill Assessments 
            GitHub repository and update the database. It may take several minutes to complete.
          </p>
        </div>
      </div>

      {/* Sync Information */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold text-gray-900 mb-4">About GitHub Sync</h2>
        
        <div className="space-y-4 text-sm text-gray-700">
          <div>
            <h3 className="font-semibold text-gray-900 mb-2">What does sync do?</h3>
            <ul className="list-disc list-inside space-y-1 ml-2">
              <li>Fetches the latest quizzes from the LinkedIn Skill Assessments repository</li>
              <li>Parses markdown files and extracts quiz questions</li>
              <li>Creates or updates topics and quizzes in the database</li>
              <li>Preserves existing custom quizzes and data</li>
            </ul>
          </div>

          <div>
            <h3 className="font-semibold text-gray-900 mb-2">Repository Information</h3>
            <p>
              Source:{' '}
              <a
                href="https://github.com/Ebazhanov/linkedin-skill-assessments-quizzes"
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 hover:underline"
              >
                linkedin-skill-assessments-quizzes
              </a>
            </p>
          </div>

          <div>
            <h3 className="font-semibold text-gray-900 mb-2">Sync Frequency</h3>
            <p>
              You can trigger manual syncs anytime. For production, consider setting up automated 
              scheduled syncs (e.g., daily or weekly) to keep quizzes up to date.
            </p>
          </div>
        </div>
      </div>

      {/* Sync History (Placeholder) */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold text-gray-900 mb-4">Sync History</h2>
        <SyncHistory history={syncHistory.slice(0, 10)} />
      </div>
    </div>
  );
};

export default SyncManagement;
