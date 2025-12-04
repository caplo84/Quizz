# Question Quality Guidelines

## Overview
This document defines quality standards for importing questions from the upstream community quiz repository and guidelines for maintaining data quality.

---

## Quality Dimensions

### 1. Content Accuracy
- Question text is clear and unambiguous
- Exactly one correct answer
- Incorrect options are plausibly wrong (not obviously incorrect)
- Technical information is accurate

### 2. Format Validity
- Markdown properly formatted
- Code blocks have language specified
- Images load successfully
- No broken links

### 3. Clarity
- Question is understandable
- No grammatical errors
- Technical terms used correctly
- Context provided when needed

### 4. Completeness
- All required fields present
- Explanations/references included (when available)
- Code examples complete and runnable
- Image alt text provided

### 5. Consistency
- Follows repository formatting standards
- Numbering is sequential
- Option count is appropriate (2-6)
- Difficulty is appropriate for topic

---

## Quality Checks

### Automated Validation

#### 1. Structural Validation

```python
def validate_question_structure(question: Dict) -> List[str]:
    """Validate question meets structural requirements."""
    errors = []
    
    # Check required fields
    if not question.get('question_text'):
        errors.append("Missing question_text")
    
    # Check question text length
    if len(question.get('question_text', '')) < 10:
        errors.append("Question text too short (< 10 chars)")
    
    # Check choices
    choices = question.get('choices', [])
    if len(choices) < 2:
        errors.append("Too few choices (minimum 2)")
    if len(choices) > 6:
        errors.append("Too many choices (maximum 6)")
    
    # Check correct answer
    correct_count = sum(1 for c in choices if c.get('is_correct'))
    if correct_count == 0:
        errors.append("No correct answer marked")
    elif correct_count > 1:
        errors.append(f"Multiple correct answers ({correct_count})")
    
    # Check choice text
    for i, choice in enumerate(choices):
        if not choice.get('choice_text'):
            errors.append(f"Choice {i+1} has no text")
    
    return errors
```

#### 2. Content Validation

```python
def validate_question_content(question: Dict) -> List[str]:
    """Validate question content quality."""
    warnings = []
    
    # Check for code language specification
    if question.get('question_code') and not question.get('question_code_language'):
        warnings.append("Code block missing language specification")
    
    # Check image URLs
    if question.get('question_image_url'):
        if not is_valid_url(question['question_image_url']):
            warnings.append("Invalid image URL")
    
    # Check for explanations on complex questions
    if has_code_blocks(question) and not question.get('explanation'):
        warnings.append("Complex question missing explanation")
    
    # Check choice length variance
    choice_lengths = [len(c['choice_text']) for c in question.get('choices', [])]
    if max(choice_lengths) / min(choice_lengths) > 10:
        warnings.append("Choice lengths vary significantly")
    
    return warnings
```

#### 3. Code Quality Validation

```python
def validate_code_blocks(question: Dict) -> List[str]:
    """Validate code block quality."""
    warnings = []
    
    if question.get('question_code'):
        code = question['question_code']
        language = question.get('question_code_language')
        
        # Check for common syntax errors
        if language == 'javascript':
            if 'function' in code and not code.strip().endswith('}'):
                warnings.append("JavaScript code may be incomplete")
        
        # Check for very long code blocks
        if len(code) > 1000:
            warnings.append("Code block very long (>1000 chars)")
        
        # Check for proper indentation
        if not has_consistent_indentation(code):
            warnings.append("Code indentation inconsistent")
    
    return warnings
```

---

## Quality Scoring

### Quality Score Algorithm

```python
def calculate_quality_score(question: Dict, errors: List[str], warnings: List[str]) -> float:
    """Calculate quality score 0-100."""
    
    score = 100.0
    
    # Deduct for errors
    score -= len(errors) * 15
    
    # Deduct for warnings
    score -= len(warnings) * 5
    
    # Bonus for completeness
    if question.get('explanation'):
        score += 5
    if question.get('source') == 'github':
        score += 2
    
    # Bonus for references
    if question.get('explanation') and 'http' in question['explanation']:
        score += 3
    
    # Cap at 0-100
    return max(0, min(100, score))
```

### Quality Levels

| Score | Level | Action |
|-------|-------|--------|
| 90-100 | Excellent | Import immediately |
| 70-89 | Good | Import with minor review |
| 50-69 | Fair | Import but flag for improvement |
| 30-49 | Poor | Review before import |
| 0-29 | Very Poor | Skip import, log for manual review |

---

## Quality Flags

### Flag Types

```sql
CREATE TABLE quality_flags (
    id SERIAL PRIMARY KEY,
    question_id INTEGER REFERENCES questions(id),
    flag_type VARCHAR(50),
    severity VARCHAR(20), -- 'error', 'warning', 'info'
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved BOOLEAN DEFAULT FALSE
);
```

### Common Flags

| Flag Type | Severity | Description |
|-----------|----------|-------------|
| `multiple_correct` | error | Multiple answers marked correct |
| `no_correct` | error | No correct answer marked |
| `missing_language` | warning | Code block without language |
| `broken_image` | warning | Image URL returns 404 |
| `too_short` | warning | Question < 10 characters |
| `no_explanation` | info | Missing reference/explanation |
| `ambiguous_options` | warning | Similar option text |
| `incomplete_code` | warning | Code appears incomplete |

---

## Import Quality Gates

### Gate 1: Critical Errors
**Block import if:**
- No question text
- No correct answer
- Multiple correct answers
- < 2 choices

**Action:** Skip question, log error

---

### Gate 2: Validation Warnings
**Allow import but flag if:**
- Missing code language
- No explanation
- Broken image URL
- Suspicious formatting

**Action:** Import with quality flag

---

### Gate 3: Content Review
**Manual review needed if:**
- Quality score < 50
- More than 5 warnings
- Unusual choice count (1 or >6)
- Duplicate detection triggered

**Action:** Queue for review

---

## Duplicate Detection

### Algorithm

```python
def detect_duplicates(new_question: Dict, existing_questions: List[Dict]) -> Optional[int]:
    """Detect if question is duplicate of existing question."""
    
    new_text = normalize_text(new_question['question_text'])
    
    for existing in existing_questions:
        existing_text = normalize_text(existing['question_text'])
        
        # Exact match
        if new_text == existing_text:
            return existing['id']
        
        # High similarity
        similarity = calculate_similarity(new_text, existing_text)
        if similarity > 0.90:  # 90% similar
            return existing['id']
    
    return None

def normalize_text(text: str) -> str:
    """Normalize text for comparison."""
    # Remove markdown formatting
    text = re.sub(r'[`*_\[\]()]', '', text)
    # Lowercase
    text = text.lower()
    # Remove extra whitespace
    text = re.sub(r'\s+', ' ', text).strip()
    return text
```

---

## Quality Improvement Workflow

### 1. Identify Low-Quality Questions

```sql
SELECT 
    q.id,
    q.question_text,
    COUNT(qf.id) as flag_count,
    q.quality_score
FROM questions q
LEFT JOIN quality_flags qf ON q.id = qf.question_id AND qf.resolved = FALSE
WHERE q.quality_score < 70
GROUP BY q.id
ORDER BY q.quality_score ASC, flag_count DESC
LIMIT 50;
```

### 2. Review and Fix

- Manual review by content team
- Fix formatting issues
- Verify correct answers
- Add missing explanations
- Update image URLs

### 3. Re-validate

```python
def revalidate_question(question_id: int):
    """Re-run validation after fixes."""
    question = get_question(question_id)
    errors = validate_question_structure(question)
    warnings = validate_question_content(question)
    
    new_score = calculate_quality_score(question, errors, warnings)
    update_quality_score(question_id, new_score)
    
    # Clear resolved flags
    if new_score > 70:
        clear_quality_flags(question_id)
```

---

## Source Quality Assessment

### GitHub Repository Quality

**Pros:**
- Community-reviewed
- Many contributors
- Regular updates
- Version controlled

**Cons:**
- Varying quality by contributor
- No formal review process
- Occasional errors/typos
- Inconsistent formatting

**Overall Rating:** 7/10 (Good)

---

## Quality Monitoring

### Metrics to Track

```sql
-- Quality score distribution
SELECT 
    CASE 
        WHEN quality_score >= 90 THEN 'Excellent'
        WHEN quality_score >= 70 THEN 'Good'
        WHEN quality_score >= 50 THEN 'Fair'
        ELSE 'Poor'
    END as quality_level,
    COUNT(*) as question_count,
    ROUND(AVG(quality_score), 2) as avg_score
FROM questions
GROUP BY quality_level
ORDER BY MIN(quality_score) DESC;

-- Most common quality flags
SELECT 
    flag_type,
    COUNT(*) as count,
    AVG(CASE WHEN resolved THEN 1 ELSE 0 END) as resolution_rate
FROM quality_flags
GROUP BY flag_type
ORDER BY count DESC;

-- Questions needing review
SELECT COUNT(*) as needs_review
FROM questions
WHERE quality_score < 50 OR 
      (SELECT COUNT(*) FROM quality_flags qf 
       WHERE qf.question_id = questions.id AND qf.severity = 'error') > 0;
```

---

## Quality Improvement Priorities

### High Priority
1. Fix questions with no correct answer
2. Fix questions with multiple correct answers
3. Resolve broken image URLs
4. Add missing code languages

### Medium Priority
5. Add explanations to complex questions
6. Fix formatting inconsistencies
7. Improve choice text clarity
8. Validate code syntax

### Low Priority
9. Optimize question text length
10. Add tags/metadata
11. Enhance images
12. Translate to more languages

---

## User Feedback Integration

### Quality Reporting

Allow users to report quality issues:

```sql
CREATE TABLE user_quality_reports (
    id SERIAL PRIMARY KEY,
    question_id INTEGER REFERENCES questions(id),
    user_id INTEGER,
    report_type VARCHAR(50), -- 'wrong_answer', 'unclear', 'broken_link', 'typo'
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending' -- 'pending', 'reviewed', 'resolved'
);
```

**API Endpoint:**
```http
POST /api/questions/{id}/report
{
  "report_type": "wrong_answer",
  "description": "Option B should be correct, not C"
}
```

---

## Continuous Quality Improvement

### Weekly Quality Review Process

1. **Monday:** Generate quality report
   - Questions with flags
   - User reports
   - Low scores

2. **Tuesday-Thursday:** Fix issues
   - Update questions
   - Add explanations
   - Fix images

3. **Friday:** Re-validate and deploy
   - Run validation suite
   - Clear resolved flags
   - Update quality scores

---

## Quality Benchmarks

### Target Metrics

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| Avg Quality Score | > 85 | TBD | - |
| Questions w/ Errors | < 1% | TBD | - |
| Questions w/ Explanations | > 60% | ~40% | ⚠️ |
| Broken Images | < 2% | TBD | - |
| Duplicate Questions | 0% | TBD | - |
| Code w/ Language | > 95% | ~90% | ⚠️ |

---

## Summary

**Quality is critical for:**
- User trust
- Learning effectiveness
- Platform reputation
- SEO rankings

**Key Actions:**
1. Implement automated validation
2. Score all questions
3. Flag low-quality content
4. Continuous improvement process
5. User feedback integration
