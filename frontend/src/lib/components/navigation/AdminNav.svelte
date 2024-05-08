<script>
	import { invalidateAll, goto } from '$app/navigation';
	import { LightSwitch, getDrawerStore } from '@skeletonlabs/skeleton';
	import { apiUri } from '$lib/utils/common';

	const drawerStore = getDrawerStore();
	/** @type {ApiResponse} */
	let formResult;

	/** @param {{ currentTarget: EventTarget & HTMLFormElement}} event */
	async function handleLogout(event) {
		const data = new FormData(event.currentTarget);

		const response = await fetch(event.currentTarget.action, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json'
			},
			body: JSON.stringify(Object.fromEntries(data))
		});

		formResult = await response.json();

		if (formResult.status === 'ok') {
			await invalidateAll();
			goto('/login/');
		}
	}
</script>

<nav class="bg-surface-100-800-token fixed top-0 z-30 w-full shadow-2xl">
	<div class="px-3 py-2 lg:px-5 lg:pl-3">
		<div class="flex items-center justify-between">
			<div class="flex items-center justify-start rtl:justify-end">
				<button
					class="btn btn-sm inline-flex items-center md:hidden"
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
				<a href="/app/prober/" class="ms-2 flex md:me-24" aria-label="title">
					<span class="hidden self-center whitespace-nowrap text-2xl font-semibold lg:block"
						>XMR Nodes</span
					>
				</a>
			</div>
			<div class="flex items-center">
				<div class="ms-3 flex items-center space-x-4">
					<LightSwitch />
					<form
						action={apiUri('/auth/logout')}
						method="POST"
						on:submit|preventDefault={handleLogout}
					>
						<input type="hidden" name="logout" value="logout" />
						<button type="submit" class="btn btn-sm variant-filled-error" role="menuitem">
							Sign out
						</button>
					</form>
				</div>
			</div>
		</div>
	</div>
</nav>
