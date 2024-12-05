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
	onClick: function (evt: any, elements: any[], chart: any) {
		const points = chart.getElementsAtEventForMode(evt, 'nearest', {intersect: true}, true)
		if (points.length) {
			const fp = points[0]
			const dataset = fp.datasetIndex
			const datapoint = fp.index
		
			console.log(`${dataset} ${datapoint}`)
			console.log("" + props.chartData.labels![datapoint])
			const obj = {
				index: datapoint,
				label: props.chartData.labels![datapoint],
				datasetLabel: props.chartData.datasets[dataset].label
			}

			console.log(JSON.stringify(obj))
			emit('dataset-clicked', obj)
		}
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