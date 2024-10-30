package main

import (
	"encoding/json"
	"errors"
	_ "errors"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"math"
	"net/http"
	_ "net/http"
	"slices"
	"sort"
	"strconv"
	_ "strconv"
)

/*
Objective:
Develop a Movers Recommendation API using Go and the Gin framework (or Go Kit), allowing users to view, add, delete, and review mover organizations. This service will support interaction with movers and provide a mechanism for recording user recommendations, updating ratings, and retrieving mover information.

API Endpoints:

1. Add a Mover

Description: Allows the addition of a new mover to the system.
Endpoint: POST /movers
Request Body: JSON object containing:
name: String, required – name of the mover organization.
rate: Float (0.0 to 5.0), required – initial rating in 0.0 format.
telephone_number: String, required – contact phone number.
jobs_done: Integer, required – total completed jobs by the mover.
Response: Returns status and the added mover information in JSON format.
Status: Completed

2. Delete a Mover

Description: Deletes a mover from the system based on their unique ID.
Endpoint: DELETE /movers/<id>
Parameters:
id: Path parameter, required – ID of the mover to delete.
Response: Returns a success status on successful deletion, or an error message if the ID is not found.
Status: Completed

3. Get All Movers (Sorted)

Description: Retrieves a list of all movers, sorted alphabetically by mover name.
Endpoint: GET /movers
Response: JSON array of mover objects, each containing:
id, name, rate, telephone_number, jobs_done
Status: Completed

4. New Recommendation

Description: Allows a user to provide a review for a mover, updating the mover's average rating and completed jobs count.
Endpoint: POST /movers/<id>/review
Request Body: JSON object containing:
rate: Float (0.0 to 5.0), required – the rating provided by the user for this mover.
Response: Returns the updated mover information with the recalculated average rating.
Calculation Logic: Each new rating will update the mover's average rating based on the previous ratings and total completed jobs.
Status: Pending

_____________________
Implementation Notes:
 - Gin Package: Utilize Gin functions for JSON handling:
	Error handling: context.JSON(http.StatusBadRequest, gin.H{"error": "<error_message>"})
	Success response: context.JSON(http.StatusCreated, <response_data>)
		JSON Parsing: context.BindJSON(&<struct>)
 - Data Storage: The list of movers is currently stored as an in-memory array but can be migrated to a database in future versions.
*/

// Struct represents our mover model:
type mover struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Rate            float64 `json:"rate"`
	TelephoneNumber string  `json:"telephone_number"`
	JobsAmount      int     `json:"jobs_done"`
}

// MarshalJSON Custom MarshalJSON to round the Rate field in JSON output only
func (m mover) MarshalJSON() ([]byte, error) {
	type Alias mover                    // Alias to prevent recursion in MarshalJSON
	m.Rate = math.Round(m.Rate*10) / 10 // Round Rate to 1 decimal place for JSON output
	return json.Marshal((Alias)(m))
}

// Database of movers:
var movers = []mover{
	{ID: 1, Name: "San Francisco MOV", Rate: 4.6, TelephoneNumber: "+15615557689", JobsAmount: 3780},
	{ID: 2, Name: "Rapid Movers", Rate: 4.2, TelephoneNumber: "+15617384568", JobsAmount: 1240},
	{ID: 3, Name: "Reliable Relocations", Rate: 4.7, TelephoneNumber: "+14155538692", JobsAmount: 2050},
	{ID: 4, Name: "City Express Movers", Rate: 4.5, TelephoneNumber: "+18025559482", JobsAmount: 1870},
	{ID: 5, Name: "Pro Mover Co.", Rate: 4.8, TelephoneNumber: "+17024457893", JobsAmount: 2500},
	{ID: 6, Name: "MoveOn Solutions", Rate: 4.4, TelephoneNumber: "+19025548765", JobsAmount: 1730},
	{ID: 7, Name: "All Star Moving", Rate: 4.3, TelephoneNumber: "+13125587612", JobsAmount: 1290},
	{ID: 8, Name: "Swift Relocation", Rate: 4.6, TelephoneNumber: "+12026758741", JobsAmount: 3100},
	{ID: 9, Name: "Speedy Transport", Rate: 4.5, TelephoneNumber: "+14027759832", JobsAmount: 1980},
	{ID: 10, Name: "Premier Movers", Rate: 4.7, TelephoneNumber: "+15022556478", JobsAmount: 2300},
	{ID: 11, Name: "Ace Relocators", Rate: 4.3, TelephoneNumber: "+16024457812", JobsAmount: 1670},
	{ID: 12, Name: "Trusted Movers Co.", Rate: 4.6, TelephoneNumber: "+17024459874", JobsAmount: 2890},
	{ID: 13, Name: "Urban Move", Rate: 4.5, TelephoneNumber: "+18024458736", JobsAmount: 3200},
	{ID: 14, Name: "FastTrack Movers", Rate: 4.7, TelephoneNumber: "+13027758495", JobsAmount: 2150},
	{ID: 15, Name: "Metro Moving Solutions", Rate: 4.4, TelephoneNumber: "+14028854721", JobsAmount: 1390},
}

func getMoverById(id int) (*mover, error) {
	for i, mover := range movers {
		if mover.ID == id {
			return &movers[i], nil
		}
	}
	return nil, errors.New("mover not found")
}

func deleteElement(slice []mover, index int) []mover {
	return slices.Delete(slice, index, index+1)
}

func findMoverIndexById(id int) (int, error) {
	for index, mover := range movers {
		if mover.ID == id {
			return index, nil
		}
	}
	return -1, errors.New("mover not found")
}

// GET request. Sort by Rate. If rates are equal, sort by ID
func getMovers(context *gin.Context) {
	if movers == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "movers list is empty"})
		return
	}
	sort.Slice(movers, func(i, j int) bool {
		if movers[i].Rate == movers[j].Rate {
			return movers[i].ID < movers[j].ID
		} else {
			return movers[i].Rate > movers[j].Rate
		}
	})

	context.JSON(http.StatusOK, movers)
}

// POST request. Add a new mover
func addMover(context *gin.Context) {
	var newMover mover

	if err := context.BindJSON(&newMover); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	movers = append(movers, newMover)
	context.JSON(http.StatusCreated, newMover)
}

// DELETE request. Delete mover by ID
func deleteMover(context *gin.Context) {
	id := context.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Conversion error"})
		return
	}

	moverIndex, err := findMoverIndexById(intId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Mover not found"})
		return
	}

	movers = deleteElement(movers, moverIndex)

	context.JSON(http.StatusOK, gin.H{"message": "Mover deleted successfully"})
}

func updateMover(context *gin.Context) {
	id := context.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Conversion error"})
		return
	}

	currMover, currErr := getMoverById(intId)

	if currErr != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Mover not found"})
		return
	}

	var updatedMover mover

	if err := context.BindJSON(&updatedMover); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if updatedMover.Rate >= 0.0 && updatedMover.Rate <= 5.0 {
		totalJobs := currMover.JobsAmount

		currMover.Rate = (currMover.Rate*float64(totalJobs) + updatedMover.Rate) / (float64(totalJobs) + 1)
		currMover.JobsAmount += 1
		// Calculate the average rate based on provided rate (if provided)
	} else {
		context.JSON(http.StatusExpectationFailed, gin.H{"error": "Provided rate should be in range between 0 and 5"})
		return
	}

	context.JSON(http.StatusOK, currMover)
}

func main() {
	router := gin.Default()

	router.GET("/movers", getMovers)
	router.POST("/movers", addMover)
	router.DELETE("/movers/:id", deleteMover)
	router.POST("/movers/:id/review", updateMover)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
