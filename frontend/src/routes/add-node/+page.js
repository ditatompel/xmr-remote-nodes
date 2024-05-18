/** @type {import('./$types').PageLoad} */
export async function load({ data }) {
	/* prettier-ignore */
	return {
		meta: {
			title: 'Add Monero Node',
			description:
				'You can use this page to add known remote node to the system so my bots can monitor it.',
			keywords: 'monero,monero node,monero public node,monero wallet,list monero node,monero node monitoring'
		}
	};
}
