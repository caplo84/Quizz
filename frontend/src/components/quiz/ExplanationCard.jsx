import React from 'react';
import { CheckCircle, XCircle, Lightbulb } from 'lucide-react';

const ExplanationCard = ({ isCorrect, explanation, darkMode }) => {
  return (
    <div className={`explanation-card ${isCorrect ? 'correct' : 'incorrect'} ${darkMode ? 'dark-mode' : ''}`}>
      <div className="explanation-header">
        {isCorrect ? (
          <>
            <CheckCircle size={24} className="text-green" />
            <h3 className="text-xl font-bold text-green">Correct!</h3>
          </>
        ) : (
          <>
            <XCircle size={24} className="text-red" />
            <h3 className="text-xl font-bold text-red">Incorrect</h3>
          </>
        )}
      </div>
      
      {explanation && (
        <div className="explanation-content">
          <div className="flex items-start gap-3">
            <Lightbulb size={20} className={darkMode ? 'text-light-bluish' : 'text-purple'} />
            <p className={`text-base leading-relaxed ${darkMode ? 'text-light-bluish' : 'text-grey-navy'}`}>
              {explanation}
            </p>
          </div>
        </div>
      )}
    </div>
  );
};

export default ExplanationCard;
