<script setup lang="ts">
import { IconCopy , IconDeviceFloppy } from '@tabler/icons-vue'
import { NButton, NDivider, NFormItem, NInput, NInputNumber, NModal, useMessage } from 'naive-ui'
import { computed, onMounted, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import Moveable from 'vue3-moveable'

import HtmlLayer from './layers/html.vue'
import HtmlLayerForm from './layers/htmlForm.vue'

import type { CustomOverlayCreateInput, CustomOverlayUpdateInput } from '@/gql/graphql'
import type { OnDrag, OnResize, OnWarp } from 'vue3-moveable'

import { useCustomOverlaysApi } from '@/api/custom-overlays'
import {
	useProfile,
} from '@/api/index.js'
import NewSelector from '@/components/registry/overlays/newSelector.vue'
import { CustomOverlayLayerType } from '@/gql/graphql.js'
import { copyToClipBoard } from '@/helpers'

const { t } = useI18n()

const router = useRouter()

const overlaysManager = useCustomOverlaysApi()
const { data: overlays } = overlaysManager.useData()
const creator = overlaysManager.useCreate()
const updater = overlaysManager.useUpdate()

type OverlayForm = CustomOverlayCreateInput | CustomOverlayUpdateInput

const formValue = ref<OverlayForm>({
	id: '',
	name: '',
	layers: [],
	width: 1920,
	height: 1080,
})

function setCurrentOverlay() {
	const id = router.currentRoute.value.params.id
	if (typeof id !== 'string' || id === 'new') {
		return
	}

	const overlay = overlays.value?.customOverlays.find((overlay) => overlay.id === id)
	if (!overlay) {
		return
	}

	const raw = toRaw(overlay)

	const input: CustomOverlayUpdateInput = {
		id: raw.id,
		height: raw.height,
		width: raw.width,
		layers: raw.layers.map((layer) => {
			return {
				type: layer.type,
				periodicallyRefetchData: layer.periodicallyRefetchData,
				width: layer.width,
				height: layer.height,
				transformString: layer.transformString,
				settings: 'html' in layer.settings
					? {
						html_html: layer.settings.html,
						html_css: layer.settings.css,
						html_js: layer.settings.js,
						html_pollSecondsInterval: layer.settings.pollSecondsInterval,
					}
					: 'url' in layer.settings
						? {
							image_url: layer.settings.url,
						}
						: {},
			}
		}),
		name: raw.name,
	}

	formValue.value = input
}

watch(overlays, (v) => {
	if (!v) return

	setCurrentOverlay()
})

onMounted(() => setCurrentOverlay())

const messages = useMessage()

async function save() {
	const data = toRaw(formValue.value)

	if (!data.name || data.name.length > 30) {
		messages.error(t('overlaysRegistry.validations.name'))
		return
	}

	if (!data.layers.length || data.layers.length > 15) {
		messages.error(t('overlaysRegistry.validations.layers'))
		return
	}

	if ('id' in data && data.id) {
		await updater.executeMutation({
			updateInput: data,
		})
	} else {
		const { data: newOverlay } = await creator.executeMutation({
			createInput: {
				...data,
				// eslint-disable-next-line ts/ban-ts-comment
				// @ts-expect-error
				id: undefined,
			},
		})

		if (!newOverlay) {
			return
		}

		router.push(`/dashboard/registry/overlays/${newOverlay.customOverlaysCreate.id}`)
	}

	messages.success(t('sharedTexts.saved'))
}

const currentlyFocused = ref(0)
function focus(index: number) {
	currentlyFocused.value = index
}

interface EventWithLayerIndex {
	index: number
}

function onDrag({ target, transform, index }: OnDrag & EventWithLayerIndex) {
	focus(index)
	target.style.transform = transform

	formValue.value.layers[index].transformString = transform
}

function onWarp({ target, transform, index }: OnWarp & EventWithLayerIndex) {
	focus(index)
	target.style.transform = transform

	formValue.value.layers[index].transformString = transform
}

function onResize({ target, width, height, transform, index }: OnResize & EventWithLayerIndex) {
	focus(index)

	target.style.width = `${width}px`
	target.style.height = `${height}px`
	target.style.transform = transform

	formValue.value.layers[index].height = height
	formValue.value.layers[index].width = width
}

function removeLayer(index: number) {
	formValue.value.layers = formValue.value.layers.filter((_, i) => i !== index)
	focus(-1)
}

const isOverlayNewModalOpened = ref(false)

const userProfile = useProfile()
async function copyUrl() {
	if (!('id' in formValue.value)) return

	await copyToClipBoard(`${window.location.origin}/overlays/${userProfile.data.value?.apiKey}/registry/overlays/${formValue.value.id}`)
}

const innerWidth = computed(() => window.innerWidth)

const mode = ref<'resize' | 'warp'>('resize')
</script>

<template>
	<div class="flex w-full relative">
		<div class="w-[85%] absolute left-[275px]">
			<div class="flex gap-2">
				<NButton :type="mode === 'resize' ? 'success' : 'info'" secondary @click="mode = 'resize'">
					Resize mode
				</NButton>
				<NButton :type="mode === 'warp' ? 'success' : 'info'" secondary @click="mode = 'warp'">
					Warp mode
				</NButton>
				<NButton type="warning" secondary @click="resetTransform">
					Reset transform
				</NButton>
			</div>
			<div
				class="grid-container"
				:style="{
					width: `${formValue.width}px`,
					height: `${formValue.height}px`,
					transform: `scale(${(innerWidth / formValue.width) * 0.7})`,
				}"
			>
				<div v-for="(layer, index) of formValue.layers" :key="index">
					<HtmlLayer
						v-if="layer.type === CustomOverlayLayerType.Html"
						:transformString="layer.transformString"
						:width="layer.width"
						:height="layer.height"
						:index="index"
						:text="layer.settings?.html_html ?? ''"
						:css="layer.settings?.html_css ?? ''"
						:js="layer.settings?.html_js ?? ''"
						:periodicallyRefetchData="layer.periodicallyRefetchData"
					/>

					<Moveable
						className="moveable"
						:target="`#layer-${index}`"
						:draggable="true"
						:resizable="mode === 'resize'"
						:rotatable="false"
						:snappable="true"
						:warpable="mode === 'warp'"
						:bounds="{ left: 0, top: 0, right: 0, bottom: 0, position: 'css' }"
						:origin="false"
						:renderDirections="currentlyFocused === index ? ['nw', 'n', 'ne', 'w', 'e', 'sw', 's', 'se'] : []"
						@drag="(opts) => onDrag({ ...opts, index })"
						@resize="(opts) => onResize({ ...opts, index })"
						@click="focus(index)"
						@warp="(opts) => onWarp({ ...opts, index })"
					>
					</Moveable>
				</div>
			</div>
		</div>
		<div class="flex flex-col gap-1 relative">
			<NButton
				:disabled="!formValue.name || !formValue.layers.length" block secondary
				type="success" @click="save"
			>
				<IconDeviceFloppy />
				{{ t('sharedButtons.save') }}
			</NButton>
			<NButton
				block secondary type="info" :disabled="!('id' in formValue)"
				@click="copyUrl"
			>
				<IconCopy />
				{{ t('overlays.copyOverlayLink') }}
			</NButton>

			<NFormItem :label="t('overlaysRegistry.name')">
				<NInput
					v-model:value="formValue.name" :placeholder="t('overlaysRegistry.name')"
					:maxlength="30"
				/>
			</NFormItem>

			<NFormItem :label="t('overlaysRegistry.customWidth')">
				<NInputNumber
					v-model:value="formValue.width" :min="50"
					:placeholder="t('overlaysRegistry.customWidth')"
				/>
			</NFormItem>

			<NFormItem :label="t('overlaysRegistry.customHeight')">
				<NInputNumber
					v-model:value="formValue.height" :min="50"
					:placeholder="t('overlaysRegistry.customHeight')"
				/>
			</NFormItem>

			<NDivider />

			<NButton
				secondary
				type="success"
				@click="isOverlayNewModalOpened = true"
			>
				{{ t('overlaysRegistry.createNewLayer') }}
			</NButton>

			<div class="flex flex-col gap-3 w-full">
				<template v-for="(layer, index) of formValue.layers">
					<HtmlLayerForm
						v-if="layer.type === CustomOverlayLayerType.Html"
						:key="index"
						v-model:html="formValue.layers[index].settings!.html_html"
						v-model:css="formValue.layers[index].settings!.html_css"
						v-model:js="formValue.layers[index].settings!.html_js"

						v-model:periodicallyRefetchData="formValue.layers[index].periodicallyRefetchData"
						:isFocused="currentlyFocused === index"
						:layerIndex="index"
						:type="layer.type"
						@remove="removeLayer"
						@focus="focus"
					/>
				</template>
			</div>
		</div>
	</div>

	<NModal
		v-model:show="isOverlayNewModalOpened" class="w-[50vw]" preset="card"
		:title="t('sharedButtons.create')"
	>
		<NewSelector
			@select="v => {
				formValue.layers.push(v)
				isOverlayNewModalOpened = false
			}"
		/>
	</NModal>
</template>

<style scoped>
.grid-container {
	background-color: rgb(18, 18, 18);
	transform-origin: 0px 0px;

	background-image: linear-gradient(45deg, rgb(34, 34, 34) 25%, transparent 25%), linear-gradient(135deg, rgb(34, 34, 34) 25%, transparent 25%), linear-gradient(45deg, transparent 75%, rgb(34, 34, 34) 75%), linear-gradient(135deg, transparent 75%, rgb(34, 34, 34) 75%);
	background-size: 20px 20px;
	background-position: 0px 0px, 10px 0px, 10px -10px, 0px 10px;
}
</style>
