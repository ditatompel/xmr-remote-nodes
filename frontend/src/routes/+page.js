/** @type {import('./$types').PageLoad} */
export async function load() {
	return {
		meta: {
			title: 'Monero Remote Node',
			description:
				'A website that helps you monitor your favourite Monero remote nodes, a device on the internet running the Monero software with copy of the Monero blockchain.',
			keywords:
				'monero,monero,xmr,monero node,xmrnode,cryptocurrency,monero remote node,monero testnet,monero stagenet'
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
			// You change donation address and qr image below if you run your own "instance"
			address:
				'8BWYe6GzbNKbxe3D8mPkfFMQA2rViaZJFhWShhZTjJCNG6EZHkXRZCKHiuKmwwe4DXDYF8KKcbGkvNYaiRG3sNt7JhnVp7D',
			qr: '/img/monerotip.png'
		}
	};
}
