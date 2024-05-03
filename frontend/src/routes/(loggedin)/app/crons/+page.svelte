<script>
	import { DataHandler } from '@vincjo/datatables/remote';
	import { format, formatDistance } from 'date-fns';
	import { loadData } from './api-handler';
	import { onMount, onDestroy } from 'svelte';
	import { DtSrThSort, DtSrThFilter, DtSrRowCount } from '$lib/components/datatables/server';

	const handler = new DataHandler([], { rowsPerPage: 1000, totalRows: 0 });
	let rows = handler.getRows();

	const reloadData = () => {
		handler.invalidate();
	};

	/** @type {string | number} */
	let filterState = -1;
	/** @type {string | number} */
	let filterEnabled = -1;

	/** @type {number | undefined} */
	let intervalId;
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
			reloadData();
			intervalId = setInterval(() => {
				reloadData();
			}, seconds * 1000);
		}
	};

	$: startInterval(); // Automatically start the interval on change

	onDestroy(() => {
		clearInterval(intervalId); // Clear the interval when the component is destroyed
	});
	onMount(() => {
		handler.onChange((state) => loadData(state));
		handler.invalidate();
	});
</script>

<div class="mb-4">
	<h1 class="h2 font-extrabold dark:text-white">Crons</h1>
</div>

<div class="dashboard-card">
	<div class="flex justify-between">
		<div class="invisible flex place-items-center md:visible">
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
		</div>
		<div class="flex place-items-center">
			<button
				id="reloadDt"
				name="reloadDt"
				class="variant-filled-primary btn"
				on:click={reloadData}>Reload</button
			>
		</div>
	</div>

	<div class="my-2 overflow-x-auto">
		<table class="table table-hover table-compact w-full table-auto">
			<thead>
				<tr>
					<DtSrThSort {handler} orderBy="id">ID</DtSrThSort>
					<th>Title</th>
					<th>Slug</th>
					<th>Description</th>
					<DtSrThSort {handler} orderBy="run_every">Run Every</DtSrThSort>
					<DtSrThSort {handler} orderBy="last_run">Last Run</DtSrThSort>
					<DtSrThSort {handler} orderBy="next_run">Next Run</DtSrThSort>
					<DtSrThSort {handler} orderBy="run_time">Run Time</DtSrThSort>
					<th>State</th>
					<th>Enabled</th>
				</tr>
				<tr>
					<DtSrThFilter {handler} filterBy="title" placeholder="Title" colspan={3} />
					<DtSrThFilter {handler} filterBy="description" placeholder="Description" colspan={5} />
					<th>
						<select
							id="fState"
							name="fState"
							class="select variant-form-material"
							bind:value={filterState}
							on:change={() => {
								handler.filter(filterState, 'cron_state');
								reloadData();
							}}
						>
							<option value={-1}>Any</option>
							<option value={1}>Running</option>
							<option value={0}>Idle</option>
						</select>
					</th>
					<th>
						<select
							id="fEnabled"
							name="fEnabled"
							class="select variant-form-material"
							bind:value={filterEnabled}
							on:change={() => {
								handler.filter(filterEnabled, 'is_enabled');
								reloadData();
							}}
						>
							<option value={-1}>Any</option>
							<option value={1}>Yes</option>
							<option value={0}>No</option>
						</select>
					</th>
				</tr>
			</thead>
			<tbody>
				{#each $rows as row (row.id)}
					<tr>
						<td>{row.id}</td>
						<td>{row.title}</td>
						<td>{row.slug}</td>
						<td>{row.description}</td>
						<td>{row.run_every}s</td>
						<td>
							{format(row.last_run * 1000, 'PP HH:mm')}<br />
							{formatDistance(row.last_run * 1000, new Date(), { addSuffix: true })}
						</td>
						<td>
							{format(row.next_run * 1000, 'PP HH:mm')}<br />
							{formatDistance(row.next_run * 1000, new Date(), { addSuffix: true })}
						</td>
						<td>{row.run_time}</td>
						<td>{row.cron_state ? 'RUNNING' : 'IDLE'}</td>
						<td>{row.is_enabled ? 'ENABLED' : 'DISABLED'}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<div class="flex justify-between mb-2">
		<DtSrRowCount {handler} />
	</div>
</div>
