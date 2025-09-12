-- Clean up GitHub-synced data from database
-- This script removes only data that was imported from GitHub

BEGIN;

-- Step 1: Delete choices for questions from GitHub quizzes
DELETE FROM choices 
WHERE question_id IN (
    SELECT q.id FROM questions q
    JOIN quizzes qz ON q.quiz_id = qz.id
    WHERE qz.source = 'github'
);

-- Step 2: Delete questions from GitHub quizzes
DELETE FROM questions 
WHERE quiz_id IN (
    SELECT id FROM quizzes WHERE source = 'github'
);

-- Step 3: Delete attempts for GitHub quizzes
DELETE FROM attempts 
WHERE quiz_id IN (
    SELECT id FROM quizzes WHERE source = 'github'
);

-- Step 4: Delete GitHub quizzes
DELETE FROM quizzes WHERE source = 'github';

-- Step 5: Delete topics that have no remaining quizzes
DELETE FROM topics 
WHERE id NOT IN (
    SELECT DISTINCT topic_id FROM quizzes WHERE topic_id IS NOT NULL
);

-- Show what's left
SELECT 'Topics remaining:' as info, COUNT(*) as count FROM topics
UNION ALL
SELECT 'Quizzes remaining:', COUNT(*) FROM quizzes
UNION ALL
SELECT 'Questions remaining:', COUNT(*) FROM questions
UNION ALL
SELECT 'Choices remaining:', COUNT(*) FROM choices
UNION ALL
SELECT 'Attempts remaining:', COUNT(*) FROM attempts;

COMMIT;