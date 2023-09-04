<template>
  <Button bgcolor="bg-neutral-50"
          bordercolor="border-neutral-400"
          :callback="callback"
          hoverbgcolor="hover:bg-neutral-100"
          hovertextcolor="hover:text-neutral-600"
          :hovereffect="true"
          :icon="icon"
          :text="text"
          textcolor="text-neutral-500" />
</template>

<script setup lang="ts">
const props = defineProps({
  callback: { type: Function as PropType<(event: MouseEvent) => void> },
  icon: { type: String, default: null },
  text: { type: String, default: "Cancel" },
})
</script>