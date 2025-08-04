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

-- Insert default topics
INSERT INTO topics (name, slug, description, icon_url) VALUES
('Programming', 'programming', 'Test your programming knowledge', '/icons/programming.svg'),
('Science', 'science', 'General science questions', '/icons/science.svg'),
('History', 'history', 'World history and events', '/icons/history.svg'),
('Mathematics', 'mathematics', 'Math problems and concepts', '/icons/math.svg'),
('Technology', 'technology', 'Modern technology and computing', '/icons/tech.svg');
