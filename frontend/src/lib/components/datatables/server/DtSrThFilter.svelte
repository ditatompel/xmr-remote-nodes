<script lang="ts">
	import type { DataHandler, Row } from '@vincjo/datatables/remote';

	type T = $$Generic<Row>;

	export let handler: DataHandler<T>;
	export let filterBy: keyof T;

	/** @type {string} */
	export let placeholder: string = 'Filter';

	/** @type {number} */
	export let colspan: number = 1;

	let value: string = '';
	let timeout: any;

	const filter = () => {
		handler.filter(value, filterBy);
		clearTimeout(timeout);
		timeout = setTimeout(() => {
			handler.invalidate();
		}, 400);
	};
</script>

<th {colspan}>
	<input
		class="input variant-form-material h-8 w-full text-sm"
		type="text"
		{placeholder}
		bind:value
		on:input={filter}
	/>
</th>
