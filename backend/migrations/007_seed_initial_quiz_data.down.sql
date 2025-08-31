-- Migration 007 DOWN: Remove initial quiz data

-- Remove choices for all seeded quizzes
DELETE FROM choices WHERE question_id IN (
    SELECT id FROM questions WHERE quiz_id IN (
        SELECT id FROM quizzes WHERE slug IN ('html-basics', 'css-basics', 'javascript-basics', 'accessibility-basics')
    )
);

-- Remove questions for all seeded quizzes
DELETE FROM questions WHERE quiz_id IN (
    SELECT id FROM quizzes WHERE slug IN ('html-basics', 'css-basics', 'javascript-basics', 'accessibility-basics')
);

-- Note: Quizzes and topics are removed in migrations 002 and 001 respectively
