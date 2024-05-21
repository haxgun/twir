<script lang="ts" setup>
import { useWebSocket } from '@vueuse/core'
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import type { TTSMessage } from '@twir/grpc/websockets/websockets'

import { generateSocketUrlWithParams } from '@/helpers.js'

declare global {
	interface Window {
		webkitAudioContext: typeof AudioContext
	}
}

const isProcessing = ref(false)
const queueMessages = ref<TTSMessage[]>([])
const currentAudioBuffer = ref<AudioBufferSourceNode | null>(null)

const route = useRoute()

const apiKey = route.params.apiKey as string
const ttsUrl = generateSocketUrlWithParams('/overlays/tts', {
	apiKey,
})

const { data } = useWebSocket(ttsUrl, {
	autoReconnect: {
		delay: 500,
	},
})

watch(data, (message) => {
	const parsedData = JSON.parse(message)
	if (parsedData.eventName === 'say') {
		queueMessages.value.push(parsedData.data)
		processQueue()
	}

	if (parsedData.eventName === 'skip') {
		currentAudioBuffer.value?.stop()
	}

	if (parsedData.eventName === 'skipall') {
		currentAudioBuffer.value?.stop()
		queueMessages.value = []
	}
})

async function processQueue() {
	if (isProcessing.value) return

	const message = queueMessages.value.shift()
	if (!message) return

	isProcessing.value = true
	await sayMessage(message)
	isProcessing.value = false

	// Process the next item in the queue
	processQueue()
}

const ttsGenerateApiUrl = `${window.location.origin}/api/v1/public/overlays/tts/generate-file`

async function sayMessage(data: TTSMessage) {
	if (!apiKey || !data.text) return

	const audioContext = new (window.AudioContext || window.webkitAudioContext)()
	const gainNode = audioContext.createGain()

	const req = await fetch(ttsGenerateApiUrl, {
		method: 'POST',
		body: JSON.stringify({
			voice: data.voice,
			voiceService: data.voiceService,
			text: data.text,
			volume: Number(data.volume),
			pitch: Number(data.pitch),
			rate: Number(data.rate),
		}),
	})
	if (!req.ok) return

	const response = await req.arrayBuffer()

	const source = audioContext.createBufferSource()
	currentAudioBuffer.value = source

	source.buffer = await audioContext.decodeAudioData(response)

	gainNode.gain.value = Number.parseInt(data.volume) / 100
	source.connect(gainNode)
	gainNode.connect(audioContext.destination)

	return new Promise<void>((resolve) => {
		source.onended = () => {
			currentAudioBuffer.value = null
			resolve()
		}

		source.start(0)
	})
}
</script>
