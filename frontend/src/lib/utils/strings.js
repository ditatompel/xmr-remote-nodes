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
 * @param {number} n
 * @param {number} p
 */
const maxPrecision = (n, p) => {
	return parseFloat(n.toFixed(p));
};

/**
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
