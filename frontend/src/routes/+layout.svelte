<script>
	import '../app.css';
	import { page } from '$app/stores';
	import { Toast, Modal } from '@skeletonlabs/skeleton';
	import { beforeNavigate, afterNavigate } from '$app/navigation';
	import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
	import { ProgressBar, initializeStores, storePopup } from '@skeletonlabs/skeleton';

	initializeStores();
	storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });
	let isLoading = false;

	beforeNavigate(() => (isLoading = true));
	afterNavigate(() => {
		isLoading = false;
	});

	/* prettier-ignore */
	const meta = {
		title: 'Monero Remote Node',
		description: 'A website that helps you monitor your favourite Monero remote nodes, a device on the internet running the Monero software with copy of the Monero blockchain.',
		keywords: 'monero,monero,xmr,monero node,xmrnode,cryptocurrency,monero remote node,monero testnet,monero stagenet'
	};

	page.subscribe((page) => {
		if (typeof page.data.meta === 'object') {
			meta.title = page.data.meta.title ?? meta.title;
			meta.description = page.data.meta.description ?? meta.description;
			meta.keywords = page.data.meta.keywords ?? meta.description;
		}
	});
</script>

<svelte:head>
	<title>{meta.title} — xmr.ditatompel.com</title>
	<!-- Meta Tags -->
	<meta name="title" content="{meta.title} — xmr.ditatompel.com" />
	<meta name="description" content={meta.description} />
	<meta name="keywords" content={meta.keywords} />
	<meta name="theme-color" content="#272b31" />
	<meta name="author" content="ditatompel" />

	<!-- Open Graph - https://ogp.me/ -->
	<meta property="og:site_name" content="xmr.ditatompel.com" />
	<meta property="og:type" content="website" />
	<meta property="og:url" content="https://xmr.ditatompel.com{$page.url.pathname}" />
	<meta property="og:locale" content="en_US" />
	<meta property="og:title" content="{meta.title} — xmr.ditatompel.com" />
	<meta property="og:description" content={meta.description} />
</svelte:head>

<Modal />
<Toast />

{#if isLoading}
	<ProgressBar
		class="fixed top-0 z-50"
		height="h-1"
		track="bg-opacity-100"
		meter="bg-gradient-to-br from-purple-600 via-pink-600 to-blue-600"
	/>
{/if}
<slot />
