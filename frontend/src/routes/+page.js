/** @type {import('./$types').PageLoad} */
export async function load() {
	return {
		meta: {
			title: 'Monero Remote Node',
			description:
				'A website that helps you monitor your favourite Monero remote nodes, but YOU BETTER RUN AND USE YOUR OWN NODE.',
			keywords:
				'monero,monero,xmr,monero node,xmrnode,cryptocurrency,monero remote node,monero testnet,monero stagenet'
		},
		links: [
			{ text: 'moneroworld.com', uri: 'https://moneroworld.com' },
			{ text: 'monero.how', uri: 'https://www.monero.how' },
			{ text: 'monero.observer', uri: 'https://www.monero.observer' },
			{ text: 'revuo-xmr.com', uri: 'https://revuo-xmr.com' },
			{ text: 'themonoeromoon.com', uri: 'https://www.themoneromoon.com' },
			{ text: 'monerotopia.com', uri: 'https://monerotopia.com' },
			{ text: 'sethforprivacy.com', uri: 'https://sethforprivacy.com' }
		],
		stagenet: [
			{ label: 'P2P', value: 'stagenet.xmr.ditatompel.com:38080', key: 'snetP2P' },
			{ label: 'RPC', value: 'stagenet.xmr.ditatompel.com:38089', key: 'snetRPC' },
			{ label: 'RPC SSL', value: 'stagenet.xmr.ditatompel.com:443', key: 'snetSSL' }
		],
		testnet: [
			{ label: 'P2P', value: 'testnet.xmr.ditatompel.com:28080', key: 'tnetP2P' },
			{ label: 'RPC', value: 'testnet.xmr.ditatompel.com:28089', key: 'tnetRPC' },
			{ label: 'RPC SSL', value: 'testnet.xmr.ditatompel.com:443', key: 'tnetSSL' }
		],
		donation: {
			// You change donation address and qr image below if you run your own "instance"
			address:
				'8BWYe6GzbNKbxe3D8mPkfFMQA2rViaZJFhWShhZTjJCNG6EZHkXRZCKHiuKmwwe4DXDYF8KKcbGkvNYaiRG3sNt7JhnVp7D',
			qr: '/img/monerotip.png'
		}
	};
}
