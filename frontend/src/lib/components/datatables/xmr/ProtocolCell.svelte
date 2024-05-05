<script>
	/** @type {string} */
	export let protocol;
	/** @type {boolean} */
	export let cors;
</script>

{#if protocol === 'http'}
	<span class="font-semibold uppercase text-sky-800 dark:text-sky-500">{protocol}</span>
{:else}
	<span class="font-semibold uppercase text-green-800 dark:text-green-500">{protocol}</span>
{/if}

{#if cors}
	<br />(CORS 💪)
{/if}
