# Fallback Strategy

## Overview
This document defines the fallback strategy when requested translations are not available.

---

## Fallback Hierarchy

### Level 1: Requested Language
User requests quiz in Spanish (`lang=es`)

```
1. Check if Spanish translation exists for this question
   ↓ Yes → Return Spanish version
   ↓ No → Go to Level 2
```

### Level 2: Default Language (English)
Spanish translation not found

```
2. Return English version (default)
   ↓ Mark response with fallback flag
   ↓ Log missing translation
```

### Level 3: Partial Fallback
Some questions translated, others not

```
3. Return mixed response
   ↓ Translated questions in Spanish
   ↓ Untranslated questions in English
   ↓ Flag indicates mixed language
```

---

## Implementation

### Database Query with Fallback

```sql
-- Get question with translation fallback
SELECT 
    q.id,
    q.question_text as english_text,
    COALESCE(qt.question_text, q.question_text) as display_text,
    CASE 
        WHEN qt.question_text IS NOT NULL THEN :requested_lang
        ELSE 'en'
    END as language_used,
    qt.question_text IS NOT NULL as is_translated
FROM questions q
LEFT JOIN question_translations qt 
    ON q.id = qt.question_id 
    AND qt.language_code = :requested_lang
WHERE q.quiz_id = :quiz_id
ORDER BY q.order_index;
```

**Parameters:**
- `:requested_lang` - e.g., `'es'`
- `:quiz_id` - Quiz ID

**Result:**
```json
{
  "id": 123,
  "english_text": "Which operator returns true?",
  "display_text": "[Spanish translation] Which operator returns true?",
  "language_used": "es",
  "is_translated": true
}
```

**Fallback Result:**
```json
{
  "id": 124,
  "english_text": "What is the correct syntax?",
  "display_text": "What is the correct syntax?",
  "language_used": "en",
  "is_translated": false
}
```

---

## API Response Format

### Complete Translation

**Request:**
```http
GET /api/questions/123?lang=es
```

**Response:**
```json
{
  "id": 123,
  "question_text": "[Spanish translation] Which operator returns true?",
  "language": "es",
  "fallback": false,
  "choices": [...]
}
```

---

### Fallback to English

**Request:**
```http
GET /api/questions/124?lang=es
```

**Response:**
```json
{
  "id": 124,
  "question_text": "What is the correct syntax?",
  "language": "en",
  "fallback": true,
  "requested_language": "es",
  "choices": [...]
}
```

**Note:** `fallback: true` indicates English was returned instead of Spanish.

---

### Mixed Language Quiz

**Request:**
```http
GET /api/quizzes/javascript-assessment/questions?lang=es
```

**Response:**
```json
{
  "quiz_slug": "javascript-assessment",
  "requested_language": "es",
  "translation_coverage": 85,
  "questions": [
    {
      "id": 1,
      "question_text": "[Spanish translation] Which operator returns true?",
      "language": "es",
      "fallback": false
    },
    {
      "id": 2,
      "question_text": "What is the correct syntax?",
      "language": "en",
      "fallback": true
    },
    {
      "id": 3,
      "question_text": "[Spanish translation] What is the difference?",
      "language": "es",
      "fallback": false
    }
  ]
}
```

---

## Frontend Handling

### Display Fallback Indicator

```jsx
function QuestionCard({ question }) {
  return (
    <div className="question-card">
      {question.fallback && (
        <div className="fallback-notice">
          ⚠️ This question is shown in English (translation not available)
        </div>
      )}
      
      <h2>{question.question_text}</h2>
      {/* ... */}
    </div>
  );
}
```

---

### Language Coverage Warning

```jsx
function LanguageSelector({ quizSlug, currentLang, coverage }) {
  const showWarning = coverage[currentLang] < 100;
  
  return (
    <div>
      <select value={currentLang} onChange={handleChange}>
        {/* language options */}
      </select>
      
      {showWarning && (
        <p className="coverage-warning">
          ⚠️ {coverage[currentLang]}% of questions available in {languageName}.
          Missing translations will be shown in English.
        </p>
      )}
    </div>
  );
}
```

---

## Fallback Preferences

### User Preference

Allow users to configure fallback behavior:

```typescript
interface FallbackPreferences {
  allowFallback: boolean;          // Allow English fallback?
  showFallbackWarning: boolean;    // Show warning indicators?
  fallbackLanguage: 'en' | 'es';   // Fallback to which language?
}
```

### Example Scenarios

**Scenario 1: Strict Language Mode**
```json
{
  "allowFallback": false,
  "preferredLanguage": "es"
}
```
**Behavior:** Skip questions without Spanish translation OR show error

**Scenario 2: Flexible Mode (Default)**
```json
{
  "allowFallback": true,
  "showFallbackWarning": true,
  "preferredLanguage": "es",
  "fallbackLanguage": "en"
}
```
**Behavior:** Show English for missing translations with warning

**Scenario 3: Silent Fallback**
```json
{
  "allowFallback": true,
  "showFallbackWarning": false,
  "preferredLanguage": "fr",
  "fallbackLanguage": "en"
}
```
**Behavior:** Show English without indicators

---

## Choice-Level Fallback

### Problem
Question translated but some choices are not.

### Solution
Fallback at choice level:

```sql
SELECT 
    c.id,
    c.choice_text as english_text,
    COALESCE(ct.choice_text, c.choice_text) as display_text,
    ct.choice_text IS NOT NULL as is_translated
FROM choices c
LEFT JOIN choice_translations ct 
    ON c.id = ct.choice_id 
    AND ct.language_code = :requested_lang
WHERE c.question_id = :question_id
ORDER BY c.order_index;
```

---

## Fallback Logging

### Track Fallback Usage

```sql
CREATE TABLE translation_fallback_log (
    id SERIAL PRIMARY KEY,
    question_id INTEGER,
    requested_language VARCHAR(5),
    fallback_language VARCHAR(5),
    user_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Log each fallback event
INSERT INTO translation_fallback_log (question_id, requested_language, fallback_language)
VALUES (124, 'es', 'en');
```

**Purpose:**
- Identify most needed translations
- Prioritize translation efforts
- Analytics on language usage

---

## Translation Priority Algorithm

Use fallback logs to prioritize translation work:

```sql
-- Questions most frequently falling back
SELECT 
    q.id,
    q.question_text,
    COUNT(*) as fallback_count,
    tfl.requested_language
FROM translation_fallback_log tfl
JOIN questions q ON tfl.question_id = q.id
WHERE tfl.created_at > NOW() - INTERVAL '30 days'
GROUP BY q.id, q.question_text, tfl.requested_language
ORDER BY fallback_count DESC
LIMIT 20;
```

**Result:** Top 20 questions that need translation most urgently.

---

## Partial Content Fallback

### Code Blocks
**Rule:** Never translate code, always show original

```json
{
  "question_text": "[Spanish translation] What does this code do?",  // Translated
  "question_code": "function test() {}",       // NOT translated
  "question_code_language": "javascript"
}
```

### Images
**Rule:** Images are language-agnostic, no fallback needed

```json
{
  "question_text": "[French translation] Which shape?",           // Translated
  "question_image_url": "image/shape.png"     // Same image
}
```

---

## Regional Language Fallbacks

### Dialect Fallback Chain

**Example:** User requests Brazilian Portuguese (`pt-BR`)

```
1. Try pt-BR (Brazilian Portuguese)
   ↓ Not found
2. Try pt (Portuguese)
   ↓ Not found
3. Fallback to en (English)
```

**Implementation:**
```python
def get_fallback_chain(language_code: str) -> List[str]:
    """Get fallback chain for a language."""
    
    fallback_map = {
        'pt-BR': ['pt-BR', 'pt', 'en'],
        'es-MX': ['es-MX', 'es', 'en'],
        'zh-CN': ['zh-CN', 'ch', 'zh', 'en'],
        'zh-TW': ['zh-TW', 'ch', 'zh', 'en'],
    }
    
    return fallback_map.get(language_code, [language_code, 'en'])
```

---

## Performance Optimization

### Caching Fallback Results

```typescript
interface CachedTranslation {
  questionId: number;
  language: string;
  text: string;
  isFallback: boolean;
  cachedAt: Date;
}
```

**Strategy:**
1. Cache translated questions (TTL: 7 days)
2. Cache fallback responses (TTL: 1 day)
3. Invalidate on translation update

---

## Error Handling

### No Translation Available

**Scenario:** Language requested but no translations exist for entire quiz.

**Response:**
```json
{
  "error": "translation_not_available",
  "message": "This quiz is not available in Spanish yet. Would you like to try it in English?",
  "available_languages": ["en"],
  "requested_language": "es"
}
```

---

## Accessibility Considerations

### Screen Reader Support

Announce language changes:

```html
<div lang="es">
  <p>[Spanish translation] Which operator returns true?</p>
</div>

<div lang="en" aria-label="This question is shown in English">
  <p>What is the correct syntax?</p>
</div>
```

---

## Summary

### Fallback Strategy Decision Tree

```
User requests question in language X
    ↓
Is question available in language X?
    ↓ Yes → Return in language X
    ↓ No → Check fallback settings
        ↓
    Allow fallback?
        ↓ Yes → Return in English with fallback flag
        ↓ No → Return error or skip question
```

### Key Principles

1. **Always provide content** - Better to show English than nothing
2. **Be transparent** - Flag fallbacks clearly
3. **Track usage** - Log fallbacks for improvement
4. **Progressive enhancement** - Start with English, add translations incrementally
5. **User choice** - Let users configure fallback behavior
