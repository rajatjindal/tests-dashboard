<template>
	<div
			 class="grid grid-cols-9 md:grid-cols-9 border border-darkplum text-darkmode-blue-contrast1 px-2 py-3 text-xs uppercase tracking-wider border-b-0 text-darkplum mt-20">
		<div class="col-span-2">Suite name</div>
		<div class="col-span-1 flex">
			<SortableDiv :key="'total' + sortBy"
									 :sortingBy="sortBy"
									 label="total"
									 v-on:upd-sortby="updateSortBy"
									 v-on:upd-sortdirection="updateSortDirection" />
		</div>
		<div class="col-span-1 flex">
			<SortableDiv :key="'passed' + sortBy"
									 :sortingBy="sortBy"
									 label="passed"
									 v-on:upd-sortby="updateSortBy"
									 v-on:upd-sortdirection="updateSortDirection" />
		</div>
		<div class="col-span-1 flex">
			<SortableDiv :key="'failed' + sortBy"
									 :sortingBy="sortBy"
									 label="failed"
									 v-on:upd-sortby="updateSortBy"
									 v-on:upd-sortdirection="updateSortDirection" />
		</div>
		<div class="col-span-1 flex">
			<SortableDiv :key="'ignored' + sortBy"
									 :sortingBy="sortBy"
									 label="ignored"
									 v-on:upd-sortby="updateSortBy"
									 v-on:upd-sortdirection="updateSortDirection" />
		</div>
		<div class="col-span-1 flex">
			<SortableDiv :key="'duration' + sortBy"
									 :sortingBy="sortBy"
									 label="duration"
									 v-on:upd-sortby="updateSortBy"
									 v-on:upd-sortdirection="updateSortDirection" />
		</div>
		<div class="col-span-1">Trends</div>
		<div class="col-span-1"></div>
	</div>


	<div v-for="summary in sortedSuites"
			 :key="summary.suiteId"
			 class="grid grid-cols-9 md:grid-cols-9 border border-darkplum px-2 py-3 text-xs text-darkmode-blue-contrast1">
		<div class="col-span-2">{{ summary.suiteName }}</div>
		<div class="col-span-1">{{ summary.passed + summary.failed + summary.ignored }}</div>
		<div class="col-span-1">{{ summary.passed }}</div>
		<div class="col-span-1">{{ summary.failed }}</div>
		<div class="col-span-1">{{ summary.ignored }}</div>
		<div class="col-span-1">{{ (Math.round(summary.duration * 100) / 100).toFixed(2) }}s</div>
		<div class="col-span-1">
			<NuxtLink target="_blank"
								:to="'/timetrends?suiteName=' + summary.suiteName">
				<GraphIcon class="w-6 h-6 text-green-600" />
			</NuxtLink>
		</div>
		<div class="col-span-1"
					v-if="!currentSuiteId"
				 v-on:click="updateCurrentSuiteRun(summary.runId, summary.suiteId, summary.suiteName)">Show Details</div>
	
	<div class="col-span-1"
					v-else
				 v-on:click="updateCurrentSuiteRun(summary.runId, '', '')">Hide</div>
				</div>

	<div class="grid grid-cols-12" v-if="currentSuiteId">
		<div class="col-span-12">
			<Tests :runId="currentRunId"
						 :suiteId="currentSuiteId"
						 :suiteName="currentSuiteName"
						 :showIgnored="showIgnored"
						 :key="currentRunId + currentSuiteId" />
		</div>
	</div>

</template>
<script lang="ts" setup>
import { formatDate, humanDuration } from '@/sdk/base/myfetch';
import { SuiteSummary, Metadata } from '@/sdk/backend/types';
import { getSuiteSummaryForRun } from '@/sdk/backend/api';

const emit = defineEmits(['hide-details'])

const props = defineProps({
	runId: { type: String },
	showIgnored: { type: Boolean, default: true },
})

const currentRunId = ref("")
const currentSuiteId = ref("")
const currentSuiteName = ref("")

const sortDirection = ref(-1)
const sortBy = ref("duration")

const updateSortDirection = function (val: number) {
	console.log("sortDir callsed => ", val)
	sortDirection.value = val
}

const updateSortBy = function (val: string) {
	console.log("sortBy callsed => ", val)
	sortBy.value = val
}


const showFailedOnly = ref(false)
const suites = ref<SuiteSummary[]>()

onBeforeMount(async () => {
	if (props.runId) {
		suites.value = await getSuiteSummaryForRun(props.runId)
	}
})

const sortedSuites = computed(() => suites.value?.filter(item => currentSuiteId.value === '' || currentSuiteId.value === item.suiteId ).sort(function (a: SuiteSummary, b: SuiteSummary): number {
	var cmp = 0
	switch (sortBy.value) {
		case "duration":
			cmp = a.duration - b.duration
			break;
		case "total":
			cmp = (a.passed + a.failed + a.ignored) - (b.passed + b.failed + b.ignored)
			break;
		case "passed":
			cmp = a.passed - b.passed
			break;
		case "failed":
			cmp = a.failed - b.failed
			break;
		case "ignored":
			cmp = a.ignored - b.ignored
			break;
	}

	// asc
	if (sortDirection.value === 2) {
		return cmp
	}

	// descending
	return cmp * -1
}))

const updateCurrentSuiteRun = function (runId: string, suiteId: string, suiteName: string) {
	currentRunId.value = runId
	currentSuiteId.value = suiteId
	currentSuiteName.value = suiteName
}

// const metadata = ref({} as Metadata)
// onBeforeMount(async () => {
// 	metadata.value = await getMetadata(props.summary.runId)
// })

// const hideDetails = function () { 
// 	emit('hide-details', '')
// }

// const getTags = function (metadata = ""): string[] { return metadata.split(",") }
</script>