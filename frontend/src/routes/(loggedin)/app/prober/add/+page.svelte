<script>
	import { invalidateAll, goto } from '$app/navigation';
	import { apiUri } from '$lib/utils/common';
	import { ProgressBar } from '@skeletonlabs/skeleton';

	/**
	 * @typedef formResult
	 * @type {object}
	 * @property {string} status
	 * @property {string} message
	 * @property {null | object} data
	 */
	/** @type {formResult} */
	export let formResult;

	let isProcessing = false;

	/** @param {{ currentTarget: EventTarget & HTMLFormElement}} event */
	async function handleSubmit(event) {
		isProcessing = true;
		const data = new FormData(event.currentTarget);

		const response = await fetch(event.currentTarget.action, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json'
			},
			body: JSON.stringify(Object.fromEntries(data))
		});

		formResult = await response.json();
		isProcessing = false;

		if (formResult.status === 'ok') {
			// rerun all `load` functions, following the successful update
			await invalidateAll();
			goto('/app/prober/');
		}
	}
</script>

<div class="mb-4">
	<h1 class="h2 font-extrabold dark:text-white">Add Prober</h1>
</div>
{#if !isProcessing}
	{#if formResult?.status === 'error'}
		<div class="p-4 mb-4 text-sm rounded-lg bg-gray-700 text-red-400" role="alert">
			<span class="font-medium">Error:</span>
			{formResult.message}!
		</div>
	{/if}
	{#if formResult?.status === 'ok'}
		<div class="p-4 mb-4 text-sm rounded-lg bg-gray-700 text-green-400" role="alert">
			<span class="font-medium">Success:</span>
			{formResult.message}!
		</div>
	{/if}
{:else}
	<ProgressBar meter="bg-secondary-500" track="bg-secondary-500/30" value={undefined} />
	<div class="mx-4 p-4 mb-4 text-sm rounded-lg bg-gray-700 text-blue-400" role="alert">
		<span class="font-medium">Processing...</span>
	</div>
{/if}
<div class="dashboard-card">
	<form
		class="space-y-4 md:space-y-6"
		action={apiUri('/api/v1/prober')}
		method="POST"
		on:submit|preventDefault={handleSubmit}
	>
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<div>
				<label for="name" class="label">
					<span>Name</span>
					<input
						type="text"
						name="name"
						id="name"
						placeholder="Prober name"
						autocomplete="off"
						class="input variant-form-material"
					/>
				</label>
			</div>
		</div>

		<button
			type="submit"
			class="w-full rounded-lg bg-primary-600 px-5 py-2.5 text-center text-sm font-medium hover:bg-primary-700"
			>Submit</button
		>
	</form>
</div>
