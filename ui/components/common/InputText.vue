<template>
	<div class="w-full">
		<label v-if="showLabel && label !== ''"
					 class="block text-xs text-darkmode-blue-contrast1">
			{{ label }}
		</label>
		<input type="text"
					 :disabled="!editable"
					 class="w-full border border-neutral-300 px-2 py-2 rounded-r-md text-xs focus:outline-none focus:border-2 focus:ring-primary-200 focus:border-primary-200 placeholder-neutral-400 focus:placeholder-neutral-50"
					 :placeholder="placeholder"
					 v-on:input="updateInput"
					 v-on:focus="toggleFocused"
					 v-on:blur="toggleFocused"
					 :value="modelValue"
					 v-bind:modelValue="modelValue" />
	</div>
</template>
<script setup lang="ts">
const emit = defineEmits(['update:modelValue', 'focus', 'blur'])

const props = defineProps({
	modelValue: { type: String, required: true, default: "" },
	editable: { type: Boolean, default: true },
	placeholder: { type: String, default: "" },
	label: { type: String, default: "" },
	showLabel: { type: Boolean, default: true }
})

const isFocused = ref(false);
const toggleFocused = function () {
	isFocused.value = !isFocused.value
	if (!isFocused.value) {
		emit('blur')
		return
	}

	emit('focus')
}

const updateInput = function (event: Event) {
	const el = event.target as HTMLInputElement;
	emit('update:modelValue', el.value)
}
</script>