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
	functions.HTTP("DeleteDoc", DeleteDoc)
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

func DeleteDoc(w http.ResponseWriter, r *http.Request) {
	// Initialize app
	firestore, ctx := InitalizeApp()
	defer firestore.Close()

	w.Header().Set("Content-Type", "application/json")
	var b struct {
		ShareDocId string `json:"share_doc_id"`
		Uid        string `json:"uid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		fmt.Fprint(w, "No uid sent")
		return
	}

	_, err := firestore.Collection("shared_docs").Doc(b.ShareDocId).Delete(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Some error occoured")
		return
	}
	_, errLink := firestore.Collection("users").Doc(b.Uid).Collection("shared_links").Doc(b.ShareDocId).Delete(ctx)
	if errLink != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Some error occoured")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Document Deleted Successfully"))
}
