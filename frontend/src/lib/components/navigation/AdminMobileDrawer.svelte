<script>
	import { page } from '$app/stores';
	import { adminNavs } from './navs';
	import { getDrawerStore } from '@skeletonlabs/skeleton';

	const drawerStore = getDrawerStore();
	$: style = (/** @type {string} */ href) =>
		$page.url.pathname.startsWith(href) ? 'bg-primary-500' : '';
</script>

<nav class="list-nav p-4">
	<ul>
		{#each adminNavs as nav}
			<li>
				<a href={nav.path} class={style(nav.path)} on:click={() => drawerStore.close()}
					>{nav.name}</a
				>
			</li>
		{/each}
	</ul>
</nav>
