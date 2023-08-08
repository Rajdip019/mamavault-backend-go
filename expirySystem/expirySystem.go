// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IDoc struct {
	DocId          string    `firestore:"doc_id"`
	DocUrl         string    `firestore:"doc_url"`
	DocType        string    `firestore:"doc_type"`
	DocFormat      string    `firestore:"doc_format"`
	DocDownloadUrl string    `firestore:"doc_download_url"`
	UploadTime     time.Time `firestore:"upload_time"`
	TimelineTime   time.Time `firestore:"timeline_time"`
}

type IUser struct {
	Name                     string    `firestore:"name"`
	Email                    string    `firestore:"email"`
	Age                      string    `firestore:"age"`
	BloodGroup               string    `firestore:"blood_group"`
	DateOfPregnancy          time.Time `firestore:"date_of_pregnancy"`
	ComplicationsDescription string    `firestore:"complications_description,omitempty"`
	Medicines                []string  `firestore:"medicines,omitempty"`
	Diseases                 []string  `firestore:"diseases,omitempty"`
	Allergies                []string  `firestore:"allergies,omitempty"`
}

type SharedDoc struct {
	UserData  IUser
	Documents []IDoc
}

func init() {
	functions.HTTP("ExpirySystem", ExpirySystem)
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

// createHTTPTask creates a new task with a HTTP target then adds it to a Queue.
func createHTTPTask(ctx context.Context, projectID, locationID, queueID, url, share_doc_id string, timer int, uid string) (*taskspb.Task, error) {

	// Create a new Cloud Tasks client instance.
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// make time
	var d time.Duration = time.Duration(timer) * time.Second
	ts := &timestamppb.Timestamp{
		Seconds: time.Now().Add(d).Unix(),
	}

	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
				},
			},
			ScheduleTime: ts,
		},
	}

	// Build the Task payload.
	payloadValues := map[string]string{"share_doc_id": share_doc_id, "uid": uid}
	jsonStr, err := json.Marshal(payloadValues)

	if err != nil {
		return nil, fmt.Errorf("unable to marshal json: %v", err)
	}

	// Add a payload message if one is present.
	req.Task.GetHttpRequest().Body = jsonStr

	// log.Println(req.Task.ScheduleTime.Seconds)
	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return createdTask, nil
}

// Duplicates the data to share an instance
func DuplicateData(firestore *firestore.Client, ctx context.Context, uid string, isProfile bool, sharedDocsId []string) (string, error) {

	// Read data
	var fetchedDocs []map[string]interface{}
	var fetchedDocsP *[]map[string]interface{} = &fetchedDocs
	var profileData map[string]interface{} = nil
	var sharedDocs []map[string]interface{} = nil
	var sharedDocsP *[]map[string]interface{} = &sharedDocs
	iter := firestore.Collection("users").Doc(html.EscapeString(uid)).Collection("documents").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "Error", err
		}
		fetchedDocs = append(fetchedDocs, doc.Data())
	}

	// Cheking is profile is shared or not
	if isProfile {
		profileSnap, err := firestore.Collection("users").Doc(html.EscapeString(uid)).Get(ctx)
		if err != nil {
			return "Error", err
		}
		profileData = profileSnap.Data()
	}

	if len(sharedDocsId) == len(fetchedDocs) {
		sharedDocsP = fetchedDocsP
	} else {
		for _, id := range sharedDocsId {
			for _, doc := range fetchedDocs {
				if id == doc["doc_id"] {
					sharedDocs = append(sharedDocs, doc)
				}
			}
		}
	}

	doc_ref, _, err := firestore.Collection("shared_docs").Add(ctx, map[string]interface{}{
		"userData":  profileData,
		"documents": *sharedDocsP,
	})
	if err != nil {
		return "Error", err
	}
	return doc_ref.ID, nil

}

func ExpirySystem(w http.ResponseWriter, r *http.Request) {
	// Initialize app
	firestore, ctx := InitializeApp()
	defer firestore.Close()

	w.Header().Set("Content-Type", "application/json")

	var b struct {
		Uid        string   `json:"uid"`
		IsProfile  bool     `json:"isProfile"`
		TTL        int      `json:"ttl"`
		SharedDocs []string `json:"shared_docs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		fmt.Fprint(w, "Wrong body sent")
		log.Fatal("Error", err)
		return
	}
	// Duplicating shared docs data
	id, err := DuplicateData(firestore, ctx, b.Uid, b.IsProfile, b.SharedDocs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error"))
		log.Fatal("Error", err)
		return
	}
	// Making share docs link
	var share_doc_link = "https://mamavault.vercel.app/" + id

	// Getting expiry time
	var d time.Duration = time.Duration(20) * time.Second
	expiryTime := time.Now().Add(d).UnixMilli()
	//  Share docs reference stored
	_, errLink := firestore.Collection("users").Doc(b.Uid).Collection("shared_links").Doc(id).Set(ctx, map[string]interface{}{
		"shared_doc_id": id,
		"shared_link":   share_doc_link,
		"expiry_time":   expiryTime,
		"views":         0,
	})
	if errLink != nil {
		return
	}

	task, err := createHTTPTask(ctx, "mamavault", "asia-south1", "shared-doc-delete-queue", "https://delete-shared-doc-s6e4vwvwlq-el.a.run.app", id, b.TTL, b.Uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error"))
		log.Fatal("Error", err)
		return
	}

	// Sending response
	responseMap := map[string]string{"shared_document_id": id, "share_doc_link": share_doc_link, "delete_task_path": task.Name}
	response, err := json.Marshal(responseMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to marshal JSON"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
