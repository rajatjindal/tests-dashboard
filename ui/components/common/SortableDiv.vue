<template>
	<div class="col-span-1 flex hover:cursor-pointer"
			 v-on:click="changeSort">
		<span>{{ label }}</span>
		<span v-if="sortingBy === label">
			<ArrowDown v-if="sortDirection === 1"
								 class="w-5 h-5 font-bold" />
			<ArrowUp v-if="sortDirection === 2"
							 class="w-5 h-5 font-bold" />
		</span>
	</div>
</template>
<script lang="ts" setup>
const props = defineProps({
	label: { type: String },
	sortingBy: { type: String, default: ""}
})

// 0 -> disable
// 1 -> desc
// 2 -> asc
const defaultSortDirection = 1;
const sortDirection = ref(defaultSortDirection) // desc
const emit = defineEmits(['upd-sortby', 'upd-sortdirection'])
const changeSort = function() {
	if (props.sortingBy !== props.label) {
		emit('upd-sortby', props.label)
		emit('upd-sortdirection', defaultSortDirection)
		return
	}

	if (sortDirection.value === 2) {
		sortDirection.value = 0
	} else {
		sortDirection.value = sortDirection.value + 1
	}
	
	emit('upd-sortby', props.label)
	emit('upd-sortdirection', sortDirection.value)
}
</script>