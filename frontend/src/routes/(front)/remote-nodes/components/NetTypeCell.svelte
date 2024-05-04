<script>
	/** @type {string} */
	export let nettype;
	/** @type {number} */
	export let height;
</script>

{#if nettype === 'stagenet'}
	<span class="font-semibold uppercase text-sky-800 dark:text-sky-500">{nettype}</span>
{:else if nettype === 'testnet'}
	<span class="font-semibold uppercase text-rose-800 dark:text-rose-400">{nettype}</span>
{:else}
	<span class="font-semibold uppercase text-green-800 dark:text-green-500">{nettype}</span>
{/if}
<br />{height}
