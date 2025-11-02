import { useState, useEffect, useMemo, useCallback } from 'react';
import { Link } from 'react-router-dom';
import {
  Plus as PlusIcon,
  Search as MagnifyingGlassIcon,
  RefreshCw,
  ArrowUpDown,
  ChevronLeft,
  ChevronRight,
  ShieldAlert,
  ListChecks,
} from 'lucide-react';
import adminApi from '../../services/adminApi';
import QuizList from './components/QuizList';

const PAGE_SIZE = 10;
const AUDIT_STORAGE_KEY = 'quiz_admin_audit_log';

const getReviewAssessment = (report) => {
  const processed = Number(report?.total_processed || 0);
  const flagged = Number(report?.total_fixed || 0);
  const failed = Number(report?.total_failed || 0);
  const ratio = processed > 0 ? flagged / processed : 0;

  if (failed > 0 || ratio >= 0.2) {
    return {
      level: 'danger',
      comment: `Nguy hiểm: ${flagged} câu có vấn đề, ${failed} câu lỗi phân tích. Cần review thủ công trước publish.`,
      canPublish: false,
    };
  }

  if (flagged > 0) {
    return {
      level: 'warning',
      comment: `Cảnh báo: phát hiện ${flagged} câu cần xem lại. Nên chỉnh tay trước publish.`,
      canPublish: false,
    };
  }

  return {
    level: 'good',
    comment: `Tốt: không thấy cảnh báo trên ${processed} câu đã quét.`,
    canPublish: true,
  };
};

const QuizManagement = () => {
  const [quizzes, setQuizzes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [topicFilter, setTopicFilter] = useState('all');
  const [statusFilter, setStatusFilter] = useState('all');
  const [sortBy, setSortBy] = useState('updated_desc');
  const [page, setPage] = useState(1);
  const [deleteConfirm, setDeleteConfirm] = useState(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [publishLoadingId, setPublishLoadingId] = useState(null);
  const [reviewLoadingId, setReviewLoadingId] = useState(null);
  const [reviewGateByQuiz, setReviewGateByQuiz] = useState({});
  const [reviewModalQuizId, setReviewModalQuizId] = useState(null);
  const [error, setError] = useState('');
  const [auditLog, setAuditLog] = useState(() => {
    try {
      const raw = localStorage.getItem(AUDIT_STORAGE_KEY);
      const parsed = raw ? JSON.parse(raw) : [];
      return Array.isArray(parsed) ? parsed : [];
    } catch {
      return [];
    }
  });

  const appendAudit = useCallback((action, message) => {
    setAuditLog((current) => {
      const next = [
        {
          id: crypto.randomUUID(),
          action,
          message,
          at: new Date().toISOString(),
        },
        ...current,
      ].slice(0, 25);
      localStorage.setItem(AUDIT_STORAGE_KEY, JSON.stringify(next));
      return next;
    });
  }, []);

  const loadQuizzes = useCallback(async () => {
    try {
      setLoading(true);
      setError('');
      const data = await adminApi.getAllQuizzes();
      setQuizzes(data || []);
      appendAudit('LOAD_SUCCESS', `Loaded ${data?.length || 0} quizzes`);
    } catch (err) {
      console.error('Failed to load quizzes:', err);
      setError('Failed to load quizzes. Please try again.');
      appendAudit('LOAD_FAILED', 'Failed to load quizzes from API');
    } finally {
      setLoading(false);
    }
  }, [appendAudit]);

  useEffect(() => {
    loadQuizzes();
  }, [loadQuizzes]);

  useEffect(() => {
    setPage(1);
  }, [searchTerm, topicFilter, statusFilter, sortBy]);

  const handleDelete = async () => {
    if (!deleteConfirm?.id) return;

    const quizId = deleteConfirm.id;
    const previousQuizzes = quizzes;

    setIsDeleting(true);
    setQuizzes((current) => current.filter((q) => q.id !== quizId));

    try {
      await adminApi.deleteQuiz(quizId);
      appendAudit('DELETE_SUCCESS', `Deleted quiz #${quizId}`);
      setDeleteConfirm(null);
    } catch (err) {
      setQuizzes(previousQuizzes);
      console.error('Failed to delete quiz:', err);
      setError('Failed to delete quiz. The list has been restored. Please try again.');
      appendAudit('DELETE_FAILED', `Failed deleting quiz #${quizId}; rolled back`);
    } finally {
      setIsDeleting(false);
    }
  };

  const handleReviewQuiz = async (quiz) => {
    if (!quiz?.slug) {
      setError('Quiz slug is missing, cannot run AI review.');
      return;
    }

    try {
      setError('');
      setReviewLoadingId(quiz.id);

      const reviewResponse = await adminApi.reviewQuizBeforePublish(quiz.slug);
      const report = reviewResponse?.report || {};
      const fixed = Number(report.total_fixed || 0);
      const failed = Number(report.total_failed || 0);
      const processed = Number(report.total_processed || 0);
      const assessment = getReviewAssessment(report);

      setReviewGateByQuiz((current) => ({
        ...current,
        [quiz.id]: {
          quizId: quiz.id,
          quizTitle: quiz.title || '',
          quizSlug: quiz.slug || '',
          isClean: assessment.canPublish,
          level: assessment.level,
          comment: assessment.comment,
          canPublish: assessment.canPublish,
          fixed,
          failed,
          processed,
          skipped: Number(report.total_skipped || 0),
          duration: report.duration || '',
          bySource: report.by_source || {},
          byConfidence: report.by_confidence || {},
          estimatedApiCost: Number(report.estimated_api_cost || 0),
          details: Array.isArray(report.details) ? report.details : [],
          reviewedAt: new Date().toISOString(),
        },
      }));
      setReviewModalQuizId(quiz.id);

      setError(
        `AI review-only: không ghi DB. ${assessment.comment}`,
      );
      appendAudit(
        'AI_REVIEW_COMPLETED',
        `Quiz #${quiz.id} => level=${assessment.level}, processed=${processed}, flagged=${fixed}, failed=${failed}`,
      );
    } catch (err) {
      console.error('AI review failed:', err);
      setError('AI review failed. Please retry.');
      appendAudit('AI_REVIEW_FAILED', `Review failed for quiz #${quiz?.id}`);
    } finally {
      setReviewLoadingId(null);
    }
  };

  const openReviewReport = (quiz) => {
    const existing = reviewGateByQuiz?.[quiz?.id];
    if (!existing) {
      setError(`Chưa có report cho "${quiz?.title || quiz?.slug || 'quiz'}". Hãy chạy AI Review trước.`);
      return;
    }
    setReviewModalQuizId(quiz.id);
  };

  const handleTogglePublish = async (quiz) => {
    const previous = quizzes;
    const nextIsActive = !quiz.isActive;

    setPublishLoadingId(quiz.id);

    if (nextIsActive) {
      const gate = reviewGateByQuiz[quiz.id];
      if (!gate) {
        setError(`Please run AI Review for "${quiz.title || quiz.slug}" before publishing.`);
        appendAudit('PUBLISH_BLOCKED_NO_REVIEW', `Blocked publish for quiz #${quiz.id} (no review)`);
        setPublishLoadingId(null);
        return;
      }

      if (!gate.canPublish) {
        setError(
          `AI Review chưa đạt mức tốt cho "${quiz.title || quiz.slug}" (${gate.level}). ${gate.comment}`,
        );
        appendAudit('PUBLISH_BLOCKED_REVIEW_ISSUES', `Blocked publish for quiz #${quiz.id} (review not clean)`);
        setPublishLoadingId(null);
        return;
      }
    }

    setQuizzes((current) =>
      current.map((item) =>
        item.id === quiz.id ? { ...item, isActive: nextIsActive } : item,
      ),
    );

    try {
      await adminApi.setQuizPublicationStatus(quiz.id, nextIsActive);

      appendAudit(
        'PUBLISH_STATUS_CHANGED',
        `${nextIsActive ? 'Published' : 'Unpublished'} quiz "${quiz.title || quiz.id}"`,
      );
    } catch (err) {
      setQuizzes(previous);
      setError(err?.message || 'Failed to update publication status. Changes were rolled back.');
      appendAudit('PUBLISH_FAILED', `Failed publication change for quiz #${quiz.id}`);
    } finally {
      setPublishLoadingId(null);
    }
  };

  const uniqueTopics = useMemo(() => {
    return Array.from(new Set(quizzes.map((quiz) => quiz.topic).filter(Boolean)));
  }, [quizzes]);

  const filteredQuizzes = useMemo(() => {
    const term = searchTerm.trim().toLowerCase();

    return quizzes.filter((quiz) => {
      const matchesSearch =
        !term ||
        [quiz.title, quiz.slug, quiz.topic]
          .filter(Boolean)
          .some((value) => String(value).toLowerCase().includes(term));

      const matchesTopic = topicFilter === 'all' || quiz.topic === topicFilter;

      const questionCount = Number(quiz.totalQuestions || quiz.questions?.length || 0);
      const isLive = Boolean(quiz.isActive) && questionCount > 0;
      const matchesStatus =
        statusFilter === 'all' ||
        (statusFilter === 'live' && isLive) ||
        (statusFilter === 'draft' && !isLive);

      return matchesSearch && matchesTopic && matchesStatus;
    });
  }, [quizzes, searchTerm, topicFilter, statusFilter]);

  const sortedQuizzes = useMemo(() => {
    const next = [...filteredQuizzes];

    const getCount = (quiz) => Number(quiz.totalQuestions || quiz.questions?.length || 0);

    switch (sortBy) {
      case 'title_asc':
        next.sort((a, b) => String(a.title || '').localeCompare(String(b.title || '')));
        break;
      case 'title_desc':
        next.sort((a, b) => String(b.title || '').localeCompare(String(a.title || '')));
        break;
      case 'questions_desc':
        next.sort((a, b) => getCount(b) - getCount(a));
        break;
      case 'questions_asc':
        next.sort((a, b) => getCount(a) - getCount(b));
        break;
      default:
        next.sort((a, b) => Number(b.id || 0) - Number(a.id || 0));
        break;
    }

    return next;
  }, [filteredQuizzes, sortBy]);

  const pageCount = Math.max(1, Math.ceil(sortedQuizzes.length / PAGE_SIZE));
  const currentPage = Math.min(page, pageCount);
  const paginatedQuizzes = useMemo(() => {
    const start = (currentPage - 1) * PAGE_SIZE;
    return sortedQuizzes.slice(start, start + PAGE_SIZE);
  }, [sortedQuizzes, currentPage]);

  const totalLive = quizzes.filter(
    (quiz) => Boolean(quiz.isActive) && Number(quiz.totalQuestions || quiz.questions?.length || 0) > 0,
  ).length;
  const totalDraft = quizzes.length - totalLive;

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Quiz Management</h1>
          <p className="text-gray-600 mt-1">Manage all quizzes in your system</p>
        </div>
        <Link
          to="/admin/quizzes/new"
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          <PlusIcon className="w-5 h-5" />
          Create Quiz
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-lg shadow-md p-4 border border-gray-200">
          <p className="text-xs uppercase tracking-wide text-gray-500">Total Quizzes</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{quizzes.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow-md p-4 border border-gray-200">
          <p className="text-xs uppercase tracking-wide text-gray-500">Live</p>
          <p className="text-2xl font-bold text-green-700 mt-1">{totalLive}</p>
        </div>
        <div className="bg-white rounded-lg shadow-md p-4 border border-gray-200">
          <p className="text-xs uppercase tracking-wide text-gray-500">Draft</p>
          <p className="text-2xl font-bold text-amber-700 mt-1">{totalDraft}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-4 space-y-3">
        <div className="flex items-center justify-between gap-3">
          <h2 className="text-sm font-semibold text-gray-700 uppercase tracking-wide">Filters</h2>
          <button
            type="button"
            onClick={loadQuizzes}
            className="inline-flex items-center gap-2 px-3 py-2 border border-gray-300 rounded-lg text-sm text-gray-700 hover:bg-gray-50"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-3">
          <div className="relative lg:col-span-2">
            <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search by title, slug, or topic"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <select
            value={topicFilter}
            onChange={(e) => setTopicFilter(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="all">All topics</option>
            {uniqueTopics.map((topic) => (
              <option key={topic} value={topic}>
                {topic}
              </option>
            ))}
          </select>

          <div className="flex gap-2">
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="all">All status</option>
              <option value="live">Live</option>
              <option value="draft">Draft</option>
            </select>

            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="updated_desc">Newest</option>
              <option value="title_asc">Title A-Z</option>
              <option value="title_desc">Title Z-A</option>
              <option value="questions_desc">Most Questions</option>
              <option value="questions_asc">Least Questions</option>
            </select>
          </div>
        </div>

        <div className="flex items-center gap-2 text-xs text-gray-500">
          <ArrowUpDown className="w-3.5 h-3.5" />
          Showing {paginatedQuizzes.length} of {sortedQuizzes.length} matched quizzes
        </div>
      </div>

      {error ? (
        <div className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-lg p-3 flex items-center justify-between gap-3">
          <span>{error}</span>
          <button
            type="button"
            onClick={loadQuizzes}
            className="px-3 py-1.5 rounded-md border border-red-300 hover:bg-red-100"
          >
            Retry
          </button>
        </div>
      ) : null}

      <QuizList
        quizzes={paginatedQuizzes}
        onDeleteRequest={setDeleteConfirm}
        onReviewQuiz={handleReviewQuiz}
        onOpenReviewReport={openReviewReport}
        onTogglePublish={handleTogglePublish}
        publishLoadingId={publishLoadingId}
        reviewLoadingId={reviewLoadingId}
        reviewGateByQuiz={reviewGateByQuiz}
      />

      {reviewModalQuizId && reviewGateByQuiz?.[reviewModalQuizId] ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 px-4 py-6">
          <div className="max-h-[90vh] w-full max-w-4xl overflow-hidden rounded-2xl bg-white shadow-2xl">
            {(() => {
              const report = reviewGateByQuiz[reviewModalQuizId];
              const levelClassName =
                report.level === 'good'
                  ? 'bg-emerald-100 text-emerald-700'
                  : report.level === 'warning'
                    ? 'bg-amber-100 text-amber-700'
                    : 'bg-red-100 text-red-700';

              return (
                <>
                  <div className="flex items-start justify-between border-b border-gray-200 px-6 py-4">
                    <div>
                      <p className="text-xs uppercase tracking-wide text-gray-500">AI Review Report</p>
                      <h3 className="mt-1 text-xl font-bold text-gray-900">
                        {report.quizTitle || report.quizSlug || `Quiz #${report.quizId}`}
                      </h3>
                      <p className="mt-1 text-sm text-gray-500">
                        Reviewed at {new Date(report.reviewedAt).toLocaleString()} • Duration {report.duration || '-'}
                      </p>
                    </div>
                    <div className="flex items-center gap-3">
                      <span className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold ${levelClassName}`}>
                        {report.level === 'good'
                          ? 'Tốt'
                          : report.level === 'warning'
                            ? 'Cảnh báo'
                            : 'Nguy hiểm'}
                      </span>
                      <button
                        type="button"
                        onClick={() => setReviewModalQuizId(null)}
                        className="rounded-lg border border-gray-300 px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50"
                      >
                        Close
                      </button>
                    </div>
                  </div>

                  <div className="space-y-4 overflow-y-auto px-6 py-4">
                    <div className="rounded-lg border border-gray-200 bg-gray-50 p-3 text-sm text-gray-700">
                      {report.comment}
                    </div>

                    <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
                      <div className="rounded-lg border border-gray-200 p-3">
                        <p className="text-xs uppercase text-gray-500">Processed</p>
                        <p className="mt-1 text-lg font-bold text-gray-900">{report.processed}</p>
                      </div>
                      <div className="rounded-lg border border-amber-200 bg-amber-50 p-3">
                        <p className="text-xs uppercase text-amber-700">Flagged</p>
                        <p className="mt-1 text-lg font-bold text-amber-800">{report.fixed}</p>
                      </div>
                      <div className="rounded-lg border border-red-200 bg-red-50 p-3">
                        <p className="text-xs uppercase text-red-700">Failed</p>
                        <p className="mt-1 text-lg font-bold text-red-800">{report.failed}</p>
                      </div>
                      <div className="rounded-lg border border-gray-200 p-3">
                        <p className="text-xs uppercase text-gray-500">Skipped</p>
                        <p className="mt-1 text-lg font-bold text-gray-900">{report.skipped}</p>
                      </div>
                    </div>

                    <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                      <div className="rounded-lg border border-gray-200 p-3">
                        <p className="text-sm font-semibold text-gray-800">By Source</p>
                        <ul className="mt-2 space-y-1 text-sm text-gray-600">
                          {Object.entries(report.bySource || {}).map(([k, v]) => (
                            <li key={k} className="flex items-center justify-between">
                              <span>{k}</span>
                              <span className="font-medium text-gray-800">{v}</span>
                            </li>
                          ))}
                        </ul>
                      </div>
                      <div className="rounded-lg border border-gray-200 p-3">
                        <p className="text-sm font-semibold text-gray-800">By Confidence</p>
                        <ul className="mt-2 space-y-1 text-sm text-gray-600">
                          {Object.entries(report.byConfidence || {}).map(([k, v]) => (
                            <li key={k} className="flex items-center justify-between">
                              <span>{k}</span>
                              <span className="font-medium text-gray-800">{v}</span>
                            </li>
                          ))}
                        </ul>
                      </div>
                    </div>

                    <div className="rounded-lg border border-gray-200">
                      <div className="border-b border-gray-200 px-3 py-2">
                        <p className="text-sm font-semibold text-gray-800">Findings</p>
                      </div>
                      {report.details?.length ? (
                        <div className="max-h-64 overflow-y-auto divide-y divide-gray-100">
                          {report.details.slice(0, 100).map((item, index) => (
                            <div key={`${item.question_id || 'q'}-${index}`} className="px-3 py-2 text-sm">
                              <div className="flex items-center justify-between gap-2">
                                <p className="font-medium text-gray-800">
                                  Question #{item.question_id || '-'} • {item.action || 'reviewed'}
                                </p>
                                {item.confidence !== undefined ? (
                                  <span className="text-xs text-gray-500">conf: {item.confidence}</span>
                                ) : null}
                              </div>
                              {item.issues?.length ? (
                                <p className="mt-1 text-xs text-amber-700">{item.issues.join(' • ')}</p>
                              ) : null}
                              {item.error ? (
                                <p className="mt-1 text-xs text-red-600">{item.error}</p>
                              ) : null}
                            </div>
                          ))}
                        </div>
                      ) : (
                        <p className="px-3 py-4 text-sm text-gray-500">Không có findings chi tiết cho lần quét này.</p>
                      )}
                    </div>
                  </div>
                </>
              );
            })()}
          </div>
        </div>
      ) : null}

      {quizzes.length > 500 && (
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 text-sm text-blue-700 flex items-center gap-2">
          <ListChecks className="w-4 h-4" />
          Large dataset mode active: filtering/sorting is memoized and paginated for performance.
        </div>
      )}

      <div className="bg-white rounded-lg shadow-sm border border-gray-200">
        <div className="px-4 py-3 border-b border-gray-200 flex items-center gap-2">
          <ShieldAlert className="w-4 h-4 text-gray-600" />
          <h3 className="text-sm font-semibold text-gray-700">Audit Log (Recent Actions)</h3>
        </div>
        <div className="max-h-56 overflow-y-auto">
          {auditLog.length === 0 ? (
            <p className="px-4 py-3 text-sm text-gray-500">No admin actions logged yet.</p>
          ) : (
            <ul className="divide-y divide-gray-100">
              {auditLog.map((entry) => (
                <li key={entry.id} className="px-4 py-2 text-sm flex items-center justify-between gap-3">
                  <div>
                    <p className="font-medium text-gray-700">{entry.action}</p>
                    <p className="text-gray-500">{entry.message}</p>
                  </div>
                  <span className="text-xs text-gray-400 whitespace-nowrap">
                    {new Date(entry.at).toLocaleString()}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>

      {pageCount > 1 && (
        <div className="flex items-center justify-between bg-white rounded-lg shadow-sm border border-gray-200 px-4 py-3">
          <p className="text-sm text-gray-600">
            Page {currentPage} of {pageCount}
          </p>
          <div className="flex items-center gap-2">
            <button
              type="button"
              disabled={currentPage === 1}
              onClick={() => setPage((prev) => Math.max(1, prev - 1))}
              className="inline-flex items-center gap-1 px-3 py-1.5 border border-gray-300 rounded-lg text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              <ChevronLeft className="w-4 h-4" /> Prev
            </button>
            <button
              type="button"
              disabled={currentPage === pageCount}
              onClick={() => setPage((prev) => Math.min(pageCount, prev + 1))}
              className="inline-flex items-center gap-1 px-3 py-1.5 border border-gray-300 rounded-lg text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              Next <ChevronRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}

      {deleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 className="text-xl font-bold text-gray-900 mb-4">Confirm Delete</h3>
            <p className="text-gray-600 mb-6">
              Are you sure you want to delete <strong>{deleteConfirm.title || 'this quiz'}</strong>? This action cannot be undone.
            </p>
            <div className="flex gap-3 justify-end">
              <button
                disabled={isDeleting}
                onClick={() => setDeleteConfirm(null)}
                className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                disabled={isDeleting}
                onClick={handleDelete}
                className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:bg-red-400 transition-colors"
              >
                {isDeleting ? 'Deleting...' : 'Delete'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default QuizManagement;
