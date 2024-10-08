<script>
	import { invalidateAll, goto } from '$app/navigation';
	import { apiUri } from '$lib/utils/common';
	import { ProgressBar } from '@skeletonlabs/skeleton';

	/** @type {import('./$types').PageData} */
	export let data;
	/** @type {ApiResponse} */
	export let formResult;

	let isProcessing = false;

	/** @param {{ currentTarget: EventTarget & HTMLFormElement}} event */
	async function handleSubmit(event) {
		isProcessing = true;
		const data = new FormData(event.currentTarget);

		const response = await fetch(event.currentTarget.action, {
			method: 'POST',
			body: data
		});

		formResult = await response.json();
		isProcessing = false;

		if (formResult.status === 'ok') {
			await invalidateAll();
			goto('/remote-nodes');
		}
	}
</script>

<header id="hero" class="hero-gradient py-7">
	<div class="section-container text-center">
		<h1 class="h1 pb-2 font-extrabold">{data.meta.title}</h1>
		<p>{data.meta.description}</p>
	</div>
	<div class="mx-auto w-full max-w-3xl px-20">
		<hr class="!border-primary-400-500-token !border-t-4 !border-double" />
	</div>
</header>

<section id="page-info" class="mx-auto w-full max-w-4xl px-4 pb-7">
	<div class="alert card shadow-xl">
		<div class="alert-message">
			<h2 class="h3 text-center">Important Note</h2>
			<ul class="list-inside list-disc">
				<li>
					As an administrator of this instance, I have full rights to delete, and blacklist any
					submitted node with or without providing any reason.
				</li>
			</ul>
		</div>
	</div>
</section>

<section id="form-add-monero-node">
	<div class="section-container text-center">
		<p>Enter your Monero node information below (IPv6 host check is experimental):</p>

		<form
			class="mx-auto w-full max-w-3xl py-2"
			action={apiUri('/api/v1/nodes')}
			method="POST"
			on:submit|preventDefault={handleSubmit}
		>
			<div class="grid grid-cols-1 gap-4 py-6 md:grid-cols-4">
				<label class="label">
					<span>Protocol *</span>
					<select name="protocol" class="select variant-form-material" disabled={isProcessing}>
						<option value="http">HTTP or TOR</option>
						<option value="https">HTTPS</option>
					</select>
				</label>
				<label class="label md:col-span-2">
					<span>Host / IP *</span>
					<input
						class="input variant-form-material"
						name="hostname"
						type="text"
						required
						placeholder="Eg: node.example.com or 172.16.17.18"
						disabled={isProcessing}
					/>
				</label>
				<label class="label">
					<span>Port *</span>
					<input
						class="input variant-form-material"
						name="port"
						type="number"
						required
						placeholder="Eg: 18081"
						disabled={isProcessing}
					/>
				</label>
			</div>
			<button class="variant-filled-success btn" disabled={isProcessing}
				>{isProcessing ? 'Processing...' : 'Submit'}</button
			>
		</form>

		<div class="mx-auto w-full max-w-3xl py-2">
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

		<p>
			Here you can find list of <a class="anchor" href="/remote-nodes/">Monero Remote Node</a>.
		</p>
	</div>
</section>
