import React from 'react';
import { AlertCircle, CheckCircle2, Clock3, RefreshCw } from 'lucide-react';

const SyncStatus = ({ statusData }) => {
  const status = statusData?.status || 'idle';

  if (status === 'running') {
    return (
      <div className="inline-flex items-center gap-2 rounded-full bg-blue-50 px-3 py-1.5 text-sm font-medium text-blue-700">
        <RefreshCw className="h-4 w-4 animate-spin" />
        Sync in progress
      </div>
    );
  }

  if (status === 'completed' || status === 'success') {
    return (
      <div className="inline-flex items-center gap-2 rounded-full bg-green-50 px-3 py-1.5 text-sm font-medium text-green-700">
        <CheckCircle2 className="h-4 w-4" />
        Completed
      </div>
    );
  }

  if (status === 'failed' || status === 'error') {
    return (
      <div className="inline-flex items-center gap-2 rounded-full bg-red-50 px-3 py-1.5 text-sm font-medium text-red-700">
        <AlertCircle className="h-4 w-4" />
        Failed
      </div>
    );
  }

  return (
    <div className="inline-flex items-center gap-2 rounded-full bg-gray-100 px-3 py-1.5 text-sm font-medium text-gray-700">
      <Clock3 className="h-4 w-4" />
      Idle
    </div>
  );
};

export default SyncStatus;
