<script>
	import { page } from '$app/stores';
	import { getDrawerStore } from '@skeletonlabs/skeleton';

	const drawerStore = getDrawerStore();

	function drawerClose() {
		drawerStore.close();
	}

	$: classesActive = (/** @type {string} */ href) =>
		$page.url.pathname.startsWith(href) ? 'bg-primary-500' : '';
</script>

<nav class="list-nav p-4">
	<ul>
		<li>
			<a href="/" class={$page.url.pathname === '/' ? 'bg-primary-500' : ''} on:click={drawerClose}
				>Home</a
			>
		</li>
		<li><a href="/dib" class={classesActive('/dib')} on:click={drawerClose}>Data</a></li>
		<li><a href="/tools" class={classesActive('/tools')} on:click={drawerClose}>Tools</a></li>
    <li><a href="/asn" class={classesActive('/asn')} on:click={drawerClose}>ASN</a></li>
		<li><a href="/proxy" class={classesActive('/proxy')} on:click={drawerClose}>Proxy</a></li>
		<li><a href="/monero" class={classesActive('/monero')} on:click={drawerClose}>Monero</a></li>
	</ul>
</nav>
