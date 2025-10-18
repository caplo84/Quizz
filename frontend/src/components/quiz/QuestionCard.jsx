
import { Code2 } from 'lucide-react';
import CodeBlock from '../../ui/CodeBlock';
import ImageDisplay from '../../ui/ImageDisplay';

const QuestionCard = ({ 
  question, 
  code, 
  codeLanguage,
  difficulty, 
  tags, 
  imageUrl, 
  imageAlt,
  topic,
  darkMode 
}) => {
  const difficultyColors = {
    easy: darkMode ? 'bg-green-900 text-green-200' : 'bg-green-100 text-green-800',
    medium: darkMode ? 'bg-yellow-900 text-yellow-200' : 'bg-yellow-100 text-yellow-800',
    hard: darkMode ? 'bg-red-900 text-red-200' : 'bg-red-100 text-red-800'
  };

  return (
    <div className="question-card">
      {/* Question Header */}
      {(difficulty || tags) && (
        <div className="question-header mb-6">
          <div className="question-meta flex gap-3 flex-wrap">
            {difficulty && (
              <span className={`difficulty-badge ${difficultyColors[difficulty]}`}>
                {difficulty}
              </span>
            )}
            {tags?.map(tag => (
              <span 
                key={tag} 
                className={`tag ${darkMode ? 'bg-navy text-light-bluish' : 'bg-light-grey text-grey-navy'}`}
              >
                #{tag}
              </span>
            ))}
          </div>
        </div>
      )}

      {/* Question Text */}
      <div className="question-text">
        <h2 className={`text-3xl md:text-4xl font-medium leading-tight mb-6 ${
          darkMode ? 'text-white' : 'text-dark-navy'
        }`}>
          {question}
        </h2>
      </div>

      {/* Code Block (if exists) */}
      {code && (
        <div className="code-block mb-6">
          <div className={`code-header flex items-center gap-2 px-4 py-3 rounded-t-lg ${
            darkMode ? 'bg-grey-navy text-light-bluish' : 'bg-light-grey text-grey-navy'
          }`}>
            <Code2 size={16} />
            <span className="text-sm font-semibold">Code Example</span>
          </div>
          <CodeBlock 
            code={code} 
            language={codeLanguage}
            className="rounded-t-none"
          />
        </div>
      )}

      {/* Image (if exists) */}
      {imageUrl && (
        <div className="question-image mb-6">
          <ImageDisplay 
            imageUrl={imageUrl} 
            altText={imageAlt}
            topic={topic}
            className="rounded-lg shadow-md"
          />
        </div>
      )}
    </div>
  );
};

export default QuestionCard;
