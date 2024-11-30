<template class="bg-darkmode-background-contrast2">
  <div class="flex text-darkmode-blue-contrast1 mt-5 ml-20">
    <div :class="{'border rounded-lg px-5 py-2 cursor-pointer': true, 'bg-darkmode-background-contrast1 ': selectedIndex === 0}" v-on:click="changeSelection(0)">Reliability</div>
    <div :class="{'border rounded-lg ml-2 px-5 py-2 cursor-pointer': true, 'bg-darkmode-background-contrast1 ': selectedIndex === 1}" v-on:click="changeSelection(1)">Duration</div>
  </div>
  <Trends v-if="selectedIndex===0"/>
  <TimeTrends v-if="selectedIndex===1"/>
</template>
<script lang="ts" setup>
const selectedIndex = ref(0)
const changeSelection = function(index: number) {
  selectedIndex.value = index
}

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
  