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

		<div class="flex">
		<div v-for="(value, index) in selectedTags" :key="index" class="text-white border py-1 mt-5 rounded flex ml-2">
			<span class="px-2">{{ value[0] }}: {{ value[1] }}</span>
			<span class="" v-on:click="removeIndex(value[0])"><CloseIcon class="w-4 h-4"/></span>
		</div>
	</div>
	</div>
</template>
<script setup lang="ts">
import { getTags } from "@/sdk/backend/api";
import type { Tag } from '~/sdk/backend/types'

const currentValue = ref("")
const currentKey = ref("")
const selectedTags = ref(new Map<string, string>());

const tags = ref([] as Tag[])
onBeforeMount(async () => {
	tags.value = await getTags()
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
	currentKey.value = ""
}

const removeIndex = function(key: string) {
	selectedTags.value.delete(key)
}
</script>