import { useState, useEffect } from 'react';
import { Wifi, WifiOff } from 'lucide-react';
import { healthApi } from '../services/api.js';

function ConnectionStatus() {
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const checkConnection = async () => {
      try {
        setIsLoading(true);
        await healthApi.check();
        setIsConnected(true);
      } catch (error) {
        setIsConnected(false);
      } finally {
        setIsLoading(false);
      }
    };
    
    checkConnection();
    const interval = setInterval(checkConnection, 30000);
    
    return () => clearInterval(interval);
  }, []);

  if (isLoading) return null;

  return (
    <div className={`flex items-center gap-1.5 text-xs px-2.5 py-1 rounded-full font-medium transition-all ${
      isConnected
        ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300'
        : 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
    }`}>
      {isConnected ? (
        <>
          <Wifi className="w-3 h-3" /> Online
        </>
      ) : (
        <>
          <WifiOff className="w-3 h-3" /> Offline
        </>
      )}
    </div>
  );
}

export default ConnectionStatus;
