## Expiry System API

**Path** : `/get-shareable-link`

```json
{
  "uid" : "SGQyh1JTc8SK4MiNY4JnAeWO5Un2", // user id
  "isprofile" : true, // shares profile
  "ttl" : 60, //in seconds
  "shared_docs" : ["w7bdh36tQdVZXqXy2DcU", "UYTp0VX3jCD2mM8nnwht"] // documents to share
}
```

### Expected Response

```jsx
{
  "delete_task_path": "projects/mamavault/locations/asia-south1/queues/shared-doc-delete-queue/tasks/4339113183054972091",
  "share_doc_link": "https://mamavault.vercel.app/xt7u7SPulK3aHYkJJVYr",
  "shared_document_id": "xt7u7SPulK3aHYkJJVYr"
}
```

## Delete Shared Docs Manually

**Path** : `/delete-shared-doc`

```json
{
  "uid" : "SGQyh1JTc8SK4MiNY4JnAeWO5Un2"
  "share_doc_id" : "" //the id of the shared doc
}
```

### Expected Response

```jsx
{
  "status" : 200,
  "message" : "Shared Documents Deleted"
}
```

## How it works ⏬
![Screenshot 2023-08-08 at 5 34 25 PM](https://github.com/Rajdip019/mamavault-backend-go/assets/91758830/806bca35-83d7-47da-aa5e-29891858e619)

### Step 1
The mobile app makes q req to secure end-point of API Gateway which triggers a cloud function.

### Step 2
The cloud function duplicates the specific shared documents from the path `users/{uid}/documents` to `/shared-docs` and gets the id of the new docuemnt.

### Step 3 :
We schedule a cloud task with the id of the shared-doc and tells it to execute another cloud function `DeleteSharedDoc` which deletes the doc when triggered.
