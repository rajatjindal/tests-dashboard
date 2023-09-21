<template class="bg-darkmode-background-contrast2">
  <Trends />
</template>
<script lang="ts" setup>
  console.log(useRoute().query["redirect-original-path"])
  const original = useRoute().query["redirect-original-path"]
  if (original !== null) {
    useRouter().push(original.toString())  
  }
  
</script>
  