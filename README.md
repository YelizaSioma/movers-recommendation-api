# Objective:
Develop a Movers Recommendation API using Go and the Gin framework (or Go Kit), allowing users to view, add, delete, and review mover organizations. This service will support interaction with movers and provide a mechanism for recording user recommendations, updating ratings, and retrieving mover information.

# API Endpoints:

1. Add a Mover

- Description: Allows the addition of a new mover to the system.
- Endpoint: POST /movers
- Request Body: JSON object containing:
name: String, required – name of the mover organization.
rate: Float (0.0 to 5.0), required – initial rating in 0.0 format.
telephone_number: String, required – contact phone number.
jobs_done: Integer, required – total completed jobs by the mover.
- Response: Returns status and the added mover information in JSON format.

2. Delete a Mover

- Description: Deletes a mover from the system based on their unique ID.
- Endpoint: DELETE /movers/<id>
- Parameters:
id: Path parameter, required – ID of the mover to delete.
- Response: Returns a success status on successful deletion, or an error message if the ID is not found.

3. Get All Movers (Sorted)

- Description: Retrieves a list of all movers, sorted alphabetically by mover name.
- Endpoint: GET /movers
- Response: JSON array of mover objects, each containing:
id, name, rate, telephone_number, jobs_done

4. New Recommendation

- Description: Allows a user to provide a review for a mover, updating the mover's average rating and completed jobs count.
- Endpoint: POST /movers/<id>/review
- Request Body: JSON object containing:
rate: Float (0.0 to 5.0), required – the rating provided by the user for this mover.
- Response: Returns the updated mover information with the recalculated average rating.
- Calculation Logic: Each new rating will update the mover's average rating based on the previous ratings and total completed jobs.

_____________________
## Implementation Notes:
 - Gin Package: Utilize Gin functions for JSON handling:
	Error handling: context.JSON(http.StatusBadRequest, gin.H{"error": "<error_message>"})
	Success response: context.JSON(http.StatusCreated, <response_data>)
	sing context.JSON() instead of less optimized context.IntendedJSON()
	JSON Parsing: context.BindJSON(&<struct>)
 - Data Storage: The list of movers is currently stored as an in-memory array but can be migrated to a database in future versions.