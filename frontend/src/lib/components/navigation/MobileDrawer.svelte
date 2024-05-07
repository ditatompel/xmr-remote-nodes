<script>
	import { page } from '$app/stores';
	import { navs } from './navs';
	import { getDrawerStore } from '@skeletonlabs/skeleton';

	const drawerStore = getDrawerStore();
	$: classes = (/** @type {string} */ href) =>
		$page.url.pathname.startsWith(href) ? 'bg-primary-500' : '';
</script>

<nav class="list-nav p-4">
	<ul>
		<li>
			<a
				href="/"
				class={$page.url.pathname === '/' ? 'bg-primary-500' : ''}
				on:click={() => drawerStore.close()}>Home</a
			>
		</li>
		{#each navs as nav}
			<li>
				<a href={nav.path} class={classes(nav.path)} on:click={() => drawerStore.close()}
					>{nav.name}</a
				>
			</li>
		{/each}
	</ul>
</nav>
