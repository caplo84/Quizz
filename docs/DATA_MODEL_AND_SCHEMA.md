# Data Model and Schema Guide

## Core Entities

## Topic

Represents a quiz domain/category.

Key fields:
- `id`
- `name`
- `slug`
- `description`
- `icon_url`
- `is_active`

Relations:
- One Topic has many Quizzes.

## Quiz

Represents a quiz container under a topic.

Key fields:
- `id`
- `title`
- `slug`
- `topic_id`
- `difficulty_level`
- `time_limit_minutes`
- `total_questions`
- `is_active`
- Source metadata: `source`, `external_reference`, `external_id`, `last_synced_at`

Relations:
- One Quiz belongs to one Topic.
- One Quiz has many Questions.
- One Quiz has many Attempts.

## Question

Represents one question item in a quiz.

Key fields:
- `id`
- `quiz_id`
- `question_text`
- `question_type`
- `points`
- `order_index`
- `is_active`
- Optional code/media: `question_code`, `question_code_language`, `question_image_url`, `question_image_alt`
- Source metadata: `source`, `external_reference`, `external_id`

Relations:
- One Question belongs to one Quiz.
- One Question has many Choices.

## Choice

Represents one answer option.

Key fields:
- `id`
- `question_id`
- `choice_text`
- `is_correct`
- `order_index`
- Optional explanation
- Correction metadata: `answer_source`, `ai_confidence`, `corrected_at`
- Optional code/media: `choice_code`, `choice_code_language`, `choice_image_url`, `choice_image_alt`

Relations:
- One Choice belongs to one Question.

## Attempt

Represents one user attempt for one quiz.

Key fields:
- `id`
- `quiz_id`
- `user_identifier`
- `user_name`
- `started_at`
- `completed_at`
- `total_score`
- `max_possible_score`
- `percentage_score`
- `time_taken_seconds`
- `status`
- `answers` (JSONB)
- `is_completed`

Relations:
- One Attempt belongs to one Quiz.

## Migration Strategy

- SQL migrations are stored in `backend/migrations`.
- Current migration set includes:
  - base tables,
  - update triggers,
  - seed data,
  - external-source fields,
  - correction metadata fields.

## Data Integrity Principles

1. Keep slugs stable once published.
2. Preserve source metadata for traceability.
3. Keep question order deterministic for quiz rendering.
4. Treat attempts as append-only records after completion.
5. Apply soft business constraints through service layer when DB constraints are not sufficient.

## Index and Query Notes

- Slug-based lookups are common for topic and quiz retrieval.
- Topic and quiz IDs are used heavily in joins.
- Attempt read/write paths should be monitored for growth and indexing needs.

## Recommended Next Improvements

1. Add explicit uniqueness and conflict policies for external IDs by source.
2. Add partial indexes for active content paths.
3. Add data retention and archiving policy for attempts.
4. Add schema documentation generation from model tags + migration metadata.
