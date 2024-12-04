CREATE TABLE metadata (
	run_id TEXT NOT NULL,
	commit_sha TEXT NOT NULL,
	job_name TEXT NOT NULL,
	repo TEXT NOT NULL,
	branch TEXT NOT NULL,
	format TEXT NOT NULL,
	link TEXT NOT NULL,
	tags TEXT NOT NULL,
	created_at TEXT NOT NULL 
);