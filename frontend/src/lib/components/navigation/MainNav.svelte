<script>
	import { page } from '$app/stores';
	import { navs } from './navs';
	import { LightSwitch, getDrawerStore } from '@skeletonlabs/skeleton';

	const drawerStore = getDrawerStore();
</script>

<nav class="fixed w-full z-20 top-0 start-0 bg-surface-100-800-token shadow-2xl">
	<div class="mx-auto flex max-w-screen-xl flex-wrap items-center justify-between px-4 py-1">
		<a href="/" class="flex items-center space-x-3 rtl:space-x-reverse" aria-label="ditatompel">
			<span
				class="txt-logo-gradient hidden self-center whitespace-nowrap text-2xl font-semibold lg:block"
				>XMR Nodes</span
			>
		</a>
		<div class="flex items-center space-x-1 md:order-2 md:space-x-0 rtl:space-x-reverse">
			<LightSwitch />
			<button
				class="btn btn-sm mr-4 md:hidden"
				aria-label="Mobile Drawer Button"
				on:click={() => drawerStore.open({})}
			>
				<span>
					<svg viewBox="0 0 100 80" class="fill-token h-4 w-4">
						<rect width="100" height="20" />
						<rect y="30" width="100" height="20" />
						<rect y="60" width="100" height="20" />
					</svg>
				</span>
			</button>
		</div>
		<div class="hidden w-full items-center justify-between md:order-1 md:flex md:w-auto">
			<ul
				class="flex flex-row space-x-1 rounded-lg bg-white p-0 dark:bg-gray-900 rtl:space-x-reverse"
			>
				<li>
					<a
						href="/"
						class={$page.url.pathname === '/' ? 'active' : 'nav-link'}
						aria-current={$page.url.pathname === '/' ? 'page' : undefined}>Home</a
					>
				</li>
				{#each navs as nav}
					<li>
						<a
							href={nav.path}
							class={$page.url.pathname.startsWith(nav.path) ? 'active' : 'nav-link'}
						>
							{nav.name}
						</a>
					</li>
				{/each}
			</ul>
		</div>
	</div>
</nav>

<style lang="postcss">
	.active {
		@apply block rounded bg-primary-500 p-2 text-black;
	}
	.nav-link {
		@apply block rounded hover:bg-secondary-500 md:p-2 hover:text-black;
	}
</style>
