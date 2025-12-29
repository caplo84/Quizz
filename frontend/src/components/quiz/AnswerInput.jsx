import { CheckSquare, Square } from "lucide-react";
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

function normalizeIndex(value) {
  if (typeof value === "number") return Number.isNaN(value) ? null : value;
  if (typeof value === "string" && value.trim() !== "") {
    const parsed = Number(value);
    return Number.isNaN(parsed) ? null : parsed;
  }
  return null;
}

function SingleChoice({ options, selected, onSelect, darkMode }) {
  return (
    <div className="space-y-3">
      {options.map((option, index) => {
        const optionText = getOptionText(option);
        const optionCode = getOptionCode(option);
        const optionCodeLanguage = getOptionCodeLanguage(option);
        const isSelected = normalizeIndex(selected) === index;

        return (
          <button
            key={index}
            type="button"
            onClick={() => onSelect(index)}
            className={`w-full text-left p-3.5 rounded-xl border-2 transition-colors duration-150 focus:outline-none ${
              isSelected
                ? darkMode
                  ? "border-violet-400 bg-violet-900/35"
                  : "border-violet-500 bg-violet-100/70"
                : darkMode
                  ? "border-slate-700 hover:border-violet-400 hover:bg-slate-800"
                  : "border-slate-300 hover:border-violet-300 hover:bg-slate-50"
            }`}
          >
            <div className="flex items-start gap-3">
              <span
                className={`shrink-0 w-6 h-6 rounded-full border-2 flex items-center justify-center text-xs font-bold mt-0.5 ${
                  isSelected
                    ? "border-violet-500 bg-violet-500 text-white"
                    : darkMode
                      ? "border-slate-500 text-slate-300"
                      : "border-slate-300 text-slate-500"
                }`}
              >
                {String.fromCharCode(65 + index)}
              </span>

              <div className="flex-1 min-w-0">
                {optionText && (
                  <p className={`text-base font-normal ${darkMode ? (isSelected ? "text-white" : "text-slate-100") : (isSelected ? "text-slate-900" : "text-slate-800")}`}>
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
  const selectedIndices = Array.isArray(selected)
    ? selected
        .map((value) => normalizeIndex(value))
        .filter((value) => value !== null)
    : [];

  const toggle = (index) => {
    const next = selectedIndices.includes(index)
      ? selectedIndices.filter((item) => item !== index)
      : [...selectedIndices, index];
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
        const isSelected = selectedIndices.includes(index);

        return (
          <button
            key={index}
            type="button"
            onClick={() => toggle(index)}
            className={`w-full text-left p-3.5 rounded-xl border-2 transition-colors duration-150 focus:outline-none ${
              isSelected
                ? darkMode
                  ? "border-violet-400 bg-violet-900/35"
                  : "border-violet-500 bg-violet-100/70"
                : darkMode
                  ? "border-slate-700 hover:border-violet-400 hover:bg-slate-800"
                  : "border-slate-300 hover:border-violet-300 hover:bg-slate-50"
            }`}
          >
            <div className="flex items-start gap-3">
              {isSelected ? (
                <CheckSquare className="w-5 h-5 text-violet-500 shrink-0 mt-0.5" />
              ) : (
                <Square className={`w-5 h-5 shrink-0 mt-0.5 ${darkMode ? "text-slate-400" : "text-slate-400"}`} />
              )}

              <div className="flex-1 min-w-0">
                {optionText && (
                  <p className={`text-base font-normal ${darkMode ? (isSelected ? "text-white" : "text-slate-100") : (isSelected ? "text-slate-900" : "text-slate-800")}`}>
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
                  ? "border-slate-600 text-slate-100 hover:border-violet-400"
                  : "border-slate-200 text-slate-800 hover:border-violet-300"
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
            ? "border-slate-600 bg-slate-800 text-white focus:border-violet-400"
            : "border-slate-200 bg-white text-slate-900 focus:border-violet-500"
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
    return <TrueFalse selected={normalizeIndex(answer)} onSelect={onAnswer} darkMode={darkMode} />;
  }

  if (type === "text_input") {
    return <TextInput value={typeof answer === "string" ? answer : ""} onChange={onAnswer} darkMode={darkMode} />;
  }

  return <SingleChoice options={options} selected={normalizeIndex(answer)} onSelect={onAnswer} darkMode={darkMode} />;
}
