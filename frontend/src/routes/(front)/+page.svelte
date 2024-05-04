<script>
	import { clipboard } from '@skeletonlabs/skeleton';
	import {
		IcnGitHub,
		IcnMonero,
		IcnReddit,
		IcnTwitter,
		IcnFacebook,
		IcnTelegram
	} from '$lib/components/svg';

	/** @type {import('./$types').PageData} */
	export let data;
	let currentVariant = 'bg-initial';
	let donationCopied = false;

	/** @param {Event & { target: HTMLInputElement }} e */
	function copyHandler(e) {
		e.target.disabled = true;
		e.target.innerText = 'Copied 👍';
		setTimeout(() => {
			e.target.innerText = 'Copy';
			e.target.disabled = false;
		}, 1000);
	}

	function copyDonationAddr() {
		donationCopied = true;
		setTimeout(() => {
			donationCopied = false;
		}, 2500);
	}

	let snetProtos = [
		{ label: 'P2P', value: 'stagenet.xmr.ditatompel.com:38080', key: 'snetP2P' },
		{ label: 'RPC', value: 'stagenet.xmr.ditatompel.com:38089', key: 'snetRPC' },
		{ label: 'RPC SSL', value: 'stagenet.xmr.ditatompel.com:443', key: 'snetSSL' }
	];

	let tnetProtos = [
		{ label: 'P2P', value: 'testnet.xmr.ditatompel.com:28080', key: 'tnetP2P' },
		{ label: 'RPC', value: 'testnet.xmr.ditatompel.com:28089', key: 'tnetRPC' },
		{ label: 'RPC SSL', value: 'testnet.xmr.ditatompel.com:443', key: 'tnetSSL' }
	];
</script>

<header id="hero" class="hero-gradient py-7">
	<div class="section-container text-center">
		<h1 class="h1 pb-2 font-extrabold">{data.meta.title}</h1>
		<p>{data.meta.description}</p>

		<div class="pt-2">
			<a
				href="https://www.getmonero.org"
				class="variant-ghost chip mt-2 hover:variant-filled"
				target="_blank"
				rel="noopener"
			>
				<span><IcnMonero class="h-4 w-4" /></span>
				<span>GetMonero.org</span>
			</a>
			<a
				href="https://github.com/monero-project"
				class="variant-ghost chip mt-2 hover:variant-filled"
				target="_blank"
				rel="noopener"
			>
				<span><IcnGitHub fill="currentColor" class="h-4 w-4" /></span>
				<span>monero-project</span>
			</a>
			<a
				href="https://www.reddit.com/r/Monero/"
				class="variant-ghost chip mt-2 hover:variant-filled"
				target="_blank"
				rel="noopener"
			>
				<span><IcnReddit fill="currentColor" class="h-4 w-4" /></span>
				<span>r/Monero</span>
			</a>
			<a
				href="https://twitter.com/monero"
				class="variant-ghost chip mt-2 hover:variant-filled"
				target="_blank"
				rel="noopener"
			>
				<span><IcnTwitter fill="currentColor" class="h-4 w-4" /></span>
				<span>@monero</span>
			</a>
			<a
				href="https://www.facebook.com/monerocurrency/"
				class="variant-ghost chip mt-2 hover:variant-filled"
				target="_blank"
				rel="noopener"
			>
				<span><IcnFacebook fill="currentColor" class="h-4 w-4" /></span>
				<span>monerocurrency</span>
			</a>
			<a
				href="https://telegram.me/monero"
				class="variant-ghost chip mt-2 hover:variant-filled"
				target="_blank"
				rel="noopener"
			>
				<span><IcnTelegram fill="currentColor" class="h-4 w-4" /></span>
				<span>monero</span>
			</a>
		</div>
	</div>
	<div class="mx-auto w-full max-w-3xl px-20">
		<hr class="!border-primary-400-500-token !border-t-4 !border-double" />
	</div>
</header>

<section id="introduction">
	<div class="section-container text-center">
		<p>If you're new to Monero, the official links above is a perfect place to start.</p>
		<p class="py-2">
			Of course, there are lots of personal and community sites which generally discusses a lot
			about Monero, such as <a
				href="https://moneroworld.com"
				class="external"
				target="_blank"
				rel="noopener">moneroworld.com</a
			>,
			<a href="https://www.monero.how" class="external" target="_blank" rel="noopener">monero.how</a
			>,
			<a href="https://www.monero.observer" class="external" target="_blank" rel="noopener"
				>monero.observer</a
			>,
			<a href="https://sethforprivacy.com" class="external" target="_blank" rel="noopener"
				>sethforprivacy.com</a
			>,
			<a
				href="https://localmonero.co
/knowledge"
				class="external"
				target="_blank"
				rel="noopener">localmonero.co</a
			>,
			<a href="https://revuo-xmr.com/" class="external" target="_blank" rel="noopener"
				>revuo-xmr.com</a
			>, etc, can be an other good reference for you.
		</p>
		<p>You can find few resources I provide related to Monero below:</p>
	</div>

	<div class="section-container text-token grid grid-cols-1 gap-2 md:grid-cols-3">
		<a
			class="card {currentVariant} card-hover overflow-hidden py-2 text-center"
			href="/remote-nodes"
		>
			<h2 class="h2 font-bold">Remote Nodes</h2>
			<div class="space-y-4 p-4">
				<p>
					<!-- cspell:disable -->
					List of Monero remote nodes you can use when you cannot (or don't want to) run your own node.
					<!-- cspell:enable -->
				</p>
			</div>
		</a>

		<a class="card {currentVariant} card-hover overflow-hidden py-2 text-center" href="/add-node">
			<h2 class="h2 font-bold">Add Node</h2>
			<div class="space-y-4 p-4">
				<p>
					<!-- cspell:disable -->
					Add your Monero public node so others can connect to your node.
					<!-- cspell:enable -->
				</p>
			</div>
		</a>

		<a
			class="card {currentVariant} card-hover overflow-hidden py-2 text-center"
			href="https://monitor.ditatompel.com/d/xmr_metrics/monero-metrics?orgId=2"
			target="_blank"
			rel="noopener"
		>
			<h2 class="h2 font-bold">Metrics</h2>
			<div class="space-y-4 p-4">
				<p>
					<!-- cspell:disable -->
					Collection of my Monero metrics (GitHub repository, blockchain, market, P2Pool) presented through
					Grafana.
					<small>* Will open in new window</small>
					<!-- cspell:enable -->
				</p>
			</div>
		</a>
	</div>
</section>

<section id="my-monero-public-nodes" class="bg-surface-100-800-token">
	<div class="section-container text-token grid grid-cols-1 gap-10 md:grid-cols-2">
		<div class="text-center">
			<h2 class="h2 pb-2 font-bold">My Stagenet Public Node</h2>
			<p>
				Stagenet is what you need to learn Monero safely. Stagenet is technically equivalent to
				mainnet, both in terms of features and consensus rules.
			</p>
			{#each snetProtos as { label, value, key }}
				<div class="input-group input-group-divider my-2 grid-cols-[auto_1fr_auto]">
					<div class="input-group-shim"><label for={key}>{label}</label></div>
					<input class="text-center" type="text" id={key} name={key} {value} data-clipboard={key} />
					<button
						class="variant-filled-secondary"
						use:clipboard={{ input: key }}
						on:click={copyHandler}>Copy</button
					>
				</div>
			{/each}
		</div>

		<div class="text-center">
			<h2 class="h2 pb-2 font-bold">My Testnet Public Node</h2>
			<p>
				Testnet is the <em>"experimental"</em> network and blockchain where things get released long
				before mainnet. As a normal user, use mainnet instead.
			</p>
			{#each tnetProtos as { label, value, key }}
				<div class="input-group input-group-divider my-2 grid-cols-[auto_1fr_auto]">
					<div class="input-group-shim"><label for={key}>{label}</label></div>
					<input class="text-center" type="text" id={key} name={key} {value} data-clipboard={key} />
					<button
						class="variant-filled-secondary"
						use:clipboard={{ input: key }}
						on:click={copyHandler}>Copy</button
					>
				</div>
			{/each}
		</div>
	</div>
</section>

<section id="privacy-quote">
	<div class="text-token mx-auto w-full max-w-4xl py-4 text-center">
		<blockquote class="blockquote">
			<p class="text-3xl">
				Since we desire privacy, we must ensure that each party to a transaction have knowledge only
				of that which is directly necessary for that transaction.
			</p>
			<p class="my-2">
				<strong>Eric Hughes</strong> in
				<a
					href="https://www.activism.net/cypherpunk/manifesto.html"
					class="external"
					target="_blank"
					rel="noopener"><cite title="Source Title">A Cypherpunk's Manifesto</cite></a
				>.
			</p>
		</blockquote>
	</div>
</section>

<section id="monero-donation" class="section-container text-token text-center">
	<div class="mx-auto flex w-full max-w-4xl flex-row items-center gap-10">
		<div class="md:basis-3/4">
			<label for="donationAddress">If you like to buy me a coffee, here is my Monero address:</label
			>
			<textarea
				class="textarea my-2"
				id="donationAddress"
				name="donationAddress"
				data-clipboard="donationAddress"
				readonly
				>8BWYe6GzbNKbxe3D8mPkfFMQA2rViaZJFhWShhZTjJCNG6EZHkXRZCKHiuKmwwe4DXDYF8KKcbGkvNYaiRG3sNt7JhnVp7D</textarea
			>
			<button
				class="variant-filled-success btn"
				use:clipboard={{ input: 'donationAddress' }}
				disabled={donationCopied}
				on:click={copyDonationAddr}
				>{donationCopied ? 'Donation Address Copied! 🤩' : 'Copy Donation Address'}</button
			>
		</div>
		<div class="md:basis-1/4">
			<img
				src="https://w3-s3-jkt1.ditatompel.com/pub/assets/img/fav/monerotip.png"
				alt="ditatompel monero address"
			/>
			<p>Thank you so much! 🥰</p>
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
</style>
