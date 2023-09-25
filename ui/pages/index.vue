<template class="bg-darkmode-background-contrast2">
  <Trends />
</template>
<script lang="ts" setup>
if (useRoute() !== null && useRoute().query !== null) {
  console.log(useRoute().query)
  const original = useRoute().query["redirect-original-path"]
  if (original && original != '') {
    let x = original + "?"
    for (const key in useRoute().query) {
      if (key === 'redirect-original-path') {
        continue
      }
      
      console.log("key is ", key)
      x += `${key}=${useRoute().query[key]}&`
    }
    useRouter().push(x.toString())
  }
}


</script>
  