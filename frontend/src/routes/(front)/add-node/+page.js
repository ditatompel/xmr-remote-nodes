/** @type {import('./$types').PageLoad} */
export async function load({ data }) {
	/* prettier-ignore */
	const metaDefaults = {
		title: 'Add Monero Node',
		description: 'You can use this page to add known remote node to the system so my bots can monitor it.',
		keywords: 'monero, monero node, monero public node, monero wallet'
	};

	return {
		meta: {
			title: metaDefaults.title,
			description: metaDefaults.description,
			keywords: metaDefaults.keywords,
			// Twitter
			twitter: {
				title: metaDefaults.title,
				description: metaDefaults.description
			}
		}
	};
}
