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
	functions.HTTP("DeleteDocs", DeleteDocs)
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

func DeleteDocs(w http.ResponseWriter, r *http.Request) {
	// Initialize app
	firestore, ctx := InitializeApp()
	defer firestore.Close()

	w.Header().Set("Content-Type", "application/json")
	var b struct {
		DocIdArr []string `json:"doc_id_arr"`
		Uid      string   `json:"uid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		fmt.Println("Wrong body sent")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": "wrong body sent",
		})
		return
	}

	for _, id := range b.DocIdArr {
		_, err := firestore.Collection("users").Doc(b.Uid).Collection("documents").Doc(id).Delete(ctx)
		if err != nil {
			fmt.Printf("An error has occurred: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": string(err.Error()),
			})
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Documents are Deleted",
	})
}
