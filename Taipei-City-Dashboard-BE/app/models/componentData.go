// Package models stores the models for the postgreSQL databases.
package models

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

/* ----- Models ----- */

// ChartDataQuery is the model for getting the chart data query.
type ChartDataQuery struct {
	QueryType  string `json:"query_type" gorm:"column:query_type"`
	QueryChart string `json:"query_chart" gorm:"column:query_chart"`
}

// HistoryDataQuery is the model for getting the history data query.
type HistoryDataQuery struct {
	QueryHistory string `json:"query_history" gorm:"column:query_history"`
}

/*
TwoDimensionalData Json Format:

	{
		"data": [
			{
				"data": [
					{ "x": "", "y": 17 },
					...
				]
			}
		]
	}
*/
type TwoDimensionalData struct {
	Xaxis string  `gorm:"column:x_axis" json:"x"`
	Data  float64 `gorm:"column:data" json:"y"`
}
type TwoDimensionalDataOutput struct {
	Data []TwoDimensionalData `json:"data"`
}

/*
ThreeDimensionalData & PercentData Json Format:

	{
		"data": [
			{
				"name": ""
				"data": [...]
			},
			...
		]
	}

>> ThreeDimensionalData is shared by 3D and percentage data
*/
type ThreeDimensionalData struct {
	Xaxis string `gorm:"column:x_axis"`
	Icon  string `gorm:"column:icon"`
	Yaxis string `gorm:"column:y_axis"`
	Data  int    `gorm:"column:data"`
}

type ThreeDimensionalDataOutput struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Data []int  `json:"data"`
}

/*
TimeSeriesData Json Format:

	{
		"data": [
			{
				"name": "",
				"data": [
					{ "x": "2023-05-25T06:29:00+08:00", "y": 17 },
					...
				]
			},
			...
		]
	}
*/
type TimeSeriesData struct {
	Xaxis time.Time `gorm:"column:x_axis"`
	Yaxis string    `gorm:"column:y_axis"`
	Data  float64   `gorm:"column:data"`
}

type TimeSeriesDataItem struct {
	X string  `json:"x"`
	Y float64 `json:"y"`
}

type TimeSeriesDataOutput struct {
	Name string               `json:"name"`
	Data []TimeSeriesDataItem `json:"data"`
}

/*
MapLegendData Json Format:
*/
type MapLegendData struct {
	Name  string  `gorm:"column:name" json:"name"`
	Type  string  `gorm:"column:type" json:"type"`
	Icon  string  `gorm:"column:icon" json:"icon"`
	Value float64 `gorm:"column:value" json:"value"`
}

/* ----- Handlers ----- */

func GetComponentChartDataQuery(id int) (queryType string, queryString string, err error) {
	var chartDataQuery ChartDataQuery

	err = DBManager.
		Table("components").
		Select("query_type, query_chart").
		Where("components.id = ?", id).
		Find(&chartDataQuery).Error
	if err != nil {
		return queryType, queryString, err
	}
	return chartDataQuery.QueryType, chartDataQuery.QueryChart, nil
}

func GetComponentHistoryDataQuery(id int, timeFrom string, timeTo string) (queryHistory string, err error) {
	var historyDataQuery HistoryDataQuery

	err = DBManager.
		Table("components").
		Select("query_history").
		Where("components.id = ?", id).
		Find(&historyDataQuery).Error
	if err != nil {
		return queryHistory, err
	}
	if historyDataQuery.QueryHistory == "" {
		return historyDataQuery.QueryHistory, err
	}

	var timeStepUnit string

	timeFromTime, err := time.Parse("2006-01-02T15:04:05+08:00", timeFrom)
	if err != nil {
		return queryHistory, err
	}
	timeToTime, err := time.Parse("2006-01-02T15:04:05+08:00", timeTo)
	if err != nil {
		return queryHistory, err
	}

	/*
			timesteps are automatically determined based on the time range:
		  - Within 24hrs: hour
		  - Within 1 month: day
		  - Within 3 months: week
		  - Within 2 years: month
		  - More than 2 years: year
	*/
	if timeToTime.Sub(timeFromTime).Hours() <= 24 {
		timeStepUnit = "hour" // Within 24hrs
	} else if timeToTime.Sub(timeFromTime).Hours() < 24*32 {
		timeStepUnit = "day" // Within 1 month
	} else if timeToTime.Sub(timeFromTime).Hours() < 24*93 {
		timeStepUnit = "week" // Within 3 months
	} else if timeToTime.Sub(timeFromTime).Hours() < 24*740 {
		timeStepUnit = "month" // Within 2 years
	} else {
		timeStepUnit = "year" // More than 2 years
	}

	// Insert the time range and timestep unit into the query
	var queryInsertStrings []any

	if strings.Count(historyDataQuery.QueryHistory, "%s")%3 != 0 {
		return queryHistory, fmt.Errorf("invalid query string")
	}

	for i := 0; i < strings.Count(historyDataQuery.QueryHistory, "%s")/3; i++ {
		queryInsertStrings = append(queryInsertStrings, timeStepUnit, timeFrom, timeTo)
	}

	historyDataQuery.QueryHistory = fmt.Sprintf(historyDataQuery.QueryHistory, queryInsertStrings...)

	return historyDataQuery.QueryHistory, nil
}

/*
Below are the parsing functions for the four data types:
two_d, three_d, percent, and time. three_d and percent data share a common handler.
*/

func GetTwoDimensionalData(query *string, timeFrom string, timeTo string) (chartDataOutput []TwoDimensionalDataOutput, err error) {
	var chartData []TwoDimensionalData
	var queryString string

	// 1. Check if query contains substring '%s'. If so, the component can be queried by time.
	if strings.Count(*query, "%s") == 2 {
		queryString = fmt.Sprintf(*query, timeFrom, timeTo)
	} else {
		queryString = *query
	}

	// 2. Get the data from the database
	err = DBDashboard.Raw(queryString).Scan(&chartData).Error
	if err != nil {
		return chartDataOutput, err
	}
	if len(chartData) == 0 {
		return chartDataOutput, err
	}

	// 3. Convert the data to the format required by the front-end
	chartDataOutput = append(chartDataOutput, TwoDimensionalDataOutput{Data: chartData})

	return chartDataOutput, nil
}

func GetThreeDimensionalData(query *string, timeFrom string, timeTo string) (chartDataOutput []ThreeDimensionalDataOutput, categories []string, err error) {
	var chartData []ThreeDimensionalData
	var queryString string

	// 1. Check if query contains substring '%s'. If so, the component can be queried by time.
	if strings.Count(*query, "%s") == 2 {
		queryString = fmt.Sprintf(*query, timeFrom, timeTo)
	} else {
		queryString = *query
	}

	// 2. Get the data from the database
	err = DBDashboard.Raw(queryString).Scan(&chartData).Error
	if err != nil {
		return chartDataOutput, categories, err
	}

	// 3. Create maps to store unique categories and yAxis values
	categoryMap := make(map[string]bool)
	yAxisMap := make(map[string]int)

	// 4. Convert the data to the format required by the front-end
	for _, data := range chartData {
		// Add unique categories to categoryMap
		categoryMap[data.Xaxis] = true

		// Add unique yAxis values to yAxisMap and initialize corresponding index in chartDataOutput
		if _, ok := yAxisMap[data.Yaxis]; !ok {
			yAxisMap[data.Yaxis] = len(chartDataOutput)
			chartDataOutput = append(chartDataOutput, ThreeDimensionalDataOutput{Name: data.Yaxis, Icon: data.Icon, Data: make([]int, 0)})
		}
	}

	// 5. Populate categories slice from categoryMap
	for category := range categoryMap {
		categories = append(categories, category)
	}

	// 6. Sort categories
	sort.Strings(categories)

	// 7. Fill chartDataOutput with data
	for _, data := range chartData {
		yIndex := yAxisMap[data.Yaxis]
		categoryIndex := findIndex(categories, data.Xaxis)
		if categoryIndex != -1 {
			// Ensure chartDataOutput has enough elements to accommodate categoryIndex
			for len(chartDataOutput[yIndex].Data) <= categoryIndex {
				chartDataOutput[yIndex].Data = append(chartDataOutput[yIndex].Data, 0)
			}
			chartDataOutput[yIndex].Data[categoryIndex] = data.Data
		}
	}

	return chartDataOutput, categories, nil
}

// Function to find index of element in slice
func findIndex(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

func GetTimeSeriesData(query *string, timeFrom string, timeTo string) (chartDataOutput []TimeSeriesDataOutput, err error) {
	var chartData []TimeSeriesData
	var queryString string

	// 1. Check if query contains substring '%s'. If so, the component can be queried by time.
	if strings.Count(*query, "%s") == 2 {
		queryString = fmt.Sprintf(*query, timeFrom, timeTo)
	} else {
		queryString = *query
	}

	// 2. Get the data from the database
	err = DBDashboard.Raw(queryString).Scan(&chartData).Error
	if err != nil {
		return chartDataOutput, err
	}
	if len(chartData) == 0 {
		return chartDataOutput, err
	}

	// 3. Convert the data to the format required by the front-end
	for _, data := range chartData {
		// Group data together by yAxis
		var foundY bool
		formattedDate := data.Xaxis.Format("2006-01-02T15:04:05+08:00")
		for i, output := range chartDataOutput {
			if output.Name == data.Yaxis {
				// Append the data to the output
				chartDataOutput[i].Data = append(output.Data, TimeSeriesDataItem{X: formattedDate, Y: data.Data})
				foundY = true
				break
			}
		}

		// If a unique yAxis is found, create a new entry in the output
		if !foundY {
			chartDataOutput = append(chartDataOutput, TimeSeriesDataOutput{Name: data.Yaxis, Data: []TimeSeriesDataItem{{X: formattedDate, Y: data.Data}}})
		}
	}

	return chartDataOutput, nil
}

func GetMapLegendData(query *string, timeFrom string, timeTo string) (chartData []MapLegendData, err error) {
	var queryString string

	// 1. Check if query contains substring '%s'. If so, the component can be queried by time.
	if strings.Count(*query, "%s") == 2 {
		queryString = fmt.Sprintf(*query, timeFrom, timeTo)
	} else {
		queryString = *query
	}

	// 2. Get the data from the database
	err = DBDashboard.Raw(queryString).Scan(&chartData).Error
	if err != nil {
		return chartData, err
	}
	if len(chartData) == 0 {
		return chartData, err
	}

	return chartData, nil
}

// QueryExperiment is the model for executing a custom query and validating the output.
type QueryExperiment struct {
	Query string `json:"query"`
}

// Bypass dangerous keywords to avoid SQL injection
func QueryValidation(query string) error {
	// List of dangerous keywords
	dangerousKeywords := []string{"DROP", "DELETE", "TRUNCATE", "INSERT", "UPDATE", "ALTER"}

	// Check if the query contains any dangerous keywords
	for _, keyword := range dangerousKeywords {
		if strings.Contains(strings.ToUpper(query), keyword) {
			return fmt.Errorf("query contains invalid keyword: %s", keyword)
		}
	}

	return nil
}

// ListTablesInComponents lists all tables in the components database using the DBManager connection.
func ListTablesInComponents() ([]string, error) {
	var tables []string

	// Query the information schema to fetch table names
	rows, err := DBDashboard.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and extract table names
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}
