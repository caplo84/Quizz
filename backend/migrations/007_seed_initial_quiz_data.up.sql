-- Seed all questions and choices from data.json for HTML, CSS, JavaScript, and Accessibility

-- HTML Quiz Questions
INSERT INTO questions (quiz_id, question_text, question_type, points, order_index) VALUES
((SELECT id FROM quizzes WHERE slug = 'html-basics'), 'What does HTML stand for?', 'multiple_choice', 1, 1),
((SELECT id FROM quizzes WHERE slug = 'html-basics'), 'Which of the following is the correct structure for an HTML document?', 'multiple_choice', 1, 2),
((SELECT id FROM quizzes WHERE slug = 'html-basics'), 'Which HTML element is used to define the title of a document?', 'multiple_choice', 1, 3),
((SELECT id FROM quizzes WHERE slug = 'html-basics'), 'What is the purpose of the <body> tag in HTML?', 'multiple_choice', 1, 4),
((SELECT id FROM quizzes WHERE slug = 'html-basics'), 'Which HTML tag is used to create a hyperlink?', 'multiple_choice', 1, 5),
((SELECT id FROM quizzes WHERE slug = 'html-basics'), 'Which tag is used to display images in HTML?', 'multiple_choice', 1, 6),
((SELECT id FROM quizzes WHERE slug = 'html-basics'), 'What attribute is used to provide the path of an image in the <img> tag?', 'multiple_choice', 1, 7),
((SELECT id FROM quizzes WHERE slug = 'html-basics'), 'Which HTML tag is used to create an unordered list?', 'multiple_choice', 1, 8),
((SELECT id FROM quizzes WHERE slug = 'html-basics'), 'What does the <br> tag do?', 'multiple_choice', 1, 9),
((SELECT id FROM quizzes WHERE slug = 'html-basics'), 'In HTML, what does the `fieldset` tag do?', 'multiple_choice', 1, 10);

-- CSS Quiz Questions
INSERT INTO questions (quiz_id, question_text, question_type, points, order_index) VALUES
((SELECT id FROM quizzes WHERE slug = 'css-basics'), 'What does CSS stand for?', 'multiple_choice', 1, 1),
((SELECT id FROM quizzes WHERE slug = 'css-basics'), 'Which HTML attribute is used to define inline styles?', 'multiple_choice', 1, 2),
((SELECT id FROM quizzes WHERE slug = 'css-basics'), 'How do you insert a comment in a CSS file?', 'multiple_choice', 1, 3),
((SELECT id FROM quizzes WHERE slug = 'css-basics'), 'Which property is used to change the background color of an element?', 'multiple_choice', 1, 4),
((SELECT id FROM quizzes WHERE slug = 'css-basics'), 'How do you apply a style to all <p> elements?', 'multiple_choice', 1, 5),
((SELECT id FROM quizzes WHERE slug = 'css-basics'), 'Which property is used to change the text color of an element?', 'multiple_choice', 1, 6),
((SELECT id FROM quizzes WHERE slug = 'css-basics'), 'How do you select an element with the class name "header"?', 'multiple_choice', 1, 7),
((SELECT id FROM quizzes WHERE slug = 'css-basics'), 'What is the default position value of an HTML element?', 'multiple_choice', 1, 8),
((SELECT id FROM quizzes WHERE slug = 'css-basics'), 'Which property is used to change the font of an element?', 'multiple_choice', 1, 9),
((SELECT id FROM quizzes WHERE slug = 'css-basics'), 'How do you make each word in a text start with a capital letter?', 'multiple_choice', 1, 10);

-- JavaScript Quiz Questions
INSERT INTO questions (quiz_id, question_text, question_type, points, order_index) VALUES
((SELECT id FROM quizzes WHERE slug = 'javascript-basics'), 'What does the ''let'' keyword do?', 'multiple_choice', 1, 1),
((SELECT id FROM quizzes WHERE slug = 'javascript-basics'), 'Which of the following is NOT a JavaScript data type?', 'multiple_choice', 1, 2),
((SELECT id FROM quizzes WHERE slug = 'javascript-basics'), 'What is the purpose of the ''return'' statement in a function?', 'multiple_choice', 1, 3),
((SELECT id FROM quizzes WHERE slug = 'javascript-basics'), 'Which method is used to add an element to the end of an array?', 'multiple_choice', 1, 4),
((SELECT id FROM quizzes WHERE slug = 'javascript-basics'), 'What will ''typeof null'' return in JavaScript?', 'multiple_choice', 1, 5),
((SELECT id FROM quizzes WHERE slug = 'javascript-basics'), 'Which operator is used to check both value and type in JavaScript?', 'multiple_choice', 1, 6),
((SELECT id FROM quizzes WHERE slug = 'javascript-basics'), 'What is the correct way to call a function named ''myFunction''?', 'multiple_choice', 1, 7),
((SELECT id FROM quizzes WHERE slug = 'javascript-basics'), 'Which method is used to remove the last element from an array?', 'multiple_choice', 1, 8),
((SELECT id FROM quizzes WHERE slug = 'javascript-basics'), 'What is the correct way to write a JavaScript array?', 'multiple_choice', 1, 9),
((SELECT id FROM quizzes WHERE slug = 'javascript-basics'), 'Which built-in method returns the length of the string?', 'multiple_choice', 1, 10);

-- Accessibility Quiz Questions
INSERT INTO questions (quiz_id, question_text, question_type, points, order_index) VALUES
((SELECT id FROM quizzes WHERE slug = 'accessibility-basics'), 'What does ''WCAG'' stand for?', 'multiple_choice', 1, 1),
((SELECT id FROM quizzes WHERE slug = 'accessibility-basics'), 'Which element is used to provide alternative text for images?', 'multiple_choice', 1, 2),
((SELECT id FROM quizzes WHERE slug = 'accessibility-basics'), 'What is the purpose of semantic HTML in accessibility?', 'multiple_choice', 1, 3),
((SELECT id FROM quizzes WHERE slug = 'accessibility-basics'), 'Which attribute is used to describe the purpose of an input field?', 'multiple_choice', 1, 4),
((SELECT id FROM quizzes WHERE slug = 'accessibility-basics'), 'What is the minimum color contrast ratio for normal text?', 'multiple_choice', 1, 5),
((SELECT id FROM quizzes WHERE slug = 'accessibility-basics'), 'Which element is inherently focusable?', 'multiple_choice', 1, 6),
((SELECT id FROM quizzes WHERE slug = 'accessibility-basics'), 'What is the purpose of the ''lang'' attribute?', 'multiple_choice', 1, 7),
((SELECT id FROM quizzes WHERE slug = 'accessibility-basics'), 'Which guideline ensures keyboard accessibility?', 'multiple_choice', 1, 8),
((SELECT id FROM quizzes WHERE slug = 'accessibility-basics'), 'What is the role of ''skip navigation'' links?', 'multiple_choice', 1, 9),
((SELECT id FROM quizzes WHERE slug = 'accessibility-basics'), 'Which tool can help check website accessibility?', 'multiple_choice', 1, 10);

-- HTML Question Choices
INSERT INTO choices (question_id, choice_text, is_correct, order_index) VALUES
-- Q1: What does HTML stand for?
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 1), 'Hyper Trainer Marking Language', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 1), 'Hyper Text Marketing Language', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 1), 'Hyper Text Markup Language', true, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 1), 'Hyper Text Markup Leveler', false, 4),

-- Q2: HTML document structure
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 2), '<html><head></head><body></body></html>', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 2), '<head><html></html><body></body></head>', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 2), '<body><head></head><html></html></body>', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 2), '<html><body></body><head></head></html>', false, 4),

-- Q3: Title element
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 3), '<head>', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 3), '<title>', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 3), '<header>', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 3), '<top>', false, 4),

-- Q4: Body tag purpose
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 4), 'It defines the document''s head section.', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 4), 'It contains all the content such as text, images, and links.', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 4), 'It is used to define the main content of an HTML document.', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 4), 'It specifies the body of the email content in HTML.', false, 4),

-- Q5: Hyperlink tag
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 5), '<hyperlink>', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 5), '<link>', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 5), '<a>', true, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 5), '<href>', false, 4),

-- Q6: Image tag
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 6), '<img>', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 6), '<image>', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 6), '<src>', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 6), '<pic>', false, 4),

-- Q7: Image path attribute
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 7), 'link', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 7), 'src', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 7), 'href', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 7), 'url', false, 4),

-- Q8: Unordered list tag
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 8), '<ul>', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 8), '<ol>', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 8), '<list>', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 8), '<li>', false, 4),

-- Q9: BR tag purpose
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 9), 'It breaks the text into two sections.', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 9), 'It creates a bold text.', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 9), 'It inserts a line break.', true, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 9), 'It adds a new row in a table.', false, 4),

-- Q10: Fieldset tag purpose
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 10), 'It is used to group related data in a form.', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 10), 'It sets the field to a fixed size.', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 10), 'It automatically validates the fields within a form.', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'html-basics') AND order_index = 10), 'It hides the fields in a form.', false, 4);

-- CSS Question Choices
INSERT INTO choices (question_id, choice_text, is_correct, order_index) VALUES
-- Q1: What does CSS stand for?
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 1), 'Colorful Style Sheets', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 1), 'Computer Style Sheets', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 1), 'Cascading Style Sheets', true, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 1), 'Creative Style Sheets', false, 4),

-- Q2: Inline styles attribute
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 2), 'styles', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 2), 'style', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 2), 'class', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 2), 'font-style', false, 4),

-- Q3: CSS comment
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 3), '// this is a comment //', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 3), '/* this is a comment */', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 3), '<!-- this is a comment -->', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 3), '** this is a comment **', false, 4),

-- Q4: Background color property
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 4), 'color', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 4), 'background-color', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 4), 'bgcolor', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 4), 'bg-color', false, 4),

-- Q5: Select all p elements
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 5), 'p', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 5), '.p', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 5), '#p', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 5), 'all.p', false, 4),

-- Q6: Text color property
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 6), 'color', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 6), 'text-color', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 6), 'font-color', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 6), 'text', false, 4),

-- Q7: Class selector
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 7), 'header', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 7), '.header', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 7), '#header', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 7), '*header', false, 4),

-- Q8: Default position value
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 8), 'relative', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 8), 'absolute', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 8), 'fixed', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 8), 'static', true, 4),

-- Q9: Font property
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 9), 'font-family', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 9), 'font-weight', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 9), 'font-style', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 9), 'font', false, 4),

-- Q10: Capitalize text
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 10), 'text-transform: capitalize', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 10), 'text-style: capital', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 10), 'font-transform: capitalize', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'css-basics') AND order_index = 10), 'transform: capitalize', false, 4);

-- JavaScript Question Choices
INSERT INTO choices (question_id, choice_text, is_correct, order_index) VALUES
-- Q1: let keyword
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 1), 'Declares a block-scoped variable', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 1), 'Declares a function-scoped variable', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 1), 'Declares a constant variable', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 1), 'Declares a global variable', false, 4),

-- Q2: NOT a JavaScript data type
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 2), 'string', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 2), 'boolean', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 2), 'float', true, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 2), 'object', false, 4),

-- Q3: return statement purpose
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 3), 'Stops the execution of a function and returns a value', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 3), 'Restarts the function', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 3), 'Throws an error', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 3), 'Prints a value to the console', false, 4),

-- Q4: Add element to end of array
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 4), 'append()', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 4), 'push()', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 4), 'add()', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 4), 'insert()', false, 4),

-- Q5: typeof null
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 5), '"null"', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 5), '"undefined"', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 5), '"object"', true, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 5), '"number"', false, 4),

-- Q6: Strict equality operator
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 6), '==', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 6), '===', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 6), '!=', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 6), '!==', false, 4),

-- Q7: Call function
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 7), 'call function myFunction()', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 7), 'call myFunction()', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 7), 'myFunction()', true, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 7), 'execute myFunction()', false, 4),

-- Q8: Remove last element from array
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 8), 'pop()', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 8), 'push()', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 8), 'shift()', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 8), 'unshift()', false, 4),

-- Q9: JavaScript array syntax
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 9), 'var colors = "red", "green", "blue"', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 9), 'var colors = ["red", "green", "blue"]', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 9), 'var colors = 1 = ("red"), 2 = ("green"), 3 = ("blue")', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 9), 'var colors = (1:"red", 2:"green", 3:"blue")', false, 4),

-- Q10: String length method
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 10), 'length()', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 10), 'length', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 10), 'getLength()', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'javascript-basics') AND order_index = 10), 'len()', false, 4);

-- Accessibility Question Choices
INSERT INTO choices (question_id, choice_text, is_correct, order_index) VALUES
-- Q1: WCAG stands for
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 1), 'Web Content Accessibility Guidelines', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 1), 'Web Compliance Accessibility Guide', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 1), 'Website Content Accessibility Guidelines', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 1), 'Web Communication Accessibility Guidelines', false, 4),

-- Q2: Alternative text for images
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 2), 'alt', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 2), 'title', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 2), 'description', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 2), 'summary', false, 4),

-- Q3: Semantic HTML purpose
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 3), 'It helps search engines understand the content better', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 3), 'It provides meaning and structure to the content', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 3), 'It makes the website load faster', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 3), 'It improves the website''s appearance', false, 4),

-- Q4: Input field description attribute
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 4), 'label', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 4), 'name', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 4), 'placeholder', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 4), 'title', false, 4),

-- Q5: Color contrast ratio
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 5), '3:1', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 5), '4.5:1', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 5), '7:1', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 5), '2:1', false, 4),

-- Q6: Inherently focusable element
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 6), '<div>', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 6), '<span>', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 6), '<a href="...">', true, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 6), '<p>', false, 4),

-- Q7: lang attribute purpose
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 7), 'To specify the scripting language', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 7), 'To define the character set', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 7), 'To indicate the language of the page content', true, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 7), 'To declare a language pack', false, 4),

-- Q8: Keyboard accessibility guideline
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 8), 'Keyboard Accessible', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 8), 'Mouse Independence', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 8), 'Device Independence', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 8), 'Operable Controls', false, 4),

-- Q9: Skip navigation links role
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 9), 'To skip over primary navigation to the main content', true, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 9), 'To provide shortcuts to different sections of the website', false, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 9), 'To help users skip unwanted sections like advertisements', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 9), 'To bypass broken links in the navigation', false, 4),

-- Q10: Accessibility checking tool
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 10), 'W3C Validator', false, 1),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 10), 'Google Lighthouse', true, 2),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 10), 'CSS Validator', false, 3),
((SELECT id FROM questions WHERE quiz_id = (SELECT id FROM quizzes WHERE slug = 'accessibility-basics') AND order_index = 10), 'JavaScript Console', false, 4);
