export const getDistinct = (items) => {
	return Array.from(getCounter(items).keys());
};

export const getDuplicates = (items) => {
	return Array.from(getCounter(items).entries())
		.filter(([, count]) => count !== 1)
		.map(([key]) => key);
};

export const getCounter = (items) => {
	const result = new Map();
	items.forEach((item) => {
		result.set(item, (result.get(item) ?? 0) + 1);
	});
	return result;
};
