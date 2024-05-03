import { apiUri } from '$lib/utils/common';

/** @param {import('@vincjo/datatables/remote/state')} state */
export const loadData = async (state) => {
	const response = await fetch(apiUri(`/api/v1/crons?${getParams(state)}`));
	const json = await response.json();
	state.setTotalRows(json.data.length ?? 0);
	return json.data ?? [];
};

const getParams = ({ pageNumber, rowsPerPage, sort, filters }) => {
	let params = `page=${pageNumber}&limit=${rowsPerPage}`;

	if (sort) {
		params += `&orderBy=${sort.orderBy}&orderDir=${sort.direction}`;
	}
	if (filters) {
		params += filters.map(({ filterBy, value }) => `&${filterBy}=${value}`).join('');
	}
	return params;
};
