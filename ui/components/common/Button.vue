<template>
  <button v-on:click="onClick"
          :disabled="loading"
          :class="{ [bgcolor]: true, [bordercolor]: true, [textcolor]: true, [hoverbgcolor]: hovereffect, [hovertextcolor]: hovereffect }"
          class="py-1 px-3 md:px-4 border border-transparent text-xs font-semibold rounded-md hover:shadow-sm flex">
    <component v-if="icon && icon !== ''"
               :is="computedComponent"
               class="w-4 h-4"></component>
    <span v-if="text && text !== ''">{{ loading ? 'Submitting...' : text }}</span>

  </button>
</template>

<script setup lang="ts">
import { defineAsyncComponent } from "vue";

const props = defineProps({
  bgcolor: { type: String, default: 'bg-primary-50' },
  bordercolor: { type: String, default: 'border-primary-500' },
  callback: { type: Function as PropType<(event: MouseEvent) => void> },
  hoverbgcolor: { type: String, default: 'hover:bg-primary-200' },
  hovertextcolor: { type: String, default: 'hover:text-neutral-1000' },
  hovereffect: { type: Boolean, default: true, },
  icon: { type: String, default: null },
  text: { type: String, default: "Save" },
  textcolor: { type: String, default: 'text-primary-800' },
})


const loading = ref(false);
const startLoading = function () {
  loading.value = true
}

const isLoading = function (): boolean {
  return loading.value
}

const resetLoading = function () {
  setTimeout(() => {
    loading.value = false
  }, 500)
}

const runCallback = async function (event: MouseEvent) {
  if (!props.callback) {
    return
  }

  try {
    await props.callback(event)
  } catch { }
}

const onClick = async function (event: MouseEvent) {
  if (isLoading()) {
    return
  }

  startLoading()
  await runCallback(event)
  resetLoading()
}

const computedComponent = computed(() => defineAsyncComponent(() => import(`../../components/icons/${props.icon}.vue`)))
</script>