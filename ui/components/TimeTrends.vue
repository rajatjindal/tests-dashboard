<template>
	<div class="mx-auto w-11/12">
		<div
				 class="w-full md:w-full gap-4 border border-darkplum px-4 py-4 mx-auto rounded bg-indigo-50 shadow-2xl mt-20">
			<BarChart v-if="timetrends && timetrends.datasets && timetrends.datasets.length > 0" class="w-full h-screen"
								key="timetrends"
                v-on:dataset-clicked="showRunWithIndex"
								:chartData="timetrends" />
		</div>
	</div>
</template>

<script setup lang="ts">
import {  getTimeTrendsForSuites } from "@/sdk/backend/api";
import type {  TimeTrendsData } from '~/sdk/backend/types'

const suiteName = useRoute().query["suiteName"]?.toString() ?? ""
const timetrends = ref({} as TimeTrendsData)

onBeforeMount(async () => {
	timetrends.value = await getTimeTrendsForSuites(suiteName)
})

const showRunWithIndex = async function (_: number, label: string) {
  // console.log("clicked at ", index, " valat:", label)
  await navigateTo(`/timetrends?suiteName=${label}`, {
    open: {
      target: '_blank',
    }
  })
}
</script>