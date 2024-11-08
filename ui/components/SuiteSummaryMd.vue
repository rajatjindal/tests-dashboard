<template>
	<div
			 class="grid grid-cols-6 border border-darkplum text-darkmode-blue-contrast1 px-2 py-3 text-xs uppercase tracking-wider border-b-0 text-darkplum mt-20">
		<div class="col-span-1">Suite name</div>
		<div class="col-span-1">Total</div>
		<div class="col-span-1">Passed</div>
		<div class="col-span-1">Failed</div>
		<div class="col-span-1">Ignored</div>
		<div class="col-span-1"></div>
	</div>


	<div v-for="summary in suites"
			 :key="summary.suiteId"
			 class="grid grid-cols-6 border border-darkplum px-2 py-3 text-xs text-darkmode-blue-contrast1">
		<div class="col-span-1">{{ summary.suiteName }}</div>
		<div class="col-span-1">{{ summary.passed + summary.failed + summary.ignored }}</div>
		<div class="col-span-1">{{ summary.passed }}</div>
		<div class="col-span-1">{{ summary.failed }}</div>
		<div class="col-span-1">{{ summary.ignored }}</div>
		<div class="col-span-1"
				 v-on:click="updateCurrentSuiteRun(summary.runId, summary.suiteId, summary.suiteName)">details</div>
	</div>

	<div class="grid grid-cols-12">
		<div class="col-span-12">
			<RunMobile :runId="currentRunId"
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

const showFailedOnly = ref(false)
const suites = ref<SuiteSummary[]>()

onBeforeMount(async () => {
	if (props.runId) {
		suites.value = await getSuiteSummaryForRun(props.runId)
	}
})

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