<template>
	<div
			 class="grid grid-cols-8 border border-darkplum text-darkmode-blue-contrast1 px-2 py-3 text-xs uppercase tracking-wider border-b-0 text-darkplum mt-20">
		<div class="col-span-1">Link</div>
		<div class="col-span-1">Total</div>
		<div class="col-span-1">Passed</div>
		<div class="col-span-1">Failed</div>
		<div class="col-span-1">Ignored</div>
		<div class="col-span-1">Duration</div>
		<div class="col-span-2 text-left">On</div>
	</div>
	<div class="grid">
		<div class="grid grid-cols-8 border border-b-0 border-darkplum px-2 py-3 text-xs text-darkmode-blue-contrast1">
			<div class="col-span-1 font-bold">
				<a target="_blank"
					 rel="noopener noreferrer"
					 :href="metadata.link">
					<div class="text-darkmode-blue-contrast1">
						<GitHubIcon />
					</div>
				</a>
			</div>

			<div class="col-span-1">{{ summary.passed + summary.failed + summary.ignored }}</div>
			<div class="col-span-1">{{ summary.passed }}</div>
			<div class="col-span-1">{{ summary.failed }}</div>
			<div class="col-span-1">{{ summary.ignored }}</div>
			<div class="col-span-1">{{ humanDuration(summary.duration) }}</div>
			<div class="col-span-2 text-left">{{ formatDate(summary.createdAt) }}</div>
			</div>
			<div class="grid grid-cols-4 text-xs text-darkmode-blue-contrast1 border border-darkplum  border-t-0 px-2 py-2">
			
			<div class="col-span-4 flex justify-between">

				<div class="flex justify-start">
					<span v-for="(value, key) in metadata.tags"
								class="px-2 py-1 rounded bg-darkmode-blue-contrast3 text-darkgreen ml-2 wrap">{{ key }}: {{ value }}</span>
				</div>

				<div class="flex justify-items-end">
					<CancelButton v-if="showHideButton" text="Hide" v-on:click="hideDetails"/>
				</div>
			</div>

		</div>
	</div>
</template>
<script lang="ts" setup>
import { formatDate, humanDuration } from '@/sdk/base/myfetch';
import { Summary, Metadata } from '@/sdk/backend/types';
import { getMetadata } from '@/sdk/backend/api';

const emit = defineEmits(['hide-details'])

const props = defineProps({
	summary: { type: Object as PropType<Summary>, required: true, default: {} },
	showHideButton: { type: Boolean, default: true}
})

const metadata = ref({} as Metadata)
onBeforeMount(async () => {
	metadata.value = await getMetadata(props.summary.runId)
})

const hideDetails = function () { 
	emit('hide-details', '')
}
</script>