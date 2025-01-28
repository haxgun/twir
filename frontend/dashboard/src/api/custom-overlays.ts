import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

export const useCustomOverlaysApi = createGlobalState(() => {
	const cacheKey = 'customOverlays'

	const useData = () => useQuery({
		query: graphql(`
			query CustomOverlays {
				customOverlays {
					id
					name
					createdAt
					updatedAt
					width
					height
					layers {
						id
						type
						overlayId
						width
						height
						createdAt
						updatedAt
						periodicallyRefetchData
						transformString
						settings {
							... on CustomOverlayLayerSettingsHTML {
								html
								css
								js
								pollSecondsInterval
							}
							... on CustomOverlayLayerSettingsImage {
								url
							}
						}
					}
				}
			}
		`),
		context: {
			additionalTypenames: [cacheKey],
		},
		variables: {},
	})

	const useDelete = () => useMutation(
		graphql(`
			mutation DeleteCustomOverlay($id: UUID!) {
				customOverlaysDelete(id: $id)
			}
		`),
		[cacheKey],
	)

	const useCreate = () => useMutation(
		graphql(`
			mutation CreateCustomOverlay($createInput: CustomOverlayCreateInput!) {
				customOverlaysCreate(input: $createInput) {
					id
				}
			}
		`),
		[cacheKey],
	)

	const useUpdate = () => useMutation(
		graphql(`
			mutation UpdateCustomOverlay($updateInput: CustomOverlayUpdateInput!) {
				customOverlaysUpdate(input: $updateInput) {
					id
				}
			}
		`),
		[cacheKey],
	)

	const useParseVariablesInHtml = () => useMutation(graphql(`
		mutation ParseVariablesInHtml($text: String!) {
			customOverlaysParseHtml(input: { text: $text})
		}
	`))

	return {
		useData,
		useDelete,
		useCreate,
		useUpdate,
		useParseVariablesInHtml,
	}
})
