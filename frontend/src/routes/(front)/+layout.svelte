<script>
	// import { writable } from 'svelte/store';
	import { page } from '$app/stores';
	import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
	import { beforeNavigate, afterNavigate } from '$app/navigation';
	import '../../app.css';
	import {
		Toast,
		Modal,
		Drawer,
		initializeStores,
		storePopup // PopUps
	} from '@skeletonlabs/skeleton';
	import { MainNav, MobileDrawer } from '$lib/components/navigation';
	import Footer from '$lib/components/Footer.svelte';

	initializeStores();

	// popups
	storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });

	/* prettier-ignore */
	const metaDefaults = {
		title: 'Monero is private, decentralized cryptocurrency that keeps your finances confidential and secure.',
		description: '',
		keywords: '',
		image:
			'https://vcl-og-img.ditatompel.com/' +
			encodeURIComponent('Monero is private, decentralized cryptocurrency that keeps your finances confidential and secure.') +
			'.png?md=0'
	};
	const meta = {
		title: metaDefaults.title,
		description: metaDefaults.description,
		keywords: metaDefaults.keywords,
		image: metaDefaults.image,
		article: { publishTime: '', modifiedTime: '', author: '' },
		// Twitter
		twitter: {
			title: metaDefaults.title,
			description: metaDefaults.description,
			image: metaDefaults.image
		}
	};

	page.subscribe((page) => {
		// Restore Page Defaults
		meta.title = metaDefaults.title;
		meta.description = metaDefaults.description;
		meta.keywords = metaDefaults.keywords;
		meta.image = metaDefaults.image;
		// Restore Twitter Defaults
		meta.twitter.title = metaDefaults.title;
		meta.twitter.description = metaDefaults.description;
		meta.twitter.image = metaDefaults.image;
		if (typeof page.data.meta === 'object') {
			meta.title = page.data.meta.title ?? metaDefaults.title;
			meta.description = page.data.meta.description ?? metaDefaults.description;
			meta.keywords = page.data.meta.keywords ?? metaDefaults.description;
			meta.image = page.data.meta.image ?? metaDefaults.image;
			// Restore Twitter Defaults
			meta.twitter.title = page.data.meta.title ?? metaDefaults.title;
			meta.twitter.description = page.data.meta.description ?? metaDefaults.description;
			meta.twitter.image = page.data.meta.image ?? metaDefaults.description;
			if (typeof page.data.meta.article === 'object') {
				meta.article.author = page.data.meta.article.author ?? '';
			}
		}
	});

	let isLoading = false;

	// progress bar show
	beforeNavigate(() => (isLoading = true));

	// scroll to top after nafigation and progress bar
	afterNavigate((/* params */) => {
		isLoading = false;
		// scroll to top when navigate
		// const isNewPage = params.from?.url.pathname !== params.to?.url.pathname;
		// const elemPage = document.querySelector('#page');
		// if (isNewPage && elemPage !== null) {
		// 	elemPage.scrollTop = 0;
		// }
	});
</script>

<svelte:head>
	<title>{meta.title} — xmr.ditatompel.com</title>
	<!-- Meta Tags -->
	<meta name="title" content="{meta.title} — ditatompel.com" />
	<meta name="description" content={meta.description} />
	<meta name="keywords" content={meta.keywords} />
	<meta name="theme-color" content="#272b31" />
	<meta name="author" content="ditatompel" />

	<!-- Open Graph - https://ogp.me/ -->
	<meta property="og:site_name" content="ditatompel.com" />
	<meta property="og:type" content="website" />
	<meta property="og:url" content="https://www.ditatompel.com{$page.url.pathname}" />
	<meta property="og:locale" content="en_US" />
	<meta property="og:title" content="{meta.title} — ditatompel.com" />
	<meta property="og:description" content={meta.description} />
	<meta property="og:image" content={meta.image} />
	<meta property="og:image:secure_url" content={meta.image} />
	<meta property="og:image:type" content="image/png" />
	<meta property="og:image:width" content="2048" />
	<meta property="og:image:height" content="1170" />

	<!-- Open Graph: Twitter -->
	<meta name="twitter:card" content="summary" />
	<meta name="twitter:site" content="@ditatompel" />
	<meta name="twitter:creator" content="@ditatompel" />
	<meta name="twitter:title" content="{meta.twitter.title} — ditatompel.com" />
	<meta name="twitter:description" content={meta.twitter.description} />
	<meta name="twitter:image" content={meta.twitter.image} />
</svelte:head>

<Modal />
<Toast />
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
