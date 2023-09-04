<template>
	<div class="mt-5">
		<div class="col-span-1 text-xs my-auto">
			<input type="checkbox"
						 class="rounded mt-2"
						 v-model="showFailedOnly" /><span class="ml-1 text-darkmode-blue-contrast1">Show Failed only</span>
		</div>
	</div>
	<div class="grid md:hidden">
		<div class="text-darkmode-blue-contrast1 grid grid-rows border border-darkplum" :class="{'border-b-0': !lastIndex, 'border-b': lastIndex}" v-for="(test, index) in filteredTests" :key="test.name">
			<div class="w-full px-1 py-2 text-xs">
				<div class="flex ml-0.5">
					<CheckIcon class="-ml-1 w-5 h-5 text-seagreen font-bold" v-if="test.result === 'ok'"/>
					<CloseIcon class="-ml-1 w-5 h-5 text-darkmode-red-dark font-bold" v-else-if="test.result === 'failed'"/>
					<div class="text-gray-400 font-bold" v-if="test.result === 'ignored'">!</div>
					<div class="ml-1">{{ test.name }}</div>
				</div>
				<div>
					<p v-for="line in test.logs.split('\n')"
							class="text-darkplum italic">
						{{ line }}
					</p>
				</div>
			</div>
		</div>
	</div>
</template>


<script setup lang="ts">
import { getTestsForRunAndSuite } from "@/sdk/backend/api";
import { Test } from "@/sdk/backend/types";

const props = defineProps({
	runId: { type: String, required: true },
	showIgnored: { type: Boolean, required: true }
})

const showFailedOnly = ref(false)
const tests = ref<Test[]>()

onBeforeMount(async () => {
	tests.value = await getTestsForRunAndSuite(props.runId, "0")
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
</script>