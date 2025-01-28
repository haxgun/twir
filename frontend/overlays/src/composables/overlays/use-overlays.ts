import { createGlobalState, useWebSocket } from '@vueuse/core'
import { ref, watch } from 'vue'

import { generateSocketUrlWithParams } from '@/helpers.js'

export interface Layer {
	id: string
	type: 'HTML' | 'IMAGE'
	overlay_id: string
	width: number
	height: number
	createdAt: string
	updatedAt: string
	periodically_refetch_data: boolean
	transform_string: string
	settings_html_html?: string
	settings_html_css?: string
	settings_html_js?: string
	settings_html_data_poll_seconds_interval?: any
	settings_image_url?: any
}

export const useOverlays = createGlobalState(() => {
	const overlayUrl = ref('')
	const overlayId = ref('')
	const layers = ref<Array<Layer>>([])

	const { data, status, send, open } = useWebSocket(
		overlayUrl,
		{
			immediate: false,
			autoReconnect: {
				delay: 500,
			},
			onConnected() {
				send(
					JSON.stringify({
						eventName: 'getLayers',
						data: {
							overlayId: overlayId.value,
						},
					}),
				)
			},
		},
	)

	const parsedLayersData = ref<Record<string, string>>({})

	watch(data, (d) => {
		const parsedData = JSON.parse(d)

		if (parsedData.eventName === 'layers') {
			layers.value = parsedData.layers as Array<Layer>
		}

		if (parsedData.eventName === 'parsedLayerVariables') {
			parsedLayersData.value[parsedData.layerId] = parsedData.data
		}

		if (parsedData.eventName === 'refreshOverlays') {
			window.location.reload()
		}
	})

	function requestLayerData(layerId: string): void {
		send(
			JSON.stringify({
				eventName: 'parseLayerVariables',
				data: {
					layerId,
				},
			}),
		)
	}

	function connectToOverlays(apiKey: string, _overlayId: string): void {
		if (status.value === 'OPEN') return

		const url = generateSocketUrlWithParams('/overlays/registry/overlays', {
			apiKey,
		})
		console.log(url)

		overlayUrl.value = url
		overlayId.value = _overlayId

		open()
	}

	return {
		layers,
		parsedLayersData,
		requestLayerData,
		connectToOverlays,
	}
})
