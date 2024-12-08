<template>
	<div class="mx-auto w-11/12">
		<div class="w-full hidden md:grid" v-if="currentRun.runId">
			<Summary :summary="currentRun" :key="currentRun.runId" :showHideButton="false"/>
			<SuiteSummary :runId="currentRun.runId" :showIgnored="showIgnored" :key="currentRun.runId" />
		</div>
	</div>
</template>
<script setup lang="ts">
import { getSummary } from "@/sdk/backend/api";
import type { Summary } from "@/sdk/backend/types";

const currentRun = ref<Summary>({} as Summary)
const showIgnored = ref(false)
const runId: string = useRoute().params["id"] ? useRoute().params["id"].toString() : ""

console.log("run id -> ", runId)
onBeforeMount(async () => {
	currentRun.value = await getSummary(runId)
})
</script>