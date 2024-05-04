<script>
	import { getModalStore } from '@skeletonlabs/skeleton';

	const modalStore = getModalStore();
	/** @type {string} */
	export let ip;
	/** @type {boolean} */
	export let is_tor;
	/** @type {string} */
	export let hostname;
	/** @type {number} */
	export let port;

	// if (is_tor) {
	// 	hostname = hostname.substring(0, 8) + '[...].onion';
	// }

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
		class="max-w-32 truncate text-orange-800 dark:text-orange-300"
		on:click={() => modalAlert(hostname, port)}
	>
		👁 {hostname}
	</button><br />.onion:<span class="text-indigo-800 dark:text-indigo-400">{port}</span>
	<span class="text-gray-700 dark:text-gray-400">(TOR)</span>
{:else}
	{hostname}:<span class="text-indigo-800 dark:text-indigo-400">{port}</span>
	{#if ip !== ''}
		<br /><span class="text-gray-700 dark:text-gray-400">{ip}</span>
	{/if}
{/if}
