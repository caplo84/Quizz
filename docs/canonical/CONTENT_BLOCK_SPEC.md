# Content Block Specification

## Overview
Content blocks represent different types of media embedded within questions or answer choices. This specification defines how to extract, transform, and store content from upstream quiz markdown files.

---

## Content Block Types

### 1. Text Content
**Description**: Plain text content (default type)

**Source Format:**
```markdown
#### Q1. Which operator returns true if values are not equal?
- [x] `!==`
```

**Storage:**
```json
{
  "question_text": "Which operator returns true if values are not equal?",
  "choice_text": "`!==`"
}
```

**Processing Rules:**
- Preserve inline code markers (backticks)
- Clean up extra whitespace
- Preserve line breaks within text
- Strip markdown formatting symbols (**, __, *, _) if needed

---

### 2. Code Block Content

**Description**: Multi-line code snippets with syntax highlighting

**Source Format:**
````markdown
#### Q3. Review the code below. Which statement calls the addTax function?

```js
function addTax(total) {
  return total * 1.05;
}
```

- [x] `addTax(50);`
````

**Storage:**
```json
{
  "question_text": "Review the code below. Which statement calls the addTax function?",
  "question_code": "function addTax(total) {\n  return total * 1.05;\n}",
  "question_code_language": "js"
}
```

**Language Detection:**
- Extract from code fence: ` ```javascript `, ` ```java `, ` ```python `
- Common languages: `js`, `javascript`, `java`, `kotlin`, `python`, `css`, `html`, `xml`, `sql`, `bash`, `json`, `yaml`
- Default to `text` if no language specified

**Processing Rules:**
1. Extract code between ` ``` ` markers
2. Preserve all whitespace and indentation
3. Store code as-is (don't escape)
4. Detect language from fence marker
5. Store in separate `*_code` field

---

### 3. Image Content

**Description**: Visual content (diagrams, screenshots, examples)

**Source Formats:**

**Local Images:**
```markdown
![img](image/43.jpeg)
```

**External URLs:**
```markdown
![quote](https://raw.githubusercontent.com/ram-sah/Community-Assessments/master/CSS/images/rm-1.png)
```

**Storage:**
```json
{
  "question_image_url": "https://raw.githubusercontent.com/your-org/quiz-content-source/main/android/image/43.jpeg",
  "question_image_alt": "img"
}
```

**URL Transformation Rules:**
1. **Local paths**: Convert to full GitHub URL
   - `image/43.jpeg` → `https://raw.githubusercontent.com/your-org/quiz-content-source/main/{topic}/image/43.jpeg`
   
2. **External URLs**: Store as-is
   - Already absolute URLs
   
3. **Relative paths**: Resolve relative to file location
   - `../assets/image.png` → Full path

**Processing Rules:**
- Extract alt text from `![alt-text](url)`
- Store `url` → `*_image_url`
- Store `alt-text` → `*_image_alt`
- Validate URL accessibility (optional)

---

### 4. Mixed Content

**Description**: Questions with multiple content types

**Source Format:**
````markdown
#### Q15. What color will the link be?

```css
.example {
  color: yellow;
}
```

```html
<ul>
  <li><a href="#" class="example">link</a></li>
</ul>
```

- [x] yellow
````

**Storage Strategy:**
```json
{
  "question_text": "What color will the link be?",
  "question_code": ".example {\n  color: yellow;\n}\n\n<ul>\n  <li><a href=\"#\" class=\"example\">link</a></li>\n</ul>",
  "question_code_language": "html"
}
```

**Processing Rules:**
1. Concatenate multiple code blocks
2. Use the last specified language
3. Separate blocks with `\n\n`
4. Preserve order of appearance

---

## Content Block Location

### Question-Level Content
Content that appears in the question itself.

**Fields:**
- `question_text` (required)
- `question_code` (optional)
- `question_code_language` (optional)
- `question_image_url` (optional)
- `question_image_alt` (optional)

### Choice-Level Content
Content that appears within answer options.

**Fields:**
- `choice_text` (required)
- `choice_code` (optional)
- `choice_code_language` (optional)
- `choice_image_url` (optional)
- `choice_image_alt` (optional)

---

## Extraction Algorithm

### Step 1: Parse Question Block
```
1. Extract question text from "#### Q[N]. [TEXT]"
2. Scan for code blocks before first option
3. Scan for images before first option
4. Store in question_* fields
```

### Step 2: Parse Options
```
For each line starting with "- [":
  1. Extract choice text
  2. Check for code block immediately after
  3. Check for image immediately after
  4. Store in choice_* fields
```

### Step 3: Handle Nested Content
```markdown
- [ ] A

```java
RecycleView
RecyclerView.Adapter
```
```

**Processing:**
1. Identify multi-line option (letter followed by blank line)
2. Extract code block following the option
3. Associate code with that specific choice
4. Continue until next option marker

---

## Content Validation

### Text Validation
- **Min length**: 10 characters for questions
- **Min length**: 1 character for choices
- **Max length**: 5000 characters
- **Special chars**: Allow all Unicode

### Code Validation
- **Languages**: Validate against known list
- **Format**: Must be valid markdown code block
- **Size**: Max 10,000 characters

### Image Validation
- **Protocols**: `http://`, `https://` only
- **Extensions**: `.png`, `.jpg`, `.jpeg`, `.gif`, `.svg`
- **Accessibility**: Check URL returns 200 status (optional)

---

## Edge Cases

### 1. Inline Code in Text
```markdown
- [x] Use the `ActivityManager.isLowRamDevice()` method
```

**Handling**: Keep inline code markers in `choice_text`

### 2. Multiple Code Blocks in Options
````markdown
- [ ] A

```java
Option A Code
```

```xml
<layout>Option A XML</layout>
```
````

**Handling**: Concatenate with separator

### 3. Images in Both Question and Options
```markdown
#### Q10. What shape is this?
![shape](shape.png)

- [ ] ![circle](circle.png)
- [x] ![square](square.png)
```

**Handling**: Store separately in question_* and choice_* fields

### 4. No Code Language Specified
````markdown
```
function test() {}
```
````

**Handling**: Default to `text` or infer from context/file extension

---

## Rendering Guidelines

### Frontend Display

**Question Code:**
```jsx
{question.question_code && (
  <CodeBlock 
    code={question.question_code}
    language={question.question_code_language}
  />
)}
```

**Question Image:**
```jsx
{question.question_image_url && (
  <ImageDisplay
    src={question.question_image_url}
    alt={question.question_image_alt}
  />
)}
```

**Choice Code:**
```jsx
{choice.choice_code ? (
  <CodeBlock code={choice.choice_code} language={choice.choice_code_language} />
) : (
  <p>{choice.choice_text}</p>
)}
```

---

## Performance Considerations

### Storage
- Store code/images separately to optimize text search
- Index text fields for full-text search
- Store images as URLs (don't embed base64)

### Caching
- Cache rendered code blocks
- Cache image proxies
- Pre-render markdown on import

### Bandwidth
- Lazy load images
- Use CDN for external images
- Compress code blocks (gzip)
