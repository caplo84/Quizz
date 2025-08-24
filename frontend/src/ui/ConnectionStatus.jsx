import { useState, useEffect } from 'react';
import { healthApi } from '../services/api.js';

function ConnectionStatus() {
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const checkConnection = async () => {
      console.log('🔍 Checking connection...'); // DEBUG
      try {
        setIsLoading(true);
        const result = await healthApi.check();
        console.log('✅ Health check response:', result); // DEBUG
        setIsConnected(true);
      } catch (error) {
        console.error('❌ Health check failed:', error); // DEBUG
        setIsConnected(false);
      } finally {
        setIsLoading(false);
      }
    };
    
    checkConnection();
    const interval = setInterval(checkConnection, 30000);
    
    return () => clearInterval(interval);
  }, []);

  console.log('🎯 ConnectionStatus state:', { isConnected, isLoading }); // DEBUG

  if (isLoading) {
    return (
      <div className="connection-status loading">
        Checking connection...
      </div>
    );
  }

  return (
    <div className={`connection-status ${isConnected ? 'success' : 'error'}`}>
      {isConnected ? (
        <>
          <span>✅ Backend connected</span>
          <button onClick={() => setIsConnected(false)}>Hide</button>
        </>
      ) : (
        <>
          <span>❌ Backend disconnected</span>
          <button onClick={() => window.location.reload()}>Retry</button>
        </>
      )}
    </div>
  );
}

export default ConnectionStatus;
