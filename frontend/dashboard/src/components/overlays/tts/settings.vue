<script setup lang="ts">
import { IconPlayerPlay } from '@tabler/icons-vue'
import {
	NAlert,
	NButton,
	NDivider,
	NForm,
	NFormItem,
	NGrid,
	NGridItem,
	NInput,
	NRow,
	NSelect,
	NSkeleton,
	NSlider,
	NSpace,
	NSwitch,
	NText,
	useMessage,
} from 'naive-ui'
import { computed, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import type { GetResponse as TTSSettings } from '@twir/api/messages/modules_tts/modules_tts'
import type { VoiceService } from '@twir/grpc/websockets/websockets'

import { useTtsOverlayManager } from '@/api/index.js'

const ttsManager = useTtsOverlayManager()
const ttsSettings = ttsManager.getSettings()
const ttsUpdater = ttsManager.updateSettings()
const ttsInfo = ttsManager.getInfo()
const ttsSay = ttsManager.useSay()

const countriesMapping: Record<string, string> = {
	ru: '🇷🇺 Russian',
	mk: '🇲🇰 Macedonian',
	uk: '🇺🇦 Ukrainian',
	ka: '🇬🇪 Georgian',
	ky: '🇰🇬 Kyrgyz',
	en: '🇺🇸 English',
	pt: '🇵🇹 Portuguese',
	eo: '🇺🇳 Esperanto',
	sq: '🇦🇱 Albanian',
	cs: '🇨🇿 Czech',
	pl: '🇵🇱 Polish',
	br: '🇧🇷 Brazilian',
}

interface Voice {
	label: string
	value: string
	key: string
}
type VoiceGroup = Omit<Voice, 'value' | 'gender'> & { children: Voice[], type: 'group' }
const voicesOptions = computed<VoiceGroup[]>(() => {
	if (!ttsInfo.data.value?.voicesInfo) return []

	const voices: Record<string, VoiceGroup> = {}

	for (const [voiceKey, voice] of Object.entries(ttsInfo.data.value.voicesInfo)) {
		let lang = voice.lang

		if (voice.lang === 'tt') {
			lang = 'ru'
		}

		if (!voices[lang]) {
			voices[lang] = {
				key: lang,
				label: `${countriesMapping[lang] ?? ''}`,
				type: 'group',
				children: [],
			}
		}

		voices[lang].children.push({
			key: lang,
			value: voiceKey,
			label: `${voice.name} (${voice.gender})`,
		})
	}

	return Object.entries(voices).map(([, group]) => group)
})

const formValue = ref<TTSSettings['data']>({
	enabled: false,
	voice: 'alan',
	disallowedVoices: [],
	pitch: 50,
	rate: 50,
	volume: 30,
	doNotReadTwitchEmotes: true,
	doNotReadEmoji: true,
	doNotReadLinks: true,
	allowUsersChooseVoiceInMainCommand: false,
	maxSymbols: 0,
	readChatMessages: false,
	readChatMessagesNicknames: false,
})

watch(ttsSettings.data, (v) => {
	if (!v?.data) return
	formValue.value = toRaw(v.data)
}, { immediate: true })

const message = useMessage()
const { t } = useI18n()

async function save() {
	await ttsUpdater.mutateAsync({ data: formValue.value })
	message.success(t('sharedTexts.saved'))
}

const previewText = ref('')

async function previewVoice() {
	if (!previewText.value || !formValue.value) return

	await ttsSay.mutateAsync({
		voice: formValue.value.voice,
		voiceService: 0 as VoiceService,
		text: previewText.value,
		volume: formValue.value.volume.toString(),
		pitch: formValue.value.pitch.toString(),
		rate: formValue.value.rate.toString(),
	})
}
</script>

<template>
	<NSpace vertical class="p-5">
		<NAlert type="info">
			{{ t('overlays.tts.eventsHint') }}
		</NAlert>

		<NSkeleton v-if="!formValue || ttsSettings.isLoading.value" :sharp="false" size="large" />

		<NForm v-else class="mt-4">
			<NGrid cols="1 s:1 m:2 l:2" responsive="screen" :x-gap="20" :y-gap="20">
				<NGridItem :span="1">
					<NSpace justify="space-between">
						<NText>{{ t('sharedTexts.enabled') }}</NText>
						<NSwitch v-model:value="formValue.enabled" />
					</NSpace>
				</NGridItem>

				<NGridItem :span="1">
					<NRow justify-content="space-between" align-items="flex-start" class="flex-nowrap">
						<NText>{{ t('overlays.tts.allowUsersChooseVoice') }}</NText>
						<NSwitch v-model:value="formValue.allowUsersChooseVoiceInMainCommand" />
					</NRow>
				</NGridItem>

				<NGridItem :span="1">
					<NSpace justify="space-between">
						<NText>{{ t('overlays.tts.doNotReadEmoji') }}</NText>
						<NSwitch v-model:value="formValue.doNotReadEmoji" />
					</NSpace>
				</NGridItem>

				<NGridItem :span="1">
					<NSpace justify="space-between">
						<NText>{{ t('overlays.tts.doNotReadTwitchEmotes') }}</NText>
						<NSwitch v-model:value="formValue.doNotReadTwitchEmotes" />
					</NSpace>
				</NGridItem>

				<NGridItem :span="1">
					<NSpace justify="space-between">
						<NText>{{ t('overlays.tts.doNotReadLinks') }}</NText>
						<NSwitch v-model:value="formValue.doNotReadLinks" />
					</NSpace>
				</NGridItem>

				<NGridItem>
					<NSpace justify="space-between">
						<NText>{{ t('overlays.tts.readChatMessages') }}</NText>
						<NSwitch v-model:value="formValue.readChatMessages" />
					</NSpace>
				</NGridItem>

				<NGridItem :span="1">
					<NSpace justify="space-between">
						<NText>{{ t('overlays.tts.readChatMessagesNicknames') }}</NText>
						<NSwitch v-model:value="formValue.readChatMessagesNicknames" />
					</NSpace>
				</NGridItem>
			</NGrid>

			<NDivider />

			<NFormItem :label="t('overlays.tts.voice')" show-require-mark>
				<NSelect
					v-model:value="formValue.voice"
					remote
					:loading="ttsInfo.isLoading.value"
					:options="voicesOptions"
				/>
			</NFormItem>

			<NFormItem :label="t('overlays.tts.disallowedVoices')">
				<NSelect
					v-model:value="formValue.disallowedVoices"
					remote
					clearable
					:loading="ttsInfo.isLoading.value"
					:options="voicesOptions"
					multiple
				/>
			</NFormItem>

			<NSpace class="w-full" vertical size="small">
				<NFormItem :label="t('overlays.tts.volume')" size="small">
					<NSlider v-model:value="formValue.volume" :step="1" />
				</NFormItem>
				<NFormItem :label="t('overlays.tts.pitch')" size="small">
					<NSlider v-model:value="formValue.pitch" :step="1" />
				</NFormItem>
				<NFormItem :label="t('overlays.tts.rate')" size="small">
					<NSlider v-model:value="formValue.rate" :step="1" />
				</NFormItem>
			</NSpace>

			<NDivider class="m-0 mb-2.5" />

			<NFormItem :label="`🎤 ${t('overlays.tts.previewText')}`">
				<div class="flex gap-1 w-full">
					<NInput
						v-model:value="previewText" :placeholder="t('overlays.tts.previewText')"
						class="w-1/2"
					/>
					<NButton text @click="previewVoice">
						<IconPlayerPlay />
					</NButton>
				</div>
			</NFormItem>
		</NForm>

		<NButton secondary type="success" block class="mt-2.5" @click="save">
			{{ t('sharedButtons.save') }}
		</NButton>
	</NSpace>
</template>
