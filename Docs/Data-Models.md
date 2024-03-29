### Collections 👇

1. `users/{uid}/documents`
2. `users/{uid}/memories`
3. `users/{uid}/shared_docs`
4. `users/{uid}/panic_info`

### Cloud storage

1. `{uid}/documents/doc_type_(date and time)`
2. `{uid}/memories/doc_name_timestamp`

## Database Models

user info. `/users/{uid}`

```json
{
	"uid" : string,
	"name" : string,
	"email" : string,
	"image" :  string,
	"age" : number,
	"blood_group" : string,
	"starting_date" : Firebase Timestamp,
	"account_created" : Firebase TimeStamp,
	"date_of_pregnancy" : "YYYY-MM-DD"
	"complications_description" : string,
	"medicines" : string[],
	"diseases" : string[],
	"allegies" : string[]
}
```

Document types.

```json
{
	"doc_id" : string, // This is the id auto generated by firebase. After saving the data gets the uid of the data given by firebase and we will update the data again with uid in the collection
	"doc_url" : string,
	"doc_type" : string,
	"doc_format" : string,
	"doc_download_url" : string,
	"upload_time" : "YYYY-MM-DD",
	"timeline_time" : "YYYY-MM-DD"
}
```

Memories type : 

```json
{
	"id" : string,
	"title" : string,
	"desctiption" : string,
	"img_url" ?: string,
	"timeline_time" : "YYYY-MM-DD"
	"upload_time" : "YYYY-MM-DD"
}
```

Shared Docs

```json
{
	"doc_id" : string, // This is the id auto generated by firebase. After saving the data gets the uid of the data given by firebase and we will update the data again with uid in the collection
	"doc_url" : string,
	"doc_type" : string,
	"doc_format" : string,
	"doc_download_url" : string,
	"starting_date" : Firebase type
}
```

Panic Info

```json
{
	"phone_number" : string,
	"number_id" : string,
	"status" : not_added | pending | approved | rejected
}
```
