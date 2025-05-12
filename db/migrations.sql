-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    github_id VARCHAR(255) UNIQUE,
    email VARCHAR(255),
    avatar_url TEXT,
    groq_api_key TEXT,
    default_model VARCHAR(50) DEFAULT 'llama3-8b-8192',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create prompts table
CREATE TABLE IF NOT EXISTS prompts (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    prompt TEXT NOT NULL,
    response TEXT,
    model VARCHAR(50) DEFAULT 'llama3-8b-8192',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
