// test data for creating component
export const testComponent = {
	map: {
		index: "garbage_truck",
		title: "給愛麗絲",
		type: "",
		source: "",
		size: "",
		icon: "",
		// paint: "",
		// property: "",
		paint: JSON.stringify(
			JSON.parse(
				JSON.stringify({
					"circle-color": [
						"case",
						["<", ["get", "抵達時間"], 1700],
						"#E6DF44",
						["<", ["get", "抵達時間"], 1900],
						"#F4633C",
						["<", ["get", "抵達時間"], 2100],
						"#D63940",
						"#9C2A4B",
					],
				})
			)
		),
		property: JSON.stringify(
			JSON.parse(
				JSON.stringify([
					{ key: "行政區", name: "行政區" },
					{ key: "抵達時間", name: "抵達時間" },
					{ key: "離開時間", name: "離開時間" },
					{ key: "地點", name: "地址" },
				])
			)
		),
	},
	chart: {
		index: "garbage_truck",
		color: ["#E6DF44", "#F4633C", "#D63940", "#9C2A4B"],
		types: ["ColumnChart", "DistrictChart", "RadarChart"],
		unit: "處",
	},
	data: {
		index: "garbage_truck",
		name: "垃圾車收運",
		history_config: "",
		map_config: "",
		map_config_ids: "",
		chart_config: "",
		map_filter: "",
		time_from: "",
		time_to: "",
		update_freq: 0,
		update_freq_unit: "",
		source: "環保局",
		short_desc: "垃圾車收運！",
		long_desc: "垃圾車收運資料；猩猩，一起，強大。",
		use_case: "垃圾車good",
		links: [
			"https://data.taipei/dataset/detail?id=6bb3304b-4f46-4bb0-8cd1-60c66dcd1cae",
		],
		contributors: ["管誰import啥"],
		query_type: "three_d",
		query_chart: `
        SELECT * FROM (
            SELECT
                行政區 as x_axis,
                CASE
                    WHEN 抵達時間 BETWEEN 0 AND 1659 THEN '1700前'
                    WHEN 抵達時間 BETWEEN 1700 AND 1859 THEN '1700-1900'
                    WHEN 抵達時間 BETWEEN 1900 AND 2059 THEN '1900-2100'
                    ELSE '2100後'
                END AS y_axis,
                COUNT(*) AS data
            FROM garbage_truck
            GROUP BY
                行政區,
                CASE
                    WHEN 抵達時間 BETWEEN 0 AND 1659 THEN '1700前'
                    WHEN 抵達時間 BETWEEN 1700 AND 1859 THEN '1700-1900'
                    WHEN 抵達時間 BETWEEN 1900 AND 2059 THEN '1900-2100'
                    ELSE '2100後'
                END
        ) as t
        
        ORDER BY
            ARRAY_POSITION(ARRAY['北投區', '士林區', '內湖區', '南港區', '松山區', '信義區', '中山區', '大同區', '中正區', '萬華區', '大安區', '文山區']::varchar[], t.x_axis), 
            ARRAY_POSITION(ARRAY['1700前', '1700-1900', '1900-2100', '2100後'], t.y_axis);
        `,
	},
};
