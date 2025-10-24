DROP TABLE IF EXISTS correction_audit_log;

DROP INDEX IF EXISTS idx_choices_answer_source;
DROP INDEX IF EXISTS idx_choices_corrected_at;

ALTER TABLE choices
    DROP COLUMN IF EXISTS answer_source,
    DROP COLUMN IF EXISTS ai_confidence,
    DROP COLUMN IF EXISTS corrected_at;
