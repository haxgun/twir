import { useMutation } from '@urql/vue'

import { graphql } from '@/gql'

export function useTelegramIntegrationCallback() {
	return useMutation(graphql(`
		mutation telegramIntegrationCallback($input: Map!) {
			telegramCallback(input: { values: $input })
		}
	`))
}
