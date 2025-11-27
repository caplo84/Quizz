# Format Mapping Guide

## Overview
This guide provides step-by-step instructions for mapping upstream quiz markdown format to our database schema.

---

## Source to Schema Mapping

### Complete Flow

```
Upstream Markdown File
    ↓
Parse Topic/Quiz Metadata
    ↓
Extract Questions (#### Q[N])
    ↓
Parse Question Content
    ↓
Extract Choices (- [x]/[ ])
    ↓
Store in Database
```

---

## Level 1: Topic Mapping

### Source
```
Repository Folder: community-quiz-content-main/javascript/
```

### Extraction
```python
topic_slug = folder_name  # "javascript"
topic_name = folder_name.replace('-', ' ').title()  # "Javascript"
```

### Database Insert
```sql
INSERT INTO topics (name, slug, description, icon_url, is_active, created_at, updated_at)
VALUES (
    'JavaScript',
    'javascript',
    'JavaScript Programming Language',
    '/icons/javascript.svg',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
```

**Field Mapping:**

| Source | Target Field | Transformation |
|--------|-------------|----------------|
| Folder name | `name` | Title case |
| Folder name | `slug` | Lowercase, normalize |
| Derived | `description` | Template: "[Name] Assessment" |
| Manual/Config | `icon_url` | Map from config file |
| Default | `is_active` | `true` |

---

## Level 2: Quiz Mapping

### Source
```
File: javascript/javascript-quiz.md
README: | Javascript | 166 | 166 |
```

### Extraction
```python
quiz_title = f"{topic_name} Assessment"
quiz_slug = f"{topic_slug}-assessment"
question_count = count_questions_in_file(file_path)
external_ref = github_url(file_path)
```

### Database Insert
```sql
INSERT INTO quizzes (
    title, slug, description, topic_id, difficulty_level,
    time_limit_minutes, total_questions, is_active,
    source, external_reference, external_id, last_synced_at,
    created_at, updated_at
)
VALUES (
    'JavaScript Assessment',
    'javascript-assessment',
    'Test your JavaScript knowledge with community questions',
    (SELECT id FROM topics WHERE slug = 'javascript'),
    'medium',
    30,
    166,
    true,
    'github',
    'https://github.com/your-org/quiz-content-source/blob/main/javascript/javascript-quiz.md',
    'javascript-quiz-main',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
```

**Field Mapping:**

| Source | Target Field | Transformation |
|--------|-------------|----------------|
| Derived | `title` | "[Topic] Assessment" |
| Derived | `slug` | "[topic-slug]-assessment" |
| Template | `description` | Standard template |
| FK | `topic_id` | Lookup from topics |
| Default | `difficulty_level` | `'medium'` |
| Default | `time_limit_minutes` | `30` |
| Count | `total_questions` | Count `#### Q` markers |
| Default | `is_active` | `true` |
| Constant | `source` | `'github'` |
| File path | `external_reference` | GitHub URL |
| Derived | `external_id` | `{topic-slug}-quiz-main` |
| Now | `last_synced_at` | `CURRENT_TIMESTAMP` |

---

## Level 3: Question Mapping

### Source (Plain Text Example)
```markdown
#### Q1. Which operator returns true if the two compared values are not equal?

- [ ] `<>`
- [ ] `~`
- [ ] `==!`
- [x] `!==`

[Reference Javascript Comparison Operators](https://www.w3schools.com/js/js_operators.asp)
```

### Parsing Logic
```python
def parse_question(question_block: str, quiz_id: int, order_index: int) -> Question:
    # Extract question number and text
    match = re.match(r'#### Q(\d+)\. (.+)', question_block)
    question_num = int(match.group(1))
    question_text = match.group(2)
    
    # Extract optional reference
    explanation = None
    ref_match = re.search(r'\[Reference[^\]]*\]\(([^)]+)\)', question_block)
    if ref_match:
        explanation = f"Reference: {ref_match.group(1)}"
    
    return Question(
        quiz_id=quiz_id,
        question_text=question_text,
        question_type='multiple_choice',
        points=1,
        explanation=explanation,
        order_index=order_index,
        is_active=True,
        source='github',
        external_reference=f"{quiz_external_ref}#q{question_num}",
        external_id=f"{quiz_slug}-q{question_num}"
    )
```

### Database Insert
```sql
INSERT INTO questions (
    quiz_id, question_text, question_type, points, explanation,
    order_index, is_active, source, external_reference, external_id,
    created_at, updated_at
)
VALUES (
    (SELECT id FROM quizzes WHERE slug = 'javascript-assessment'),
    'Which operator returns true if the two compared values are not equal?',
    'multiple_choice',
    1,
    'Reference: https://www.w3schools.com/js/js_operators.asp',
    1,
    true,
    'github',
    'https://github.com/.../javascript-quiz.md#q1',
    'javascript-q1',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
```

**Field Mapping:**

| Source | Target Field | Transformation |
|--------|-------------|----------------|
| FK | `quiz_id` | Lookup from quizzes |
| `#### Q1. TEXT` | `question_text` | Extract after `Q1.` |
| Constant | `question_type` | `'multiple_choice'` |
| Default | `points` | `1` |
| `[Reference]` | `explanation` | Extract URL from markdown link |
| Question # | `order_index` | Use question number |
| Default | `is_active` | `true` |
| Constant | `source` | `'github'` |
| Derived | `external_reference` | Quiz URL + `#qN` |
| Derived | `external_id` | `{quiz-slug}-qN` |

---

## Level 4: Question with Code

### Source
````markdown
#### Q3. Review the code below. Which statement calls the addTax function?

```js
function addTax(total) {
  return total * 1.05;
}
```

- [ ] `addTax = 50;`
- [x] `addTax(50);`
````

### Parsing Logic
```python
def extract_code_blocks(text: str) -> List[Tuple[str, str]]:
    """Extract code blocks and their languages."""
    pattern = r'```(\w+)?\n(.*?)\n```'
    matches = re.findall(pattern, text, re.DOTALL)
    return [(lang or 'text', code.strip()) for lang, code in matches]

def parse_question_with_code(block: str, quiz_id: int, order: int) -> Question:
    # Split question text and options
    parts = block.split('\n- [')
    question_part = parts[0]
    
    # Extract text and code
    question_text = re.match(r'#### Q\d+\. (.+?)(?=\n```|\n-|$)', question_part, re.DOTALL).group(1).strip()
    
    code_blocks = extract_code_blocks(question_part)
    question_code = code_blocks[0][1] if code_blocks else None
    question_code_language = code_blocks[0][0] if code_blocks else None
    
    return Question(
        quiz_id=quiz_id,
        question_text=question_text,
        question_code=question_code,
        question_code_language=question_code_language,
        # ... other fields
    )
```

### Database Insert
```sql
INSERT INTO questions (
    quiz_id, question_text, question_type, points,
    question_code, question_code_language,
    order_index, is_active, source, external_reference, external_id,
    created_at, updated_at
)
VALUES (
    (SELECT id FROM quizzes WHERE slug = 'javascript-assessment'),
    'Review the code below. Which statement calls the addTax function?',
    'multiple_choice',
    1,
    E'function addTax(total) {\n  return total * 1.05;\n}',
    'js',
    3,
    true,
    'github',
    'https://github.com/.../javascript-quiz.md#q3',
    'javascript-q3',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
```

**Additional Fields:**

| Source | Target Field | Transformation |
|--------|-------------|----------------|
| ` ```lang ` | `question_code_language` | Extract language identifier |
| Code content | `question_code` | Extract text between ` ``` ` markers |

---

## Level 5: Choice Mapping

### Source
```markdown
- [ ] `<>`
- [ ] `~`
- [ ] `==!`
- [x] `!==`
```

### Parsing Logic
```python
def parse_choices(question_block: str, question_id: int) -> List[Choice]:
    # Extract option lines
    option_pattern = r'^- \[(x| )\] (.+)$'
    matches = re.findall(option_pattern, question_block, re.MULTILINE)
    
    choices = []
    for order, (marker, text) in enumerate(matches, start=1):
        choices.append(Choice(
            question_id=question_id,
            choice_text=text.strip(),
            is_correct=(marker == 'x'),
            order_index=order
        ))
    
    return choices
```

### Database Insert
```sql
INSERT INTO choices (
    question_id, choice_text, is_correct, order_index,
    created_at, updated_at
)
VALUES
    ((SELECT id FROM questions WHERE external_id = 'javascript-q1'), '`<>`', false, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM questions WHERE external_id = 'javascript-q1'), '`~`', false, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM questions WHERE external_id = 'javascript-q1'), '`==!`', false, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM questions WHERE external_id = 'javascript-q1'), '`!==`', true, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
```

**Field Mapping:**

| Source | Target Field | Transformation |
|--------|-------------|----------------|
| FK | `question_id` | Lookup from questions |
| After `] ` | `choice_text` | Extract text after marker |
| `[x]` vs `[ ]` | `is_correct` | `x` → `true`, space → `false` |
| Position | `order_index` | Sequential (1, 2, 3, 4) |

---

## Level 6: Choice with Code

### Source
````markdown
- [ ] A

```java
RecyclerView.Adapter<T extends BaseAdapter>
```

- [x] B

```java
RecyclerView.Adapter<VH extends ViewHolder>
```
````

### Parsing Logic
```python
def parse_choice_with_code(option_block: str) -> Choice:
    # Extract label
    label_match = re.match(r'- \[(x| )\] (.+)', option_block)
    is_correct = label_match.group(1) == 'x'
    choice_text = label_match.group(2).strip()
    
    # Extract code block if present
    code_match = re.search(r'```(\w+)?\n(.*?)\n```', option_block, re.DOTALL)
    choice_code = code_match.group(2).strip() if code_match else None
    choice_code_language = code_match.group(1) if code_match else None
    
    return Choice(
        choice_text=choice_text,
        is_correct=is_correct,
        choice_code=choice_code,
        choice_code_language=choice_code_language
    )
```

### Database Insert
```sql
INSERT INTO choices (
    question_id, choice_text, is_correct, order_index,
    choice_code, choice_code_language,
    created_at, updated_at
)
VALUES
    ((SELECT id FROM questions WHERE external_id = 'android-q5'), 'A', false, 1,
     'RecyclerView.Adapter<T extends BaseAdapter>', 'java',
     CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM questions WHERE external_id = 'android-q5'), 'B', true, 2,
     'RecyclerView.Adapter<VH extends ViewHolder>', 'java',
     CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
```

**Additional Fields:**

| Source | Target Field | Transformation |
|--------|-------------|----------------|
| Code block | `choice_code` | Extract between ` ``` ` |
| Code language | `choice_code_language` | Extract language identifier |

---

## Image URL Transformation

### Local Image Path
**Source:** `![img](image/43.jpeg)`

**Transform:**
```python
def transform_image_url(markdown_image: str, topic: str, file_path: str) -> str:
    # Parse markdown image
    match = re.match(r'!\[([^\]]*)\]\(([^)]+)\)', markdown_image)
    alt_text = match.group(1)
    image_path = match.group(2)
    
    # Check if relative path
    if not image_path.startswith('http'):
        # Build GitHub raw URL
        base_url = 'https://raw.githubusercontent.com/your-org/quiz-content-source/main'
        full_url = f"{base_url}/{topic}/{image_path}"
        return full_url, alt_text
    
    return image_path, alt_text
```

**Result:** `https://raw.githubusercontent.com/your-org/quiz-content-source/main/android/image/43.jpeg`

---

## Complete Example: End-to-End

### Input
````markdown
## JavaScript

#### Q1. Which operator returns true if values are not equal?

- [ ] `<>`
- [x] `!==`

[Reference](https://www.w3schools.com/js/js_operators.asp)
````

### Output
```python
{
    "topic": {
        "name": "JavaScript",
        "slug": "javascript"
    },
    "quiz": {
        "title": "JavaScript Assessment",
        "slug": "javascript-assessment",
        "topic_slug": "javascript",
        "external_id": "javascript-quiz-main"
    },
    "question": {
        "question_text": "Which operator returns true if values are not equal?",
        "question_type": "multiple_choice",
        "points": 1,
        "explanation": "Reference: https://www.w3schools.com/js/js_operators.asp",
        "order_index": 1,
        "external_id": "javascript-q1"
    },
    "choices": [
        {
            "choice_text": "`<>`",
            "is_correct": false,
            "order_index": 1
        },
        {
            "choice_text": "`!==`",
            "is_correct": true,
            "order_index": 2
        }
    ]
}
```

---

## Validation Checklist

After mapping, validate:

- [ ] One topic per folder
- [ ] One quiz per topic
- [ ] Sequential question numbering (Q1, Q2, Q3...)
- [ ] Each question has 2-6 choices
- [ ] Exactly ONE correct answer per question
- [ ] All `external_id` values are unique
- [ ] Code languages are valid
- [ ] Image URLs are accessible
- [ ] No null/empty required fields

---

## Error Handling

| Error Condition | Action |
|----------------|--------|
| Missing question text | Skip question, log error |
| No correct answer | Skip question, log error |
| Multiple correct answers | Skip question, log error |
| Invalid markdown | Try best-effort parse, flag for review |
| Duplicate `external_id` | Skip (already imported) |
| Missing code language | Default to `'text'` |
| Broken image URL | Import anyway, flag for review |
