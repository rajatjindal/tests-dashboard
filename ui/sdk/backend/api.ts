import { Metadata, Summary, Test} from './types'
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

export const getTestsForRunAndSuite = async (runId: string, suiteId: string): Promise<Test[]> => {
	const path = `/api/runs/${runId}/suites/${suiteId}/tests`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getTestsForLogLine = async (logLine: string): Promise<Test[]> => {
	const path = `/api/history/log?logLine=${logLine}`
	return await myfetch(path, {
		method: 'GET'
	});
}

export const getTestsForTestcase = async (testname: string): Promise<Test[]> => {
	const path = `/api/history/test?name=${testname}`
	return await myfetch(path, {
		method: 'GET'
	});
}