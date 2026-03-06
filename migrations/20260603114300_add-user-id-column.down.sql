DROP INDEX IF EXISTS idx_outcomes_user_id;
ALTER TABLE outcomes DROP COLUMN user_id;

DROP INDEX IF EXISTS idx_incomes_user_id;
ALTER TABLE incomes DROP COLUMN user_id;

DROP INDEX IF EXISTS idx_categories_user_id;
ALTER TABLE categories DROP COLUMN user_id;