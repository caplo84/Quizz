import { CheckSquare, Circle, Square } from "lucide-react";
import CodeBlock from "../../ui/CodeBlock";

function getOptionText(option) {
  if (typeof option === "string") return option;
  return option?.text || option?.choice_text || option?.option || "";
}

function getOptionCode(option) {
  if (typeof option === "string") return null;
  return option?.code || option?.choice_code || null;
}

function getOptionCodeLanguage(option) {
  if (typeof option === "string") return null;
  return option?.code_language || option?.choice_code_language || "javascript";
}

function SingleChoice({ options, selected, onSelect, darkMode }) {
  return (
    <div className="space-y-3">
      {options.map((option, index) => {
        const optionText = getOptionText(option);
        const optionCode = getOptionCode(option);
        const optionCodeLanguage = getOptionCodeLanguage(option);
        const isSelected = selected === index;

        return (
          <button
            key={index}
            onClick={() => onSelect(index)}
            className={`w-full text-left p-3.5 rounded-xl border-2 transition-all duration-150 ${
              isSelected
                ? darkMode
                  ? "border-purple-400 bg-purple-900/30"
                  : "border-purple-500 bg-purple-50"
                : darkMode
                  ? "border-slate-600 hover:border-purple-400 hover:bg-slate-700"
                  : "border-slate-200 hover:border-purple-300 hover:bg-slate-50"
            }`}
          >
            <div className="flex items-start gap-3">
              <span
                className={`shrink-0 w-6 h-6 rounded-full border-2 flex items-center justify-center text-xs font-bold mt-0.5 ${
                  isSelected
                    ? "border-purple-500 bg-purple-500 text-white"
                    : darkMode
                      ? "border-slate-500 text-slate-300"
                      : "border-slate-300 text-slate-500"
                }`}
              >
                {String.fromCharCode(65 + index)}
              </span>

              <div className="flex-1 min-w-0">
                {optionText && (
                  <p className={`text-base ${darkMode ? "text-slate-100" : "text-slate-800"}`}>
                    {optionText}
                  </p>
                )}
                {optionCode && (
                  <CodeBlock code={optionCode} language={optionCodeLanguage} className="mt-2" />
                )}
              </div>
            </div>
          </button>
        );
      })}
    </div>
  );
}

function MultipleChoice({ options, selected = [], onSelect, darkMode }) {
  const toggle = (index) => {
    const next = selected.includes(index)
      ? selected.filter((item) => item !== index)
      : [...selected, index];
    onSelect(next);
  };

  return (
    <div className="space-y-3">
      <p className={`text-sm font-medium ${darkMode ? "text-slate-300" : "text-slate-500"}`}>
        Select all that apply
      </p>
      {options.map((option, index) => {
        const optionText = getOptionText(option);
        const optionCode = getOptionCode(option);
        const optionCodeLanguage = getOptionCodeLanguage(option);
        const isSelected = selected.includes(index);

        return (
          <button
            key={index}
            onClick={() => toggle(index)}
            className={`w-full text-left p-3.5 rounded-xl border-2 transition-all duration-150 ${
              isSelected
                ? darkMode
                  ? "border-purple-400 bg-purple-900/30"
                  : "border-purple-500 bg-purple-50"
                : darkMode
                  ? "border-slate-600 hover:border-purple-400 hover:bg-slate-700"
                  : "border-slate-200 hover:border-purple-300 hover:bg-slate-50"
            }`}
          >
            <div className="flex items-start gap-3">
              {isSelected ? (
                <CheckSquare className="w-5 h-5 text-purple-500 shrink-0 mt-0.5" />
              ) : (
                <Square className={`w-5 h-5 shrink-0 mt-0.5 ${darkMode ? "text-slate-400" : "text-slate-400"}`} />
              )}

              <div className="flex-1 min-w-0">
                {optionText && (
                  <p className={`text-base ${darkMode ? "text-slate-100" : "text-slate-800"}`}>
                    {optionText}
                  </p>
                )}
                {optionCode && (
                  <CodeBlock code={optionCode} language={optionCodeLanguage} className="mt-2" />
                )}
              </div>
            </div>
          </button>
        );
      })}
    </div>
  );
}

function TrueFalse({ selected, onSelect, darkMode }) {
  return (
    <div className="grid grid-cols-2 gap-4">
      {["True", "False"].map((label, index) => {
        const isSelected = selected === index;
        return (
          <button
            key={label}
            onClick={() => onSelect(index)}
            className={`py-4 rounded-xl border-2 font-semibold text-lg transition-all ${
              isSelected
                ? index === 0
                  ? "border-green-500 bg-green-50 text-green-700"
                  : "border-red-500 bg-red-50 text-red-700"
                : darkMode
                  ? "border-slate-600 text-slate-100 hover:border-purple-400"
                  : "border-slate-200 text-slate-800 hover:border-purple-300"
            }`}
          >
            {label}
          </button>
        );
      })}
    </div>
  );
}

function TextInput({ value = "", onChange, darkMode }) {
  return (
    <div>
      <p className={`text-sm font-medium mb-2 ${darkMode ? "text-slate-300" : "text-slate-500"}`}>
        Type your answer below
      </p>
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Your answer..."
        className={`w-full px-4 py-2.5 rounded-lg border-2 text-sm focus:outline-none transition-colors ${
          darkMode
            ? "border-slate-600 bg-slate-800 text-white focus:border-purple-400"
            : "border-slate-200 bg-white text-slate-900 focus:border-purple-500"
        }`}
      />
    </div>
  );
}

export default function AnswerInput({ question, answer, onAnswer, darkMode }) {
  const type = question?.question_type || "single";
  const options = Array.isArray(question?.options) ? question.options : [];

  if (type === "multiple") {
    return <MultipleChoice options={options} selected={Array.isArray(answer) ? answer : []} onSelect={onAnswer} darkMode={darkMode} />;
  }

  if (type === "true_false") {
    return <TrueFalse selected={typeof answer === "number" ? answer : null} onSelect={onAnswer} darkMode={darkMode} />;
  }

  if (type === "text_input") {
    return <TextInput value={typeof answer === "string" ? answer : ""} onChange={onAnswer} darkMode={darkMode} />;
  }

  return <SingleChoice options={options} selected={typeof answer === "number" ? answer : null} onSelect={onAnswer} darkMode={darkMode} />;
}
