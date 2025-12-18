CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    label TEXT NOT NULL UNIQUE
);

CREATE TABLE outcomes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    amount INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_outcomes_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_outcomes_category_id ON outcomes(category_id);

CREATE TABLE incomes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);