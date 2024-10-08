/**
 * Modifies the input string based on whether it is an IPv6 address.
 * If the input is an IPv6 address, it wraps it in square brackets `[ ]`.
 * Otherwise, it returns the input string as-is (for domain names or
 * IPv4 addresses). AND I'M SORRY USING REGEX FOR THIS!
 *
 * @param {string} hostname
 * @returns {string} - The modified string, IPv6 addresses wrapped in `[ ]`.
 */
export const formatHostname = (hostname) => {
	// const ipv6Pattern = /^(?:[a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}$/; // full
	// pattern for both full and compressed IPv6 addresses.
	// source: https://regex101.com/library/cP9mH9?filterFlavors=dotnet&filterFlavors=javascript&orderBy=RELEVANCE&search=ip
	// This may be incorrect, but let's assume it's correct. xD
	const ipv6Pattern =
		/^(([0-9A-Fa-f]{1,4}:){7})([0-9A-Fa-f]{1,4})$|(([0-9A-Fa-f]{1,4}:){1,6}:)(([0-9A-Fa-f]{1,4}:){0,4})([0-9A-Fa-f]{1,4})$/;
	if (ipv6Pattern.test(hostname)) {
		return `[${hostname}]`;
	}

	return hostname;
};

/**
 * @param {number} bytes
 * @param {number} decimals
 * @returns {string}
 */
export const formatBytes = (bytes, decimals = 2) => {
	if (!+bytes) return '0 Bytes';

	const k = 1024;
	const dm = decimals < 0 ? 0 : decimals;
	const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];

	const i = Math.floor(Math.log(bytes) / Math.log(k));

	return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
};

/**
 * Returns a number with a maximum precision.
 *
 * This function was copied from jtgrassie/monero-pool project.
 * Source: https://github.com/jtgrassie/monero-pool/blob/master/src/webui-embed.html
 *
 * Copyright (c) 2018, The Monero Project
 *
 * @param {number} n
 * @param {number} p
 */
const maxPrecision = (n, p) => {
	return parseFloat(n.toFixed(p));
};

/**
 * Formats a hash value (h) into human readable format.
 *
 * This function was copied from jtgrassie/monero-pool project.
 * Source: https://github.com/jtgrassie/monero-pool/blob/master/src/webui-embed.html
 *
 * Copyright (c) 2018, The Monero Project
 *
 * @param {number} h
 */
export const formatHashes = (h) => {
	if (h < 1e-12) return '0 H';
	else if (h < 1e-9) return maxPrecision(h * 1e12, 0) + ' pH';
	else if (h < 1e-6) return maxPrecision(h * 1e9, 0) + ' nH';
	else if (h < 1e-3) return maxPrecision(h * 1e6, 0) + ' μH';
	else if (h < 1) return maxPrecision(h * 1e3, 0) + ' mH';
	else if (h < 1e3) return h + ' H';
	else if (h < 1e6) return maxPrecision(h * 1e-3, 2) + ' KH';
	else if (h < 1e9) return maxPrecision(h * 1e-6, 2) + ' MH';
	else return maxPrecision(h * 1e-9, 2) + ' GH';
};
