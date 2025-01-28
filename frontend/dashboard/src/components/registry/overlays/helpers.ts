import type { CustomOverlayLayerSettingsHtml, CustomOverlayLayerSettingsImage } from '@/gql/graphql.js'

export function isSettingsHtml(settings: CustomOverlayLayerSettingsHtml | CustomOverlayLayerSettingsImage): settings is CustomOverlayLayerSettingsHtml {
	return (settings as CustomOverlayLayerSettingsHtml).html !== undefined
}

export function isSettingsImage(settings: CustomOverlayLayerSettingsHtml | CustomOverlayLayerSettingsImage): settings is CustomOverlayLayerSettingsImage {
	return 'image_url' in settings
}
