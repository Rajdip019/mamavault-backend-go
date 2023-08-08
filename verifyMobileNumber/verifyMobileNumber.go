package helloworld

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// setting environment variables
var accountSid = os.Getenv("TWILIO_ACCOUNT_SID")
var authToken = os.Getenv("TWILIO_AUTH_TOKEN")
var twilioNumber = os.Getenv("TWILIO_NUMBER")

// Response is the response from SendMessage function
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	SentTo  string `json:"sent_to"`
}

func init() {
	functions.HTTP("VerifyMobileNumber", VerifyMobileNumber) // entry point for the function
}

// InitializeApp initializes a firebase app
func InitializeApp() (*firestore.Client, context.Context) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "mamavault"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	firestore, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return firestore, ctx
}

// SendMessage sends a message to the given number
func SendMessage(number string, verification_link string, name string) (Response, error) {

	fmt.Println("accountSid", accountSid)
	fmt.Println("authToken", authToken)
	// exit if environment variables are not set
	if accountSid == "" || authToken == "" {
		fmt.Println("environment variables not set")
		return Response{}, fmt.Errorf("environment variables not set")
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(number)
	params.SetFrom(twilioNumber)
	params.SetBody(name + " has request you to be an emergency contact. Click this link below to verify your number \n " + verification_link + "sent by MamaVault")

	_, err := client.Api.CreateMessage(params) // sending message using twilio
	if err != nil {
		fmt.Println(err)
		return Response{}, err
	}
	response := Response{
		Status:  http.StatusOK,
		Message: name + " has request you to be an emergency contact. Click this link below to verify your number \n " + verification_link + "\nsent by MamaVault",
		SentTo:  number,
	}
	return response, nil
}

func VerifyMobileNumber(w http.ResponseWriter, r *http.Request) {
	// Initialize app
	client, ctx := InitializeApp()
	defer client.Close()

	w.Header().Set("Content-Type", "application/json")
	var b struct {
		Uid    string `json:"uid"`
		Number string `json:"number"`
		Name   string `json:"name"`
	}

	// checking if the request body is correct
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": "wrong body sent",
		})
		return
	}

	// adding the number to the database
	panicContactRef, _, err := client.Collection("users").Doc(b.Uid).Collection("panic_info").Add(ctx, map[string]interface{}{
		"number": b.Number,
		"status": "Waiting for confirmation",
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": string(err.Error()),
		})
		return
	}

	// updating the number_id field
	panic, err := client.Collection("users").Doc(b.Uid).Collection("panic_info").Doc(panicContactRef.ID).Update(ctx, []firestore.Update{
		{
			Path:  "number_id",
			Value: panicContactRef.ID,
		},
	})

	fmt.Printf("%v", panic)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": string(err.Error()),
		})
		return
	}
	// sending the verification link
	verificationLink := "https://mamavault.vercel.app/verify-mobile-number?panic_id=" + panicContactRef.ID + "&uid=" + b.Uid
	res, err := SendMessage(b.Number, verificationLink, b.Name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": string(err.Error()),
		})
		return
	}
	// sending the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
