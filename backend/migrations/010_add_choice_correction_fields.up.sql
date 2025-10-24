-- Add correction metadata fields to choices table
ALTER TABLE choices
    ADD COLUMN IF NOT EXISTS answer_source VARCHAR(50) DEFAULT 'parsed',
    ADD COLUMN IF NOT EXISTS ai_confidence DECIMAL(3,2) DEFAULT 0.00,
    ADD COLUMN IF NOT EXISTS corrected_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_choices_answer_source ON choices(answer_source);
CREATE INDEX IF NOT EXISTS idx_choices_corrected_at ON choices(corrected_at);

-- Audit table for correction runs
CREATE TABLE IF NOT EXISTS correction_audit_log (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    old_data JSONB,
    new_data JSONB,
    confidence DECIMAL(3,2),
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT 'system'
);

CREATE INDEX IF NOT EXISTS idx_correction_audit_question_id ON correction_audit_log(question_id);
CREATE INDEX IF NOT EXISTS idx_correction_audit_created_at ON correction_audit_log(created_at);
