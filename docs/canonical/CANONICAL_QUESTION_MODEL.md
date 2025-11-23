# Canonical Question Model

## Overview
This document defines the canonical data model for quiz questions imported from the upstream community quiz repository.

## Data Source
- **Repository**: [your-org/quiz-content-source](https://github.com/your-org/quiz-content-source)
- **Format**: Markdown files organized by topic
- **Update Frequency**: Community-maintained, irregular updates
- **License**: Educational use (see DISCLAIMER.md)

---

## Entity Model

### Topic
Represents a skill category (e.g., JavaScript, Python, CSS).

```json
{
  "name": "JavaScript",
  "slug": "javascript",
  "description": "JavaScript Programming Language Assessment",
  "icon_url": null,
  "is_active": true,
  "source": "github",
  "external_reference": "https://github.com/your-org/quiz-content-source/blob/main/javascript/javascript-quiz.md"
}
```

**Mapping Rules:**
- Folder name → `slug` (normalized)
- Folder name (title case) → `name`
- File path → `external_reference`

---

### Quiz
One quiz per topic file.

```json
{
  "title": "JavaScript Assessment",
  "slug": "javascript-assessment",
  "description": "Test your JavaScript knowledge with questions from upstream community quiz dataset",
  "topic_id": "[TOPIC_ID]",
  "difficulty_level": "medium",
  "time_limit_minutes": 30,
  "total_questions": 166,
  "is_active": true,
  "source": "github",
  "external_reference": "https://github.com/your-org/quiz-content-source/blob/main/javascript/javascript-quiz.md",
  "external_id": "javascript-quiz-main",
  "last_synced_at": "2026-01-22T00:00:00Z"
}
```

**Mapping Rules:**
- Default difficulty: `medium`
- Default time limit: `30 minutes`
- `total_questions`: Counted from markdown
- `source`: Always `"github"`

---

### Question
Individual quiz question with optional content blocks.

```json
{
  "quiz_id": "[QUIZ_ID]",
  "question_text": "Which operator returns true if the two compared values are not equal?",
  "question_type": "multiple_choice",
  "points": 1,
  "explanation": "Reference: https://www.w3schools.com/js/js_operators.asp",
  "order_index": 1,
  "is_active": true,
  "source": "github",
  "external_reference": "https://github.com/your-org/quiz-content-source/blob/main/javascript/javascript-quiz.md#q1",
  "external_id": "javascript-q1",
  "question_code": null,
  "question_code_language": null,
  "question_image_url": null,
  "question_image_alt": null
}
```

**Content Block Fields:**
- `question_code`: Code snippet in the question (if present)
- `question_code_language`: Programming language (e.g., `javascript`, `java`, `python`)
- `question_image_url`: Image URL (local or external)
- `question_image_alt`: Alt text for image

**Mapping Rules:**
- Question number from `#### Q[N]` → `order_index`
- Reference links → `explanation`
- Default points: `1`
- `question_type`: Always `"multiple_choice"` (upstream format)

---

### Choice
Answer options for a question.

```json
{
  "question_id": "[QUESTION_ID]",
  "choice_text": "`!==`",
  "is_correct": true,
  "order_index": 4,
  "explanation": null,
  "choice_code": null,
  "choice_code_language": null,
  "choice_image_url": null,
  "choice_image_alt": null
}
```

**Content Block Fields:**
- `choice_code`: Code snippet in the option (if present)
- `choice_code_language`: Programming language
- `choice_image_url`: Image URL
- `choice_image_alt`: Alt text

**Mapping Rules:**
- `[x]` marker → `is_correct = true`
- `[ ]` marker → `is_correct = false`
- Option position (1-6) → `order_index`
- Code blocks within option → `choice_code`

---

## Question Type Classification

### Standard Question Type
All upstream questions are `multiple_choice` with 2-6 options.

**Validation Rules:**
- Must have at least 2 choices
- Must have exactly ONE correct answer (`[x]`)
- Maximum 6 choices (most have 4)

---

## External ID Convention

Format: `{topic-slug}-q{question-number}`

**Examples:**
- `javascript-q1`
- `css-q15`
- `android-q42`

**Purpose:**
- Prevent duplicate imports
- Track question updates
- Enable incremental syncing

---

## Sync Strategy

### Initial Import
1. Parse markdown file
2. Extract all questions sequentially
3. Create topic → quiz → questions → choices
4. Set `last_synced_at` timestamp

### Incremental Update
1. Check `last_synced_at` timestamp
2. Compare question count
3. Re-import if count differs
4. Update `external_id` based questions

### Deduplication
- Use `external_id` as unique constraint
- Skip questions with existing `external_id`
- Update if content hash differs

---

## Data Integrity Rules

1. **Topic Uniqueness**: `slug` must be unique
2. **Quiz Uniqueness**: One quiz per topic file
3. **Question Order**: `order_index` must match source order
4. **Answer Correctness**: Exactly one `is_correct = true` per question
5. **External References**: Store GitHub file path for traceability

---

## Content Preservation

### What to Preserve
- Original question text (exact copy)
- Original code formatting
- Original explanations/references
- Question numbering

### What to Transform
- Markdown image syntax → `image_url` field
- Markdown code blocks → separate `code` fields
- Backticks in options → clean text + code field
- Relative paths → absolute URLs

---

## Metadata Tracking

Track import metadata for each question:

```json
{
  "source": "github",
  "external_reference": "https://github.com/.../javascript-quiz.md",
  "external_id": "javascript-q1",
  "last_synced_at": "2026-01-22T00:00:00Z"
}
```

**Purpose:**
- Audit trail
- Update detection
- Source attribution
- Quality tracking
