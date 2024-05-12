<!-- Developed by Taipei Urban Intelligence Center 2023-2024-->

<script setup>
import { ref, defineProps, watch, onBeforeMount } from "vue";
import { storeToRefs } from "pinia";
import { DashboardComponent } from "city-dashboard-component";
import { useDialogStore } from "../../../store/dialogStore";
import { useAdminStore } from "../../../store/adminStore";
import http from "../../../router/axios";

import DialogContainer from "../DialogContainer.vue";
import InputTags from "../../utilities/forms/InputTags.vue";
import SelectButtons from "../../utilities/forms/SelectButtons.vue";
import HistoryChart from "../../charts/HistoryChart.vue";

import { chartsPerDataType } from "../../../assets/configs/apexcharts/chartTypes";
import { timeTerms } from "../../../assets/configs/AllTimes";
import { mapTypes } from "../../../assets/configs/mapbox/mapConfig";

const dialogStore = useDialogStore();
const adminStore = useAdminStore();

const props = defineProps(["searchParams"]);

const currentSettings = ref("sql");
const tempInputStorage = ref({
	link: "",
	contributor: "",
	chartColor: "#000000",
	historyColor: "#000000",
});

const maindata = ref({
	name: "",
	index: "",
	isUseMap: false,
	alllist: [],
	csvStore: [],
	afterSql: [],
});
function handleConfirm() {
	newInputStorage.value.map.index = maindata.value.index;
	newInputStorage.value.map.title = maindata.value.name;
	newInputStorage.value.chart.index = maindata.value.index;
	newInputStorage.value.data.index = maindata.value.index;
	newInputStorage.value.data.name = maindata.value.name;

	handleTestNewComponents();
	handleClose();
}

function handleClose() {
	currentSettings.value = "sql";
	dialogStore.hideAllDialogs();
}
function handleSelectOthers(input) {
	currentSettings.value = input;
}

const newInputStorage = ref({
	map: {
		index: "",
		title: "",
		type: "circle",
		source: "geojson",
		size: "big",
		icon: "",
		paint: "",
		property: "",
	},
	chart: {
		index: "",
		color: ["#E6DF44", "#F4633C", "#D63940", "#9C2A4B"],
		types: ["ColumnChart", "DistrictChart", "RadarChart"],
		unit: "處",
	},
	data: {
		index: "garbage_truck",
		name: "垃圾車收運",
		history_config: null,
		map_config: "",
		map_config_ids: "",
		chart_config: "",
		map_filter: "",
		time_from: "static",
		time_to: "",
		update_freq: 0,
		update_freq_unit: "",
		source: "",
		short_desc: "",
		long_desc: "",
		use_case: "",
		links: [
			"https://data.taipei/dataset/detail?id=6bb3304b-4f46-4bb0-8cd1-60c66dcd1cae",
		],
		contributors: ["管誰import啥"],
		query_type: "three_d",
		query_chart: "",
	},
});

// Trim the input query on the frontend side
function trimQuery(query) {
	// Remove leading and trailing whitespace
	query = query.trim();
	// Replace multiple consecutive whitespace characters with a single space
	query = query.replace(/\s+/g, " ");
	return query;
}

async function handleGetCsv() {
	const datasetIndex = maindata.value.index;
	const csvData = await http.get(`/helper/data-csv/${datasetIndex}`);
	maindata.value.csvStore = csvData.data.data;
}

async function handlePushSqlData() {
	// Validate query chart and query type
	const queryType = newInputStorage.value.data.query_type,
		queryChart = newInputStorage.value.data.query_chart;

	if (!queryType || !queryChart) {
		console.error("Query chart and query type cannot be empty");
		return;
	}

	// Trim query chart
	const trimmedQuery = trimQuery(queryChart);

	// Prepare request body
	const requestBody = {
		queryType: queryType,
		queryString: trimmedQuery,
	};

	try {
		// Make HTTP request
		const response = await http.post(
			"/helper/query",
			JSON.stringify(JSON.parse(JSON.stringify(requestBody)))
		);
		console.log(response);
		// Show notification on success
		dialogStore.showNotification("success", response.data);

		// Update newInputStorage chart field if necessary
		// (Assuming there's a method to update the chart field)
		maindata.value.afterSql = response.data;
	} catch (error) {
		// Handle errors
		console.error("Error:", error);
		dialogStore.showNotification(
			"error",
			"An error occurred while processing the request."
		);
	}
}

function handleTestNewComponents() {
	// no dialog or UI for this testing function
	adminStore.newComponentMap = JSON.parse(
		JSON.stringify(newInputStorage.value.map)
	);
	adminStore.newComponentChart = JSON.parse(
		JSON.stringify(newInputStorage.value.chart)
	);
	adminStore.newComponent = JSON.parse(
		JSON.stringify(newInputStorage.value.data)
	);

	adminStore.createComponent(props.searchParams.value,maindata.value.isUseMap);
}

watch(
	() => dialogStore.dialogs.adminCreator,
	async (newValue, oldValue) => {
		if (newValue === true) {
			const response = await http.get("/helper/list-tables");
			maindata.value.alllist = response.data.tables;
		}
	}
);
</script>

<template>
	<DialogContainer :dialog="`adminCreator`" @on-close="handleClose">
		<div class="adminCreator">
			<div class="adminCreator-header">
				<h2>組件設定</h2>
				<button
					v-if="maindata.afterSql.length != 0"
					@click="handleConfirm"
				>
					確定新增
				</button>
			</div>
			<div class="adminCreator-tabs">
				<button
					:class="{ active: currentSettings === 'sql' }"
					@click="handleSelectOthers('sql')"
				>
					基設
				</button>
				<button
					v-if="maindata.afterSql.length != 0"
					:class="{ active: currentSettings === 'all' }"
					@click="handleSelectOthers('all')"
				>
					整體
				</button>
				<button
					v-if="maindata.afterSql.length != 0"
					:class="{ active: currentSettings === 'chart' }"
					@click="handleSelectOthers('chart')"
				>
					圖表
				</button>
			</div>
			<div class="adminCreator-content">
				<div class="adminCreator-settings">
					<div
						v-if="currentSettings === 'sql'"
						class="adminCreator-settings-items"
					>
						<label>
							組件名稱* ({{ maindata.name.length }}/10)
						</label>
						<input
							v-model="maindata.name"
							type="text"
							:minlength="1"
							:maxlength="15"
							required
						/>
						<label>組件 Index</label>

						<select
							v-model="maindata.index"
							@change="handleGetCsv()"
						>
							<option
								:value="item"
								v-for="item in maindata.alllist"
								:key="item"
							>
								{{ item }}
							</option>
						</select>

						<label>圖表資料型態 Query Type</label>
						<select v-model="newInputStorage.data.query_type">
							<option value="two_d">二維資料</option>
							<option value="three_d">三維資料</option>
							<option value="time">時間序列資料</option>
							<option value="percent">百分比資料</option>
							<option value="map_legend">圖例資料</option>
						</select>
						<label>sql設定</label>
						<textarea v-model="newInputStorage.data.query_chart" />
						<div class="adminCreator-runsql">
							<button @click="handlePushSqlData()">
								匯入資料
							</button>
						</div>
					</div>

					<div
						v-else-if="currentSettings === 'all'"
						class="adminCreator-settings-items"
					>
						<label>資料來源*</label>
						<input
							v-model="newInputStorage.data.source"
							type="text"
							:minlength="1"
							:maxlength="12"
							required
						/>
						<label>更新頻率* (0 = 不定期更新)</label>
						<div class="two-block">
							<input
								v-model="newInputStorage.data.update_freq"
								type="number"
								:min="0"
								:max="31"
								required
							/>
							<select
								v-model="newInputStorage.data.update_freq_unit"
							>
								<option value="minute" />
								<option value="hour">時</option>
								<option value="day">天</option>
								<option value="week">週</option>
								<option value="month">月</option>
								<option value="year">年</option>
							</select>
						</div>
						<label>資料區間</label>
						<!-- eslint-disable no-mixed-spaces-and-tabs -->
						<input
							:value="`${
								timeTerms[newInputStorage.data.time_from]
							}${
								timeTerms[newInputStorage.data.time_to]
									? ' ~ ' +
									  timeTerms[newInputStorage.data.time_to]
									: ''
							}`"
							disabled
						/>
						<label required
							>組件簡述* ({{
								newInputStorage.data.short_desc.length
							}}/50)</label
						>
						<textarea
							v-model="newInputStorage.data.short_desc"
							:minlength="1"
							:maxlength="50"
							required
						/>
						
						
						<div  class="adminCreator-settings-checkbox">
						<input type="checkbox" v-model="maindata.isUseMap">

							<label>
							載入地圖
						</label>
						</div>
						<label
							>組件詳述* ({{
								newInputStorage.data.long_desc.length
							}}/100)</label
						>
						<textarea
							v-model="newInputStorage.data.long_desc"
							:minlength="1"
							:maxlength="100"
							required
						/>
						<label
							>範例情境* ({{
								newInputStorage.data.use_case.length
							}}/100)</label
						>
						<textarea
							v-model="newInputStorage.data.use_case"
							:minlength="1"
							:maxlength="100"
							required
						/>
						<label>資料連結</label>
						<InputTags
							:tags="newInputStorage.data.links"
							@deletetag="
								(index) => {
									newInputStorage.data.links.splice(index, 1);
								}
							"
							@updatetagorder="
								(updatedTags) => {
									newInputStorage.data.links = updatedTags;
								}
							"
						/>
						<input
							v-model="tempInputStorage.link"
							type="text"
							:minlength="1"
							@keypress.enter="
								() => {
									if (tempInputStorage.link.length > 0) {
										newInputStorage.data.links.push(
											tempInputStorage.link
										);
										tempInputStorage.link = '';
									}
								}
							"
						/>
						<label>貢獻者</label>
						<InputTags
							:tags="newInputStorage.data.contributors"
							@deletetag="
								(index) => {
									newInputStorage.data.contributors.splice(
										index,
										1
									);
								}
							"
							@updatetagorder="
								(updatedTags) => {
									newInputStorage.data.contributors =
										updatedTags;
								}
							"
						/>
						<input
							v-model="tempInputStorage.contributor"
							type="text"
							@keypress.enter="
								() => {
									if (
										tempInputStorage.contributor.length > 0
									) {
										newInputStorage.data.contributors.push(
											tempInputStorage.contributor
										);
										tempInputStorage.contributor = '';
									}
								}
							"
						/>
					</div>

					<div
						v-else-if="currentSettings === 'chart'"
						class="adminCreator-settings-items"
					>
						<label>圖表資料型態</label>
						<select
							:value="newInputStorage.data.query_type"
							disabled
						>
							<option value="two_d">二維資料</option>
							<option value="three_d">三維資料</option>
							<option value="time">時間序列資料</option>
							<option value="percent">百分比資料</option>
							<option value="map_legend">圖例資料</option>
						</select>
						<label>資料單位*</label>
						<input
							v-model="newInputStorage.chart.unit"
							type="text"
							:minlength="1"
							:maxlength="6"
							required
						/>
						<label>圖表類型*（限3種，依點擊順序排列）</label>
						<SelectButtons
							:tags="
								chartsPerDataType[
									newInputStorage.data.query_type
								]
							"
							:selected="newInputStorage.chart.types"
							:limit="3"
							@updatetagorder="
								(updatedTags) => {
									newInputStorage.chart.types = updatedTags;
								}
							"
						/>
						<label>圖表顏色</label>
						<InputTags
							:tags="newInputStorage.chart.color"
							:color-data="true"
							@deletetag="
								(index) => {
									newInputStorage.chart.color.splice(
										index,
										1
									);
								}
							"
							@updatetagorder="
								(updatedTags) => {
									newInputStorage.chart.color = updatedTags;
								}
							"
						/>
						<input
							v-model="tempInputStorage.chartColor"
							type="color"
							class="adminCreator-settings-inputcolor"
							@focusout="
								() => {
									if (
										tempInputStorage.chartColor.length === 7
									) {
										newInputStorage.chart.color.push(
											tempInputStorage.chartColor
										);
										tempInputStorage.chartColor = '#000000';
									}
								}
							"
						/>
						<div v-if="newInputStorage.data.map_config[0] !== null">
							<label>地圖篩選</label>
							<textarea
								v-model="newInputStorage.data.map_filter"
							/>
						</div>
					</div>
				</div>
				<div class="adminCreator-preview">
					<div
						class="adminCreator__preData"
						v-if="currentSettings === 'all'"
					>
						<div class="adminCreator__preData-box">
							<pre><code>{{maindata.afterSql}}</code></pre>
						</div>
					</div>
					<div
						v-else-if="currentSettings === 'sql'"
						class="adminCreator__preData"
					>
						<p v-if="maindata.csvStore.length == 0">尚未匯入資料</p>
						<div
							class="adminCreator__preData-box"
							v-if="maindata.csvStore.length != 0"
						>
							<div
								class="adminCreator__preData-title"
								:style="`width:${
									Object.keys(maindata.csvStore[0]).length *
									40
								}%`"
							>
								<div
									v-for="item in Object.keys(
										maindata.csvStore[0]
									)"
									:style="`width:${Math.floor(
										100 /
											Object.keys(maindata.csvStore[0])
												.length
									)}%`"
								>
									{{ item }}
								</div>
							</div>
							<div
								class="adminCreator__preData-item"
								v-for="item in Object.keys(maindata.csvStore)
									.length > 10
									? 10
									: Object.keys(maindata.csvStore)"
								:style="`width:${
									Object.keys(maindata.csvStore[0]).length *
									40
								}%`"
							>
								<div
									v-for="item in maindata.csvStore[item]"
									:style="`width:${Math.floor(
										100 /
											Object.keys(maindata.csvStore[0])
												.length
									)}%`"
								>
									{{ item }}
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</DialogContainer>
</template>

<style scoped lang="scss">
.adminCreator {
	width: 750px;
	height: 500px;

	@media (max-width: 770px) {
		display: none;
	}
	@media (max-height: 520px) {
		display: none;
	}

	&-header {
		display: flex;
		justify-content: space-between;

		button {
			display: flex;
			align-items: center;
			justify-self: baseline;
			padding: 2px 4px;
			border-radius: 5px;
			background-color: var(--color-highlight);
			font-size: var(--font-ms);
		}
	}

	&-content {
		height: calc(100% - 70px);
		display: grid;
		grid-template-columns: 1fr 350px;
	}

	&-tabs {
		height: 30px;
		display: flex;
		align-items: center;
		margin-top: var(--font-s);

		button {
			width: 70px;
			height: 30px;
			border-radius: 5px 5px 0px 0px;
			background-color: var(--color-border);
			font-size: var(--font-m);
			color: var(--color-text);
			cursor: pointer;
			transition: background-color 0.2s;

			&:hover {
				background-color: var(--color-complement-text);
			}
		}
		.active {
			background-color: var(--color-complement-text);
		}
	}
	&-runsql {
		display: flex;
		justify-content: flex-end;
		margin-top: var(--font-m);
		button {
			background-color: var(--color-highlight);
			font-size: var(--font-ms);
			padding: 2px 8px;
			border-radius: 2px;
			text-align: end;
		}
	}
	&-settings {
		padding: 0 0.5rem 0.5rem 0.5rem;
		margin-right: var(--font-ms);
		border-radius: 0px 5px 5px 5px;
		border: solid 1px var(--color-border);
		overflow-y: scroll;

		label {
			margin: 8px 0 4px;
			font-size: var(--font-s);
			color: var(--color-complement-text);
		}

		.two-block {
			display: grid;
			grid-template-columns: 1fr 1fr;
			column-gap: 0.5rem;
			select {
				width: 100%;
			}
		}
		.three-block {
			display: grid;
			grid-template-columns: 1fr 2rem 1fr;
			column-gap: 0.5rem;
		}
		&-checkbox{
			display: flex;
			input{
				margin-right: 4px;
			}
		}
		&-items {
			display: flex;
			flex-direction: column;

			hr {
				margin: var(--font-ms) 0 0.5rem;
				border: none;
				border-bottom: dashed 1px var(--color-complement-text);
			}
		}

		&-inputcolor {
			width: 140px;
			height: 40px;
			appearance: none;
			display: flex;
			justify-content: center;
			align-items: center;
			padding: 0;
			outline: none;
			cursor: pointer;

			&::-webkit-color-swatch {
				border: none;
				border-radius: 5px;
			}
			&::-moz-color-swatch {
				border: none;
			}
			&:before {
				content: "選擇顏色";
				position: absolute;
				display: block;
				border-radius: 5px;
				font-size: var(--font-ms);
				color: var(--color-complement-text);
			}
			&:focus:before {
				content: "點擊空白處確認";
				text-shadow: 0px 0px 1px black;
			}
		}

		&::-webkit-scrollbar {
			width: 4px;
		}
		&::-webkit-scrollbar-thumb {
			background-color: rgba(136, 135, 135, 0.5);
			border-radius: 4px;
		}
		&::-webkit-scrollbar-thumb:hover {
			background-color: rgba(136, 135, 135, 1);
		}
	}

	&-preview {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		border-radius: 5px;
		border: solid 1px var(--color-border);
	}
	&__preData {
		width: 100%;
		height: 100%;
		&-box {
			display: block;
			width: 100%;
			height: 100%;
			overflow: auto !important;

			&::-webkit-scrollbar {
				width: 4px;
				height: 4px;
				-webkit-appearance: none;
				-webkit-overflow-scrolling: auto;
			}
			&::-webkit-scrollbar-thumb {
				background-color: rgba(136, 135, 135, 0.5);
				border-radius: 4px;
			}
			&::-webkit-scrollbar-thumb:hover {
				height: 4px;

				background-color: rgba(136, 135, 135, 1);
			}
		}
		&-item {
			display: flex;

			& > div {
				border-bottom: solid 1px var(--color-border);
				background-color: var(--color-component-background);
				padding: 4px 8px;
			}
		}
		&-title {
			display: flex;
			width: 120%;
			& > div {
				position: sticky;
				top: 0;
				background-color: var(--color-border);
				padding: 4px 8px;
			}
		}
	}
}
</style>
