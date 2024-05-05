<script lang="ts">
	export let uptime: number;
</script>

{#if uptime > 98}
	<span class="text-green-800 dark:text-green-500">{uptime}%</span>
{:else if uptime < 98 && uptime > 80}
	<span class="text-sky-800 dark:text-sky-500">{uptime}%</span>
{:else if uptime < 80 && uptime > 75}
	<span class="text-orange-800 dark:text-orange-300">{uptime}%</span>
{:else}
	<span class="text-rose-800 dark:text-rose-400">{uptime}%</span>
{/if}
