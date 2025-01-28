<script setup lang="ts">
import { NGrid, NGridItem } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import Card from '@/components/card/card.vue'
import { type CustomOverlayLayer, CustomOverlayLayerType } from '@/gql/graphql'

defineEmits<{
	select: [CustomOverlayLayer]
}>()

const { t } = useI18n()
</script>

<template>
	<NGrid responsive="screen" cols="1 s:2 m:3 l:4">
		<NGridItem :span="1">
			<Card
				class="cursor-pointer"
				title="HTML"
				@click="() => {
					$emit('select', {
						transformString: '',
						width: 200,
						height: 200,
						settings: {
							html_css: '.text { color: red }',
							html_html: `<span class='text'>$(stream.uptime)</span>`,
							html_pollSecondsInterval: 5,
							html_js: `
// will be triggered, when new overlay data comes from backend
function onDataUpdate() {
	console.log('updated')
}
							`,
						},
						type: CustomOverlayLayerType.Html,
						periodicallyRefetchData: true,
					})
				}"
			>
				<template #content>
					{{ t('overlaysRegistry.html.description') }}
				</template>
			</Card>
		</NGridItem>
	</NGrid>
</template>
