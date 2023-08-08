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
)

func init() {
	functions.HTTP("UpdateVerificationStatus", UpdateVerificationStatus)
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

func UpdateVerificationStatus(w http.ResponseWriter, r *http.Request) {
	// Initialize app
	client, ctx := InitializeApp()
	defer client.Close()

	w.Header().Set("Content-Type", "application/json")
	var b struct {
		Uid            string `json:"uid"`
		VerificationId string `json:"verification_id"`
		Action         string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		fmt.Fprint(w, "No number sent")
		return
	}
	if b.Action == "Confirm" {
		_, err := client.Collection("users").Doc(b.Uid).Collection("panic_info").Doc(b.VerificationId).Update(ctx, []firestore.Update{
			{
				Path:  "status",
				Value: "Verified",
			},
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Some error occurred")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mobile Number verified"))
	} else if b.Action == "Delete" {
		_, err := client.Collection("users").Doc(b.Uid).Collection("panic_info").Doc(b.VerificationId).Delete(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Some error occurred")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mobile Number verified"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Wrong action")
		return
	}
}
