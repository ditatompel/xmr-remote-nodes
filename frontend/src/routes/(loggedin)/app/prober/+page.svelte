<script>
	import { DataHandler } from '@vincjo/datatables/remote';
	import { format, formatDistance } from 'date-fns';
	import { loadData, createProber, editProber, deleteProber } from './api-handler';
	import { onMount } from 'svelte';
	import { getModalStore, getToastStore } from '@skeletonlabs/skeleton';
	import {
		DtSrRowsPerPage,
		DtSrThSort,
		DtSrThFilter,
		DtSrRowCount,
		DtSrPagination,
		DtSrAutoRefresh
	} from '$lib/components/datatables/server';
	const modalStore = getModalStore();
	const toastStore = getToastStore();

	function showAddModal() {
		/** @type {import('@skeletonlabs/skeleton').ModalSettings} */
		const modal = {
			type: 'prompt',
			// Data
			title: 'Enter Name',
			body: 'Enter a name for the prober',
			valueAttr: { type: 'text', minlength: 3, maxlength: 50, required: true },
			response: (r) => {
				if (!r) return;
				createProber(r)
					.then((res) => {
						if (res.status !== 'ok') {
							toastStore.trigger({ message: 'Failed to create prober' });
						} else {
							toastStore.trigger({
								message: 'Prober created',
								background: 'variant-filled-success'
							});
							handler.invalidate();
						}
					})
					.catch(() => {
						toastStore.trigger({ message: 'Failed to create prober' });
					});
			}
		};
		modalStore.trigger(modal);
	}

	/**
	 * @param {string} proberId
	 * @param {string} proberName
	 */
	function showEditModal(proberId, proberName) {
		/** @type {import('@skeletonlabs/skeleton').ModalSettings} */
		const modal = {
			type: 'prompt',
			// Data
			title: 'Enter Name',
			body: 'Enter a new name for the prober',
			value: proberName,
			valueAttr: { type: 'text', minlength: 3, maxlength: 50, required: true },
			response: (r) => {
				if (!r) return;
				editProber(proberId, r)
					.then((res) => {
						if (res.status !== 'ok') {
							toastStore.trigger({ message: 'Failed to edit prober' });
						} else {
							toastStore.trigger({
								message: 'Prober edited',
								background: 'variant-filled-success'
							});
							handler.invalidate();
						}
					})
					.catch(() => {
						toastStore.trigger({ message: 'Failed to edit prober' });
					});
			}
		};
		modalStore.trigger(modal);
	}

	/** @param {number} id */
	const handleDelete = (id) => {
		modalStore.trigger({
			type: 'confirm',
			title: 'Please Confirm',
			body: 'Are you sure you wish to proceed?',
			/** @param {boolean} r */
			response: async (r) => {
				if (r) {
					deleteProber(id)
						.then((res) => {
							if (res.status !== 'ok') {
								toastStore.trigger({ message: 'Failed to delete prober' });
							} else {
								toastStore.trigger({
									message: 'Prober deleted',
									background: 'variant-filled-success'
								});
								handler.invalidate();
							}
						})
						.catch(() => {
							toastStore.trigger({ message: 'Prober could not be deleted' });
						});
				}
			}
		});
	};

	const handler = new DataHandler([], { rowsPerPage: 10, totalRows: 0 });
	let rows = handler.getRows();

	onMount(() => {
		handler.onChange((state) => loadData(state));
		handler.invalidate();
	});
</script>

<div class="mb-4">
	<h1 class="h2 font-extrabold dark:text-white">Prober</h1>
</div>

<div class="dashboard-card">
	<button class="variant-filled-success btn btn-sm mb-4" on:click={showAddModal}>Add Prober</button>
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

	<div class="my-2 overflow-x-auto">
		<table class="table table-hover table-compact w-full table-auto">
			<thead>
				<tr>
					<DtSrThSort {handler} orderBy="id">ID</DtSrThSort>
					<th>Name</th>
					<th>API Key</th>
					<DtSrThSort {handler} orderBy="last_submit_ts">Last Submit</DtSrThSort>
				</tr>
				<tr>
					<DtSrThFilter {handler} filterBy="name" placeholder="Name" colspan={2} />
					<DtSrThFilter {handler} filterBy="api_key" placeholder="API Key" colspan={2} />
				</tr>
			</thead>
			<tbody>
				{#each $rows as row (row.id)}
					<tr>
						<td>
							{row.id}
							<button
								class="variant-filled-secondary btn btn-sm mr-1"
								on:click={() => showEditModal(row.id, row.name)}>Edit</button
							>
							<button
								class="variant-filled-error btn btn-sm"
								name="delete_{row.id}"
								on:click={() => {
									handleDelete(row.id);
								}}>Delete</button
							>
						</td>
						<td>{row.name}</td>
						<td>{row.api_key}</td>
						<td>
							{format(row.last_submit_ts * 1000, 'PP HH:mm')}<br />
							{formatDistance(row.last_submit_ts * 1000, new Date(), { addSuffix: true })}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<div class="flex justify-between mb-2">
		<DtSrRowCount {handler} />
		<DtSrPagination {handler} />
	</div>
</div>
