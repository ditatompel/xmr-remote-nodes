import { apiUri } from '$lib/utils/common';

/** @param {import('@vincjo/datatables/remote/state')} state */
export const loadData = async (state) => {
	const response = await fetch(apiUri(`/api/v1/nodes/logs?${getParams(state)}`));
	const json = await response.json();
	state.setTotalRows(json.data.total_rows ?? 0);
	return json.data.items ?? [];
};

export const loadNodeInfo = async (nodeId) => {
	const response = await fetch(apiUri(`/api/v1/nodes/id/${nodeId}`));
	const json = await response.json();
	return json.data;
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

const getParams = ({ pageNumber, rowsPerPage, sort, filters }) => {
	let params = `page=${pageNumber}&limit=${rowsPerPage}`;

	if (sort) {
		params += `&sort_by=${sort.orderBy}&sort_direction=${sort.direction}`;
	}
	if (filters) {
		params += filters.map(({ filterBy, value }) => `&${filterBy}=${value}`).join('');
	}
	return params;
};
