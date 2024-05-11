<!-- Developed by Taipei Urban Intelligence Center 2023-2024-->

<script setup>
import { ref, defineProps , watch,onBeforeMount } from "vue";
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


function handleConfirm() {
	handleClose();
}

function handleClose() {
	currentSettings.value = "sql";
	dialogStore.hideAllDialogs();
}
function handleSelectOthers(input){
	currentSettings.value = input;
}


const newInputStorage = ref({
    "id": null,
    "index": "",
    "name": "",
    "history_config": null,
    "map_config": [
        {
            "id": 64,
            "index": "socl_welfare_organization_plc",
            "title": "社福機構",
            "type": "circle",
            "source": "geojson",
            "size": "big",
            "icon": null,
            "paint": "{\n  \"circle-color\": [\n    \"match\",\n    [\n      \"get\",\n      \"main_type\"\n    ],\n    \"銀髮族服務\",\n    \"#F49F36\",\n    \"身障機構\",\n    \"#F65658\",\n    \"兒童與少年服務\",\n    \"#F5C860\",\n    \"社區服務、NPO\",\n    \"#9AC17C\",\n    \"婦女服務\",\n    \"#4CB495\",\n    \"貧困危機家庭服務\",\n    \"#569C9A\",\n    \"保護性服務\",\n    \"#60819C\",\n    \"#2F8AB1\"\n  ]\n}",
            "property": "[\n  {\n    \"key\": \"main_type\",\n    \"name\": \"主要類別\"\n  },\n  {\n    \"key\": \"sub_type\",\n    \"name\": \"次要分類\"\n  },\n  {\n    \"key\": \"name\",\n    \"name\": \"名稱\"\n  },\n  {\n    \"key\": \"address\",\n    \"name\": \"地址\"\n  }\n]"
        }
    ],
    "chart_config": {
        "index": "welfare_institutions",
        "color": [],
        "types": [],
        "unit": ""
    },
    "map_filter": "",
    "time_from": "static",
    "time_to": "",
    "update_freq": null,
    "update_freq_unit": "",
    "source": "",
    "short_desc": "",
    "long_desc": "",
    "use_case": "",
    "links": [],
    "contributors": [],
    "updated_at": "",
    "query_type": "",
    "chart_data": [],
	  "query_chart": ""
})

// Trim the input query on the frontend side
function trimQuery(query) {
    // Remove leading and trailing whitespace
    query = query.trim();
    // Replace multiple consecutive whitespace characters with a single space
    query = query.replace(/\s+/g, ' ');
    return query;
}

async function handlePushSqlData() {
    // Validate query chart and query type
    const queryType = newInputStorage.value.query_type,
          queryChart = newInputStorage.value.query_chart;

    if (!queryType || !queryChart) {
        console.error("Query chart and query type cannot be empty");
        return;
    }

    // Trim query chart
    const trimmedQuery = trimQuery(queryChart);

    // Prepare request body
    const requestBody = {
        queryType:  queryType,
        queryString: trimmedQuery
    };

    try {
        // Make HTTP request
		console.log(requestBody)
        const response = await http.post(
            "/helper/query", 
            JSON.stringify(JSON.parse(JSON.stringify(requestBody)))
        );
        console.log(response);
        // Show notification on success
        dialogStore.showNotification("success", response.data);
        
        // Update newInputStorage chart field if necessary
        // (Assuming there's a method to update the chart field)
        newInputStorage.chart_data = response.data;
    } catch (error) {
        // Handle errors
        console.error("Error:", error);
        dialogStore.showNotification("error", "An error occurred while processing the request.");
    }
}

const alllist = ref([])

watch(
    () => dialogStore.dialogs.adminCreator, 
    async (newValue, oldValue) => {
      if (newValue === true) {
        const response = await http.get("/helper/list-tables");
        alllist.value = response.data.tables;
      }
    }
);


</script>

<template>
  <DialogContainer
    :dialog="`adminCreator`"
    @on-close="handleClose"
  >
    <div class="adminCreator">
      <div class="adminCreator-header">
        <h2>組件設定</h2>
        <button 
		v-if="newInputStorage.chart_data.length!=0"
		 @click="handleConfirm" >
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
		  v-if="newInputStorage.chart_data.length!=0"
          :class="{ active: currentSettings === 'all' }"
          @click="handleSelectOthers('all')"
        >
          整體
        </button>
        <button
		v-if="newInputStorage.chart_data.length!=0"

          :class="{ active: currentSettings === 'chart' }"
          @click="handleSelectOthers('chart')"
        >
          圖表
        </button>
      <button
          v-if="newInputStorage.history_config !== null&&newInputStorage.chart_data.length!=0"
          :class="{ active: currentSettings === 'history' }"
          @click="handleSelectOthers('history')"
        >
          歷史軸
        </button> 
         <button
          v-if="newInputStorage.map_config[0] !== null&&newInputStorage.chart_data.length!=0"
          :class="{ active: currentSettings === 'map' }"
          @click="handleSelectOthers('map')"
        >
          地圖
        </button> 
      </div>
      <div class="adminCreator-content">
        <div class="adminCreator-settings">
		<div v-if="currentSettings === 'sql'" class="adminCreator-settings-items">
			      <label>
              組件名稱* ({{newInputStorage.name.length}}/10)
            </label>
            <input
              v-model="newInputStorage.name"
              type="text"
              :minlength="1"
              :maxlength="15"
              required
            >
            <div class="two-block">
              <label>組件 ID</label>
              <label>組件 Index</label>
            </div>
            <div class="two-block">
              <input
                type="text"
                :value="newInputStorage.id"
                disabled
              >
			  <select v-model="newInputStorage.index">
				<option :value="item" v-for="item in alllist" :key="item">
				{{item}}
				</option>
			</select>
           
            </div>
            <label>圖表資料型態 Query Type</label>
            <select v-model="newInputStorage.query_type">
              <option value="two_d">
                二維資料
              </option>
              <option value="three_d">
                三維資料
              </option>
              <option value="time">
                時間序列資料
              </option>
              <option value="percent">
                百分比資料
              </option>
              <option value="map_legend">
                圖例資料
              </option>
            </select>
			      <label>sql設定</label>
            <textarea v-model="newInputStorage.query_chart"/>
			  <div class="adminCreator-runsql">
        
				<button @click="handlePushSqlData()">匯入資料</button>
			  </div>
		</div>

          <div
             v-else-if="currentSettings === 'all'"
            class="adminCreator-settings-items"
          >
        
            <label>資料來源*</label>
            <input
              v-model="newInputStorage.source"
              type="text"
              :minlength="1"
              :maxlength="12"
              required
            >
            <label>更新頻率* (0 = 不定期更新)</label>
            <div class="two-block">
              <input
                v-model="newInputStorage.update_freq"
                type="number"
                :min="0"
                :max="31"
                required
              >
              <select v-model="newInputStorage.update_freq_unit">
                <option value="minute" />
                <option value="hour">
                  時
                </option>
                <option value="day">
                  天
                </option>
                <option value="week">
                  週
                </option>
                <option value="month">
                  月
                </option>
                <option value="year">
                  年
                </option>
              </select>
            </div>
            <label>資料區間</label>
            <!-- eslint-disable no-mixed-spaces-and-tabs -->
            <input
              :value="`${timeTerms[newInputStorage.time_from]}${
                timeTerms[newInputStorage.time_to]
                  ? ' ~ ' +
                    timeTerms[newInputStorage.time_to]
                  : ''
              }`"
              disabled
            >
            <label required>組件簡述* ({{
              newInputStorage.short_desc.length
            }}/50)</label>
            <textarea
              v-model="newInputStorage.short_desc"
              :minlength="1"
              :maxlength="50"
              required
            />
            <label>組件詳述* ({{
              newInputStorage.long_desc.length
            }}/100)</label>
            <textarea
              v-model="newInputStorage.long_desc"
              :minlength="1"
              :maxlength="100"
              required
            />
            <label>範例情境* ({{
              newInputStorage.use_case.length
            }}/100)</label>
            <textarea
              v-model="newInputStorage.use_case"
              :minlength="1"
              :maxlength="100"
              required
            />
            <label>資料連結</label>
            <InputTags
              :tags="newInputStorage.links"
              @deletetag="
                (index) => {
                  newInputStorage.links.splice(index, 1);
                }
              "
              @updatetagorder="
                (updatedTags) => {
                  newInputStorage.links = updatedTags;
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
                    newInputStorage.links.push(
                      tempInputStorage.link
                    );
                    tempInputStorage.link = '';
                  }
                }
              "
            >
            <label>貢獻者</label>
            <InputTags
              :tags="newInputStorage.contributors"
              @deletetag="
                (index) => {
                  newInputStorage.contributors.splice(
                    index,
                    1
                  );
                }
              "
              @updatetagorder="
                (updatedTags) => {
                  newInputStorage.contributors = updatedTags;
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
                    newInputStorage.contributors.push(
                      tempInputStorage.contributor
                    );
                    tempInputStorage.contributor = '';
                  }
                }
              "
            >
          </div>
          <div
            v-else-if="currentSettings === 'chart'"
            class="adminCreator-settings-items"
          >
            <label>資料單位*</label>
            <input
              v-model="newInputStorage.chart_config.unit"
              type="text"
              :minlength="1"
              :maxlength="6"
              required
            >
            <label>圖表類型*（限3種，依點擊順序排列）</label>
            <SelectButtons
              :tags="
                chartsPerDataType[newInputStorage.query_type]
              "
              :selected="newInputStorage.chart_config.types"
              :limit="3"
              @updatetagorder="
                (updatedTags) => {
                  newInputStorage.chart_config.types =
                    updatedTags;
                }
              "
            />
            <label>圖表顏色</label>
            <InputTags
              :tags="newInputStorage.chart_config.color"
              :color-data="true"
              @deletetag="
                (index) => {
                  newInputStorage.chart_config.color.splice(
                    index,
                    1
                  );
                }
              "
              @updatetagorder="
                (updatedTags) => {
                  newInputStorage.chart_config.color =
                    updatedTags;
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
                    newInputStorage.chart_config.color.push(
                      tempInputStorage.chartColor
                    );
                    tempInputStorage.chartColor = '#000000';
                  }
                }
              "
            >
            <div v-if="newInputStorage.map_config[0] !== null">
              <label>地圖篩選</label>
              <textarea
                v-model="newInputStorage.map_filter"
              />
            </div>
          </div>
          <div
            v-else-if="currentSettings === 'history'"
            class="adminCreator-settings-items"
          >
            <label>歷史軸時間區間
              (依點擊順序排列，資料無法預覽)</label>
            <SelectButtons
              :tags="[
                'month_ago',
                'quarter_ago',
                'halfyear_ago',
                'year_ago',
                'twoyear_ago',
                'fiveyear_ago',
                'tenyear_ago',
              ]"
              :selected="newInputStorage.history_config.range"
              :limit="5"
              @updatetagorder="
                (updatedTags) => {
                  newInputStorage.history_config.range =
                    updatedTags;
                }
              "
            />
            <label>歷史軸顏色 (若無提供沿用圖表顏色)</label>
            <InputTags
              :tags="newInputStorage.history_config.color"
              :color-data="true"
              @deletetag="
                (index) => {
                  newInputStorage.history_config.color.splice(
                    index,
                    1
                  );
                }
              "
              @updatetagorder="
                (updatedTags) => {
                  newInputStorage.history_config.color =
                    updatedTags;
                }
              "
            />
            <input
              v-model="tempInputStorage.historyColor"
              type="color"
              class="adminCreator-settings-inputcolor"
              @focusout="
                () => {
                  if (
                    tempInputStorage.historyColor.length ===
                    7
                  ) {
                    newInputStorage.history_config.color.push(
                      tempInputStorage.historyColor
                    );
                    tempInputStorage.historyColor =
                      '#000000';
                  }
                }
              "
            >
          </div>
          <div v-else-if="currentSettings === 'map'">
            <div
              v-for="(
                map_config, index
              ) in newInputStorage.map_config"
              :key="map_config.index"
              class="adminCreator-settings-items"
            >
              <hr v-if="index > 0">
              <label>地圖{{ index + 1 }} ID / Index</label>
              <div class="two-block">
                <input
                  :value="
                    newInputStorage.map_config[index].id
                  "
                  disabled
                >
                <input
                  v-model="
                    newInputStorage.map_config[index].index
                  "
                  :maxlength="30"
                  :minlength="1"
                  required
                >
              </div>

              <label>地圖{{ index + 1 }} 名稱* ({{
                newInputStorage.map_config[index].title
                  .length
              }}/10)</label>
              <input
                v-model="
                  newInputStorage.map_config[index].title
                "
                type="text"
                :minlength="1"
                :maxlength="10"
                required
              >
              <label>地圖{{ index + 1 }} 類型*</label>
              <select
                v-model="
                  newInputStorage.map_config[index].type
                "
              >
                <option
                  v-for="(value, key) in mapTypes"
                  :key="key"
                  :value="key"
                >
                  {{ value }}
                </option>
              </select>
              <label>地圖{{
                index + 1
              }}
                預設變形（大小/圖示）</label>
              <div class="two-block">
                <select
                  v-model="
                    newInputStorage.map_config[index].size
                  "
                >
                  <option :value="''">
                    無
                  </option>
                  <option value="small">
                    small (點圖)
                  </option>
                  <option value="big">
                    big (點圖)
                  </option>
                  <option value="wide">
                    wide (線圖)
                  </option>
                </select>
                <select
                  v-model="
                    newInputStorage.map_config[index].icon
                  "
                >
                  <option :value="''">
                    無
                  </option>
                  <option value="heatmap">
                    heatmap (點圖)
                  </option>
                  <option value="dash">
                    dash (線圖)
                  </option>
                  <option value="metro">
                    metro (符號圖)
                  </option>
                  <option value="metro-density">
                    metro-density (符號圖)
                  </option>
                  <option value="triangle_green">
                    triangle_green (符號圖)
                  </option>
                  <option value="triangle_white">
                    triangle_white (符號圖)
                  </option>
                  <option value="youbike">
                    youbike (符號圖)
                  </option>
                  <option value="bus">
                    bus (符號圖)
                  </option>
                </select>
              </div>
              <label>地圖{{ index + 1 }} Paint屬性</label>
              <textarea
                v-model="
                  newInputStorage.map_config[index].paint
                "
              />
              <label>地圖{{ index + 1 }} Popup標籤</label>
              <textarea
                v-model="
                  newInputStorage.map_config[index].property
                "
              />
            </div>
          </div>
        </div>
        <div class="adminCreator-preview">
          <DashboardComponent
            v-if="
              currentSettings === 'all' ||
                currentSettings === 'chart'
            "
            :key="`${newInputStorage.index}-${newInputStorage.chart_config.color}-${newInputStorage.chart_config.types}`"
            :config="JSON.parse(JSON.stringify(newInputStorage))"
            mode="large"
          />
          <div
            v-else-if="currentSettings === 'history'"
            :style="{ width: '300px' }"
          >
            <HistoryChart
              :key="`${newInputStorage.index}-${newInputStorage.history_config.color}`"
              :chart_config="newInputStorage.chart_config"
              :series="newInputStorage.history_data"
              :history_config="
                JSON.parse(
                  JSON.stringify(
                    newInputStorage.history_config
                  )
                )
              "
            />
          </div>
		  <div v-else-if="currentSettings === 'sql'" class="adminCreator__preData">
			<p v-if="newInputStorage.chart_data.length!=0">尚未匯入資料</p>
			<div class="adminCreator__preData-box">
				<div class="adminCreator__preData-title">
					<div>test1</div>
					<div>test2</div>
					<div>test3</div>
					<div>test4</div>
					<div>test5</div>
					<div>test5</div>

				</div>
				<div class="adminCreator__preData-item">
					<div>test1</div>
					<div>test2</div>
					<div>test3</div>
					<div>test4</div>
					<div>test5</div>
					<div>test5</div>
			
					   
					  
				</div>
				<div class="adminCreator__preData-item">
					<div>test1</div>
					<div>test2</div>
					<div>test3</div>
					<div>test4</div>
					<div>test5</div>
					<div>test5</div>
				
				</div>
			</div>
		</div>
          <div
            v-else-if="currentSettings === 'map'"
            index="componentsettings"
          >
            預覽功能 Coming Soon
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
		button{
			
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
			select{
				width: 100%;
			}
		}
		.three-block {
			display: grid;
			grid-template-columns: 1fr 2rem 1fr;
			column-gap: 0.5rem;
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
	&__preData{
		width: 100%;
		height: 100%;
		&-box{
			display: block;
			width: 100%;
			height: 100%;
			overflow-x: auto !important;
			
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
		&-item{
			display: flex;
			
			width: 120%;
			&>div{
				min-width: calc(100% / 5);
				border-bottom: solid 1px var(--color-border);
    background-color: var(--color-component-background);
	padding: 4px 8px;
			}
		}
		&-title{
			display: flex;
			width: 120%;
			&>div{
				min-width: calc(100% / 5);
				position: sticky;
    			top: 0;
  				  background-color: var(--color-border);
	padding: 4px 8px;
		
			}

		}
	}

}
</style>
