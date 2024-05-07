<script>
	// import { base } from '$app/paths';
	import '../app.css';
	import { beforeNavigate, afterNavigate } from '$app/navigation';
	import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
	import {
		ProgressBar,
		initializeStores,
		storePopup // PopUps
	} from '@skeletonlabs/skeleton';

	initializeStores();

	// popups
	storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });

	let isLoading = false;

	// progress bar show
	beforeNavigate(() => (isLoading = true));

	afterNavigate(() => {
		isLoading = false;
	});
</script>

{#if isLoading}
	<ProgressBar
		class="fixed top-0 z-50"
		height="h-1"
		track="bg-opacity-100"
		meter="bg-gradient-to-br from-purple-600 via-pink-600 to-blue-600"
	/>
{/if}
<slot />
