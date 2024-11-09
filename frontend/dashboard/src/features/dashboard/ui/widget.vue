<script setup lang="ts">
import { GripVertical, XIcon } from 'lucide-vue-next'
</script>

<template>
	<div class="flex flex-col w-full h-full">
		<div class="flex items-center bg-[#343434] px-2 h-8 rounded-t-md justify-between">
			<div class="flex items-center">
				<GripVertical class="widgets-draggable-handle size-4 cursor-move" />
				<span class="text-xs uppercase font-light ml-1">
					<slot name="title" />
				</span>
			</div>

			<div class="p-1">
				<XIcon class="size-4" />
			</div>
		</div>

		<div class="h-full bg-accent rounded-b-md">
			<slot name="content" />
		</div>
	</div>
</template>
