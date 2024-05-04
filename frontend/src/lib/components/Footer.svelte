<script>
	import { IcnGitHub } from '$lib/components/svg';
</script>

<div class="flex w-full items-end border-t border-surface-500/10 bg-surface-50 dark:bg-surface-900">
	<footer class="w-full">
		<div
			class="md:flex-no-wrap container mx-auto flex max-w-screen-xl flex-col flex-wrap px-5 py-10 md:flex-row md:items-center lg:items-start"
		>
			<div class="mx-auto w-64 flex-shrink-0 text-center lg:text-left">
				<a
					href="/"
					class="title-font flex items-center justify-center font-medium md:justify-start"
				>
					<span class="txt-logo-gradient text-5xl font-bold">XMR Nodes</span>
				</a>
				<p class="mt-2">
					By <a href="https://www.ditatompel.com">ditatompel</a>,
					<a href="https://github.com/ditatompel/xmr-nodes" target="_blank" rel="noopener"
						>source code</a
					>.
				</p>
				<div class="mt-4">
					<span class="mt-2 inline-flex justify-center sm:ml-auto sm:mt-0 sm:justify-start">
						<a
							href="https://github.com/ditatompel/xmr-nodes"
							target="_blank"
							class="cursor-pointer text-gray-500 hover:text-gray-700"
							aria-label="Link to ditatompel's GitHub profile"
						>
							<IcnGitHub class="h-5 w-5" />
						</a>
					</span>
				</div>
			</div>
		</div>
		<div class="bg-surface-500/5">
			<div class="container mx-auto px-5 py-4">
				<p class="text-center text-sm">XMR Nodes</p>
			</div>
		</div>
	</footer>
</div>
