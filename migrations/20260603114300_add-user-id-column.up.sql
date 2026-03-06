-- Add nullable user_id columns first
ALTER TABLE categories ADD COLUMN user_id INTEGER;
ALTER TABLE incomes ADD COLUMN user_id INTEGER;
ALTER TABLE outcomes ADD COLUMN user_id INTEGER;

-- Step 2: Populate user_id values for existing data
-- Check if any of the tables have existing data
DO $$
DECLARE
    has_data BOOLEAN := FALSE;
    default_user_id INTEGER;
BEGIN
    -- Check if any table has existing data
    SELECT EXISTS(
        SELECT 1 FROM categories WHERE user_id IS NULL
        UNION
        SELECT 1 FROM incomes WHERE user_id IS NULL
        UNION
        SELECT 1 FROM outcomes WHERE user_id IS NULL
    ) INTO has_data;
    
    -- If there is existing data, create a default user and populate user_id
    IF has_data THEN
        -- Create default admin user (password = password123)
        INSERT INTO users (first_name, last_name, email, password_hash)
        VALUES ('Admin', 'ADMIN', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi')
        RETURNING id INTO default_user_id;
        
        -- Update existing records with the default user_id
        UPDATE categories SET user_id = default_user_id WHERE user_id IS NULL;
        UPDATE incomes SET user_id = default_user_id WHERE user_id IS NULL;
        UPDATE outcomes SET user_id = default_user_id WHERE user_id IS NULL;
    END IF;
END $$;

-- Add NOT NULL constraints after data is populated
ALTER TABLE categories ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE incomes ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE outcomes ALTER COLUMN user_id SET NOT NULL;

-- Add foreign key constraints
ALTER TABLE categories ADD CONSTRAINT fk_categories_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE incomes ADD CONSTRAINT fk_incomes_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE outcomes ADD CONSTRAINT fk_outcomes_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add indexes
CREATE INDEX idx_categories_user_id ON categories (user_id);
CREATE INDEX idx_incomes_user_id ON incomes (user_id);
CREATE INDEX idx_outcomes_user_id ON outcomes (user_id);