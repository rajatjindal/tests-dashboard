<template>
	<div class="mx-auto w-11/12">
		<div class="w-full hidden md:grid" v-if="currentRun.runId">
			<SummaryMd :summary="currentRun" :key="currentRun.runId" :showHideButton="false"/>
			<RunMd :runId="currentRun.runId" :showIgnored="showIgnored" :key="currentRun.runId" />
		</div>

		<div class="w-full grid md:hidden" v-if="currentRun.runId">
			<SummaryMobile :summary="currentRun" :key="currentRun.runId" />
			<RunMobile :runId="currentRun.runId" :showIgnored="showIgnored" :key="currentRun.runId" />
		</div>
	</div>
</template>
<script setup lang="ts">
import { getSummary } from "@/sdk/backend/api";
import { Summary } from "@/sdk/backend/types";

const currentRun = ref<Summary>({} as Summary)
const showIgnored = ref(false)
const runId: string = useRoute().params["id"].toString()

onBeforeMount(async () => {
	currentRun.value = await getSummary(runId)
})
</script>