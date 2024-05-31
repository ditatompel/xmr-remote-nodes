/** @type {import('./$types').PageLoad} */
export async function load() {
	return {
		// prettier-ignore
		meta: {
			title: 'Public Monero Remote Nodes List',
			description: 'List of public Monero remote nodes that you can use with your favourite Monero wallet. You can filter by country, protocol, or CORS capable nodes.',
			keywords: 'monero remote nodes,public monero nodes,monero public nodes,monero wallet,tor monero node,monero cors rpc'
		}
	};
}
