/** @type {import('./$types').PageLoad} */
export async function load({ data }) {
	/* prettier-ignore */
	const metaDefaults = {
		title: 'Probe Logs',
		description: 'Monero is private, decentralized cryptocurrency that keeps your finances confidential and secure.',
		keywords: 'monero,xmr,monero node,xmrnode,cryptocurrency'
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
