<script>
	import '../../app.css';
	import { Drawer } from '@skeletonlabs/skeleton';
	import { MainNav, MobileDrawer } from '$lib/components/navigation';
	import Footer from '$lib/components/Footer.svelte';
</script>

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
