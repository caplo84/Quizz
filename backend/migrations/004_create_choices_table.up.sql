-- Create choices table (for multiple choice questions)
CREATE TABLE choices (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    choice_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT false,
    order_index INTEGER NOT NULL,
    explanation TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure unique ordering within each question
    UNIQUE(question_id, order_index)
);

-- Create indexes
CREATE INDEX idx_choices_question_id ON choices(question_id);
CREATE INDEX idx_choices_is_correct ON choices(is_correct);
CREATE INDEX idx_choices_order ON choices(question_id, order_index);
