<template>
	<div class="mt-5">
		<div class="text-darkmode-blue-contrast1 mb-10" v-if="useRoute().query['logLine']">
			Showing history of log line <span class="italic font-bold">{{ useRoute().query["logLine"] }}</span>
		</div>
		<div class="text-darkmode-blue-contrast1 mb-10" v-if="useRoute().query['testname']">
			Showing history of testcase <span class="italic font-bold">{{ useRoute().query["testname"] }}</span>
		</div>
		<div class="col-span-1 text-xs my-auto">
			<input type="checkbox"
						 class="rounded mt-2"
						 v-model="showFailedOnly" /><span class="ml-1 text-darkmode-blue-contrast1">Show Failed only</span>
		</div>
	</div>
	<div class="text-darkmode-blue-contrast1 hidden md:grid grid-cols-4 gap-4 border border-darkplum px-3 py-3 text-xs uppercase tracking-wider border-b-0 text-darmplum">
		<div class="col-span-2">Name</div>
		<div class="col-span-1">Result</div>
		<div class="col-span-1">Duration</div>
	</div>

	<div class="hidden md:grid text-darkmode-blue-contrast1">
		<div v-for="(test, index) in filteredTests"
				 :key="test.name"
				 class="grid grid-cols-4 gap-4 border border-darkplum px-3 py-3 text-xs"
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
			<div class="col-span-4">
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
	runId: { type: String },
	showIgnored: { type: Boolean, default: true },
	forlogs: { type: Boolean, default: false},
	fortestname: { type: Boolean, default: false}
})

const showFailedOnly = ref(false)
const tests = ref<Test[]>()

onBeforeMount(async () => {
	if (props.forlogs) {
		const logLine = useRoute().query["logLine"];
		if (logLine) {
			tests.value = await getTestsForLogLine(logLine.toString())
		}
		
		return
	}

	if (props.fortestname) {
		const testname = useRoute().query["testname"];
		if (testname) {
			tests.value = await getTestsForLogLine(testname.toString())
		}
		
		return
	}

	if (props.runId) {
		tests.value = await getTestsForRunAndSuite(props.runId, "0")
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

const highlight = function(line: string): string {
	const logline = useRoute().query["logLine"];
	if (!logline) {
		return line
	}

  var check = new RegExp(logline.toString(), "ig");
  return line.toString().replace(check, function(matchedText,a,b){
      return ('<span class="bg-seagreen font-bold px-1 py-0.5 rounded">' + matchedText + '</span>');
  });
}
</script>