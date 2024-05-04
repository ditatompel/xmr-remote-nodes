<script>
	import { applyAction, enhance } from '$app/forms';
	import { slide } from 'svelte/transition';
	import { ProgressBar } from '@skeletonlabs/skeleton';
	/** @type {import('./$types').PageData} */
	export let data;

	/** @type {import('./$types').ActionData} */
	export let form;

	let isProcessing = false;

	/** @type {import('./$types').SubmitFunction} */
	const handleForm = async () => {
		isProcessing = true;
		return async ({ result }) => {
			isProcessing = false;
			if (result.type === 'success' || result.type === 'redirect') {
				close();
			}
			await applyAction(result);
		};
	};
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

<section id="form-add-monero-node">
	<div class="section-container text-center">
		<p>Enter your Monero node information below (IPv4 host only):</p>

		<form class="mx-auto w-full max-w-3xl py-2" method="POST" use:enhance={handleForm}>
			<div class="grid grid-cols-1 gap-4 py-6 md:grid-cols-4">
				<label class="label">
					<span>Protocol *</span>
					<select name="protocol" class="select variant-form-material" disabled={isProcessing}>
						<option value="http">HTTP / TOR</option>
						<option value="https">HTTPS</option>
					</select>
				</label>
				<label class="label md:col-span-2">
					<span>Host / IP *</span>
					<input
						class="input variant-form-material"
						name="host"
						type="text"
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
				{#if form?.status === 'error'}
					<div class="alert variant-ghost-error" transition:slide={{ duration: 500 }}>
						<div class="alert-message">
							<h3 class="h3">Error!</h3>
							<p>{form.message}!</p>
						</div>
					</div>
				{/if}
				{#if form?.status === 'ok'}
					<div class="alert variant-ghost-success" transition:slide={{ duration: 500 }}>
						<div class="alert-message">
							<h3 class="h3">Success!</h3>
							<p>{form.message}!</p>
						</div>
					</div>
				{/if}
			{:else}
				<ProgressBar meter="bg-secondary-500" track="bg-secondary-500/30" value={undefined} />
			{/if}
		</div>

		<p>
			Here you can find list of <a class="anchor" href="/remote-nodes">Monero Remote Node</a>.
		</p>
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
</style>
