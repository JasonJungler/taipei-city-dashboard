// Package controllers stores all the controllers for the Gin router.
package controllers

import (
	"log"
	"net/http"
	"strconv"

	"TaipeiCityDashboardBE/app/models"
	"TaipeiCityDashboardBE/app/util"

	"github.com/gin-gonic/gin"
)

type QueryExperimentRequest struct {
	QueryType   string `json:"queryType"`
	QueryString string `json:"queryString"`
}

/*
GetComponentChartData retrieves the chart data for a component.
/api/v1/components/:id/chart

header: time_from, time_to (optional)
*/
func GetComponentChartData(c *gin.Context) {
	// 1. Get the component id from the URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid component ID"})
		return
	}

	// 2. Get the chart data query and chart data type from the database
	queryType, queryString, err := models.GetComponentChartDataQuery(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if (queryString == "") || (queryType == "") {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "No chart data available"})
		return
	}

	timeFrom, timeTo := util.GetTime(c)

	// 3. Get and parse the chart data based on chart data type
	if queryType == "two_d" {
		chartData, err := models.GetTwoDimensionalData(&queryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})
	} else if queryType == "three_d" || queryType == "percent" {
		chartData, categories, err := models.GetThreeDimensionalData(&queryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData, "categories": categories})
	} else if queryType == "time" {
		chartData, err := models.GetTimeSeriesData(&queryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})
	} else if queryType == "map_legend" {
		chartData, err := models.GetMapLegendData(&queryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})
	}
}

/*
GetComponentHistoryData retrieves the history data for a component.
/api/v1/components/:id/history

header: time_from, time_to (mandatory)
timesteps are automatically determined based on the time range:
  - Within 24hrs: hour
  - Within 1 month: day
  - Within 3 months: week
  - Within 2 years: month
  - More than 2 years: year
*/
func GetComponentHistoryData(c *gin.Context) {
	// 1. Get the component id from the URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid component ID"})
		return
	}

	timeFrom, timeTo := util.GetTime(c)

	// 2. Get the history data query from the database
	queryHistory, err := models.GetComponentHistoryDataQuery(id, timeFrom, timeTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if queryHistory == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "No history data available"})
		return
	}

	// 3. Get and parse the history data
	chartData, err := models.GetTimeSeriesData(&queryHistory, timeFrom, timeTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})
}

// Define the route handler for executing the query experiment
func QueryChartData(c *gin.Context) {
	// Parse the JSON request body
	var request QueryExperimentRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the queryType and queryString
	if request.QueryType == "" || request.QueryString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "queryType and queryString are required"})
		return
	}

	// Execute the query experiment
	timeFrom, timeTo := util.GetTime(c)

	// 3. Get and parse the chart data based on chart data type
	if request.QueryType == "two_d" {
		chartData, err := models.GetTwoDimensionalData(&request.QueryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		log.Println(chartData)
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})

	} else if request.QueryType == "three_d" || request.QueryType == "percent" {
		chartData, categories, err := models.GetThreeDimensionalData(&request.QueryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		log.Println(chartData)
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData, "categories": categories})

	} else if request.QueryType == "time" {
		chartData, err := models.GetTimeSeriesData(&request.QueryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		log.Println(chartData)
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})
	} else if request.QueryType == "map_legend" {
		chartData, err := models.GetMapLegendData(&request.QueryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		log.Println(chartData)
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid queryType"})
		return
	}
}

// Define the route handler for listing tables in the components database
func ListTablesInComponentsHandler(c *gin.Context) {
	// Call the function to list tables in components
	tables, err := models.ListTablesInComponents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	log.Println(tables)
	c.JSON(http.StatusOK, gin.H{"status": "success", "tables": tables})
}
