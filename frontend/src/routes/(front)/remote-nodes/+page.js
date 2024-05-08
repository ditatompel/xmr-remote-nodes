/** @type {import('./$types').PageLoad} */
export async function load({ data }) {
	return {
		// prettier-ignore
		meta: {
			title: 'Public Monero Remote Nodes List',
			description:
				'Monero is private, decentralized cryptocurrency that keeps your finances confidential and secure.',
			keywords: 'monero remote nodes,public monero nodes,monero public nodes,monero wallet'
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
				estimate_fee: 58000
			},
			{
				nettype: 'testnet',
				estimate_fee: 20000
			}
		]
	};
}
