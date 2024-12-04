<template>
	<div class="mx-auto w-11/12">
		<div class="mt-5 text-darkmode-blue-contrast1 border px-6 py-6 mx-auto rounded">
			<CommonQuery v-on:updrepo="updateRepo" v-on:updtags="updateTags" :repo="repo" :tags="tags" />
		</div>

		<BarChart class="w-full h-40 mt-10"
							:key="JSON.stringify(lineChartData)"
							v-on:dataset-clicked="showRunWithIndex"
							:chartData="lineChartData" />

		<div class="w-full hidden md:grid"
				 v-if="currentRun.runId">
			<Summary :summary="currentRun"
								 v-on:hide-details="resetCurrentRun"
								 :key="currentRun.runId" />
			<SuiteSummary :runId="currentRun.runId"
						 :showIgnored="showIgnored"
						 :key="currentRun.runId" />
		</div>
	</div>
</template>
  
<script setup lang="ts">
import { formatDate } from "@/sdk/base/myfetch";
import { ChartData } from "chart.js";
import { getAllRuns } from "@/sdk/backend/api";
import type { Summary } from '~/sdk/backend/types'
import type {  TimeTrendsData } from '~/sdk/backend/types'

const defaultRepo = ""
const qrepo = useRoute().query["repo"]
const repo = ref(qrepo ? qrepo.toString() : defaultRepo)

const branch = useRoute().query["branch"]?.toString() ?? ""
const commitSha = useRoute().query["commitSha"]?.toString() ?? ""

const tags = ref(new Map<string, string>())
const suiteName = useRoute().query["suiteName"]?.toString() ?? ""
const timetrends = ref({} as TimeTrendsData)

const maxRuns = ref(20)
const showIgnored = ref(false)
const runs = ref([] as Summary[])
const currentRun = ref({} as Summary)

const lastXruns = computed(() => runs.value ? runs.value.slice(-1 * maxRuns.value) : [])
const lineChartData = computed(() => lastXruns.value ? toLineChart(lastXruns.value) : undefined)

watch(() => repo, async (currentValue, oldValue) => {
	runs.value = await getAllRuns(currentValue.value, branch, commitSha, tags.value)
})

onBeforeMount(async () => {
	runs.value = await getAllRuns(repo.value, branch, commitSha, tags.value)
})

const showRunWithIndex = function (index: number) {
	const valAt = lastXruns.value.at(index)
	if (!valAt) {
		return
	}

	currentRun.value = valAt
}

const resetCurrentRun = function () {
	currentRun.value = {} as Summary
}

const toLineChart = function (runs: Summary[]): ChartData<"bar", (number | [number, number] | null)[], unknown> {
	const labels = runs.reduce(function (result: string[], item: Summary): string[] {
		result.push(formatDate(item.createdAt));
		return result
	}, [])

	const passed = runs.reduce(function (result: number[], item: Summary): number[] {
		result.push(item.passed);
		return result
	}, [])

	const failed = runs.reduce(function (result: number[], item: Summary): number[] {
		result.push(item.failed);
		return result
	}, [])

	const ignored = runs.reduce(function (result: number[], item: Summary): number[] {
		result.push(item.ignored);
		return result
	}, [])

	return {
		labels: labels,
		datasets: [
			{
				label: 'passed',
				backgroundColor: '#34e8bd',
				data: passed,
			},
			{
				label: 'failed',
				backgroundColor: '#CE5050',
				data: failed,
			},
			{
				label: 'ignored',
				backgroundColor: '#9ca3af',
				data: ignored,
				hidden: !showIgnored.value,
			},
		],
	}
}

const updateRepo = function(val: string) {
	repo.value = val
}

const updateTags = function(val: Map<string, string>) {
	console.log("final upd tags", val)
	tags.value = val
}

watch(repo, async (currentValue, oldValue) => {
	console.log("repo watch")
	runs.value = await getAllRuns(currentValue, branch, commitSha, tags.value)
}, { deep: true })

watch(tags, async (currentValue, oldValue) => {
	console.log("tags watch")
	runs.value = await getAllRuns(repo.value,branch, commitSha, currentValue)
}, { deep: true })
</script>