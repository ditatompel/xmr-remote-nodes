<script>
	import { invalidateAll, goto } from '$app/navigation';
	import { apiUri } from '$lib/utils/common';
	import { ProgressBar, LightSwitch } from '@skeletonlabs/skeleton';

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
			goto('/app/dashboard/');
		}
	}
</script>

<section class="bg-gray-50 dark:bg-gray-900">
	<div class="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0">
		<a href="/" class="flex items-center mb-6 text-2xl font-semibold text-gray-900 dark:text-white"
			>XMR Nodes</a
		>
		<div
			class="w-full rounded-lg shadow border md:mt-0 sm:max-w-md xl:p-0 bg-white border-gray-700 dark:bg-gray-800"
		>
			<div class="p-6 space-y-4 md:space-y-6 sm:p-8">
				<h1
					class="text-xl font-bold leading-tight tracking-tight tmd:text-2xl text-gray-900 dark:text-white"
				>
					Sign in to your account
				</h1>
				<form
					class="space-y-4 md:space-y-6"
					action={apiUri('/auth/login')}
					method="POST"
					on:submit|preventDefault={handleSubmit}
				>
					<div>
						<label
							for="username"
							class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Username</label
						>
						<input
							type="text"
							name="username"
							id="username"
							class="input"
							placeholder="username"
							required
						/>
					</div>
					<div>
						<label
							for="password"
							class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label
						>
						<input
							type="password"
							name="password"
							id="password"
							placeholder="••••••••"
							class="input"
							required
						/>
					</div>
					<button type="submit" class="btn variant-filled-primary w-full">Sign in</button>

					<LightSwitch />
				</form>
			</div>

			{#if !isProcessing}
				{#if formResult?.status === 'error'}
					<div class="mx-4 p-4 mb-4 text-sm rounded-lg bg-gray-700 text-red-400" role="alert">
						<span class="font-medium">Error:</span>
						{formResult.message}!
					</div>
				{/if}
				{#if formResult?.status === 'ok'}
					<div class="mx-4 p-4 mb-4 text-sm rounded-lg bg-gray-700 text-green-400" role="alert">
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
		</div>
	</div>
</section>
