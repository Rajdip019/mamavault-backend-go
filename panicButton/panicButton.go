package helloWorld

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
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("PanicButton", PanicButton)
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

func SendMessage(number string, location_link string, name string) int {
	accountSid := ""
	authToken := ""

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(number)
	params.SetFrom("+15154977791")
	params.SetBody(name + " is in problem at " + location_link)

	res, err := client.Api.CreateMessage(params)
	if err != nil {
		return http.StatusInternalServerError
	} else {
		response, _ := json.Marshal(*res)
		fmt.Println(string(response))
		return http.StatusAccepted
	}

}

func MakeCall(number string) int {

	accountSid := ""
	authToken := ""

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	callParams := &twilioApi.CreateCallParams{}
	callParams.SetTo(number)
	callParams.SetFrom("+15154977791")
	callParams.SetUrl("http://twimlets.com/holdmusic?Bucket=com.twilio.music.ambient")

	resp, err := client.Api.CreateCall(callParams)
	if err != nil {
		return http.StatusInternalServerError
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println(string(response))
		return http.StatusAccepted
	}

}

func FetchMobileNumbers(uid string) ([]string, error) {
	// Initialize app
	firebase, ctx := InitalizeApp()
	defer firebase.Close()

	var mobileNumbersUnfiltered []map[string]interface{} = nil
	var mobileNumber []string = nil

	iter := firebase.Collection("users").Doc(uid).Collection("panic_info").Where("status", "==", "Verified").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil
		}
		mobileNumbersUnfiltered = append(mobileNumbersUnfiltered, doc.Data())
	}
	for _, data := range mobileNumbersUnfiltered {
		var new string = fmt.Sprintf("%v", data["number"])
		mobileNumber = append(mobileNumber, new)
	}
	return mobileNumber, nil
}

func PanicButton(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var b struct {
		Uid          string `json:"uid"`
		Name         string `json:"name"`
		LocationLink string `json:"location_link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		fmt.Fprint(w, "body is invalid")
		return
	}
	// Initialize app
	firebase, ctx := InitalizeApp()
	defer firebase.Close()

	var mobileNumbersUnfiltered []map[string]interface{} = nil
	var mobileNumbers []string = nil

	iter := firebase.Collection("users").Doc(b.Uid).Collection("panic_info").Where("status", "==", "Verified").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}
		mobileNumbersUnfiltered = append(mobileNumbersUnfiltered, doc.Data())
	}
	for _, data := range mobileNumbersUnfiltered {
		var new string = fmt.Sprintf("%v", data["number"])
		mobileNumbers = append(mobileNumbers, new)
	}
	for _, num := range mobileNumbers {
		messageRes := SendMessage(num, b.LocationLink, b.Name)
		callRes := MakeCall(num)
		if messageRes == 500 || callRes == 500 {
			return
		}
	}

}
