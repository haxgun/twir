<script setup lang="ts">
import { transform } from 'nested-css-to-flat'
import { computed, nextTick, watch } from 'vue'

import type { Layer } from '@/composables/overlays/use-overlays.js'

const props = defineProps<{
	layer: Layer
	parsedData?: string
}>()

const executeFunc = computed(() => {
	if (!props.layer.settings_html_js) return

	// eslint-disable-next-line no-new-func
	return new Function(`${props.layer.settings_html_js}; onDataUpdate();`)
})

watch(() => props.parsedData, async () => {
	await nextTick()
	executeFunc.value?.()
})
</script>

<template>
	<component is="style">
		{{
			transform(`#layer${layer.id} {
					${layer.settings_html_css}
				}`,
			)
		}}
	</component>
	<div
		:id="`layer${layer.id}`"
		style="position: absolute; overflow: hidden; text-wrap: nowrap"
		:style="{
			transform: layer.transform_string,
		}"
		v-html="parsedData"
	/>
</template>
