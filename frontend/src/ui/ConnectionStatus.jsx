import { useState, useEffect } from 'react';
import { api } from '../services/api';

function ConnectionStatus() {
  const [status, setStatus] = useState('checking');
  const [error, setError] = useState(null);

  useEffect(() => {
    checkConnection();
  }, []);

  const checkConnection = async () => {
    try {
      setStatus('checking');
      setError(null);
      
      // Test the health endpoint
      await api.healthCheck();
      setStatus('connected');
    } catch (err) {
      setStatus('disconnected');
      setError(err.message);
    }
  };

  const getStatusColor = () => {
    switch (status) {
      case 'checking': return 'text-yellow-600';
      case 'connected': return 'text-green-600';
      case 'disconnected': return 'text-red-600';
      default: return 'text-gray-600';
    }
  };

  const getStatusText = () => {
    switch (status) {
      case 'checking': return 'Checking connection...';
      case 'connected': return 'Backend connected ✅';
      case 'disconnected': return 'Backend disconnected ❌';
      default: return 'Unknown status';
    }
  };

  if (import.meta.env.PROD) {
    return null; // Don't show in production
  }

  return (
    <div className="fixed bottom-4 right-4 bg-white shadow-lg rounded-lg p-3 border">
      <div className={`text-sm font-medium ${getStatusColor()}`}>
        {getStatusText()}
      </div>
      {error && (
        <div className="text-xs text-red-500 mt-1">
          {error}
        </div>
      )}
      <button
        onClick={checkConnection}
        className="text-xs text-blue-600 hover:text-blue-800 mt-1 underline"
        disabled={status === 'checking'}
      >
        Retry
      </button>
    </div>
  );
}

export default ConnectionStatus;
