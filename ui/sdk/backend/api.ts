import { Metadata, Summary, SuiteSummary, Test, TimeTrendsData, Tag } from './types'
import { myfetch } from '@/sdk/base/myfetch';
export const getAllRuns = async (repo: string): Promise<Summary[]> => {
	const path = `/api/runs?repo=${repo}`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getMetadata = async (runId: string): Promise<Metadata> => {
	const path = `/api/runs/${runId}/metadata`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getSummary = async (runId: string): Promise<Summary> => {
	const path = `/api/runs/${runId}/summary`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getSuitesForRun = async (runId: string): Promise<string[]> => {
	const path = `/api/runs/${runId}/suites`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getSuiteSummaryForRun = async (runId: string): Promise<SuiteSummary[]> => {
	const path = `/api/runs/${runId}/suites-summary`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getTestsForRunAndSuite = async (runId: string, suiteId: string): Promise<Test[]> => {
	const path = `/api/runs/${runId}/suites/${suiteId}/tests`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getTestsForLogLine = async (repo: string, logLine: string, tags: Map<string, string>): Promise<Test[]> => {
	const path = `/api/history/log?logLine=${logLine}&repo=${repo}&${tagsToQueryString(tags)}`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getTestsForTestcase = async (repo: string, testname: string, tags: Map<string, string>): Promise<Test[]> => {
	const path = `/api/history/test?name=${testname}&repo=${repo}&${tagsToQueryString(tags)}`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getTimeTrendsForSuites = async (repo: string, suiteName: string, tags: Map<string, string>): Promise<TimeTrendsData> => {
	const path = `/api/trends/suites/time?suiteName=${suiteName}&repo=${repo}&${tagsToQueryString(tags)}`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getTags = async (repo: string): Promise<Tag[]> => {
	return await myfetch(`/api/tags?repo=${repo}`, {
		method: 'GET'
	});
}

const tagsToQueryString = function(tags: Map<string, string>): string {
	var slice: string[] = []
	for (let [key, value] of tags) {
		slice.push(`tag-${key}=${value}`)
	}

	return slice.join("&")
}