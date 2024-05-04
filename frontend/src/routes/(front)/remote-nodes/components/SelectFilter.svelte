<script>
	import { getDistinct } from '$lib/utils/arrays';

	/** @type {string} */
	export let filterName;
	export let filterValue; //: Writable<string>;
	export let preFilteredValues; //: Readable<unknown[]>;
	$: uniqueValues = getDistinct($preFilteredValues);
</script>

<div class="pt-2">
	<select name={filterName} class="select" bind:value={$filterValue} on:click|stopPropagation>
		<option value={undefined}>All</option>
		{#each uniqueValues as value}
			{#if value === ''}
				<option {value}>UNKNOWN</option>
			{:else}
				<option {value}>{value.toUpperCase()}</option>
			{/if}
		{/each}
	</select>
</div>
