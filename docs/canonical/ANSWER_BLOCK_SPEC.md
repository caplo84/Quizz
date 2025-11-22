# Answer Block Specification

## Overview
This specification defines how answer choices (options) are structured, extracted, and validated from upstream quiz markdown files.

---

## Answer Block Structure

### Basic Format

**Source:**
```markdown
- [ ] Option A
- [x] Option B (Correct)
- [ ] Option C
- [ ] Option D
```

**Canonical Model:**
```json
[
  {
    "choice_text": "Option A",
    "is_correct": false,
    "order_index": 1
  },
  {
    "choice_text": "Option B (Correct)",
    "is_correct": true,
    "order_index": 2
  },
  {
    "choice_text": "Option C",
    "is_correct": false,
    "order_index": 3
  },
  {
    "choice_text": "Option D",
    "is_correct": false,
    "order_index": 4
  }
]
```

---

## Answer Correctness Markers

### Standard Marker
- `[x]` - Correct answer
- `[ ]` - Incorrect answer

**Regular Expression:**
```regex
^- \[(x| )\] (.+)$
```

### Validation Rules
1. **Exactly ONE correct answer** per question
2. **At least TWO choices** minimum
3. **Maximum SIX choices** (rare, usually 4)
4. **No duplicate correct markers** (`[x]`)

---

## Answer Types

### 1. Simple Text Answers

**Source:**
```markdown
- [ ] green
- [x] yellow
- [ ] blue
- [ ] red
```

**Storage:**
```json
{
  "choice_text": "yellow",
  "is_correct": true,
  "choice_code": null
}
```

---

### 2. Inline Code Answers

**Source:**
```markdown
- [ ] `let 100 = rate;`
- [x] `let rate = 100;`
- [ ] `100 = let rate;`
- [ ] `rate = 100;`
```

**Storage:**
```json
{
  "choice_text": "`let rate = 100;`",
  "is_correct": true,
  "choice_code": null
}
```

**Note**: Keep inline backticks in `choice_text`

---

### 3. Code Block Answers

**Source:**
````markdown
- [ ] A

```java
RecycleView
RecyclerView.Adapter<T extends BaseAdapter>
```

- [x] B

```java
RecycleView
RecyclerView.Adapter<VH extends ViewHolder>
```
````

**Storage:**
```json
{
  "choice_text": "B",
  "is_correct": true,
  "choice_code": "RecycleView\nRecyclerView.Adapter<VH extends ViewHolder>",
  "choice_code_language": "java"
}
```

---

### 4. Multi-Part Answers

**Source:**
```markdown
- [x] A

```
  1. Larger z-index values appear on top.
  2. Negative and positive numbers can be used.
  3. The z-index can be used only on positioned elements.
```

- [ ] B

```
  1. Smaller z-index values appear on top.
  2. Negative and positive numbers can be used.
```
```

**Storage:**
```json
{
  "choice_text": "A",
  "is_correct": true,
  "choice_code": "  1. Larger z-index values appear on top.\n  2. Negative and positive numbers can be used.\n  3. The z-index can be used only on positioned elements.",
  "choice_code_language": "text"
}
```

---

### 5. Image Answers

**Source:**
```markdown
- [ ] ![img](image/00.jpeg)
- [x] ![img](image/01.jpeg)
- [ ] ![img](image/02.jpeg)
```

**Storage:**
```json
{
  "choice_text": "img",
  "is_correct": true,
  "choice_image_url": "https://raw.githubusercontent.com/.../image/01.jpeg",
  "choice_image_alt": "img"
}
```

---

## Extraction Algorithm

### Step-by-Step Process

```python
def extract_choices(question_block: str) -> List[Choice]:
    choices = []
    lines = question_block.split('\n')
    current_choice = None
    order_index = 0
    
    for i, line in enumerate(lines):
        # Detect choice marker
        if line.startswith('- ['):
            # Save previous choice if exists
            if current_choice:
                choices.append(current_choice)
            
            # Start new choice
            order_index += 1
            is_correct = '[x]' in line
            choice_text = line.split('] ', 1)[1].strip()
            
            current_choice = {
                'choice_text': choice_text,
                'is_correct': is_correct,
                'order_index': order_index
            }
            
        # Check for code block after choice
        elif current_choice and line.startswith('```'):
            code, language = extract_code_block(lines[i:])
            current_choice['choice_code'] = code
            current_choice['choice_code_language'] = language
            
        # Check for image after choice
        elif current_choice and line.startswith('!['):
            url, alt = extract_image(line)
            current_choice['choice_image_url'] = url
            current_choice['choice_image_alt'] = alt
    
    # Add last choice
    if current_choice:
        choices.append(current_choice)
    
    return choices
```

---

## Validation Rules

### 1. Correctness Validation
```python
def validate_correctness(choices: List[Choice]) -> bool:
    correct_count = sum(1 for c in choices if c.is_correct)
    return correct_count == 1
```

**Error Cases:**
- ❌ No correct answer (`correct_count == 0`)
- ❌ Multiple correct answers (`correct_count > 1`)

---

### 2. Count Validation
```python
def validate_count(choices: List[Choice]) -> bool:
    return 2 <= len(choices) <= 6
```

**Error Cases:**
- ❌ Only 1 choice (not a valid multiple choice)
- ❌ More than 6 choices (rare, might be parsing error)

---

### 3. Order Validation
```python
def validate_order(choices: List[Choice]) -> bool:
    expected_order = list(range(1, len(choices) + 1))
    actual_order = [c.order_index for c in choices]
    return expected_order == actual_order
```

---

### 4. Content Validation
```python
def validate_content(choice: Choice) -> bool:
    # Must have either text or code
    has_text = choice.choice_text and len(choice.choice_text.strip()) > 0
    has_code = choice.choice_code and len(choice.choice_code.strip()) > 0
    return has_text or has_code
```

---

## Special Cases

### Case 1: Option Label as Text
```markdown
- [ ] A

[Additional content below]
```

**Handling:** 
- Store "A" as `choice_text`
- Store additional content as `choice_code` or part of `choice_text`

---

### Case 2: Special Characters in Options
```markdown
- [x] Use the "clearfix hack" on the parent element or use the overflow property with a value other than "visible."
```

**Handling:**
- Preserve quotes, punctuation
- Escape for JSON: `"visible.\""`

---

### Case 3: Multi-line Text Options
```markdown
- [ ] By default, block elements span the entire width of their container;
      inline elements are the same height and width as the content.
```

**Handling:**
- Combine lines into single `choice_text`
- Preserve space separation

---

### Case 4: Nested Lists in Options
```markdown
- [ ] A
  * Sub-item 1
  * Sub-item 2
```

**Handling:**
- Treat entire block as `choice_text` or `choice_code`
- Preserve formatting

---

## Order Index Convention

### Standard Ordering
```
Option 1 → order_index: 1
Option 2 → order_index: 2
Option 3 → order_index: 3
Option 4 → order_index: 4
```

### Why Order Matters
1. **Reproducible display** - Same order every time
2. **Answer tracking** - "User selected option 2"
3. **Shuffle detection** - Can shuffle while tracking original position
4. **Audit trail** - Changes in order indicate content updates

---

## Answer Explanation Handling

### Inline Explanations
Some options have explanations:

```markdown
- [x] Option B

**Explanation:** This is why B is correct.
```

**Storage:**
```json
{
  "choice_text": "Option B",
  "is_correct": true,
  "explanation": "This is why B is correct."
}
```

---

## Rendering Strategy

### Display Format

**Simple Text:**
```html
<div class="choice">
  <label>
    <input type="radio" name="q1" value="1">
    <span>Option A</span>
  </label>
</div>
```

**With Code:**
```html
<div class="choice">
  <label>
    <input type="radio" name="q1" value="2">
    <pre><code class="language-java">
    RecyclerView.Adapter<VH>
    </code></pre>
  </label>
</div>
```

**With Image:**
```html
<div class="choice">
  <label>
    <input type="radio" name="q1" value="3">
    <img src="image-url.jpg" alt="Option C">
  </label>
</div>
```

---

## Answer Shuffling

### Shuffle Strategy
```python
def shuffle_choices(choices: List[Choice], seed: int) -> List[Choice]:
    # Store original order
    for choice in choices:
        choice.original_order = choice.order_index
    
    # Shuffle
    random.seed(seed)
    shuffled = random.sample(choices, len(choices))
    
    # Re-assign display order
    for i, choice in enumerate(shuffled):
        choice.display_order = i + 1
    
    return shuffled
```

**Use Cases:**
- Prevent answer pattern memorization
- Different order per quiz attempt
- Maintain tracking via `original_order`

---

## Error Handling

### Import Errors

| Error | Handling |
|-------|----------|
| No correct answer | Skip question + log warning |
| Multiple correct answers | Skip question + log error |
| < 2 choices | Skip question + log error |
| > 6 choices | Import all + log warning |
| Empty choice text | Skip choice + log warning |
| Malformed markdown | Try to parse + log error |

### Quality Flags

Tag questions with quality issues:
```json
{
  "quality_flags": [
    "multiple_correct_answers",
    "ambiguous_options",
    "missing_code_language"
  ]
}
```
