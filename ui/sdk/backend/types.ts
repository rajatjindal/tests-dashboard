export interface Test {
	runId: string;
	suiteId: string;
	name: string;
	result: string;
	duration: number;
	logs: string;
	createdAt: string;
}

export interface Summary {
	runId: string;
	result: string;
	passed: number;
	failed: number;
	ignored: number;
	duration: number;
	createdAt: string;
};

export interface SuiteSummary {
	runId: string;
	suiteId: string;
	suiteName: string;
	result: string;
	passed: number;
	failed: number;
	ignored: number;
	duration: number;
	createdAt: string;
}

export interface SuiteSummary {
	runId: string;
	suiteId: string;
	suiteName: string;
	result: string;
	passed: number;
	failed: number;
	ignored: number;
	duration: number;
	createdAt: string;
}

export interface Metadata {
	runId: string;
	commitSha: string;
	jobName: string;
	repo: string;
	branch: string;
	format: string;
	link: string;
	tags: Object; //TODO: change it to key=value specific type
	createdAt: string;
}

export interface Dataset {
	label: string;
	data: number[];
	stack: string;
	backgroundColor: string;
}

export interface TimeTrendsData {
	labels: string[][];
	ids: string[];
	datasets: Dataset[];
}

export interface Tag {
	key: string;
	values: string[]
}
