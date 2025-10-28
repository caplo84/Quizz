import { useState } from 'react';
import adminApi from '../../services/adminApi';
import BulkActions from './components/BulkActions';

const BulkOperations = () => {
  const [running, setRunning] = useState({
    sync: false,
    images: false,
    correction: false,
  });
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');

  const runOperation = async (key, fn) => {
    try {
      setError('');
      setResult(null);
      setRunning((prev) => ({ ...prev, [key]: true }));
      const response = await fn();
      setResult({ operation: key, response, time: new Date().toISOString() });
    } catch (e) {
      setError(e?.message || 'Operation failed.');
    } finally {
      setRunning((prev) => ({ ...prev, [key]: false }));
    }
  };

  const handleRunSync = () => runOperation('sync', () => adminApi.triggerGitHubSync());

  const handleDownloadImages = () =>
    runOperation('images', () => adminApi.downloadAllTopicImages());

  const handleCorrectQuestions = () =>
    runOperation('correction', () =>
      adminApi.correctQuestions({
        dryRun: false,
      }),
    );

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Bulk Operations</h1>
        <p className="mt-1 text-gray-600">
          Run maintenance and bulk tasks for admin-managed quiz data.
        </p>
      </div>

      {error ? (
        <div className="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div>
      ) : null}

      <BulkActions
        running={running}
        onRunSync={handleRunSync}
        onDownloadImages={handleDownloadImages}
        onCorrectQuestions={handleCorrectQuestions}
      />

      {result ? (
        <div className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
          <h2 className="text-lg font-semibold text-gray-900">Last Operation Result</h2>
          <p className="mt-1 text-sm text-gray-600">
            <span className="font-medium">Operation:</span> {result.operation}
          </p>
          <p className="text-sm text-gray-600">
            <span className="font-medium">Time:</span> {new Date(result.time).toLocaleString()}
          </p>
          <pre className="mt-3 max-h-64 overflow-auto rounded-lg bg-gray-900 p-3 text-xs text-gray-100">
            {JSON.stringify(result.response, null, 2)}
          </pre>
        </div>
      ) : null}
    </div>
  );
};

export default BulkOperations;