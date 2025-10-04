import { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';

function CodeBlock({ code, language, className = '' }) {
  const { darkMode } = useSelector((state) => state.home);
  const [highlightedCode, setHighlightedCode] = useState(code);

  useEffect(() => {
    // For now, we'll display plain code with formatting
    // Later we can integrate Prism.js for syntax highlighting
    setHighlightedCode(code);
  }, [code, language]);

  if (!code || code.trim() === '') {
    return null;
  }

  // Language display mapping
  const languageDisplay = {
    javascript: 'JavaScript',
    js: 'JavaScript',
    typescript: 'TypeScript',
    ts: 'TypeScript',
    css: 'CSS',
    html: 'HTML',
    python: 'Python',
    py: 'Python',
    java: 'Java',
    go: 'Go',
    cpp: 'C++',
    c: 'C',
    sql: 'SQL',
    json: 'JSON',
    xml: 'XML',
    bash: 'Bash',
    shell: 'Shell',
    yaml: 'YAML',
    yml: 'YAML',
    php: 'PHP',
    ruby: 'Ruby',
    swift: 'Swift',
    kotlin: 'Kotlin',
    rust: 'Rust',
    vue: 'Vue',
    react: 'React'
  };

  const displayLanguage = languageDisplay[language?.toLowerCase()] || language || 'Code';

  return (
    <div className={`relative rounded-lg overflow-hidden border ${className} ${
      darkMode 
        ? 'bg-gray-900 border-gray-700' 
        : 'bg-gray-50 border-gray-200'
    }`}>
      {/* Language label */}
      {language && (
        <div className={`px-4 py-2 text-sm font-medium border-b ${
          darkMode 
            ? 'bg-gray-800 border-gray-700 text-gray-300' 
            : 'bg-gray-100 border-gray-200 text-gray-600'
        }`}>
          {displayLanguage}
        </div>
      )}
      
      {/* Code content */}
      <div className={`p-4 overflow-x-auto ${
        darkMode ? 'bg-gray-900' : 'bg-gray-50'
      }`}>
        <pre className={`text-sm leading-relaxed ${
          darkMode ? 'text-gray-100' : 'text-gray-800'
        }`}>
          <code className="font-mono whitespace-pre-wrap break-words">
            {highlightedCode}
          </code>
        </pre>
      </div>
      
      {/* Copy button */}
      <button
        onClick={() => navigator.clipboard?.writeText(code)}
        className={`absolute top-2 right-2 p-2 rounded transition-all duration-200 ${
          darkMode
            ? 'bg-gray-700 hover:bg-gray-600 text-gray-300 hover:text-white'
            : 'bg-gray-200 hover:bg-gray-300 text-gray-600 hover:text-gray-800'
        }`}
        title="Copy code"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
      </button>
    </div>
  );
}

export default CodeBlock;