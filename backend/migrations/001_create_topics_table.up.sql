-- Create topics table
CREATE TABLE topics (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon_url VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_topics_slug ON topics(slug);
CREATE INDEX idx_topics_is_active ON topics(is_active);

-- Insert topics
INSERT INTO topics (name, slug, description, icon_url) VALUES
('HTML', 'html', 'HyperText Markup Language fundamentals', '/icon-html.svg'),
('CSS', 'css', 'Cascading Style Sheets fundamentals', '/icon-css.svg'),
('JavaScript', 'javascript', 'JavaScript programming fundamentals', '/icon-js.svg'),
('Accessibility', 'accessibility', 'Web accessibility best practices', '/icon-accessibility.svg');
