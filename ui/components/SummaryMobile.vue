<template>
	<div class="text-darkmode-blue-contrast1 block md:hidden mt-5 border rounded-lg px-4 py-4 w-full">
		<div class="grid grid-rows w-full">
			<div class="grid grid-cols-4">
				<div class="col-span-1 text-blue-500 font-bold"><a target="_blank"
						 rel="noopener noreferrer"
						 :href="metadata.link">
						<div class="text-darkmode-blue-contrast1">
							<GitHubIcon />
						</div>
					</a></div>

				<div class="col-span-3">
					<div class="flex justify-end">
						<span v-for="tag in getTags(metadata.tags)"
									class="px-2 py-0.5 rounded bg-darkmode-blue-contrast4 text-xs text-green-200 ml-1 wrap">{{ tag }}</span>
					</div>

				</div>
			</div>

			<div class="grid grid-cols-2 mt-5 text-sm border-b border-gray-100">
				<div class="text-left">Passed</div>
				<div class="text-right">{{ summary.passed }}</div>
			</div>

			<div class="grid grid-cols-2 mt-5 text-sm border-b border-gray-100">
				<div class="text-left">Failed</div>
				<div class="text-right">{{ summary.failed }}</div>
			</div>

			<div class="grid grid-cols-2 mt-5 text-sm border-b border-gray-100">
				<div class="text-left">Ignored</div>
				<div class="text-right">{{ summary.ignored }}</div>
			</div>

		</div>
	</div>
</template>
<script lang="ts" setup>
import { Summary, Metadata } from '@/sdk/backend/types';
import { getMetadata } from '@/sdk/backend/api';

const props = defineProps({
	summary: { type: Object as PropType<Summary>, required: true, default: {} }
})

const metadata = ref({} as Metadata)
onBeforeMount(async () => {
	metadata.value = await getMetadata(props.summary.runId)
})

const getTags = function (metadata = ""): string[] { return metadata.split(",") }
</script>