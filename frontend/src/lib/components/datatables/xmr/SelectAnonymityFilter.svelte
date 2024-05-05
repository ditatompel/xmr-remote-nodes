<script lang="ts">
	import { getDistinct } from '$lib/utils/arrays';

	export let filterValue; //: Writable<string>;
	export let preFilteredValues; //: Readable<unknown[]>;
	$: uniqueValues = getDistinct($preFilteredValues);
</script>

<div class="pt-2">
	<select name="filterAnonymity" class="select" bind:value={$filterValue} on:click|stopPropagation>
		<option value={undefined}>All</option>
		{#each uniqueValues as value}
			{#if value === true}
				<option {value}>TOR</option>
			{:else}
				<option {value}>CLEARNET</option>
			{/if}
		{/each}
	</select>
</div>
