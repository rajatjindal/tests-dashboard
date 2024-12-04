<template>
	<div class="text-darkmode-blue-contrast1 w-full grid grid-cols-4 gap-4 ">
		<InputText class="col-span-1"
							 v-bind:model-value="repo"
							 v-on:change="updateRepo"
							 label="Repo" />
		<InputText class="col-span-1"
							 v-bind:model-value="branch"
							 v-on:change="updateBranch"
							 label="Branch" />
		<TagsQuery class="col-span-2"
							 :repo="repo"
							 :key="repo"
							 :value="tags"
							 v-on:updtags="updateTags" />
	</div>
	<div class="text-darkmode-blue-contrast1 w-full flex">
		<div v-for="(value, key) in selectedTags"
				 :key="key"
				 class="border py-1 mt-5 rounded flex ml-2 text-darkmode-blue-contrast1">
			<span class="px-2">{{ value[0] }}: {{ value[1] }}</span>
			<span class=""
						v-on:click="removeIndex(value[0])">
				<CloseIcon class="w-4 h-4" />
			</span>
		</div>
	</div>
</template>
<script setup lang="ts">
const emit = defineEmits(['updtags', 'updrepo', 'updbranch'])

const props = defineProps({
	repo: { type: String },
	branch: { type: String },
	tags: { type: Map as PropType<Map<string, string>> },
})

const updateRepo = function (event: Event) {
	console.log("updateRepo called ", (event.target as HTMLInputElement).value)
	emit('updrepo', (event.target as HTMLInputElement).value)
}

const updateBranch = function (event: Event) {
	console.log("updateBranch called ", (event.target as HTMLInputElement).value)
	emit('updbranch', (event.target as HTMLInputElement).value)
}

const selectedTags = ref(new Map<string, string>());
const updateTags = function (value: Map<string, string>) {
	console.log("updateTags called ", value)
	selectedTags.value = value
	emit('updtags', value)
}

const removeQueryParam = function (k: string) {
	const query = { ...useRoute().query }; // Get the current query params
	delete query[k];

	useRouter().push({
		path: useRoute().path,
		query: {
			...query,
		},
	});
}

const removeIndex = function (key: string) {
	selectedTags.value.delete(key)
	removeQueryParam("tag-" + key)
}
</script>