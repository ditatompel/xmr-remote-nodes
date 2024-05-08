import { apiUri } from '$lib/utils/common';
import { goto } from '$app/navigation';

/** @param {import('@vincjo/datatables/remote/state')} state */
export const loadData = async (state) => {
	const response = await fetch(apiUri(`/api/v1/crons?${getParams(state)}`));
	const json = await response.json();
	if (json.data === null) {
		goto('/login');
		return;
	}
	state.setTotalRows(json.data.length ?? 0);
	return json.data.items ?? [];
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
