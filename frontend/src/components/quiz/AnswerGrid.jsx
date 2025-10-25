
import { Check, X } from 'lucide-react';
import CodeBlock from '../../ui/CodeBlock';

const AnswerOption = ({ 
  id, 
  code, 
  text, 
  codeLanguage,
  isSelected, 
  isCorrect, 
  showCorrect, 
  onClick,
  darkMode 
}) => {
  const getOptionClass = () => {
    let baseClass = 'answer-option';
    
    if (!showCorrect) {
      return isSelected ? `${baseClass} selected` : baseClass;
    }
    if (isCorrect) return `${baseClass} correct`;
    if (isSelected && !isCorrect) return `${baseClass} incorrect`;
    return `${baseClass} disabled`;
  };

  const getOptionLabelClass = () => {
    if (showCorrect) {
      if (isCorrect) return 'bg-green text-white';
      if (isSelected && !isCorrect) return 'bg-red text-white';
    }
    if (isSelected) return 'bg-purple text-white';
    return darkMode ? 'bg-navy text-light-bluish' : 'bg-light-grey text-grey-navy';
  };

  return (
    <div 
      className={`${getOptionClass()} ${darkMode ? 'dark-mode' : ''}`}
      onClick={() => !showCorrect && onClick()}
    >
      <div className="answer-option-header">
        <div className={`option-label ${getOptionLabelClass()}`}>
          {String(id || '').toUpperCase()}
        </div>
        {showCorrect && isCorrect && (
          <div className="result-icon correct-icon">
            <Check size={20} />
          </div>
        )}
        {showCorrect && isSelected && !isCorrect && (
          <div className="result-icon incorrect-icon">
            <X size={20} />
          </div>
        )}
      </div>
      
      <div className="answer-content">
        {code ? (
          <CodeBlock 
            code={code} 
            language={codeLanguage || 'css'}
            className="answer-code-block"
          />
        ) : (
          <p className={`text-lg ${darkMode ? 'text-light-bluish' : 'text-dark-navy'}`}>
            {text}
          </p>
        )}
      </div>
    </div>
  );
};

const AnswerGrid = ({ 
  options, 
  selectedAnswer, 
  onAnswerSelect, 
  showCorrect, 
  correctAnswer,
  hasCode,
  codeLanguage,
  darkMode 
}) => {
  const getOptionLabel = (index) => {
    // a-z, then fallback to numeric labels (27, 28, ...)
    if (index < 26) {
      return String.fromCharCode(97 + index);
    }

    return String(index + 1);
  };
  
  return (
    <div className="answer-grid">
      {options.map((option, index) => {
        const optionId = getOptionLabel(index);
        return (
          <AnswerOption
            key={`${optionId}-${index}`}
            id={optionId}
            code={hasCode ? option : null}
            text={!hasCode ? option : null}
            codeLanguage={codeLanguage}
            isSelected={selectedAnswer === option}
            isCorrect={option === correctAnswer}
            showCorrect={showCorrect}
            onClick={() => onAnswerSelect(option)}
            darkMode={darkMode}
          />
        );
      })}
    </div>
  );
};

export default AnswerGrid;
