// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
	interface ImportMetaEnv {
		VITE_API_URL: string;
	}

	interface MoneroNode {
		id: number;
		hostname: string;
		ip: string;
		port: number;
		protocol: string;
		is_tor: boolean;
		is_available: boolean;
		nettype: string;
		ip_addresses: string;
	}

	interface ApiResponse {
		status: string;
		message: string;
		data: null | object | object[];
	}
}

export {};
