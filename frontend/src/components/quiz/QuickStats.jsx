import React from 'react';
import { CheckCircle2, Circle, Flag } from 'lucide-react';

const QuickStats = ({ answered, remaining, flagged, score, darkMode }) => {
  return (
    <div className={`quick-stats ${darkMode ? 'bg-dark-navy' : 'bg-white'}`}>
      <h3 className={`text-lg font-bold mb-4 ${darkMode ? 'text-white' : 'text-dark-navy'}`}>
        Quick Stats
      </h3>
      
      <div className="stats-list space-y-3">
        <div className="stat-item">
          <div className="flex items-center gap-2">
            <CheckCircle2 size={20} className="text-green" />
            <span className={darkMode ? 'text-light-bluish' : 'text-grey-navy'}>Answered</span>
          </div>
          <span className={`font-bold ${darkMode ? 'text-white' : 'text-dark-navy'}`}>
            {answered}
          </span>
        </div>
        
        <div className="stat-item">
          <div className="flex items-center gap-2">
            <Circle size={20} className={darkMode ? 'text-light-grey' : 'text-grey-navy'} />
            <span className={darkMode ? 'text-light-bluish' : 'text-grey-navy'}>Remaining</span>
          </div>
          <span className={`font-bold ${darkMode ? 'text-white' : 'text-dark-navy'}`}>
            {remaining}
          </span>
        </div>
        
        {flagged > 0 && (
          <div className="stat-item">
            <div className="flex items-center gap-2">
              <Flag size={20} className="text-red" />
              <span className={darkMode ? 'text-light-bluish' : 'text-grey-navy'}>Flagged</span>
            </div>
            <span className={`font-bold ${darkMode ? 'text-white' : 'text-dark-navy'}`}>
              {flagged}
            </span>
          </div>
        )}
        
        {score !== undefined && (
          <div className="stat-item">
            <div className="flex items-center gap-2">
              <span className={darkMode ? 'text-light-bluish' : 'text-grey-navy'}>Current Score</span>
            </div>
            <span className={`font-bold text-purple`}>
              {score}
            </span>
          </div>
        )}
      </div>
    </div>
  );
};

export default QuickStats;
