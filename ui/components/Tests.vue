<template>
	<div class="mt-5">
		<div class="col-span-1 text-xs my-auto">
			<input type="checkbox"
						 class="rounded mt-2"
						 v-model="showFailedOnly" /><span class="ml-1 text-darkmode-blue-contrast1">Show Failed only</span>
		</div>
	</div>
	<div
			 class="text-darkmode-blue-contrast1 hidden md:grid grid-cols-6 gap-4 border border-darkplum px-3 py-3 text-xs uppercase tracking-wider border-b-0 text-darmplum">
		<div class="col-span-2">Name</div>
		<div class="col-span-1">Result</div>
		<div class="col-span-1">Duration</div>
		<div class="col-span-1">History</div>
		<div class="col-span-1"></div>
	</div>

	<div class="hidden md:grid text-darkmode-blue-contrast1">
		<div v-for="(test, index) in filteredTests"
				 :key="test.name"
				 class="grid grid-cols-6 gap-4 border border-darkplum px-3 py-3 text-xs"
				 :class="{ 'border-b': index === lastIndex, 'border-b-0': index !== lastIndex }">
			<div class="col-span-2">{{ test.name }}</div>
			<div class="col-span-1 flex text-left"
					 v-if="test.result === 'ok'">
				<CheckIcon class="-ml-1 w-4 h-4 text-seagreen font-bold" />
				<div class="text-seagreen">{{ test.result }}</div>
			</div>
			<div class="col-span-1 flex text-left"
					 v-if="test.result === 'failed'">
				<CloseIcon class="-ml-1 w-4 h-4 text-darkmode-red-dark font-bold" />
				<div class="text-darkmode-red-dark">{{ test.result }}</div>
			</div>
			<div class="col-span-1 flex text-left"
					 v-if="test.result === 'ignored'">
				<div class="text-darkplum font-bold">!</div>
				<div class="ml-1 mr-1 text-darkplum">{{ test.result }}</div>
			</div>
			<div class="col-span-1">{{ humanDuration(test.duration) }}</div>
			<div class="col-span-1"><TestHistory :name="test.name" :runId="test.runId"/></div>
			<div class="col-span-1 underline" v-on:click="showLogs(index)">Show Logs</div>
			<div class="col-span-4" v-if="showingLogs === index">
				<p v-for="line in test.logs.split('\n')"
					 class="text-darkplum italic">
					<span v-html="highlight(line)"></span>
				</p>
			</div>
		</div>
	</div>
</template>


<script setup lang="ts">
import { humanDuration } from "@/sdk/base/myfetch";
import { getTestsForRunAndSuite, getTestsForLogLine } from "@/sdk/backend/api";
import { Test } from "@/sdk/backend/types";

const props = defineProps({
	repo: { type: String, required: false, default: "" },
	runId: { type: String, required: false, default: "" },
	suiteId: { type: String, required: false, default: "" },
	suiteName: { type: String, required: false, default: "" },
	tags: { type: Object as PropType<Map<string, string>>, required: false, default: new Map<string, string>() },
	showIgnored: { type: Boolean, default: true },
	forlogs: { type: Boolean, default: false },
	fortestname: { type: Boolean, default: false }
})

const showFailedOnly = ref(true)
const tests = ref<Test[]>()

const showingLogs = ref(-1)
const showLogs = function(index: number) {
	if (showingLogs.value > 0) {
		showingLogs.value = -1
		return
	}

	showingLogs.value = index
}

onBeforeMount(async () => {
	if (props.forlogs) {
		const logLine = useRoute().query["logLine"];
		if (logLine) {
			tests.value = await getTestsForLogLine(props.repo, "", "", logLine.toString(), props.tags)
		}

		return
	}

	if (props.fortestname) {
		const testname = useRoute().query["testname"];
		if (testname) {
			tests.value = await getTestsForLogLine(props.repo, "", "", testname.toString(), props.tags)
		}

		return
	}

	if (props.runId) {
		tests.value = await getTestsForRunAndSuite(props.runId, props.suiteId)
	}
})

const filteredTests = computed(() => tests.value ? tests.value.filter(item => {
	if (showFailedOnly.value && item.result !== 'failed') {
		return false
	}

	if (item.result === 'ignored' && !props.showIgnored) {
		return false
	}

	return true
}) : [])

const lastIndex = computed(() => filteredTests.value ? filteredTests.value.length - 1 : 0)

const highlight = function (line: string): string {
	//remove terminal colors
	const loglineWoutColors = line.toString().replace(/[\u001b\u009b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]/g, "")
	const logline = useRoute().query["logLine"];
	if (!logline) {
		return loglineWoutColors
	}

	var check = new RegExp(logline.toString(), "ig");
	return loglineWoutColors.toString().replace(check, function (matchedText, a, b) {
		return ('<span class="bg-seagreen font-bold px-1 py-0.5 rounded">' + matchedText + '</span>');
	});
}
</script>