package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

const healthzEndpoint = "http://localhost:80/v1/healthz"
const errorEndpoint = "http://localhost:80/v1/err"
const usersEndpoint = "http://localhost:80/v1/users"
const feedsEndpoint = "http://localhost:80/v1/feeds"

func cleanUp(userId string) {
	// cleanup by deleting the created test user from DB
	// create a delete request and send it to /users endpoint
	delReq, err := http.NewRequest("DELETE", usersEndpoint+"/"+userId, nil)
	if err != nil {
		log.Printf("Error creating delete request for the cleanup: %v", err)
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	delResp, err := client.Do(delReq)
	if err != nil {
		log.Printf("Could not get response from the delete request: %v", err)
	}
	if delResp.StatusCode != 204 {
		log.Printf("Could not delete the test user created during this test. Got resp code: %v", delResp.StatusCode)
	}
}

func TestHealthzEndpoint(t *testing.T) {
	// create HTTP client to send a request
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	// send HTTP request to target endpoint
	resp, err := client.Get(healthzEndpoint)
	if err != nil {
		t.Fatalf("Failed to get a response from endpoint %v", healthzEndpoint)
	}
	// check for response status code
	if resp.StatusCode != 200 {
		t.Errorf("Failed to get correct response, got: %v want: 200", resp.StatusCode)
	}
}

func TestErrorEndpoint(t *testing.T) {
	// create HTTP client to send a request
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	// send HTTP request to the target endpoint
	resp, err := client.Get(errorEndpoint)
	if err != nil {
		t.Fatalf("Failed to get a response from endpoint %v", errorEndpoint)
	}
	// check for response status code
	if resp.StatusCode != 400 {
		t.Errorf("Failed to get correct response, got: %v want: 400", resp.StatusCode)
	}
}

func TestCreateUser(t *testing.T) {
	// create HTTP client to send a request
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	// create request body
	var jsonReq = []byte(`{
		"name": "Test User"
	}`)
	// send POST request to the endpoint
	resp, err := client.Post(usersEndpoint, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatalf("Failed to get a response from endpoint %v", usersEndpoint)
	}
	// check for correct response status code
	if resp.StatusCode != 201 {
		t.Errorf("Failed to get correct response, got: %v want: 201", resp.StatusCode)
	}
	// read User ID from the response body
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading create user response: %v", err)
	}
	var jsonRespUser map[string]string
	err = json.Unmarshal(dat, &jsonRespUser)
	if err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
	}
	userId := jsonRespUser["id"]
	// cleanup
	cleanUp(userId)
}

func TestDeleteUser(t *testing.T) {
	// create HTTP client to send a request
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	// create request body
	var jsonReq = []byte(`{
			"name": "Test User for Delete User Test"
		}`)
	// send POST request to the endpoint
	resp, err := client.Post(usersEndpoint, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatalf("Failed to get a response from endpoint %v", usersEndpoint)
	}
	// check if the user has been created
	if resp.StatusCode != 201 {
		log.Printf("Failed to create user, got: %v want: 201", resp.StatusCode)
	}
	// get user id of the created test user
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading create user response: %v", err)
	}
	var jsonRespUser map[string]string
	err = json.Unmarshal(dat, &jsonRespUser)
	if err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
	}
	userId := jsonRespUser["id"]
	// create a delete request and send it to /users endpoint
	delReq, err := http.NewRequest("DELETE", usersEndpoint+"/"+userId, nil)
	if err != nil {
		log.Printf("Error creating delete request: %v", err)
	}

	delResp, err := client.Do(delReq)
	if err != nil {
		log.Printf("Could not get response from the delete request: %v", err)
	}
	// check for correct response status code
	if delResp.StatusCode != 204 {
		t.Errorf("Failed to get correct response, got: %v want: 204", delResp.StatusCode)
	}
}

func TestCreateFeed(t *testing.T) {
	// create HTTP client to send a request
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	// create a user first
	var jsonReqUser = []byte(`{
		"name": "Test User for Create Feed Test"
	}`)
	resp, err := client.Post(usersEndpoint, "application/json", bytes.NewBuffer(jsonReqUser))
	if err != nil {
		t.Fatalf("Failed to get a response from endpoint %v", errorEndpoint)
	}
	// check if the user was created
	if resp.StatusCode != 201 {
		log.Printf("Test user not created, got: %v want: 201", resp.StatusCode)
	}
	// read user ID and API key from the response body
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading create user response: %v", err)
	}
	var jsonRespUser map[string]string
	err = json.Unmarshal(dat, &jsonRespUser)
	if err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
	}
	userId := jsonRespUser["id"]
	apiKey := jsonRespUser["apiKey"]
	authzVal := "ApiKey " + apiKey
	// create a feed with user's API key
	var jsonReqFeed = []byte(`{
		"name": "Test User's Test Feed",
		"url": "https://test.com/testcreatefeed"
	}`)
	httpFeedReq, err := http.NewRequest("POST", feedsEndpoint, bytes.NewBuffer(jsonReqFeed))
	if err != nil {
		log.Printf("Error creating request for create feed test: %v", err)
	}
	httpFeedReq.Header.Set("Authorization", authzVal)
	feedResp, err := client.Do(httpFeedReq)
	if err != nil {
		t.Fatalf("Failed to get a response from endpoint %v", feedsEndpoint)
	}
	// check for correct response status code
	if feedResp.StatusCode != 201 {
		t.Errorf("Failed to get correct response, got: %v want: 201", feedResp.StatusCode)
	}
	// cleanup
	cleanUp(userId)
}

func TestDeleteFeed(t *testing.T) {
	// create HTTP client to send a request
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	// create a user first
	var jsonReqUser = []byte(`{
		"name": "Test User for Delete Feed Test"
	}`)
	resp, err := client.Post(usersEndpoint, "application/json", bytes.NewBuffer(jsonReqUser))
	if err != nil {
		t.Fatalf("Failed to get a response from endpoint %v", errorEndpoint)
	}
	// check for correct response status code
	if resp.StatusCode != 201 {
		log.Printf("Failed to create Test user, got: %v want: 201", resp.StatusCode)
	}
	// read userId and API key from the response body
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading create user response: %v", err)
	}
	var jsonRespUser map[string]string
	err = json.Unmarshal(dat, &jsonRespUser)
	if err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
	}
	userId := jsonRespUser["id"]
	apiKey := jsonRespUser["apiKey"]
	authzVal := "ApiKey " + apiKey
	// create a feed with user's API key
	var jsonReqFeed = []byte(`{
		"name": "Test User's Test Feed for deletion",
		"url": "https://test.com/testdeletefeed"
	}`)
	httpFeedReq, err := http.NewRequest("POST", feedsEndpoint, bytes.NewBuffer(jsonReqFeed))
	if err != nil {
		log.Printf("Error creating request for create feed test: %v", err)
	}
	httpFeedReq.Header.Set("Authorization", authzVal)
	feedResp, err := client.Do(httpFeedReq)
	if err != nil {
		t.Fatalf("Failed to get a response from endpoint %v", feedsEndpoint)
	}
	// check if feed was created
	if feedResp.StatusCode != 201 {
		log.Printf("Failed to create feed, got: %v want: 201", feedResp.StatusCode)
	}
	// get feed id from the response
	defer feedResp.Body.Close()
	dat, err = io.ReadAll(feedResp.Body)
	if err != nil {
		log.Printf("failed to read create feed response: %v", err)
	}
	var jsonRespFeed map[string]string
	err = json.Unmarshal(dat, &jsonRespFeed)
	if err != nil {
		log.Printf("could not unmarshal feed response: %v", err)
	}
	feedId := jsonRespFeed["id"]
	// create delete feed request
	delReq, err := http.NewRequest("DELETE", feedsEndpoint+"/"+feedId, nil)
	if err != nil {
		log.Printf("failed to create delete feed request: %v", err)
	}
	delReq.Header.Set("Authorization", authzVal)
	delFeedResp, err := client.Do(delReq)
	if err != nil {
		t.Fatalf("failed to get response from delete feed endpoint: %v", err)
	}
	// check for correct response status code
	if delFeedResp.StatusCode != 204 {
		t.Errorf("Failed to get correct response, got: %v want: 204", delFeedResp.StatusCode)
	}
	// cleanup
	cleanUp(userId)
}
