<script>
	import { page } from '$app/stores';
	import { navs } from './navs';
</script>

<aside
	id="logo-sidebar"
	class="bg-surface-100-800-token fixed left-0 top-0 z-20 h-screen w-64 -translate-x-full pt-20 shadow-2xl transition-transform sm:translate-x-0"
	aria-label="Sidebar"
>
	<div class="h-full overflow-y-auto px-3 pb-4">
		<ul class="space-y-2 font-medium list-none" data-sveltekit-preload-data="false">
			{#each navs as nav}
				<li>
					<a
						href={nav.path}
						class={$page.url.pathname.startsWith(nav.path) ? 'active' : 'nav-link'}
					>
						<span class="ms-3">{nav.name}</span>
					</a>
				</li>
			{/each}
		</ul>
	</div>
</aside>

<style lang="postcss">
	.active {
		@apply flex items-center rounded-lg bg-primary-500 p-2;
	}
	.nav-link {
		@apply flex items-center rounded-lg p-2 hover:bg-secondary-500 hover:text-white;
	}
</style>
