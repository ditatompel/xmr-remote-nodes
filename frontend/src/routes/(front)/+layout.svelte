<script>
	import { page } from '$app/stores';
	import '../../app.css';
	import { Drawer } from '@skeletonlabs/skeleton';
	import { MainNav, MobileDrawer } from '$lib/components/navigation';
	import Footer from '$lib/components/Footer.svelte';

	/* prettier-ignore */
	const metaDefaults = {
		title: 'Monero is private, decentralized cryptocurrency that keeps your finances confidential and secure.',
		description: '',
		keywords: '',
	};
	const meta = {
		title: metaDefaults.title,
		description: metaDefaults.description,
		keywords: metaDefaults.keywords,
		// Twitter
		twitter: {
			title: metaDefaults.title,
			description: metaDefaults.description
		}
	};

	page.subscribe((page) => {
		meta.title = metaDefaults.title;
		meta.description = metaDefaults.description;
		meta.keywords = metaDefaults.keywords;
		meta.twitter.title = metaDefaults.title;
		meta.twitter.description = metaDefaults.description;
		if (typeof page.data.meta === 'object') {
			meta.title = page.data.meta.title ?? metaDefaults.title;
			meta.description = page.data.meta.description ?? metaDefaults.description;
			meta.keywords = page.data.meta.keywords ?? metaDefaults.description;
			meta.twitter.title = page.data.meta.title ?? metaDefaults.title;
			meta.twitter.description = page.data.meta.description ?? metaDefaults.description;
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

	<!-- Open Graph: Twitter -->
	<meta name="twitter:card" content="summary" />
	<meta name="twitter:site" content="@ditatompel" />
	<meta name="twitter:creator" content="@ditatompel" />
	<meta name="twitter:title" content="{meta.twitter.title} — xmr.ditatompel.com" />
	<meta name="twitter:description" content={meta.twitter.description} />
</svelte:head>

<Drawer>
	<h2 class="p-4">Navigation</h2>
	<hr />
	<MobileDrawer />
	<hr />
</Drawer>

<MainNav />

<div class="pt-10 md:pt-12 min-h-screen">
	<slot />
</div>

<Footer />
