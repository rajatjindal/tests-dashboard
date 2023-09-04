<template>
  <div class="grid grid-cols-12 bg-darkmode-background-contrast2">
    <top-navigation class="col-span-12" />
  </div>
  <div class="w-full mx-auto bg-darkmode-background-contrast2 h-screen">
    <slot />
  </div>

</template>

<style>
body {
  background-color: #061632;
}
</style>