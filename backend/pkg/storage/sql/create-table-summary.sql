
CREATE TABLE summary (
	run_id TEXT NOT NULL,
	result TEXT NOT NULL,
	passed integer NOT NULL,
	failed integer NOT NULL,
	ignored integer NOT NULL,
	duration READ NOT NULL,
	created_at TEXT NOT NULL 
);
