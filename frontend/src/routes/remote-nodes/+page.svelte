<script>
	import { DataHandler } from '@vincjo/datatables/remote';
	import { format, formatDistance } from 'date-fns';
	import { loadData, loadFees, loadCountries } from './api-handler';
	import { onMount } from 'svelte';
	import {
		DtSrRowsPerPage,
		DtSrThSort,
		DtSrThFilter,
		DtSrRowCount,
		DtSrPagination,
		DtSrAutoRefresh
	} from '$lib/components/datatables/server';
	import {
		HostPortCell,
		NetTypeCell,
		ProtocolCell,
		CountryCellWithAsn,
		StatusCell,
		UptimeCell,
		EstimateFeeCell
	} from '$lib/components/datatables/xmr';
	import News from '$lib/components/News.svelte';

	export let data;
	let filterNettype = 'any';
	let filterProtocol = 'any';
	let filterCc = 'any';
	let filterStatus = -1;
	let checkboxCors = false;

	/** @type {{total_nodes: number, cc: string, name: string}[]} */
	let countries = [];
	let fees = [];

	const handler = new DataHandler([], { rowsPerPage: 10, totalRows: 0 });
	let rows = handler.getRows();

	/** @type {Object.<string, number>} */
	let majorityFee;

	onMount(() => {
		loadFees().then((data) => {
			fees = data;
			majorityFee = fees.reduce(
				/**
				 * @param {Object.<string, number>} o
				 * @param {{ nettype: string, estimate_fee: number }} key
				 * @returns {Object.<string, number>}
				 */
				(o, key) => ({
					...o,
					[key.nettype]: key.estimate_fee
				}),
				{}
			);
			handler.onChange((state) => loadData(state));
			handler.invalidate();
		});
		loadCountries().then((data) => {
			countries = data;
		});
	});
</script>

<header id="hero" class="hero-gradient py-7">
	<div class="section-container text-center">
		<h1 class="h1 pb-2 font-extrabold">{data.meta.title}</h1>
		<!-- prettier-ignore -->
		<p class="mx-auto max-w-3xl">
			<strong>Monero remote node</strong> is a device on the internet running the Monero software with full copy of the Monero blockchain that doesn't run on the same local machine where the Monero wallet is located.
		</p>
	</div>
	<div class="mx-auto w-full max-w-3xl px-20">
		<hr class="!border-primary-400-500-token !border-t-4 !border-double" />
	</div>
</header>

<!-- prettier-ignore -->
<section id="introduction">
	<div class="section-container text-center !max-w-4xl">
		<p>Remote node can be used by people who, for their own reasons (usually because of hardware requirements, disk space, or technical abilities), cannot/don't want to run their own node and prefer to relay on one publicly available on the Monero network.</p>
		<p>Using an open node will allow to make a transaction instantaneously, without the need to download the blockchain and sync to the Monero network first, but at the cost of the control over your privacy. the <strong>Monero community suggests to <span class="font-extrabold text-2xl underline decoration-double decoration-2 decoration-pink-500">always run and use your own node</span></strong> to obtain the maximum possible privacy and to help decentralize the network.</p>
	</div>
</section>

<News />

<section id="monero-remote-node">
	<div class="section-container">
		<div class="space-y-2 overflow-x-auto">
			<div class="flex justify-between">
				<DtSrRowsPerPage {handler} />
				<div class="invisible flex place-items-center md:visible">
					<DtSrAutoRefresh {handler} />
				</div>
				<div class="flex place-items-center">
					<button
						id="reloadDt"
						name="reloadDt"
						class="variant-filled-primary btn"
						on:click={() => handler.invalidate()}>Reload</button
					>
				</div>
			</div>

			<table class="table table-hover table-compact w-full table-auto">
				<thead>
					<tr>
						<th>Host:Port</th>
						<th>Nettype</th>
						<th>Protocol</th>
						<th>Country</th>
						<th>Status</th>
						<th>Est. Fee</th>
						<DtSrThSort {handler} orderBy="uptime">Uptime</DtSrThSort>
						<DtSrThSort {handler} orderBy="last_checked">Check</DtSrThSort>
					</tr>
					<tr>
						<DtSrThFilter {handler} filterBy="host" placeholder="Filter Host / IP" />
						<th>
							<select
								id="nettype"
								name="nettype"
								class="select variant-form-material"
								bind:value={filterNettype}
								on:change={() => {
									handler.filter(filterNettype, 'nettype');
									handler.invalidate();
								}}
							>
								<option value="any">Any</option>
								<option value="mainnet">MAINNET</option>
								<option value="stagenet">STAGENET</option>
								<option value="testnet">TESTNET</option>
							</select>
						</th>
						<th>
							<select
								id="protocol"
								name="protocol"
								class="select variant-form-material"
								bind:value={filterProtocol}
								on:change={() => {
									handler.filter(filterProtocol, 'protocol');
									handler.invalidate();
								}}
							>
								<option value="any">Any</option>
								<option value="tor">TOR</option>
								<option value="http">HTTP</option>
								<option value="https">HTTPS</option>
							</select>
						</th>
						<th>
							<select
								id="cc"
								name="cc"
								class="select variant-form-material"
								bind:value={filterCc}
								on:change={() => {
									handler.filter(filterCc, 'cc');
									handler.invalidate();
								}}
							>
								<option value="any">Any</option>
								{#each countries as country}
									{#if country.cc === ''}
										<option value="UNKNOWN">UNKNOWN ({country.total_nodes})</option>
									{:else}
										<option value={country.cc}
											>{country.name === '' ? country.cc : country.name} ({country.total_nodes})</option
										>
									{/if}
								{/each}
							</select>
						</th>
						<th colspan="2">
							<select
								id="status"
								name="status"
								class="select variant-form-material"
								bind:value={filterStatus}
								on:change={() => {
									handler.filter(filterStatus, 'status');
									handler.invalidate();
								}}
							>
								<option value={-1}>Any</option>
								<option value="0">Offline</option>
								<option value="1">Online</option>
							</select>
						</th>
						<th colspan="2">
							<label for="cors" class="flex items-center justify-center space-x-2">
								<input
									id="cors"
									name="cors"
									class="checkbox"
									type="checkbox"
									bind:checked={checkboxCors}
									on:change={() => {
										handler.filter(checkboxCors === true ? 1 : -1, 'cors');
										handler.invalidate();
									}}
								/>
								<p>CORS</p>
							</label>
						</th>
					</tr>
				</thead>
				<tbody>
					{#each $rows as row (row.id)}
						<tr>
							<td
								><HostPortCell
									ip_addresses={row.ip_addresses}
									is_tor={row.is_tor}
									hostname={row.hostname}
									port={row.port}
									ipv6_only={row.ipv6_only}
								/>
							</td>
							<td><NetTypeCell nettype={row.nettype} height={row.height} /></td>
							<td><ProtocolCell protocol={row.protocol} cors={row.cors} /></td>
							<td
								><CountryCellWithAsn
									cc={row.cc}
									country_name={row.country_name}
									city={row.city}
									asn={row.asn}
									asn_name={row.asn_name}
								/></td
							>
							<td
								><StatusCell
									is_available={row.is_available}
									statuses={row.last_check_statuses}
								/></td
							>
							<td>
								<EstimateFeeCell
									estimate_fee={row.estimate_fee}
									majority_fee={majorityFee[row.nettype]}
								/>
							</td>
							<td
								><UptimeCell uptime={row.uptime} /><br />
								<a
									class="anchor !text-purple-800 dark:!text-purple-400"
									href="/remote-nodes/logs/?node_id={row.id}">[Logs]</a
								>
							</td>
							<td>
								{format(row.last_checked * 1000, 'PP HH:mm')}<br />
								{formatDistance(row.last_checked * 1000, new Date(), { addSuffix: true })}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>

			<div class="flex justify-between mb-2">
				<DtSrRowCount {handler} />
				<DtSrPagination {handler} />
			</div>
		</div>
	</div>
</section>

<section id="page-info" class="mx-auto w-full max-w-4xl px-4 pb-7">
	<div class="alert card shadow-xl">
		<div class="alert-message">
			<h2 class="h3">Info</h2>
			<ul class="list-inside list-disc">
				<li>
					If you find any remote nodes that are strange or suspicious, please <a
						class="external"
						href="https://github.com/ditatompel/xmr-remote-nodes/issues"
						target="_blank"
						rel="noopener">open an issue on GitHub</a
					> for removal.
				</li>
				<li>
					Uptime percentage calculated is the <strong>last 1 month</strong> uptime.
				</li>
				<li>
					<strong>Est. Fee</strong> here is just fee estimation / byte from
					<code class="code text-rose-900 font-bold">get_fee_estimate</code> RPC call method.
				</li>
				<li>
					Malicious actors who running remote nodes <a
						class="link"
						href="/img/node-tx-fee.jpg"
						rel="noopener">still can return high fee only if you about to create a transactions</a
					>.
				</li>
				<li>
					<strong
						class="font-extrabold text-2xl underline decoration-double decoration-2 decoration-pink-500"
						>The best and safest way is running your own node!</strong
					>
				</li>
				<li>
					Nodes with 0% uptime within 1 month with more than 300 check attempt will be removed. You
					can always add your node again latter.
				</li>
				<li>
					You can filter remote node by selecting on <strong>nettype</strong>,
					<strong>protocol</strong>, <strong>country</strong>,
					<strong>tor</strong>, and <strong>online status</strong> option.
				</li>
				<li>
					If you know one or more remote node that we don't currently monitor, please add them using <a
						href="/add-node">this form</a
					>.
				</li>
				<li>
					I deliberately cut the long Tor addresses, click the <span
						class="text-orange-800 dark:text-orange-300">👁 torhostname...</span
					> to see the full Tor address.
				</li>
				<li>
					You can found larger remote nodes database from <a
						class="external"
						href="https://monero.fail/"
						role="button"
						target="_blank"
						rel="noopener">monero.fail</a
					>.
				</li>
				<li>
					If you are developer or power user who like to fetch Monero remote node above in JSON
					format, you can read <a
						class="external"
						href="https://insights.ditatompel.com/en/blog/2022/01/public-api-monero-remote-node-list/"
						>Public API Monero Remote Node List</a
					> blog post for more detailed information.
				</li>
			</ul>
		</div>
	</div>
</section>
