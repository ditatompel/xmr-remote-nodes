import { apiUri } from '$lib/utils/common';

/**
 * @typedef {import('@vincjo/datatables/remote').State} State
 * @param {State} state - The state object from the data table.
 */
export const loadData = async (state) => {
	const response = await fetch(apiUri(`/api/v1/nodes/logs?${getParams(state)}`));
	const json = await response.json();
	state.setTotalRows(json.data.total_rows ?? 0);
	return json.data.items ?? [];
};

/** @param {string} nodeId */
export const loadNodeInfo = async (nodeId) => {
	const response = await fetch(apiUri(`/api/v1/nodes/id/${nodeId}`));
	const json = await response.json();
	return json.data;
};

/** @param {State} state - The state object from the data table. */
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
