ALTER TABLE dishes
DROP COLUMN preparation_time;

ALTER TABLE dishes
DROP COLUMN is_featured;

ALTER TABLE dishes
DROP COLUMN is_vegetarian;

ALTER TABLE dishes
RENAME COLUMN is_available TO available;
