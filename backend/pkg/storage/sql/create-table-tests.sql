CREATE TABLE tests (
	run_id TEXT NOT NULL,
	suite_id TEXT NOT NULL,
	suite_name TEXT NOT NULL,
	name TEXT NOT NULL,
	result TEXT NOT NULL,
	duration REAL NOT NULL,
	logs TEXT NOT NULL,
	created_at TEXT NOT NULL 
);
