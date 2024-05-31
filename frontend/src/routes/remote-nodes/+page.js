/** @type {import('./$types').PageLoad} */
export async function load() {
	return {
		// prettier-ignore
		meta: {
			title: 'Public Monero Remote Nodes List',
			description: 'List of public Monero remote nodes that you can use with your favourite Monero wallet. You can filter by country, protocol, or CORS capable nodes.',
			keywords: 'monero remote nodes,public monero nodes,monero public nodes,monero wallet,tor monero node,monero cors rpc'
		},
		/**
		 * Array containing network fees.
		 * For now, I use static data to reduce the amount of API calls.
		 * See the values from `/api/v1/fees`
		 * @type {{ nettype: string, estimate_fee: number }[]}
		 */
		netFees: [
			{
				nettype: 'mainnet',
				estimate_fee: 20000
			},
			{
				nettype: 'stagenet',
				estimate_fee: 56000
			},
			{
				nettype: 'testnet',
				estimate_fee: 20000
			}
		]
	};
}
