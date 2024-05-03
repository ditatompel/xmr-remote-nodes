/** @param {string} path */
export const apiUri = (path) => {
	return `${import.meta.env.VITE_API_URL || ''}${path}`;
};
