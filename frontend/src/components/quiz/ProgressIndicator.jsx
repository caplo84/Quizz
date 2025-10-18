

const ProgressIndicator = ({ current, total, percentage, darkMode }) => {
  return (
    <div className="progress-indicator">
      <div className={`progress-bar-container ${darkMode ? 'bg-navy' : 'bg-light-grey'}`}>
        <div 
          className="progress-bar-fill"
          style={{ width: `${percentage}%` }}
        />
      </div>
      <div className="progress-dots">
        {Array.from({ length: Math.min(total, 20) }, (_, i) => {
          const step = Math.ceil(total / 20);
          const isCompleted = i * step < current;
          const isActive = i * step === current - 1;
          
          return (
            <div 
              key={i}
              className={`progress-dot ${isCompleted ? 'completed' : ''} ${isActive ? 'active' : ''}`}
            />
          );
        })}
      </div>
    </div>
  );
};

export default ProgressIndicator;
