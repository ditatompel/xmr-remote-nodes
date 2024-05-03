<script lang="ts">
	import type { DataHandler, Row } from '@vincjo/datatables/remote';
	type T = $$Generic<Row>;

	export let handler: DataHandler<T>;
	export let orderBy: keyof T;

	const update = () => {
		handler.sort(orderBy);
		handler.invalidate();
	};
</script>

<th on:click={update} class="cursor-pointer select-none p-2 px-5">
	<div class="flex h-full items-center justify-start gap-x-2">
		<slot /> ↕️
	</div>
</th>
