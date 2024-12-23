<template>
	<div class="mx-auto w-11/12">
		<div class="mt-5 text-darkmode-blue-contrast1 border px-6 py-6 mx-auto rounded">
			<CommonQuery v-on:updrepo="updateRepo" v-on:updtags="updateTags" :repo="repo" :tags="tags" />
		</div>
		
		<div class="w-full md:w-full gap-4 border border-darkplum px-4 py-4 mx-auto rounded bg-indigo-50 shadow-2xl mt-10">
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

const defaultRepo = ""
const qrepo = useRoute().query["repo"]
const repo = ref(qrepo ? qrepo.toString() : defaultRepo)
const branch = useRoute().query["branch"]?.toString() ?? ""
const commitSha = useRoute().query["commitSha"]?.toString() ?? ""

const tags = ref(new Map<string, string>())
const suiteName = useRoute().query["suiteName"]?.toString() ?? ""
const timetrends = ref({} as TimeTrendsData)

onBeforeMount(async () => {
	timetrends.value = await getTimeTrendsForSuites(repo.value, branch, commitSha, suiteName, tags.value)
})

const clone = function <T>(item: T): T {
	return JSON.parse(JSON.stringify(item))
}

watch(repo, async (currentValue, oldValue) => {
	console.log("repo watch")
	timetrends.value = await getTimeTrendsForSuites(currentValue, branch, commitSha, suiteName, tags.value)
}, { deep: true })

watch(tags, async (currentValue, oldValue) => {
	console.log("tags watch")
	timetrends.value = await getTimeTrendsForSuites(repo.value, branch, commitSha, suiteName, currentValue)
}, { deep: true })

const showRunWithIndex = async function (obj: {index: string, label: string, datasetLabel: string}) {
  // console.log("clicked at ", index, " valat:", label)
  await navigateTo(`/timetrends?suiteName=${obj.datasetLabel}`, {
    open: {
      target: '_blank',
    }
  })
}

const updateRepo = function(val: string) {
	repo.value = val
}

const updateTags = function(val: Map<string, string>) {
	console.log("final upd tags", val)
	tags.value = val
}
</script>