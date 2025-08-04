-- Create attempts table
CREATE TABLE attempts (
    id SERIAL PRIMARY KEY,
    quiz_id INTEGER NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    user_identifier VARCHAR(255), -- For anonymous users (IP, session, etc.)
    user_name VARCHAR(100),
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    total_score INTEGER DEFAULT 0,
    max_possible_score INTEGER DEFAULT 0,
    percentage_score DECIMAL(5,2) DEFAULT 0.00,
    time_taken_seconds INTEGER,
    status VARCHAR(20) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'abandoned')),
    answers JSONB, -- Store user answers as JSON
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_attempts_quiz_id ON attempts(quiz_id);
CREATE INDEX idx_attempts_user_identifier ON attempts(user_identifier);
CREATE INDEX idx_attempts_status ON attempts(status);
CREATE INDEX idx_attempts_completed_at ON attempts(completed_at);
CREATE INDEX idx_attempts_percentage_score ON attempts(percentage_score);

-- Create GIN index for JSONB answers column for efficient querying
CREATE INDEX idx_attempts_answers ON attempts USING GIN (answers);
