<template>
	<div class="mx-auto w-11/12 text-darkmode-background-contrast2">
		
		<div class="mt-5 text-darkmode-blue-contrast1 border px-6 py-6 mx-auto rounded">
			repo -> {{ repo }}
			<CommonQuery v-on:updrepo="updateRepo" v-on:updtags="updateTags" :repo="repo" :tags="tags" />
		</div>
		{{ suites }}
		</div>

</template>
<script setup lang="ts">
import { getTopNSlowestTestSuites } from '~/sdk/backend/api';
import { SuiteSummary } from '~/sdk/backend/types';

// only suite-name and duration is populated for now by the bckend
const suites = ref([] as SuiteSummary[])
const tags = ref(new Map<string, string>())
const repo = ref( useRoute().query["repo"] ? useRoute().query["repo"]!.toString() : "")
const branch: string = useRoute().query["branch"] ? useRoute().query["branch"]!.toString() : ""
const commitSha: string = useRoute().query["commitSha"] ? useRoute().query["commitSha"]!.toString() : ""

// onBeforeMount(async () => {
// 	suites.value = await getTopNSlowestTestSuites(repo.value, branch, commitSha, tags.value)
// })

watch(repo, async (currentValue, oldValue) => {
	console.log("repo watch")
	suites.value = await getTopNSlowestTestSuites(currentValue, branch, commitSha, tags.value)
}, { deep: true })

watch(tags, async (currentValue, oldValue) => {
	console.log("tags watch")
	suites.value = await getTopNSlowestTestSuites(repo.value, branch, commitSha, currentValue)
}, { deep: true })

const updateRepo = function(val: string) {
	console.log("final repo val", val)
	repo.value = val
}

const updateTags = function(val: Map<string, string>) {
	console.log("final upd tags", val)
	tags.value = val
}
</script>