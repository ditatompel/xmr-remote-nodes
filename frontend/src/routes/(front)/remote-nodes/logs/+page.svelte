<script>
	import { DataHandler } from '@vincjo/datatables/remote';
	import { format, formatDistance } from 'date-fns';
	import { loadData, formatBytes } from './api-handler';
	import { onMount, onDestroy } from 'svelte';
	import {
		DtSrRowsPerPage,
		DtSrThSort,
		DtSrThFilter,
		DtSrRowCount,
		DtSrPagination
	} from '$lib/components/datatables/server';

	/**
	 * @param {number} n
	 * @param {number} p
	 */
	function maxPrecision(n, p) {
		return parseFloat(n.toFixed(p));
	}

	/**
	 * @param {number} h
	 */
	function formatHashes(h) {
		if (h < 1e-12) return '0 H';
		else if (h < 1e-9) return maxPrecision(h * 1e12, 0) + ' pH';
		else if (h < 1e-6) return maxPrecision(h * 1e9, 0) + ' nH';
		else if (h < 1e-3) return maxPrecision(h * 1e6, 0) + ' μH';
		else if (h < 1) return maxPrecision(h * 1e3, 0) + ' mH';
		else if (h < 1e3) return h + ' H';
		else if (h < 1e6) return maxPrecision(h * 1e-3, 2) + ' KH';
		else if (h < 1e9) return maxPrecision(h * 1e-6, 2) + ' MH';
		else return maxPrecision(h * 1e-9, 2) + ' GH';
	}

	/** @param {number | null } runtime */
	function parseRuntime(runtime) {
		return runtime === null ? '' : runtime.toLocaleString(undefined) + 's';
	}

	export let data;

	let pageId = '0';
	let filterProberId = 0;
	let filterStatus = -1;

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

	$: startInterval(); // Automatically start the interval on change

	onDestroy(() => {
		clearInterval(intervalId); // Clear the interval when the component is destroyed
	});
	onMount(() => {
		pageId = new URLSearchParams(window.location.search).get('node_id') || '0';
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
		<p class="mx-auto max-w-3xl">
			<strong>Monero remote node</strong> is a device on the internet running the Monero software with
			full copy of the Monero blockchain that doesn't run on the same local machine where the Monero
			wallet is located.
		</p>
	</div>
	<div class="mx-auto w-full max-w-3xl px-20">
		<hr class="!border-primary-400-500-token !border-t-4 !border-double" />
	</div>
</header>

<section id="introduction ">
	<div class="section-container text-center !max-w-4xl">
		<p>
			Remote node can be used by people who, for their own reasons (usually because of hardware
			requirements, disk space, or technical abilities), cannot/don't want to run their own node and
			prefer to relay on one publicly available on the Monero network.
		</p>
		<p>
			Using an open node will allow to make a transaction instantaneously, without the need to
			download the blockchain and sync to the Monero network first, but at the cost of the control
			over your privacy. the <strong>Monero community suggests to always run your own node</strong> to
			obtain the maximum possible privacy and to help decentralize the network.
		</p>
	</div>
</section>

<section id="monero-remote-node">
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
						<th><label for="prober_id">Prober</label></th>
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
						<th colspan="2">
							<select
								id="prober_id"
								name="prober_id"
								class="select variant-form-material"
								bind:value={filterProberId}
								on:change={() => {
									handler.filter(filterProberId, 'prober_id');
									handler.invalidate();
								}}
							>
								<option value={0}>Any</option>
							</select>
						</th>
						<th colspan="2">
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
							colspan={6}
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
