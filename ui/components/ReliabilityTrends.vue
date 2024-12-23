<template>
	<div class="mx-auto w-11/12">
		<div class="mt-5 text-darkmode-blue-contrast1 border px-6 py-6 mx-auto rounded">
			<CommonQuery v-on:updrepo="updateRepo"
									 v-on:updtags="updateTags"
									 v-on:updbranch="updateBranch"
									 :repo="repo"
									 :branch="branch"
									 :tags="tags" />
		</div>

		<div class="w-full md:w-full  gap-4 border border-darkplum px-4 py-4 mx-auto rounded bg-indigo-50 shadow-2xl mt-10">
			<BarChart v-if="timetrends && timetrends.datasets && timetrends.datasets.length > 0"
								class="w-full h-96"
								key="timetrends"
								v-on:dataset-clicked="showRunWithIndex"
								:chartData="timetrends" />
		</div>
	</div>
</template>

<script setup lang="ts">
import { getReliabilityTrendsForSuites } from "@/sdk/backend/api";
import type { TimeTrendsData } from '~/sdk/backend/types'

const repo = ref(useRoute().query["repo"]?.toString() ?? "")
const branch = ref(useRoute().query["branch"]?.toString() ?? "")
const commitSha = useRoute().query["commitSha"]?.toString() ?? ""
const suiteName = useRoute().query["suiteName"]?.toString() ?? ""

const tags = ref(new Map<string, string>())
const timetrends = ref({} as TimeTrendsData)

onBeforeMount(async () => {
	timetrends.value = await getReliabilityTrendsForSuites(repo.value, branch.value, commitSha, suiteName, tags.value)
})

// watch fns
watch(repo, async (currentValue, oldValue) => {
	timetrends.value = await getReliabilityTrendsForSuites(currentValue, branch.value, commitSha, suiteName, tags.value)
}, { deep: true })

watch(branch, async (currentValue, oldValue) => {
	timetrends.value = await getReliabilityTrendsForSuites(repo.value, currentValue, commitSha, suiteName, tags.value)
}, { deep: true })

watch(tags, async (currentValue, oldValue) => {
	console.log("tags watch")
	timetrends.value = await getReliabilityTrendsForSuites(repo.value, branch.value, commitSha, suiteName, currentValue)
}, { deep: true })

const showRunWithIndex = async function (obj: { index: number, label: string, datasetLabel: string }) {
	const id = timetrends.value.ids[obj.index]
	// console.log("clicked at ", index, " valat:", label)
	if (commitSha) {
		await navigateTo(`/run/${id}`, {
			// open: {
			// 	target: '_blank',
			// }
		})
		return
	}

	await navigateTo(`/reliability/commit?repo=${repo.value}&branch=${branch.value}&commitSha=${id}`, {
		// open: {
		// 	target: '_blank',
		// }
	})
}

// update fns
const updateRepo = function (val: string) {
	repo.value = val
}

const updateTags = function (val: Map<string, string>) {
	tags.value = val
}

const updateBranch = function (val: string) {	
	branch.value = val
}
</script>