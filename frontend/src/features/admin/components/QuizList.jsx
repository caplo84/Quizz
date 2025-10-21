import React from 'react';
import { Link } from 'react-router-dom';
import { Pencil, Trash2, HelpCircle } from 'lucide-react';

const QuizList = ({ quizzes, searchTerm, onDeleteRequest }) => {
  const term = searchTerm.trim().toLowerCase();

  const filtered = quizzes.filter((quiz) => {
    if (!term) return true;

    return [quiz.title, quiz.slug, quiz.topic]
      .filter(Boolean)
      .some((value) => String(value).toLowerCase().includes(term));
  });

  if (!filtered.length) {
    return (
      <div className="rounded-lg border border-gray-200 bg-white p-6 text-center text-gray-600 shadow-sm">
        No quizzes found.
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                Title
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                Slug
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                Topic
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                Questions
              </th>
              <th className="px-4 py-3 text-right text-xs font-semibold uppercase tracking-wider text-gray-500">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100 bg-white">
            {filtered.map((quiz) => (
              <tr key={quiz.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 text-sm font-medium text-gray-900">{quiz.title || '-'}</td>
                <td className="px-4 py-3 text-sm text-gray-600">{quiz.slug || '-'}</td>
                <td className="px-4 py-3 text-sm text-gray-600">{quiz.topic || '-'}</td>
                <td className="px-4 py-3 text-sm text-gray-600">
                  <span className="inline-flex items-center gap-1 rounded-full bg-blue-50 px-2 py-1 text-xs font-medium text-blue-700">
                    <HelpCircle className="h-3.5 w-3.5" />
                    {quiz.totalQuestions || quiz.questions?.length || 0}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <div className="flex justify-end gap-2">
                    <Link
                      to={`/admin/quizzes/${quiz.id}/edit`}
                      className="inline-flex items-center gap-1 rounded-lg border border-gray-300 px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-100"
                    >
                      <Pencil className="h-3.5 w-3.5" />
                      Edit
                    </Link>
                    <button
                      type="button"
                      onClick={() => onDeleteRequest(quiz.id)}
                      className="inline-flex items-center gap-1 rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-700 hover:bg-red-50"
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default QuizList;