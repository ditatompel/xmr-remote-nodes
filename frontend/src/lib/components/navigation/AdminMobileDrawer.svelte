<script>
	import { page } from '$app/stores';
	import { navs } from './navs';
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
		{#each navs as nav}
			<li>
				<a href={nav.path} class={classesActive(nav.path)} on:click={drawerClose}>{nav.name}</a>
			</li>
		{/each}
	</ul>
</nav>
