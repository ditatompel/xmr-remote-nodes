import { PUBLIC_API_ENDPOINT } from '$env/static/public';

/** @param {import('@vincjo/datatables/remote/state')} state */
export async function loadApiData(state) {
	const response = await fetch(`${PUBLIC_API_ENDPOINT}/monero/remote-node-dt?${getParams(state)}`);
	const json = await response.json();

	state.setTotalRows(json.data.total ?? 0);

	return json.data.nodes ?? [];
}

const getParams = ({ pageNumber, offset, rowsPerPage, sort, filters }) => {
	let params = `page=${pageNumber}&limit=${rowsPerPage}`;
	if (sort) {
		params += `&sort=${sort.orderBy}&dir=${sort.direction}`;
	}
	if (filters) {
		params += filters.map(({ filterBy, value }) => `&${filterBy}=${value}`).join('');
	}
	return params;
};
