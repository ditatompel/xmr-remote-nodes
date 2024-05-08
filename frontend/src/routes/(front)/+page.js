/** @type {import('./$types').PageLoad} */
export async function load({ data }) {
	return {
		meta: {
			title: 'Monero (XMR)',
			description:
				'Monero is private, decentralized cryptocurrency that keeps your finances confidential and secure.',
			keywords: 'monero,xmr,monero node,xmrnode,cryptocurrency'
		},
		links: [
			{ text: 'moneroworld.com', uri: 'https://moneroworld.com' },
			{ text: 'monero.how', uri: 'https://www.monero.how' },
			{ text: 'monero.observer', uri: 'https://www.monero.observer' },
			{ text: 'sethforprivacy.com', uri: 'https://sethforprivacy.com' },
			{ text: 'localmonero.co', uri: 'https://localmonero.co/knowledge' },
			{ text: 'revuo-xmr.com', uri: 'https://revuo-xmr.com/' }
		],
		stagenet: [
			{ label: 'P2P', value: 'testnet.xmr.ditatompel.com:28080', key: 'tnetP2P' },
			{ label: 'RPC', value: 'testnet.xmr.ditatompel.com:28089', key: 'tnetRPC' },
			{ label: 'RPC SSL', value: 'testnet.xmr.ditatompel.com:443', key: 'tnetSSL' }
		],
		testnet: [
			{ label: 'P2P', value: 'testnet.xmr.ditatompel.com:28080', key: 'tnetP2P' },
			{ label: 'RPC', value: 'testnet.xmr.ditatompel.com:28089', key: 'tnetRPC' },
			{ label: 'RPC SSL', value: 'testnet.xmr.ditatompel.com:443', key: 'tnetSSL' }
		],
		donation: {
			address:
				'8BWYe6GzbNKbxe3D8mPkfFMQA2rViaZJFhWShhZTjJCNG6EZHkXRZCKHiuKmwwe4DXDYF8KKcbGkvNYaiRG3sNt7JhnVp7D',
			qr: '/img/monerotip.png'
		}
	};
}
