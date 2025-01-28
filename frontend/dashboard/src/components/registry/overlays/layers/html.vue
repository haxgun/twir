<!-- eslint-disable vue/no-v-html -->
<!-- eslint-disable no-undef -->
<script setup lang="ts">
import { useIntervalFn } from '@vueuse/core'
import { transform as transformNested } from 'nested-css-to-flat'
import { computed, nextTick, ref, watch } from 'vue'

import { useCustomOverlaysApi } from '@/api/custom-overlays'

const props = defineProps<{
	index: number | string
	transformString: string
	width: number
	height: number
	text: string
	css: string
	js: string
	periodicallyRefetchData: boolean
}>()

const api = useCustomOverlaysApi()
const parser = api.useParseVariablesInHtml()

const exampleValue = ref('')

const { pause, resume } = useIntervalFn(async () => {
	const { data } = await parser.executeMutation({ text: props.text })
	exampleValue.value = data?.customOverlaysParseHtml ?? ''
}, 1000, { immediate: true, immediateCallback: true })

const executeFunc = computed(() => {
	// eslint-disable-next-line no-new-func
	return new Function(`${props.js}; onDataUpdate();`)
})

watch(exampleValue, async () => {
	await nextTick()
	// calling user defined function
	executeFunc.value?.()
})

watch(() => props.periodicallyRefetchData, (v) => {
	if (!v) pause()
	else resume()
}, { immediate: true })
</script>

<template>
	<div
		:id="`layer-${index}`"
		class="absolute overflow-hidden text-nowrap"
		:style="{
			transform: transformString,
			width: `${width}px`,
			height: `${height}px`,
		}"
	>
		<component is="style">
			{{
				transformNested(`#layersExampleRender${index} {
					${css}
				}`)
			}}
		</component>

		<div :id="`layersExampleRender${index}`" class="w-full h-full" v-html="exampleValue" />
	</div>
</template>
