import { Link } from 'react-router-dom';
import { Pencil, Trash2, HelpCircle, FileText, Power } from 'lucide-react';

const QuizList = ({ quizzes, onDeleteRequest, onTogglePublish, publishLoadingId }) => {
  if (!quizzes.length) {
    return (
      <div className="rounded-lg border border-gray-200 bg-white p-8 text-center text-gray-600 shadow-sm">
        <FileText className="mx-auto mb-3 h-10 w-10 text-gray-400" />
        <p className="font-medium text-gray-700">No quizzes found.</p>
        <p className="text-sm text-gray-500 mt-1">Try changing filters or create a new quiz.</p>
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
              <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                Status
              </th>
              <th className="px-4 py-3 text-right text-xs font-semibold uppercase tracking-wider text-gray-500">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100 bg-white">
            {quizzes.map((quiz) => (
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
                <td className="px-4 py-3 text-sm text-gray-600">
                  {quiz.isActive && Number(quiz.totalQuestions || 0) > 0 ? (
                    <span className="inline-flex items-center rounded-full bg-green-100 px-2 py-1 text-xs font-semibold text-green-700">
                      Live
                    </span>
                  ) : (
                    <span className="inline-flex items-center rounded-full bg-amber-100 px-2 py-1 text-xs font-semibold text-amber-700">
                      Draft
                    </span>
                  )}
                </td>
                <td className="px-4 py-3">
                  <div className="flex justify-end gap-2">
                    <button
                      type="button"
                      disabled={publishLoadingId === quiz.id}
                      onClick={() => onTogglePublish(quiz)}
                      className="inline-flex items-center gap-1 rounded-lg border border-indigo-200 px-3 py-1.5 text-xs font-medium text-indigo-700 hover:bg-indigo-50 disabled:opacity-60 disabled:cursor-not-allowed"
                    >
                      <Power className="h-3.5 w-3.5" />
                      {publishLoadingId === quiz.id
                        ? 'Saving...'
                        : quiz.isActive
                          ? 'Unpublish'
                          : 'Publish'}
                    </button>
                    <Link
                      to={`/admin/quizzes/${quiz.id}/edit`}
                      className="inline-flex items-center gap-1 rounded-lg border border-gray-300 px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-100"
                    >
                      <Pencil className="h-3.5 w-3.5" />
                      Edit
                    </Link>
                    <button
                      type="button"
                      onClick={() => onDeleteRequest({ id: quiz.id, title: quiz.title })}
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