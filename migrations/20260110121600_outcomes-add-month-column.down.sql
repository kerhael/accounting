DROP TRIGGER IF EXISTS trg_set_outcome_month ON outcomes;
DROP FUNCTION IF EXISTS set_outcome_month();
DROP INDEX IF EXISTS idx_outcomes_category_month;
ALTER TABLE outcomes DROP COLUMN month;