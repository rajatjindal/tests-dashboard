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
	repo: string;
	branch: string;
	format: string;
	link: string;
	tags: string;
	createdAt: string;
}

export interface Dataset {
	label: string;
	data: number[];
	stack: string;
	backgroundColor: string;
}

export interface TimeTrendsData {
	labels: string[];
	datasets: Dataset[];
}
