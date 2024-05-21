import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { unref } from 'vue'

import type {
	GetInfoResponse,
	GetResponse,
	GetUsersSettingsResponse,
	PostRequest,
} from '@twir/api/messages/modules_tts/modules_tts'
import type { TTSMessage } from '@twir/grpc/websockets/websockets'
import type { Ref } from 'vue'

import { protectedApiClient } from '@/api/twirp.js'

declare global {
	interface Window {
		webkitAudioContext: typeof AudioContext
	}
}

export function useTtsOverlayManager() {
	const queryClient = useQueryClient()
	const queryKey = ['ttsSettings']
	const usersQueryKey = ['ttsUsersSettings']

	return {
		getSettings: () => useQuery({
			queryKey,
			queryFn: async (): Promise<GetResponse> => {
				const call = await protectedApiClient.modulesTTSGet({})
				return call.response
			},
		}),
		updateSettings: () => useMutation({
			mutationKey: ['ttsUpdate'],
			mutationFn: async (opts: PostRequest) => {
				await protectedApiClient.modulesTTSUpdate(opts)
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(queryKey)
			},
		}),
		getInfo: () => useQuery({
			queryKey: ['ttsInfo'],
			queryFn: async (): Promise<GetInfoResponse> => {
				const call = await protectedApiClient.modulesTTSGetInfo({})
				return call.response
			},
		}),
		getUsersSettings: () => useQuery({
			queryKey: usersQueryKey,
			queryFn: async (): Promise<GetUsersSettingsResponse> => {
				const call = await protectedApiClient.modulesTTSGetUsersSettings({})
				return call.response
			},
		}),
		deleteUsersSettings: () => useMutation({
			mutationKey: ['ttsUsersSettingsDelete'],
			mutationFn: async (ids: string[] | Ref<string[]>) => {
				const usersIds = unref(ids)

				await protectedApiClient.modulesTTSUsersDelete({ usersIds })
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(usersQueryKey)
			},
		}),
		useSay: () => useMutation({
			mutationKey: ['ttsSay'],
			mutationFn: async (opts: Omit<TTSMessage, 'channelId'>) => {
				const audioContext = new (window.AudioContext || window!.webkitAudioContext)()
				const gainNode = audioContext.createGain()

				const ttsGenerateApiUrl = `${window.location.origin}/api/v1/public/overlays/tts/generate-file`

				const req = await fetch(ttsGenerateApiUrl, {
					method: 'POST',
					body: JSON.stringify({
						voice: opts.voice,
						text: opts.text,
						volume: Number(opts.volume),
						pitch: Number(opts.pitch),
						rate: Number(opts.rate),
					}),
				})
				if (!req.ok) return

				const response = await req.arrayBuffer()

				const source = audioContext.createBufferSource()

				source.buffer = await audioContext.decodeAudioData(response)

				gainNode.gain.value = opts.volume / 100
				source.connect(gainNode)
				gainNode.connect(audioContext.destination)

				return new Promise((resolve) => {
					source.onended = () => {
						resolve(null)
					}

					source.start(0)
				})
			},
		}),
	}
}
