package helloworld

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func init() {
	functions.HTTP("VerifyMobileNumber", VerifyMobileNumber)
}

func InitalizeApp() (*firestore.Client, context.Context) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "mamavault-019"}
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

func VerifyMobileNumber(w http.ResponseWriter, r *http.Request) {
	// Initialize app
	client, ctx := InitalizeApp()
	defer client.Close()

	w.Header().Set("Content-Type", "application/json")
	var b struct {
		Uid    string `json:"uid"`
		Number string `json:"number"`
		Name   string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		fmt.Fprint(w, "No number sent")
		return
	}
	panicContactRef, _, err := client.Collection("users").Doc(b.Uid).Collection("panic_info").Add(ctx, map[string]interface{}{
		"number": b.Number,
		"status": "Waiting for confirmation",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Some error occoured")
		return
	}

	panic, err := client.Collection("users").Doc(b.Uid).Collection("panic_info").Doc(panicContactRef.ID).Update(ctx, []firestore.Update{
		{
			Path:  "number_id",
			Value: panicContactRef.ID,
		},
	})
	log.Print(panic)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Some error occoured")
		return
	}
	verificationLink := "https://mamavault.vercel.app/verify-mobile-number?panic_id=" + panicContactRef.ID + "&uid=" + b.Uid
	res := SendMessage(b.Number, verificationLink, b.Name)
	if res != nil {
		fmt.Println(res)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Meassage sent for verification"))
}

func SendMessage(number string, verification_link string, name string) error {
	accountSid := "AC8f632b6f0ccba9b44557fcfd3e996f18"
	authToken := "ffe6dc6e9afc38ba33573ffa8e3c5ac6"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(number)
	params.SetFrom("+15154977791")
	params.SetBody(name + " has request you to be an emergency contact. Click this link below to verify your number \n " + verification_link + "sent by MamaVault")

	res, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		response, _ := json.Marshal(*res)
		fmt.Println(string(response))
		return nil
	}
}
