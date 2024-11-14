package main

import (
	"encoding/json"
	"errors"
	_ "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"log"
	"math"
	"net/http"
	_ "net/http"
	"os"
	_ "os"
	"slices"
	"sort"
	"strconv"
	_ "strconv"
)

// Struct represents our mover model:
type mover struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Rating          float64 `json:"rating"`
	TelephoneNumber string  `json:"telephone_number"`
	JobsAmount      int     `json:"jobs_done"`
}

// MarshalJSON Custom MarshalJSON to round the Rating field in JSON output only
func (m mover) MarshalJSON() ([]byte, error) {
	type Alias mover                        // Alias to prevent recursion in MarshalJSON
	m.Rating = math.Round(m.Rating*10) / 10 // Round Rating to 1 decimal place for JSON output
	return json.Marshal((Alias)(m))
}

// Database of movers:
var movers = []mover{
	{ID: 1, Name: "San Francisco MOV", Rating: 4.6, TelephoneNumber: "+15615557689", JobsAmount: 3780},
	{ID: 2, Name: "Rapid Movers", Rating: 4.2, TelephoneNumber: "+15617384568", JobsAmount: 1240},
	{ID: 3, Name: "Reliable Relocations", Rating: 4.7, TelephoneNumber: "+14155538692", JobsAmount: 2050},
	{ID: 4, Name: "City Express Movers", Rating: 4.5, TelephoneNumber: "+18025559482", JobsAmount: 1870},
	{ID: 5, Name: "Pro Mover Co.", Rating: 4.8, TelephoneNumber: "+17024457893", JobsAmount: 2500},
	{ID: 6, Name: "MoveOn Solutions", Rating: 4.4, TelephoneNumber: "+19025548765", JobsAmount: 1730},
	{ID: 7, Name: "All Star Moving", Rating: 4.3, TelephoneNumber: "+13125587612", JobsAmount: 1290},
	{ID: 8, Name: "Swift Relocation", Rating: 4.6, TelephoneNumber: "+12026758741", JobsAmount: 3100},
	{ID: 9, Name: "Speedy Transport", Rating: 4.5, TelephoneNumber: "+14027759832", JobsAmount: 1980},
	{ID: 10, Name: "Premier Movers", Rating: 4.7, TelephoneNumber: "+15022556478", JobsAmount: 2300},
	{ID: 11, Name: "Ace Relocators", Rating: 4.3, TelephoneNumber: "+16024457812", JobsAmount: 1670},
	{ID: 12, Name: "Trusted Movers Co.", Rating: 4.6, TelephoneNumber: "+17024459874", JobsAmount: 2890},
	{ID: 13, Name: "Urban Move", Rating: 4.5, TelephoneNumber: "+18024458736", JobsAmount: 3200},
	{ID: 14, Name: "FastTrack Movers", Rating: 4.7, TelephoneNumber: "+13027758495", JobsAmount: 2150},
	{ID: 15, Name: "Metro Moving Solutions", Rating: 4.4, TelephoneNumber: "+14028854721", JobsAmount: 1390},
}

// Helper functions
func extractId(context *gin.Context) (int, error) {
	idParam := context.Param("id")
	MoverId, err := strconv.Atoi(idParam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Conversion error"})
		return -1, err
	} else {
		return MoverId, nil
	}
}

func getMoverById(id int) (*mover, error) {
	for i, mover := range movers {
		if mover.ID == id {
			return &movers[i], nil
		}
	}
	return nil, errors.New("mover not found")
}

func findMoverIndexById(id int) (int, error) {
	for index, mover := range movers {
		if mover.ID == id {
			return index, nil
		}
	}
	return -1, errors.New("mover not found")
}

func deleteElement(slice []mover, index int) []mover {
	return slices.Delete(slice, index, index+1)
}

func sortMoversByRatingAndId(movers []mover) []mover {
	moversCopy := make([]mover, len(movers))
	copy(moversCopy, movers)

	sort.Slice(moversCopy, func(i, j int) bool {
		if moversCopy[i].Rating == moversCopy[j].Rating {
			return moversCopy[i].ID < moversCopy[j].ID
		} else {
			return moversCopy[i].Rating > moversCopy[j].Rating
		}
	})
	return moversCopy
}

func checkMoverExists(newMover mover) bool {
	for _, existingMover := range movers {
		if existingMover.ID == newMover.ID || existingMover.Name == newMover.Name {
			return true
		}
	}
	return false
}

func checkMoverTelNumber(newMover mover) bool {
	for _, existingMover := range movers {
		if existingMover.TelephoneNumber == newMover.TelephoneNumber {
			return true
		}
	}
	return false
}

func initializeRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/movers", getMovers)
	router.POST("/movers", addMover)
	router.DELETE("/movers/:id", deleteMover)
	router.POST("/movers/:id/review", recommendMover)

	return router
}

// Main Functions
// GET request. Sort by Rating. If rates are equal, sort by ID
func getMovers(context *gin.Context) {
	if len(movers) == 0 {
		context.JSON(http.StatusNotFound, gin.H{"error": "movers list is empty"})
		return
	}

	sortedMovers := sortMoversByRatingAndId(movers)

	context.JSON(http.StatusOK, sortedMovers)
}

// POST request. Add a new mover
func addMover(context *gin.Context) {

	var newMover mover
	if err := context.BindJSON(&newMover); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	//checks if mover already exists
	if checkMoverExists(newMover) {
		context.JSON(http.StatusNotFound, gin.H{"error": "Mover already exists"})
		return
	}

	//Checks if the tel. number is occupied
	if checkMoverTelNumber(newMover) {
		context.JSON(http.StatusNotFound, gin.H{"error": "Tel. number is occupied"})
		return
	}

	movers = append(movers, newMover)
	context.JSON(http.StatusCreated, newMover)
}

// DELETE request. Delete mover by ID
func deleteMover(context *gin.Context) {
	MoverId, err := extractId(context)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Extracting ID error"})
		return
	}

	//!create function DeleteMoverById that will implement delete logic

	moverIndex, err := findMoverIndexById(MoverId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Mover not found"})
		return
	}

	movers = deleteElement(movers, moverIndex)

	context.JSON(http.StatusOK, gin.H{"message": "Mover deleted successfully"})
}

// POST request. Recommendation from users, updating average mover rate
func recommendMover(context *gin.Context) {
	MoverId, err := extractId(context)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Extracting ID error"})
		return
	}

	existingMover, getErr := getMoverById(MoverId)

	if getErr != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Mover not found"})
		return
	}

	var updatedRating mover

	if err := context.BindJSON(&updatedRating); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if updatedRating.Rating >= 0.0 && updatedRating.Rating <= 5.0 {
		totalJobs := existingMover.JobsAmount

		existingMover.Rating = (existingMover.Rating*float64(totalJobs) + updatedRating.Rating) / (float64(totalJobs) + 1)
		existingMover.JobsAmount += 1
		// Calculate the average rate based on provided rate (if provided)
	} else {
		context.JSON(http.StatusExpectationFailed, gin.H{"error": "Provided rate should be in range between 0 and 5"})
		return
	}
	context.JSON(http.StatusOK, existingMover)
}

func main() {
	//load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// getting env variables HOST and PORT
	serverHost := os.Getenv("HOST")
	serverPort := os.Getenv("PORT")

	router := initializeRouter()

	routerErr := router.Run(fmt.Sprintf("%s:%s", serverHost, serverPort))
	if routerErr != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
