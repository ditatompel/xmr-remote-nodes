<script lang="ts">
	export let is_available: boolean;
	export let statuses: number[];
</script>

{#if is_available}
	<span class="font-semibold text-green-800 dark:text-green-500">Online</span>
{:else}
	<span class="text-rose-800 dark:text-rose-400">Offline</span>
{/if}
<br />
{#each statuses as status}
	{#if status === 1}
		<span class="text-success-700 dark:text-success-400 mr-1">•</span>
	{:else if status === 0}
		<span class="text-error-700 dark:text-error-400 mr-1">•</span>
	{:else}
		<span class="text-surface-400 dark:text-surface-600 mr-1">•</span>
	{/if}
{/each}
