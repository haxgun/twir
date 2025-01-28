<script lang="ts" setup>
import { IconSettings, IconTrash } from '@tabler/icons-vue'
import { NButton, NCard, useThemeVars } from 'naive-ui'
import { computed } from 'vue'

import type { CustomOverlayLayerType } from '@/gql/graphql'

defineProps<LayerProps>()
defineEmits<{
	focus: [index: number]
	remove: [index: number]
	openSettings: []
}>()
const theme = useThemeVars()
const activeLayourCardColor = computed(() => theme.value.infoColor)

export interface LayerProps {
	isFocused: boolean
	layerIndex: number
	type: CustomOverlayLayerType
}
</script>

<template>
	<NCard
		:title="type"
		class="cursor-pointer"
		:style="{
			border: isFocused ? `1px solid ${activeLayourCardColor}` : undefined,
		}"
		@click="$emit('focus', layerIndex)"
	>
		<slot />

		<div class="flex gap-3 w-full">
			<NButton
				secondary
				class="flex-1"
				@click="$emit('openSettings')"
			>
				<IconSettings />
			</NButton>
			<NButton
				secondary
				class="flex-1"
				type="error"
				@click="$emit('remove', layerIndex)"
			>
				<IconTrash />
			</NButton>
		</div>
	</NCard>
</template>
