<template>
	<div class="mt-5">
		<div class="w-full hidden md:grid"
				 v-if="runId">
			<Summary :runId="runId"
							 :key="runId" />
		</div>
	</div>
</template>

<script setup lang="ts">
const runId: string = useRoute().params["id"] ? useRoute().params["id"].toString() : ""
console.log(runId)
</script>