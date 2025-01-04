<script setup lang="ts">
import { LoginWidget, type LoginWidgetUser } from 'vue-tg'

import { useTelegramIntegrationCallback } from '@/api/integrations/telegram'

const mutation = useTelegramIntegrationCallback()
async function onAuth(user: LoginWidgetUser) {
	await mutation.executeMutation({ input: user })
}
</script>

<template>
	<LoginWidget
		bot-username="satonttestsbot"
		user-photo
		@auth="onAuth"
	/>
</template>
