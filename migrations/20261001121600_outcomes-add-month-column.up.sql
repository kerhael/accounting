ALTER TABLE outcomes
ADD COLUMN month date;

UPDATE outcomes SET month = date_trunc('month', created_at)::date;

CREATE OR REPLACE FUNCTION set_outcome_month()
RETURNS trigger AS $$
BEGIN
    NEW.month := date_trunc('month', NEW.created_at)::date;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_outcome_month
BEFORE INSERT OR UPDATE OF created_at
ON outcomes
FOR EACH ROW
EXECUTE FUNCTION set_outcome_month();

CREATE INDEX idx_outcomes_category_month
ON outcomes (category_id, month);