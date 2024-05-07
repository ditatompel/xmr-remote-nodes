<script lang="ts">
	import { onDestroy } from 'svelte';
	import type { DataHandler, Row } from '@vincjo/datatables/remote';

	type T = $$Generic<Row>;
	export let handler: DataHandler<T>;

	let intervalId: number | undefined;
	let intervalValue = 0;

	const intervalOptions = [
		{ value: 0, label: 'No' },
		{ value: 5, label: '5s' },
		{ value: 10, label: '10s' },
		{ value: 30, label: '30s' },
		{ value: 60, label: '1m' }
	];

	const startInterval = () => {
		const seconds = intervalValue;
		if (isNaN(seconds) || seconds < 0) {
			return;
		}

		if (!intervalOptions.some((option) => option.value === seconds)) {
			return;
		}

		if (intervalId) {
			clearInterval(intervalId);
		}

		if (seconds > 0) {
			handler.invalidate();
			intervalId = setInterval(() => {
				handler.invalidate();
			}, seconds * 1000);
		}
	};

	$: startInterval();

	onDestroy(() => {
		clearInterval(intervalId);
	});
</script>

<label for="autoRefreshInterval">Auto Refresh:</label>
<select
	class="select ml-2"
	id="autoRefreshInterval"
	bind:value={intervalValue}
	on:change={startInterval}
>
	{#each intervalOptions as { value, label }}
		<option {value}>{label}</option>
	{/each}
</select>
