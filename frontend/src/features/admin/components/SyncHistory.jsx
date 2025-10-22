import React from 'react';

const statusClass = (status) => {
  if (status === 'completed' || status === 'success') return 'bg-green-50 text-green-700';
  if (status === 'running') return 'bg-blue-50 text-blue-700';
  if (status === 'failed' || status === 'error') return 'bg-red-50 text-red-700';
  return 'bg-gray-100 text-gray-700';
};

const SyncHistory = ({ history }) => {
  if (!history?.length) {
    return <p className="text-sm text-gray-600">No sync history yet.</p>;
  }

  return (
    <div className="space-y-3">
      {history.map((item) => (
        <div key={item.id || item.startTime} className="rounded-lg border border-gray-200 p-3">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <span className={`rounded-full px-2 py-1 text-xs font-medium ${statusClass(item.status)}`}>
              {item.status || 'unknown'}
            </span>
            <span className="text-xs text-gray-500">
              {item.startTime ? new Date(item.startTime).toLocaleString() : '-'}
            </span>
          </div>
          {item.message ? <p className="mt-2 text-sm text-gray-700">{item.message}</p> : null}
          {item.endTime ? (
            <p className="mt-1 text-xs text-gray-500">
              Ended: {new Date(item.endTime).toLocaleString()}
            </p>
          ) : null}
        </div>
      ))}
    </div>
  );
};

export default SyncHistory;
