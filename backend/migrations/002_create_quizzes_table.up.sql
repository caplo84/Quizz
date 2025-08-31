-- Create quizzes table
CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL UNIQUE,
    description TEXT,
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    difficulty_level VARCHAR(20) DEFAULT 'medium' CHECK (difficulty_level IN ('easy', 'medium', 'hard')),
    time_limit_minutes INTEGER DEFAULT 30,
    total_questions INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_quizzes_slug ON quizzes(slug);
CREATE INDEX idx_quizzes_topic_id ON quizzes(topic_id);
CREATE INDEX idx_quizzes_difficulty ON quizzes(difficulty_level);
CREATE INDEX idx_quizzes_is_active ON quizzes(is_active);

-- Insert initial quizzes
INSERT INTO quizzes (title, slug, description, topic_id, difficulty_level, total_questions) VALUES
('HTML Basics', 'html-basics', 'Test your knowledge of HTML fundamentals', 
 (SELECT id FROM topics WHERE slug = 'html'), 'medium', 10),
('CSS Basics', 'css-basics', 'Test your knowledge of CSS fundamentals', 
 (SELECT id FROM topics WHERE slug = 'css'), 'medium', 10),
('JavaScript Basics', 'javascript-basics', 'Test your knowledge of JavaScript fundamentals', 
 (SELECT id FROM topics WHERE slug = 'javascript'), 'medium', 10),
('Accessibility Basics', 'accessibility-basics', 'Test your knowledge of web accessibility', 
 (SELECT id FROM topics WHERE slug = 'accessibility'), 'medium', 10);
