<template>
	<div class="mx-auto w-11/12">
		<div
				 class="w-full md:w-1/2 grid grid-cols-3 gap-4 border border-darkplum px-4 py-4 mx-auto rounded bg-indigo-50 shadow-2xl mt-20">
			<div class="col-span-1 text-xs my-auto">
				<input type="checkbox"
							 class="rounded mt-2"
							 v-model="showIgnored" />
				<span class="ml-1 text-darkmode-blue-contrast1">Show Ignored</span>
			</div>

			<InputNumber v-model="maxRuns"
									 label="Number of runs"
									 class="col-span-1" />

			<InputText v-model="service"
									 label="service to fetch results for"
									 class="col-span-1" />
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

const defaultService = "dagger/ci-tests"

const maxRuns = ref(20)
const showIgnored = ref(false)
const runs = ref([] as Summary[])
const currentRun = ref({} as Summary)
const qservice = useRoute().query["service"]
const service = ref(qservice ? qservice.toString() : defaultService)

const lastXruns = computed(() => runs.value ? runs.value.slice(-1 * maxRuns.value) : [])
const lineChartData = computed(() => lastXruns.value ? toLineChart(lastXruns.value) : undefined)

watch(() => service, async (currentValue, oldValue) => {
	runs.value = await getAllRuns(currentValue.value)
})

onBeforeMount(async () => {
	const s = service.value ? service.value : defaultService
	runs.value = await getAllRuns(s)
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
</script>