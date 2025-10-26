
import { Code2 } from 'lucide-react';
import CodeBlock from '../../ui/CodeBlock';
import ImageDisplay from '../../ui/ImageDisplay';

const QuestionCard = ({ 
  question, 
  code, 
  codeLanguage,
  imageUrl, 
  imageAlt,
  topic,
  darkMode 
}) => {
  return (
    <div className="question-card">
      {/* Question Text */}
      <div className="question-text">
        <h2 className={`text-xl md:text-2xl font-semibold leading-snug mb-5 ${
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
