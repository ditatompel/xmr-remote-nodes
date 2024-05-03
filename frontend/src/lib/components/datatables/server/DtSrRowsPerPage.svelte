<script lang="ts">
	import type { DataHandler, Row } from '@vincjo/datatables/remote';

	type T = $$Generic<Row>;

	export let handler: DataHandler<T>;
	export let options = [5, 10, 20, 50, 100];
	export let labelId = 'rowsPerPage';

	const rowsPerPage = handler.getRowsPerPage();

	const setRowsPerPage = () => {
		handler.setPage(1);
		handler.invalidate();
	};
</script>

<div class="flex place-items-center">
	<label for={labelId}>Show</label>
	<select
		class="select ml-2"
		id={labelId}
		name="rowsPerPage"
		bind:value={$rowsPerPage}
		on:change={setRowsPerPage}
	>
		{#each options as option}
			<option value={option}>
				{option}
			</option>
		{/each}
	</select>
</div>
