import { Link } from 'react-router-dom';
import { Pencil, Trash2, HelpCircle, FileText, Power, Sparkles, Eye } from 'lucide-react';

const QuizList = ({
  quizzes,
  onDeleteRequest,
  onReviewQuiz,
  onOpenReviewReport,
  onTogglePublish,
  publishLoadingId,
  reviewLoadingId,
  reviewGateByQuiz,
}) => {
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
                {(() => {
                  const reviewState = reviewGateByQuiz?.[quiz.id];
                  const canPublish = Boolean(quiz.isActive) || Boolean(reviewState?.canPublish);
                  const publishDisabled = publishLoadingId === quiz.id || !canPublish;

                  return (
                    <>
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
                  {reviewState ? (
                    <div className="mt-1 space-y-1">
                      {reviewState.level === 'good' ? (
                        <span className="inline-flex items-center rounded-full bg-emerald-100 px-2 py-1 text-[11px] font-semibold text-emerald-700">
                          AI: Good
                        </span>
                      ) : reviewState.level === 'warning' ? (
                        <span className="inline-flex items-center rounded-full bg-amber-100 px-2 py-1 text-[11px] font-semibold text-amber-700">
                          AI: Warning
                        </span>
                      ) : (
                        <span className="inline-flex items-center rounded-full bg-red-100 px-2 py-1 text-[11px] font-semibold text-red-700">
                          AI: Critical
                        </span>
                      )}
                      <p className="text-[11px] text-gray-500 max-w-xs">{reviewState.comment}</p>
                    </div>
                  ) : null}
                </td>
                <td className="px-4 py-3">
                  <div className="flex justify-end gap-2">
                    <button
                      type="button"
                      disabled={reviewLoadingId === quiz.id || publishLoadingId === quiz.id}
                      onClick={() => onReviewQuiz(quiz)}
                      className="inline-flex items-center gap-1 rounded-lg border border-violet-200 px-3 py-1.5 text-xs font-medium text-violet-700 hover:bg-violet-50 disabled:opacity-60 disabled:cursor-not-allowed"
                    >
                      <Sparkles className="h-3.5 w-3.5" />
                      {reviewLoadingId === quiz.id ? 'Reviewing...' : 'AI Review'}
                    </button>
                    <button
                      type="button"
                      disabled={!reviewState}
                      onClick={() => onOpenReviewReport(quiz)}
                      className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-3 py-1.5 text-xs font-medium text-slate-700 hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
                      title={reviewState ? 'View AI review report' : 'Run AI Review first to view the report'}
                    >
                      <Eye className="h-3.5 w-3.5" />
                      Report
                    </button>
                    <button
                      type="button"
                      disabled={publishDisabled}
                      onClick={() => onTogglePublish(quiz)}
                      className="inline-flex items-center gap-1 rounded-lg border border-indigo-200 px-3 py-1.5 text-xs font-medium text-indigo-700 hover:bg-indigo-50 disabled:opacity-60 disabled:cursor-not-allowed"
                      title={!canPublish && !quiz.isActive ? 'Run AI Review first. Publish only when the rating is Good.' : ''}
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
                    </>
                  );
                })()}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default QuizList;