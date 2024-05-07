<script>
	import { DataHandler } from '@vincjo/datatables/remote';
	import { format, formatDistance } from 'date-fns';
	import { loadData, loadNodeInfo } from './api-handler';
	import { onMount, onDestroy } from 'svelte';
	import { formatHashes, formatBytes } from '$lib/utils/strings';
	import {
		DtSrRowsPerPage,
		DtSrThSort,
		DtSrThFilter,
		DtSrRowCount,
		DtSrPagination
	} from '$lib/components/datatables/server';

	/** @param {number | null } runtime */
	function parseRuntime(runtime) {
		return runtime === null ? '' : runtime.toLocaleString(undefined) + 's';
	}

	export let data;

	let pageId = '0';
	let filterStatus = -1;

	/** @type {MoneroNode | null} */
	let nodeInfo;

	const handler = new DataHandler([], { rowsPerPage: 10, totalRows: 0 });
	let rows = handler.getRows();

	const reloadData = () => {
		handler.invalidate();
	};

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

	$: startInterval();

	onDestroy(() => {
		clearInterval(intervalId);
	});
	onMount(() => {
		pageId = new URLSearchParams(window.location.search).get('node_id') || '0';
		loadNodeInfo(pageId).then((data) => {
			nodeInfo = data;
		});
		handler.filter(pageId, 'node_id');
		handler.onChange((state) => loadData(state));
		handler.invalidate();
	});
</script>

<header id="hero" class="hero-gradient py-7">
	<div class="card text-token mx-auto flex w-fit justify-center p-4">
		<ol class="breadcrumb">
			<li class="crumb"><a class="link underline" href="/remote-nodes">Remote Nodes</a></li>
			<li class="crumb-separator" aria-hidden>/</li>
			<li>Logs</li>
		</ol>
	</div>
	<div class="section-container text-center">
		<h1 class="h1 pb-2 font-extrabold">{data.meta.title}</h1>
	</div>
	<div class="mx-auto w-full max-w-3xl px-20">
		<hr class="!border-primary-400-500-token !border-t-4 !border-double" />
	</div>
</header>

{#if nodeInfo === undefined}
	<div class="section-container mx-auto w-full max-w-3xl text-center">
		<p>Loading...</p>
	</div>
{:else if nodeInfo === null}
	<div class="section-container mx-auto w-full max-w-3xl text-center">
		<p>Node ID does not exist</p>
	</div>
{:else}
	<div class="section-container">
		<div class="table-container mx-auto w-full max-w-3xl">
			<table class="table">
				<tbody>
					<tr>
						<td class="font-bold">Hostname:Port</td>
						<td>{nodeInfo?.hostname}:{nodeInfo?.port}</td>
					</tr>
					<tr>
						<td class="font-bold">Public IP</td>
						<td>{nodeInfo?.ip}</td>
					</tr>
					<tr>
						<td class="font-bold">Net Type</td>
						<td>{nodeInfo?.nettype.toUpperCase()}</td>
					</tr></tbody
				>
			</table>
		</div>
	</div>

	<section id="node-logs">
		<div class="section-container">
			<div class="space-y-2 overflow-x-auto">
				<div class="flex justify-between">
					<DtSrRowsPerPage {handler} />
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

				<table class="table table-hover table-compact w-full table-auto">
					<thead>
						<tr>
							<th>#ID</th>
							<th>Prober ID</th>
							<th><label for="status">Status</label></th>
							<th>Height</th>
							<th>Adjusted Time</th>
							<th>DB Size</th>
							<th>Difficulty</th>
							<DtSrThSort {handler} orderBy="estimate_fee">Est. Fee</DtSrThSort>
							<DtSrThSort {handler} orderBy="date_checked">Date Checked</DtSrThSort>
							<DtSrThSort {handler} orderBy="fetch_runtime">Runtime</DtSrThSort>
						</tr>
						<tr>
							<th colspan="3">
								<select
									id="status"
									name="status"
									class="select variant-form-material"
									bind:value={filterStatus}
									on:change={() => {
										handler.filter(filterStatus, 'status');
										handler.invalidate();
									}}
								>
									<option value={-1}>Any</option>
									<option value="1">Online</option>
									<option value="0">Offline</option>
								</select>
							</th>
							<DtSrThFilter
								{handler}
								filterBy="failed_reason"
								placeholder="Filter reason"
								colspan={7}
							/>
						</tr>
					</thead>
					<tbody>
						{#each $rows as row (row.id)}
							<tr>
								<td>{row.id}</td>
								<td>{row.prober_id}</td>
								<td>{row.status === 1 ? 'OK' : 'ERR'}</td>
								{#if row.status !== 1}
									<td colspan="5">{row.failed_reason ?? ''}</td>
								{:else}
									<td class="text-right">{row.height.toLocaleString(undefined)}</td>
									<td>{format(row.adjusted_time * 1000, 'yyyy-MM-dd HH:mm')}</td>
									<td class="text-right">{formatBytes(row.database_size, 2)}</td>
									<td class="text-right">{formatHashes(row.difficulty)}</td>
									<td class="text-right">{row.estimate_fee.toLocaleString(undefined)}</td>
								{/if}
								<td>
									{format(row.date_checked * 1000, 'PP HH:mm')}<br />
									{formatDistance(row.date_checked * 1000, new Date(), { addSuffix: true })}
								</td>
								<td class="text-right">{parseRuntime(row.fetch_runtime)}</td>
							</tr>
						{/each}
					</tbody>
				</table>

				<div class="flex justify-between mb-2">
					<DtSrRowCount {handler} />
					<DtSrPagination {handler} />
				</div>
			</div>
		</div>
	</section>
{/if}

<style lang="postcss">
	.section-container {
		@apply mx-auto w-full max-w-7xl p-4;
	}
	/* Hero Gradient */
	/* prettier-ignore */
	.hero-gradient {
  background-image:
    radial-gradient(at 0% 0%, rgba(242, 104, 34, .4) 0px, transparent 50%),
    radial-gradient(at 98% 1%, rgba(var(--color-warning-900) / 0.33) 0px, transparent 50%);
  }
	/*
  td:nth-child(1) {
  @apply max-w-20;
  }
  */
</style>
