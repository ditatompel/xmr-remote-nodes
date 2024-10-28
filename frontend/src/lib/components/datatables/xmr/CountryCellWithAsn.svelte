<script lang="ts">
	export let cc: string;
	export let country_name: string;
	export let city: string;
	export let asn: number;
	export let asn_name: string;

	$: lowerCc = cc.toLowerCase();
</script>

{#if cc != ''}
	{#if city !== ''}
		{city},
	{/if}
	{country_name}
	<img class="inline-block" src="/img/cf/{lowerCc}.svg" alt="{cc} Flag" width="22px" />
{/if}

{#if asn !== 0}
	<br /><a
		class="external !text-purple-800 dark:!text-purple-400"
		href="https://www.ditatompel.com/asn/{asn}"
		target="_blank"
		rel="noopener">AS{asn}</a
	>
	(<span class="font-semibold text-green-800 dark:text-green-500">{asn_name}</span>)
{/if}

<style lang="postcss">
	a {
		@apply font-semibold text-sky-800 underline dark:text-sky-500;
	}
</style>
