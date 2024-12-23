<template>
	<div class="text-darkmode-blue-contrast1 flex">
		<div v-for="test in filteredTests" :key="test.name" class="ml-1 flex">
			<span class="border rounded border-text-seagreen" v-if="test.result === 'passed'"><NuxtLink target="_blank" :to="'/run/' + test.runId"><CheckIcon class="w-4 h-4 text-seagreen"/></NuxtLink></span>
			<span class="border rounded border-darkmode-red-dark" v-if="test.result === 'failed'"><NuxtLink target="_blank" :to="'/run/' + test.runId"><CloseIcon class="w-4 h-4 text-darkmode-red-dark"/></NuxtLink></span>
		</div>
	</div>
</template>


<script setup lang="ts">
import { getTestsForTestcase } from "@/sdk/backend/api";
import { Test } from "@/sdk/backend/types";

const props = defineProps({
	runId: { type: String, required: false, default: "" },
	repo: { type: String, required: false, default: "dagger/dagger" },
	name: { type: String, required: false, default: "" },
	tags: { type: Object as PropType<Map<string, string>>, required: false, default: new Map<string, string>() },
})

const tests = ref<Test[]>()
onBeforeMount(async () => {
	tests.value = await getTestsForTestcase(props.repo, "", "", props.name, props.tags)
})

const filteredTests = computed(() => tests.value?.filter(item => item.runId !== props.runId))
</script>