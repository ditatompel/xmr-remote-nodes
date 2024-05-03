<script lang="ts">
	import type { DataHandler, Row } from '@vincjo/datatables/remote';

	type T = $$Generic<Row>;

	export let handler: DataHandler<T>;
	const rowCount = handler.getRowCount();
</script>

{#if $rowCount === undefined}
	<div />
{:else}
	<div class={$$props.class ?? 'mr-6 leading-8 lg:leading-10'}>
		{#if $rowCount.total > 0}
			<b>{$rowCount.start}</b>
			- <b>{$rowCount.end}</b>
			/ <b>{$rowCount.total}</b>
		{:else}
			No entries found
		{/if}
	</div>
{/if}
