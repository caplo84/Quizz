-- Add external source tracking fields to quizzes table
ALTER TABLE quizzes ADD COLUMN source VARCHAR(50);
ALTER TABLE quizzes ADD COLUMN external_reference VARCHAR(500);
ALTER TABLE quizzes ADD COLUMN external_id VARCHAR(200);
ALTER TABLE quizzes ADD COLUMN last_synced_at TIMESTAMP;

-- Add indexes for better performance
CREATE INDEX idx_quizzes_source ON quizzes(source);
CREATE INDEX idx_quizzes_external_id ON quizzes(external_id);

-- Add external source tracking fields to questions table
ALTER TABLE questions ADD COLUMN source VARCHAR(50);
ALTER TABLE questions ADD COLUMN external_reference VARCHAR(500);
ALTER TABLE questions ADD COLUMN external_id VARCHAR(200);