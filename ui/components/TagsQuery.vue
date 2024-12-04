<template>
	<div class="">
		<div class="flex w-1/4">
			<StaticSelect class="w-1/2"
										v-bind:modelValue="currentKey"
										key="key"
										v-on:change="updateKey"
										:options="tagKeys"
										label="Tag Key" />

			<StaticSelect class="ml-2 w-1/2"
										v-bind:modelValue="currentValue"
										key="value"
										:editable="currentKey !== ''"
										v-on:change="updateValue"
										:options="tagValues"
										label="Tag Value" />
		</div>
	</div>
</template>
<script setup lang="ts">
import { getTags } from "@/sdk/backend/api";
import type { Tag } from '~/sdk/backend/types'

const emit = defineEmits(['updtags'])

const props = defineProps({
	repo: { type: String, required: false, default: "dagger/ci-tests" }
})

const currentValue = ref("")
const currentKey = ref("")
const selectedTags = ref(new Map<string, string>());

const tags = ref([] as Tag[])
onBeforeMount(async () => {
	if (!props.repo) {
		return
	}
	
	tags.value = await getTags(props.repo)
})

const tagKeys = computed(() => tags.value.reduce((accumulator: string[], current: Tag) => {
	if (!selectedTags.value.has(current.key)) {
		accumulator.push(current.key);
	}

	return accumulator;
}, [""]))

const findValuesOrDefault = function (key: string): string[] {
	const tag = tags.value.find(current => current.key === key)
	if (tag) {
		return ["", ...tag.values]
	}

	return []
}

const tagValues = computed(() => currentKey.value === "" ? [""] : findValuesOrDefault(currentKey.value))
const updateKey = function (event: Event) {
	currentKey.value = (event.target as HTMLInputElement).value
}

const updateValue = function (event: Event) {
	selectedTags.value.set(currentKey.value, (event.target as HTMLInputElement).value)
	emit('updtags', selectedTags.value)
	currentKey.value = ""
	console.log("selectedTags inside TagsQuery", selectedTags.value)
}
</script>