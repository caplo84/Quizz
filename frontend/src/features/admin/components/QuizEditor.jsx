import React from 'react';
import { Plus, Trash2 } from 'lucide-react';

const EMPTY_QUESTION = {
  question: '',
  choices: ['', '', '', ''],
  correctAnswers: [],
};

const QuizEditor = ({ questions, onQuestionsChange }) => {
  const updateQuestion = (index, partial) => {
    const next = questions.map((q, i) => (i === index ? { ...q, ...partial } : q));
    onQuestionsChange(next);
  };

  const addQuestion = () => {
    onQuestionsChange([...(questions || []), { ...EMPTY_QUESTION }]);
  };

  const removeQuestion = (index) => {
    onQuestionsChange(questions.filter((_, i) => i !== index));
  };

  const updateChoice = (qIndex, cIndex, value) => {
    const q = questions[qIndex];
    const updatedChoices = q.choices.map((choice, index) => (index === cIndex ? value : choice));
    updateQuestion(qIndex, { choices: updatedChoices });
  };

  const toggleCorrectAnswer = (qIndex, choiceIndex) => {
    const q = questions[qIndex];
    const exists = q.correctAnswers?.includes(choiceIndex);
    const nextAnswers = exists
      ? q.correctAnswers.filter((value) => value !== choiceIndex)
      : [...(q.correctAnswers || []), choiceIndex];
    updateQuestion(qIndex, { correctAnswers: nextAnswers });
  };

  return (
    <div className="rounded-lg bg-white p-6 shadow-md">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-xl font-bold text-gray-900">Questions</h2>
        <button
          type="button"
          onClick={addQuestion}
          className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-3 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          <Plus className="h-4 w-4" />
          Add Question
        </button>
      </div>

      {!questions.length ? (
        <div className="rounded-lg border border-dashed border-gray-300 p-6 text-center text-sm text-gray-600">
          No questions added yet.
        </div>
      ) : null}

      <div className="space-y-6">
        {questions.map((q, qIndex) => (
          <div key={`question-${qIndex}`} className="rounded-lg border border-gray-200 p-4">
            <div className="mb-3 flex items-start justify-between gap-3">
              <label className="w-full">
                <span className="mb-1 block text-sm font-semibold text-gray-700">
                  Question {qIndex + 1}
                </span>
                <textarea
                  value={q.question || ''}
                  onChange={(e) => updateQuestion(qIndex, { question: e.target.value })}
                  rows={2}
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Enter question text"
                />
              </label>

              <button
                type="button"
                onClick={() => removeQuestion(qIndex)}
                className="mt-6 inline-flex items-center gap-1 rounded-lg border border-red-200 px-3 py-2 text-xs font-medium text-red-700 hover:bg-red-50"
              >
                <Trash2 className="h-3.5 w-3.5" />
                Remove
              </button>
            </div>

            <div className="space-y-2">
              <p className="text-sm font-medium text-gray-700">Choices (check correct answers)</p>
              {q.choices?.map((choice, cIndex) => (
                <div key={`choice-${qIndex}-${cIndex}`} className="flex items-center gap-3">
                  <input
                    type="checkbox"
                    checked={Boolean(q.correctAnswers?.includes(cIndex))}
                    onChange={() => toggleCorrectAnswer(qIndex, cIndex)}
                    className="h-4 w-4 rounded border-gray-300 text-blue-600"
                  />
                  <input
                    type="text"
                    value={choice || ''}
                    onChange={(e) => updateChoice(qIndex, cIndex, e.target.value)}
                    className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder={`Choice ${cIndex + 1}`}
                  />
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default QuizEditor;