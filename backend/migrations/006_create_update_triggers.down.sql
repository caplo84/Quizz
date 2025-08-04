-- Drop triggers
DROP TRIGGER IF EXISTS update_topics_updated_at ON topics;
DROP TRIGGER IF EXISTS update_quizzes_updated_at ON quizzes;
DROP TRIGGER IF EXISTS update_questions_updated_at ON questions;
DROP TRIGGER IF EXISTS update_choices_updated_at ON choices;
DROP TRIGGER IF EXISTS update_attempts_updated_at ON attempts;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();
