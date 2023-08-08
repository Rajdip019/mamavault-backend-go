package helloWorld

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"google.golang.org/api/iterator"
)

// setting environment variables
var azureAPIKey = os.Getenv("AZURE_API_KEY")

func init() {
	functions.HTTP("PanicButton", PanicButton)
}

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

func SendMessage(number string, location_link string, name string) int {
	accountSid := ""
	authToken := ""

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(number)
	params.SetFrom("+13253357019")
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
	callParams.SetFrom("+13253357019")
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
	firebase, ctx := InitializeApp()
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

func GetNearbyHospitals(lat string, lon string) ([]byte, error) {
	resp, err := http.Get("https://atlas.microsoft.com/search/nearby/json?subscription-key=" + azureAPIKey + "&api-version=1.0&lat=" + lat + "&lon=" + lon + "&categorySet=7321002&radius=10000")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return body, nil
}

func PanicButton(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var b struct {
		Uid          string `json:"uid"`
		Name         string `json:"name"`
		LocationLink string `json:"location_link"`
		Location     struct {
			Lat string `json:"lat"`
			Lon string `json:"lon"`
		} `json:"location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		fmt.Fprint(w, "body is invalid")
		return
	}
	// Initialize app
	firebase, ctx := InitializeApp()
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
		if messageRes == 500 {
			fmt.Println("Error while sending message")
		}
		if callRes == 500 {
			fmt.Println("Error while making call")
		}
	}
	res, err := GetNearbyHospitals(b.Location.Lat, b.Location.Lon)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Some error occurred while getting nearby hospitals"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
