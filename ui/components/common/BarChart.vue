<script lang="ts" setup>
import { Bar } from 'vue-chartjs'
import { ChartData } from 'chart.js';

const emit = defineEmits(['dataset-clicked'])

const props = defineProps({
	chartData: { type: Object as PropType<ChartData<"bar", (number | [number, number] | null)[], unknown>>, required: true, default: null },
})

const chartOptions = ref({
	responsive: true,
	scales: {
		x: {
			stacked: true,
			showLabels: false,
		},
		y: {
			stacked: true
		}
	},
	maintainAspectRatio: false,
	onClick: function (e: any, elements: any[], chart: any) {
		emit('dataset-clicked', elements[0].index)
	}
})
</script>
<template>
	<div>
		<Bar v-if="chartData"
				 :data="chartData"
				 :options="chartOptions" />
	</div>
</template>