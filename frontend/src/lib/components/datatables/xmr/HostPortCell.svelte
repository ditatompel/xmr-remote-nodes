<script>
	import { getModalStore } from '@skeletonlabs/skeleton';
	import { formatHostname } from '$lib/utils/strings';

	const modalStore = getModalStore();

	/**
	 * @type {{
	 *  is_tor: boolean,
	 *  hostname: string,
	 *  port: number,
	 *  ipv6_only: boolean
	 * }}
	 */
	export let is_tor;
	export let hostname;
	export let port;
	export let ipv6_only;
	/** @type {string} */
	export let ip_addresses;

	/**
	 * @param {string} onionAddr
	 * @param {number} port
	 */
	function modalAlert(onionAddr, port) {
		/** @typedef {import('@skeletonlabs/skeleton').ModalSettings} ModalSettings */
		/** @type {ModalSettings} */
		const modal = {
			type: 'alert',
			title: 'Hostname:',
			body: '<code class="code">' + onionAddr + ':' + port + '</code>'
		};
		modalStore.trigger(modal);
	}
</script>

{#if is_tor}
	<button
		class="max-w-40 truncate text-orange-800 dark:text-orange-300"
		on:click={() => modalAlert(hostname, port)}
	>
		👁 {hostname}
	</button><br />.onion:<span class="text-indigo-800 dark:text-indigo-400">{port}</span>
	<span class="text-gray-700 dark:text-gray-400">(TOR)</span>
{:else}
	{formatHostname(hostname)}:<span class="text-indigo-800 dark:text-indigo-400">{port}</span><br />
	<div class="max-w-40 text-ellipsis overflow-x-auto md:overflow-hidden hover:overflow-visible">
		<span class="whitespace-break-spaces text-gray-700 dark:text-gray-400"
			>{ip_addresses.replace(/,/g, ' ')}</span
		>
		{#if ipv6_only}
			<span class="text-rose-800 dark:text-rose-400">(IPv6 only)</span>
		{/if}
	</div>
{/if}
