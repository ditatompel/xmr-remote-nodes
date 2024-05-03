<script lang="ts">
	import type { DataHandler, Row } from '@vincjo/datatables/remote';

	type T = $$Generic<Row>;

	export let handler: DataHandler<T>;

	const pageNumber = handler.getPageNumber();
	const pageCount = handler.getPageCount();
	const pages = handler.getPages({ ellipsis: true });

	const setPage = (value: 'previous' | 'next' | number) => {
		handler.setPage(value);
		handler.invalidate();
	};
</script>

<section class={$$props.class ?? ''}>
	{#if $pages === undefined}
		<button type="button" class="sm-btn" on:click={() => setPage('previous')}> &#10094; </button>
		<button class="mx-4">page <b>{$pageNumber}</b></button>
		<button type="button" class="sm-btn" on:click={() => setPage('next')}>&#10095;</button>
	{:else}
		<div class="lg:hidden">
			<button type="button" class="sm-btn" on:click={() => setPage('previous')}> &#10094; </button>
			<button class="mx-4">page <b>{$pageNumber}</b></button>
			<button
				class="sm-btn"
				class:disabled={$pageNumber === $pageCount}
				on:click={() => setPage('next')}
			>
				&#10095;
			</button>
		</div>

		<div class="btn-group variant-ghost-surface hidden lg:block">
			<button
				type="button"
				class="hover:variant-soft-secondary"
				class:disabled={$pageNumber === 1}
				on:click={() => setPage('previous')}>&#10094;</button
			>
			{#each $pages as page}<button
					type="button"
					class="hover:variant-filled-secondary"
					class:!variant-filled-primary={$pageNumber === page}
					class:ellipse={page === null}
					on:click={() => setPage(page)}>{page ?? '...'}</button
				>{/each}
			<button
				type="button"
				class="hover:variant-soft-secondary"
				class:disabled={$pageNumber === $pageCount}
				on:click={() => setPage('next')}
			>
				&#10095;
			</button>
		</div>
	{/if}
</section>

<style lang="postcss">
	.sm-btn {
		@apply btn btn-sm variant-ghost-surface hover:variant-soft-secondary;
	}
	.disabled {
		@apply cursor-not-allowed;
	}
</style>
