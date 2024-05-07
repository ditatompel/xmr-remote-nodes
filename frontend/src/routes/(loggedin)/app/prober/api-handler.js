import { apiUri } from '$lib/utils/common';

/** @param {import('@vincjo/datatables/remote/state')} state */
export const loadData = async (state) => {
	const response = await fetch(apiUri(`/api/v1/prober?${getParams(state)}`));
	const json = await response.json();
	state.setTotalRows(json.data.total_rows ?? 0);
	return json.data.items ?? [];
};

export const deleteData = async (id) => {
	const response = await fetch(apiUri(`/api/v1/prober/${id}`), {
		method: 'DELETE',
		credentials: 'include'
	});
	const json = await response.json();
	return json;
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
