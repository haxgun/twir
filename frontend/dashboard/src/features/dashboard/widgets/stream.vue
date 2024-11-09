<script setup lang="ts">
import { FilePenLine } from 'lucide-vue-next'

import BaseWidget from '../ui/widget.vue'

import { useProfile } from '@/api'

const { data: profile } = useProfile()
</script>

<template>
	<BaseWidget>
		<template #title>
			Stream
		</template>

		<template #content>
			<div v-if="profile" class="h-full">
				<div class="h-16 bg-[#252525] flex items-center px-2">
					<div class="flex w-3/5 justify-between h-full items-center">
						<div class="flex flex-col">
							<span class="font-medium">Coding and chill</span>
							<span class="text-xs font-light">qqweqwe</span>
						</div>
						<button class="mt-2 mr-2 bg-[#343434] rounded p-1 place-self-start">
							<FilePenLine class="size-4" />
						</button>
					</div>

					<div class="header-stats-item">
						<span class="text-xs">Uptime</span>
						<span class="font-medium text-md">00:01:52</span>
					</div>

					<div class="header-stats-item">
						<span class="text-xs">Viewers</span>
						<span class="font-medium text-md">32</span>
					</div>
				</div>

				<img :src="profile.selectedDashboardTwitchUser.offlineImageUrl || profile.selectedDashboardTwitchUser.profileImageUrl" class="w-full h-2/4 aspect-[16/9]" />

				<div class="h-1/4 bg-[#252525] px-2">
					<div class="grid grid-cols-3 gap-2 px-2 items-center">
						<div class="stats-item">
							<FilePenLine class="size-6" />
							<div class="flex flex-col gap-0.5">
								<span class="text-zinc-500 text-sm">Title</span>
								<span class="font-medium text-white">Content</span>
							</div>
						</div>
						<div class="stats-item">
							<FilePenLine class="size-6" />
							<div class="flex flex-col gap-0.5">
								<span class="text-zinc-500 text-sm">Title</span>
								<span class="font-medium text-white">Content</span>
							</div>
						</div>
						<div class="stats-item">
							<FilePenLine class="size-6" />
							<div class="flex flex-col gap-0.5">
								<span class="text-zinc-500 text-sm">Title</span>
								<span class="font-medium text-white">Content</span>
							</div>
						</div>
						<div class="stats-item">
							<FilePenLine class="size-6" />
							<div class="flex flex-col gap-0.5">
								<span class="text-zinc-500 text-sm">Title</span>
								<span class="font-medium text-white">Content</span>
							</div>
						</div>
						<div class="stats-item">
							<FilePenLine class="size-6" />
							<div class="flex flex-col gap-0.5">
								<span class="text-zinc-500 text-sm">Title</span>
								<span class="font-medium text-white">Content</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</template>
	</BaseWidget>
</template>

<style scoped>
.header-stats-item {
	@apply flex flex-col w-1/5 h-full items-center place-content-center border-l border-l-[#343434]
}

.stats-item {
	@apply w-2/6 flex gap-4 items-center
}
</style>
