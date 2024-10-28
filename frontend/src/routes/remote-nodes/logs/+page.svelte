<script>
	import { DataHandler } from '@vincjo/datatables/remote';
	import { format, formatDistance } from 'date-fns';
	import { loadData, loadNodeInfo } from './api-handler';
	import { onMount } from 'svelte';
	import { formatHostname, formatHashes, formatBytes } from '$lib/utils/strings';
	import {
		DtSrRowsPerPage,
		DtSrThSort,
		DtSrThFilter,
		DtSrRowCount,
		DtSrPagination,
		DtSrAutoRefresh
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
			<li class="crumb-separator" aria-hidden="true">/</li>
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
						<td>{formatHostname(nodeInfo?.hostname)}:{nodeInfo?.port}</td>
					</tr>
					<tr>
						<td class="font-bold">Public IP</td>
						<td>{nodeInfo?.ip_addresses.replace(/,/g, ', ')}</td>
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
						<DtSrAutoRefresh {handler} />
					</div>
					<div class="flex place-items-center">
						<button
							id="reloadDt"
							name="reloadDt"
							class="variant-filled-primary btn"
							on:click={() => handler.invalidate()}>Reload</button
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
