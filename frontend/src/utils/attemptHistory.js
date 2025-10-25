const ATTEMPTS_KEY = "quiz_attempt_history";
const USER_KEY = "quiz_user_profile";

export function getAttempts() {
  try {
    const raw = localStorage.getItem(ATTEMPTS_KEY);
    const parsed = raw ? JSON.parse(raw) : [];
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

export function saveAttempts(attempts) {
  localStorage.setItem(ATTEMPTS_KEY, JSON.stringify(attempts));
}

export function addAttempt(attempt) {
  const attempts = getAttempts();
  const id = typeof crypto !== "undefined" && crypto.randomUUID
    ? crypto.randomUUID()
    : `${Date.now()}-${Math.random().toString(36).slice(2)}`;
  const next = [{
    id,
    createdAt: new Date().toISOString(),
    ...attempt,
  }, ...attempts].slice(0, 300);
  saveAttempts(next);
  return next;
}

export function getCurrentUser() {
  try {
    const raw = localStorage.getItem(USER_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

export function setCurrentUser(user) {
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearCurrentUser() {
  localStorage.removeItem(USER_KEY);
}

export function getLeaderboard(topicSlug = "all") {
  const attempts = getAttempts().filter((item) => {
    if (topicSlug === "all") return true;
    return item.topicSlug === topicSlug;
  });

  return attempts.sort((a, b) => {
    if (b.percentage !== a.percentage) return b.percentage - a.percentage;
    return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
  });
}
