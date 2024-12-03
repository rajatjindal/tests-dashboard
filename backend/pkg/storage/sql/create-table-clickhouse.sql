CREATE TABLE tests
(
    run_id String,
		suite_id String,
		suite_name String,
		name String,
		result String,
		duration Float64,
		logs String,
		created_at DateTime
)
ENGINE = MergeTree
PRIMARY KEY (run_id, suite_id, name);

CREATE TABLE summary
(
    run_id String,
		result String,		
		passed UInt32,
		failed UInt32,
		ignored UInt32,
		duration Float64,
		created_at DateTime
)
ENGINE = MergeTree
PRIMARY KEY (run_id);

CREATE TABLE suite_summary
(
    run_id String,
		suite_id String,
		suite_name String,
		result String,		
		passed UInt32,
		failed UInt32,
		ignored UInt32,
		duration Float64,
		created_at DateTime
)
ENGINE = MergeTree
PRIMARY KEY (run_id, suite_id, suite_name);

CREATE TABLE metadata
(
    run_id String,
		commit_sha String,
		job_name String,
		repo String,
		branch String,
		format String,		
		link String,		
		tags String,
		created_at DateTime
)
ENGINE = MergeTree
PRIMARY KEY (run_id);