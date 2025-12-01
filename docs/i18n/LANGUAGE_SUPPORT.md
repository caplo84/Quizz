# Language Support

## Overview
The upstream community quiz repository includes quiz translations in multiple languages. This document outlines supported languages and implementation strategy.

---

## Supported Languages

### Primary Language
- **English (en)** - Default, most complete

### Translation Languages Available

| Language | Code | Flag | Coverage | Examples |
|----------|------|------|----------|----------|
| Spanish | es | 🇪🇸 | ~70% | javascript-quiz-es.md |
| French | fr | 🇫🇷 | ~50% | css-quiz-fr.md |
| Italian | it | 🇮🇹 | ~50% | android-quiz-it.md |
| Chinese | ch | 🇨🇳 | ~40% | python-quiz-ch.md |
| German | de | 🇩🇪 | ~30% | angular-quiz-de.md |
| Hindi | hi | 🇮🇳 | ~15% | java-quiz-hi.md |
| Ukrainian | ua | 🇺🇦 | ~10% | javascript-quiz-ua.md |
| Vietnamese | vi | 🇻🇳 | ~5% | accounting-quiz-vi.md |

**Note:** Coverage varies by topic. Popular topics (JavaScript, CSS, HTML) have more translations.

---

## File Naming Convention

### Pattern
```
{topic-name}/{topic-name}-quiz-{language-code}.md
```

### Examples
```
javascript/javascript-quiz.md       # English (default)
javascript/javascript-quiz-es.md    # Spanish
javascript/javascript-quiz-fr.md    # French
css/css-quiz.md                     # English
css/css-quiz-ch.md                  # Chinese
android/android-quiz-hi.md          # Hindi
```

---

## Translation Coverage by Topic

### High Coverage Topics (3+ Languages)
- JavaScript: en, es, fr, it, ch, ua
- CSS: en, es, fr, it, ch, ua
- HTML: en, es, ua
- Angular: en, es, fr, it, ch, hi, ua
- Android: en, es, fr, it, ch, hi, ua
- Git: en, fr, ua
- C++: en, es, fr, it, ch
- Python: en, (various partial)

### Medium Coverage (1-2 Languages)
- Most programming languages
- Popular frameworks

### Low Coverage (English Only)
- Specialized tools
- New topics
- Domain-specific topics

---

## Database Schema for Translations

### Option 1: Separate Translation Tables (Recommended)

```sql
CREATE TABLE question_translations (
    id SERIAL PRIMARY KEY,
    question_id INTEGER REFERENCES questions(id) ON DELETE CASCADE,
    language_code VARCHAR(5) NOT NULL,
    question_text TEXT NOT NULL,
    explanation TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(question_id, language_code)
);

CREATE TABLE choice_translations (
    id SERIAL PRIMARY KEY,
    choice_id INTEGER REFERENCES choices(id) ON DELETE CASCADE,
    language_code VARCHAR(5) NOT NULL,
    choice_text TEXT NOT NULL,
    explanation TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(choice_id, language_code)
);
```

**Advantages:**
- Clean separation of concerns
- Easy to add/remove languages
- Efficient queries
- No schema changes needed for new languages

---

### Option 2: JSONB Column

```sql
ALTER TABLE questions 
ADD COLUMN translations JSONB DEFAULT '{}';

-- Example data
{
  "es": {
    "question_text": "[Spanish translation] Which operator returns true when values are not equal?",
    "explanation": "[Spanish translation] Reference: ..."
  },
  "fr": {
    "question_text": "[French translation] Which operator returns true when values are not equal?",
    "explanation": "[French translation] Reference: ..."
  }
}
```

**Advantages:**
- Simpler schema
- Single query for all translations
- Flexible structure

**Disadvantages:**
- Harder to query specific languages
- JSON indexing complexity

---

## Import Strategy

### Phase 1: English Only (Current)
1. Import English questions as default
2. Store in main tables
3. Set `language = 'en'` (optional metadata)

### Phase 2: Add Translations
1. Parse translated markdown files
2. Match questions by `order_index` and topic
3. Store translations in `*_translations` tables
4. Link via `question_id` / `choice_id`

### Phase 3: Sync Translations
1. Detect new translations in repository
2. Match by question order/external_id
3. Update or insert translations

---

## Translation Matching Algorithm

### Challenge
Translated questions may have slight variations in numbering or content.

### Matching Strategy

```python
def match_translation(english_question: Question, translated_file: str, language: str):
    """Match English question to translated version."""
    
    # Parse translated file
    translated_questions = parse_markdown_file(translated_file)
    
    # Try to match by question number
    for trans_q in translated_questions:
        if trans_q.order_index == english_question.order_index:
            # Found potential match
            
            # Validate choice count matches
            if len(trans_q.choices) == len(english_question.choices):
                # Store translation
                store_translation(
                    question_id=english_question.id,
                    language=language,
                    question_text=trans_q.question_text,
                    choices=trans_q.choices
                )
                return True
    
    # No match found
    return False
```

---

## API Design for Translations

### Endpoint 1: Get Question with Translation
```http
GET /api/questions/{id}?lang=es
```

**Response:**
```json
{
  "id": 123,
  "question_text": "[Spanish translation] Which operator returns true?",
  "language": "es",
  "choices": [
    {"id": 1, "choice_text": "`<>`", "is_correct": false},
    {"id": 2, "choice_text": "`!==`", "is_correct": true}
  ],
  "fallback": false
}
```

### Endpoint 2: Get Available Languages
```http
GET /api/quizzes/{slug}/languages
```

**Response:**
```json
{
  "quiz_slug": "javascript-assessment",
  "available_languages": ["en", "es", "fr", "it", "ua"],
  "default_language": "en",
  "coverage": {
    "en": 100,
    "es": 100,
    "fr": 100,
    "it": 85,
    "ua": 75
  }
}
```

---

## Frontend Implementation

### Language Selector Component

```jsx
function LanguageSelector({ quizSlug, currentLang, onLanguageChange }) {
  const { data: languages } = useQuery(['languages', quizSlug], () =>
    fetch(`/api/quizzes/${quizSlug}/languages`).then(r => r.json())
  );

  return (
    <select value={currentLang} onChange={(e) => onLanguageChange(e.target.value)}>
      {languages.available_languages.map(lang => (
        <option key={lang} value={lang}>
          {getLanguageName(lang)} ({languages.coverage[lang]}%)
        </option>
      ))}
    </select>
  );
}
```

### Load Questions with Language

```jsx
function QuizPage() {
  const [language, setLanguage] = useState('en');
  const { slug } = useParams();

  const { data: questions } = useQuery(
    ['questions', slug, language],
    () => fetch(`/api/quizzes/${slug}/questions?lang=${language}`).then(r => r.json())
  );

  return (
    <div>
      <LanguageSelector 
        quizSlug={slug} 
        currentLang={language}
        onLanguageChange={setLanguage}
      />
      <QuizQuestions questions={questions} />
    </div>
  );
}
```

---

## Translation Quality

### Quality Indicators

| Quality Level | Criteria |
|--------------|----------|
| ⭐⭐⭐ High | Professional translation, all questions, verified |
| ⭐⭐ Medium | Community translation, most questions, reviewed |
| ⭐ Low | Machine translation, partial coverage, unverified |

**Source Quality:** Upstream translations are community-contributed, quality varies.

---

## Translation Metadata

Store metadata about each translation:

```sql
CREATE TABLE translation_metadata (
    id SERIAL PRIMARY KEY,
    quiz_id INTEGER REFERENCES quizzes(id),
    language_code VARCHAR(5),
    coverage_percent DECIMAL(5,2),
    last_updated TIMESTAMP,
    source VARCHAR(50), -- 'community', 'professional', 'machine'
    verified BOOLEAN DEFAULT false,
    UNIQUE(quiz_id, language_code)
);
```

---

## Translation Update Strategy

### Detecting Updates

```python
def check_translation_updates(topic: str, language: str):
    """Check if translation file has been updated."""
    
    # Get last sync timestamp
    last_sync = get_last_sync_timestamp(topic, language)
    
    # Check GitHub file last modified date
    github_last_modified = get_github_file_date(f"{topic}/{topic}-quiz-{language}.md")
    
    if github_last_modified > last_sync:
        # Translation has been updated
        return True
    
    return False
```

### Update Process

1. Detect updated translation files
2. Re-parse updated file
3. Match questions by order/external_id
4. Update translations in database
5. Update `last_synced_at` timestamp

---

## Special Considerations

### Code Blocks
**Do NOT translate:**
- Code syntax
- Programming keywords
- Variable names (unless in comments)

**Example:**
```javascript
// English
function addTax(total) {
  return total * 1.05;
}

// Spanish - code stays the same
function addTax(total) {
  return total * 1.05;
}
```

### Technical Terms
Some terms should remain in English:
- API names (e.g., `RecyclerView`, `useState`)
- File names
- URLs

### Images
Image URLs remain the same across translations (images are language-agnostic or include text overlays).

---

## Priority Languages for Implementation

### Phase 1: High Priority
1. **Spanish (es)** - High coverage, large user base
2. **French (fr)** - Good coverage, popular

### Phase 2: Medium Priority
3. **Italian (it)** - Decent coverage
4. **Chinese (ch)** - Growing coverage

### Phase 3: Lower Priority
5. German, Hindi, Ukrainian - Partial coverage
6. Others as available

---

## Translation Contribution Workflow

For future community translations:

1. Use English as source
2. Translate question_text and choice_text
3. Keep code blocks unchanged
4. Keep image URLs unchanged
5. Translate explanations/references if possible
6. Submit via pull request to GitHub
7. Auto-sync to database
