# Question Type Guide

## Overview
This guide documents the different question formats found in the upstream community quiz repository and how to classify them in our system.

---

## Question Type Classification

### All Questions Are Multiple Choice
The upstream quiz format **only uses multiple choice questions** with 2-6 options.

**Database Value:** `question_type = 'multiple_choice'`

---

## Question Format Categories

While all questions are technically "multiple choice", they vary significantly in format. We categorize them for parsing and rendering purposes.

### 1. **Plain Text Questions**
Pure text question with text answers.

**Frequency:** ~40% of all questions

**Example:**
```markdown
#### Q1. Which operator returns true if values are not equal?

- [ ] `<>`
- [ ] `~`
- [ ] `==!`
- [x] `!==`
```

**Characteristics:**
- Question: Plain text only
- Options: Short text (may include inline code)
- No images or code blocks

---

### 2. **Code-Based Questions**
Question includes code snippet(s) for context.

**Frequency:** ~35% of all questions

**Example:**
````markdown
#### Q3. Review the code below. Which statement calls the addTax function?

```js
function addTax(total) {
  return total * 1.05;
}
```

- [ ] `addTax = 50;`
- [x] `addTax(50);`
- [ ] `return addTax 50;`
````

**Characteristics:**
- Question: Text + code block(s)
- Options: Usually text
- Code provides context for the question

**Sub-types:**
- **Single language:** One code block
- **Multi-language:** Multiple code blocks (e.g., HTML + CSS)
- **Code debugging:** Find the error in code

---

### 3. **Code Option Questions**
Options are code snippets instead of text.

**Frequency:** ~15% of all questions

**Example:**
````markdown
#### Q5. What is the correct syntax?

- [ ] A
```java
RecyclerView.Adapter<T extends BaseAdapter>
```

- [x] B
```java
RecyclerView.Adapter<VH extends ViewHolder>
```
````

**Characteristics:**
- Question: Text or text + code
- Options: Labeled (A, B, C, D) with code blocks
- Each option is a code alternative

---

### 4. **Image-Based Questions**
Question includes an image for visual context.

**Frequency:** ~5% of all questions

**Example:**
```markdown
#### Q45. Which drawable definition creates this shape?

![img](image/43.jpeg)

- [ ] Option A
- [x] Option B
- [ ] Option C
```

**Characteristics:**
- Question: Text + image
- Options: Text or code
- Image shows the desired output/result

---

### 5. **Complex Multi-Part Questions**
Options contain detailed explanations or lists.

**Frequency:** ~5% of all questions

**Example:**
````markdown
#### Q17. Which set of statements is true?

- [x] A
```
  1. Larger z-index values appear on top.
  2. Negative and positive numbers can be used.
  3. The z-index can be used only on positioned elements.
```

- [ ] B
```
  1. Smaller z-index values appear on top.
  2. Only positive numbers can be used.
```
````

**Characteristics:**
- Question: Text
- Options: Multi-line explanations or numbered lists
- Each option is a complete statement set

---

## Parsing Strategy by Type

### Type Detection Algorithm

```python
def detect_question_format(question_md: str) -> str:
    """Detect question format for parsing strategy."""
    
    # Extract question and options sections
    question_part, options_part = split_question_and_options(question_md)
    
    # Check question characteristics
    has_question_code = '```' in question_part
    has_question_image = '![' in question_part
    
    # Check option characteristics
    has_option_code = '```' in options_part
    has_lettered_options = re.search(r'^- \[[x ]\] [A-F]$', options_part, re.MULTILINE)
    
    # Classify
    if has_question_image:
        return 'image_based'
    elif has_option_code and has_lettered_options:
        return 'code_options'
    elif has_question_code:
        return 'code_question'
    elif has_option_code:
        return 'complex_multipart'
    else:
        return 'plain_text'
```

---

## Rendering Strategy by Type

### Plain Text Questions
**Rendering:** Simple, minimal UI

```jsx
<QuestionCard>
  <QuestionText>{question.question_text}</QuestionText>
  <AnswerGrid>
    {choices.map(choice => (
      <AnswerOption text={choice.choice_text} />
    ))}
  </AnswerGrid>
</QuestionCard>
```

---

### Code-Based Questions
**Rendering:** Question code block + text options

```jsx
<QuestionCard>
  <QuestionText>{question.question_text}</QuestionText>
  {question.question_code && (
    <CodeBlock 
      code={question.question_code}
      language={question.question_code_language}
    />
  )}
  <AnswerGrid>
    {choices.map(choice => (
      <AnswerOption text={choice.choice_text} />
    ))}
  </AnswerGrid>
</QuestionCard>
```

---

### Code Option Questions
**Rendering:** Compact question + large code option blocks

```jsx
<QuestionCard>
  <QuestionText>{question.question_text}</QuestionText>
  <AnswerGrid layout="vertical"> {/* More space for code */}
    {choices.map(choice => (
      <AnswerOption>
        <span className="option-label">{choice.choice_text}</span>
        <CodeBlock code={choice.choice_code} />
      </AnswerOption>
    ))}
  </AnswerGrid>
</QuestionCard>
```

---

### Image-Based Questions
**Rendering:** Prominent image display + options below

```jsx
<QuestionCard>
  <QuestionText>{question.question_text}</QuestionText>
  <ImageDisplay 
    src={question.question_image_url}
    alt={question.question_image_alt}
    maxWidth="600px"
  />
  <AnswerGrid>
    {choices.map(choice => (
      <AnswerOption text={choice.choice_text} />
    ))}
  </AnswerGrid>
</QuestionCard>
```

---

### Complex Multi-Part Questions
**Rendering:** Expanded option space for readability

```jsx
<QuestionCard>
  <QuestionText>{question.question_text}</QuestionText>
  <AnswerGrid layout="vertical-expanded">
    {choices.map(choice => (
      <AnswerOption>
        <strong>{choice.choice_text}</strong>
        <div className="choice-details">
          {choice.choice_code ? (
            <pre>{choice.choice_code}</pre>
          ) : (
            <p>{choice.choice_text}</p>
          )}
        </div>
      </AnswerOption>
    ))}
  </AnswerGrid>
</QuestionCard>
```

---

## Future Question Types

### Potential Additions (Not in Current Source Format)

These types are **NOT currently supported** but could be added:

#### True/False Questions
```json
{
  "question_type": "true_false",
  "choices": [
    {"choice_text": "True", "is_correct": true},
    {"choice_text": "False", "is_correct": false}
  ]
}
```

#### Fill-in-the-Blank
```json
{
  "question_type": "text",
  "question_text": "The capital of France is _____",
  "answer_text": "Paris"
}
```

#### Multiple Answer (Select All That Apply)
```json
{
  "question_type": "multiple_answer",
  "choices": [
    {"choice_text": "A", "is_correct": true},
    {"choice_text": "B", "is_correct": true},
    {"choice_text": "C", "is_correct": false}
  ]
}
```

**Note:** These would require schema changes and are not part of current source data.

---

## Metadata Tags

Tag questions with format metadata for better filtering/search:

```json
{
  "question_id": 123,
  "format_tags": [
    "has_code",
    "code_language_javascript",
    "text_options"
  ],
  "complexity": "medium"
}
```

**Suggested Tags:**
- `has_code` - Question includes code
- `code_options` - Options are code blocks
- `has_image` - Question includes image
- `multi_language` - Multiple programming languages
- `long_question` - Question text > 200 chars
- `many_options` - More than 4 options

---

## Difficulty Inference

While the source data doesn't provide difficulty ratings, we can infer:

### Difficulty Indicators

**Easy:**
- Short question text (< 100 chars)
- No code blocks
- Common keywords (basic, simple, what is)

**Medium:**
- Moderate length (100-300 chars)
- One code block
- Analytical questions

**Hard:**
- Long question text (> 300 chars)
- Multiple code blocks
- Complex scenarios
- Debugging questions
- Multi-language questions

**Algorithm:**
```python
def infer_difficulty(question: Question) -> str:
    score = 0
    
    if len(question.question_text) > 300:
        score += 2
    if question.question_code:
        score += 1
    if question.question_image:
        score += 1
    if len(question.choices) > 4:
        score += 1
    if any(c.choice_code for c in question.choices):
        score += 2
    
    if score >= 4:
        return 'hard'
    elif score >= 2:
        return 'medium'
    else:
        return 'easy'
```

---

## Summary

| Type | % | Question Format | Option Format | Parsing Complexity |
|------|---|-----------------|---------------|-------------------|
| Plain Text | 40% | Text | Text | Low |
| Code Question | 35% | Text + Code | Text | Medium |
| Code Options | 15% | Text | Lettered + Code | High |
| Image Question | 5% | Text + Image | Text | Medium |
| Multi-Part | 5% | Text | Lists/Blocks | High |

**Key Takeaway:** All questions are `multiple_choice` type, but format varies significantly. Use format detection to apply appropriate parsing and rendering strategies.
