package handlers

import (
	"golang_microservice/plant-api/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// Initialize struct type Plant with properties
type Plant struct {
	logger *log.Logger
}

// Create struct type Plant with properties
func NewPlant(logger *log.Logger) *Plant {
	return &Plant{
		logger,
	}
}

// ServeHTTP is the main entry point for the handler and statify http.Handler interface
func (plant *Plant) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	//Handle the get Request to fetch list of plants data
	if request.Method == http.MethodGet {
		plant.getPlants(response, request)
		return
	}

	//Handle the post Request to insert new plant data into the datastore
	if request.Method == http.MethodPost {
		plant.createPlant(response, request)
		return
	}

	//Handle the put Request to update the available plant data into the datastore based on id
	if request.Method == http.MethodPut {
		plant.logger.Println(request.URL.Path)
		path := request.URL.Path
		reg := `/([0-9]+)`
		regC := regexp.MustCompile(reg)
		idMatchList := regC.FindAllStringSubmatch(path, -1)

		plant.logger.Println(idMatchList)
		if len(idMatchList) != 1 {
			plant.logger.Printf("Request URL IDList : %#v", idMatchList)
			http.Error(response, "Request URL have more than one id", http.StatusBadRequest)
			return
		}

		if len(idMatchList[0]) != 2 {
			plant.logger.Printf("Request URL caturedIDList : %#v", idMatchList[0])
			http.Error(response, "Regex found more id capture in the Request URL", http.StatusBadRequest)
			return
		}

		idString := idMatchList[0][1]
		id, err := strconv.Atoi(idString)
		plant.logger.Printf("ID: %d", id)

		if err != nil {
			plant.logger.Printf("Unable to convert id to int %v", err)
			http.Error(response, "Unable to convert id to int", http.StatusBadRequest)
			return
		}

		plant.updatePlant(id, response, request)

	}

	// Handle the deltet Request to delete the plant data into the datastore based on id
	if request.Method == http.MethodDelete {

		path := request.URL.Path
		reg := `/([0-9]+)`
		regC := regexp.MustCompile(reg)
		idMatchList := regC.FindAllStringSubmatch(path, -1)

		if len(idMatchList) != 1 {
			plant.logger.Printf("Request URL IDList : %#v", idMatchList)
			http.Error(response, "Request URL have more than one id", http.StatusBadRequest)
			return
		}

		if len(idMatchList[0]) != 2 {
			plant.logger.Printf("Request URL caturedIDList : %#v", idMatchList[0])
			http.Error(response, "Regex found more id capture in the Request URL", http.StatusBadRequest)
			return
		}

		idString := idMatchList[0][1]
		id, err := strconv.Atoi(idString)
		plant.logger.Println("ID", id)

		if err != nil {
			plant.logger.Printf("Unable to convert id to int %v", err)
			http.Error(response, "Unable to convert id to int", http.StatusBadRequest)
			return
		}

		plant.deletePlant(id, response, request)
	}

	//Catch all
	//If no method is statisfied return an error
	response.WriteHeader(http.StatusMethodNotAllowed)
}

//getPlants is used to fetch all the plants data from the datastore
func (plant *Plant) getPlants(response http.ResponseWriter, request *http.Request) {
	plantsList := data.GetAllPlants()
	marshalError := plantsList.ToJSON(response)

	if marshalError != nil {
		plant.logger.Printf("While, Marshaling the plant data. Reason : %s", marshalError)
		http.Error(response, "JSON marshaling failed.", http.StatusInternalServerError)
	}
}

//createPlants used to insert the new plant data into the datastore
func (plant *Plant) createPlant(response http.ResponseWriter, request *http.Request) {

	plantData := &data.Plant{}
	unMarshalError := plantData.FromJSON(request.Body)

	if unMarshalError != nil {
		plant.logger.Printf("While, UnMarshaling the plant data. Reason : %s", unMarshalError)
		http.Error(response, "JSON Unmarshaling failed.", http.StatusBadRequest)
	}

	data.AddPlant(plantData)
	plantsList := &data.Plants{plantData}
	plantsList.ToJSON(response)
}

//updatePlant used to update the plant data into the datastore based on id.
func (plant *Plant) updatePlant(id int, response http.ResponseWriter, request *http.Request) {

	plantData := &data.Plant{}
	marshalError := plantData.FromJSON(request.Body)

	if marshalError != nil {
		plant.logger.Printf("While, Marshaling the plant data. Reason : %s", marshalError)
		http.Error(response, "JSON Unmarshaling failed.", http.StatusBadRequest)
	}

	plant.logger.Printf("Plant : %#v", plantData)
	updateError := data.UpdatePlant(id, plantData)

	if updateError != nil {
		plant.logger.Printf("While, Update the plant data. Reason : %s", updateError)
		http.Error(response, "Plant not found.", http.StatusNotFound)
	}
	response.WriteHeader(http.StatusOK)
}

//deletePlant used to delete the plant data into the datastore based on id.
func (plant *Plant) deletePlant(id int, response http.ResponseWriter, request *http.Request) {

	deleteError := data.DeletePlant(id)

	if deleteError != nil {
		plant.logger.Printf("While, Delete the plant data. Reason : %s", deleteError)
		http.Error(response, "Plant not found.", http.StatusNotFound)
	}
	response.WriteHeader(http.StatusOK)
}
