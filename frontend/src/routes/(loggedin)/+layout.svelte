<script>
	import { Drawer } from '@skeletonlabs/skeleton';
	import { AdminNav, AdminSidebar, AdminMobileDrawer } from '$lib/components/navigation';
</script>

<svelte:head>
	<meta name="robots" content="noindex,nofollow" />
</svelte:head>
<Drawer>
	<h2 class="p-4">Navigation</h2>
	<hr />
	<AdminMobileDrawer />
	<hr />
</Drawer>

<AdminNav />
<AdminSidebar />

<div class="min-h-screen bg-gray-100/80 p-4 pt-14 dark:bg-gray-900/80 sm:ml-64">
	<slot />
</div>
