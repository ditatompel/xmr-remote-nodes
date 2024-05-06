/** @type {import('./$types').PageLoad} */
export async function load({ data }) {
	/* prettier-ignore */
	const metaDefaults = {
		title: 'Monero (XMR)',
		description: 'Monero is private, decentralized cryptocurrency that keeps your finances confidential and secure.',
		keywords: 'monero,xmr,monero node,xmrnode,cryptocurrency'
	};

	return {
		meta: {
			title: metaDefaults.title,
			description: metaDefaults.description,
			keywords: metaDefaults.keywords,
			image:
				'https://vcl-og-img.ditatompel.com/' + encodeURIComponent(metaDefaults.title) + '.png?md=0',
			// Article
			article: { publishTime: '', modifiedTime: '', author: '' },
			// Twitter
			twitter: {
				title: metaDefaults.title,
				description: metaDefaults.description,
				image: metaDefaults.image
			}
		},
	};
}
