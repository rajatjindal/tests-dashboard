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

CREATE TABLE summary (
	run_id TEXT NOT NULL,
	result TEXT NOT NULL,
	passed integer NOT NULL,
	failed integer NOT NULL,
	ignored integer NOT NULL,
	duration READ NOT NULL,
	created_at TEXT NOT NULL 
);

CREATE TABLE suite_summary (
	run_id TEXT NOT NULL,
	suite_id TEXT NOT NULL,
	suite_name TEXT NOT NULL,
	result TEXT NOT NULL,
	passed integer NOT NULL,
	failed integer NOT NULL,
	ignored integer NOT NULL,
	duration READ NOT NULL,
	created_at TEXT NOT NULL 
);


CREATE TABLE metadata (
	run_id TEXT NOT NULL,
	repo TEXT NOT NULL,
	branch TEXT NOT NULL,
	format TEXT NOT NULL,
	link TEXT NOT NULL,
	tags TEXT NOT NULL,
	created_at TEXT NOT NULL 
);