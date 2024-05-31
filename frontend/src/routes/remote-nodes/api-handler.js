import { apiUri } from '$lib/utils/common';

/**
 * @typedef {import('@vincjo/datatables/remote').State} State
 * @param {State} state - The state object from the data table.
 */
export const loadData = async (state) => {
	const response = await fetch(apiUri(`/api/v1/nodes?${getParams(state)}`));
	const json = await response.json();
	state.setTotalRows(json.data.total_rows ?? 0);
	return json.data.items ?? [];
};

export const loadCountries = async () => {
	const response = await fetch(apiUri('/api/v1/countries'));
	const json = await response.json();
	return json.data ?? [];
};

export const loadFees = async () => {
	const response = await fetch(apiUri('/api/v1/fees'));
	const json = await response.json();
	return json.data ?? [];
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
