<script lang="ts">
	export let estimate_fee: number;
	export let majority_fee: number;
</script>

{#if estimate_fee !== majority_fee}
	<span class="text-orange-800 dark:text-orange-300">{estimate_fee}<br />(CAUTION!)</span>
{:else}
	{estimate_fee}
{/if}
