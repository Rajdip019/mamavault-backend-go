## Panic Button API

**Path** : `/panic`

```json
{
	"uid" : "SGQyh1JTc8SK4MiNY4JnAeWO5Un2",
  "name":"Roshni", // actual name of the user
  "location_link":"google maps location link", // get the location from the user
  "location" : {
    "lat" : "", // latitte of user
    "lon : "" //longitude of user
  }
}
```

### Expected Response

```jsx
// A super big response will all the hospitals nearby 10km

```

See this for a reference response : [Reference Response](https://learn.microsoft.com/en-us/rest/api/maps/search/get-search-nearby?tabs=HTTP)

## How it works ⏬

![Screenshot 2023-08-08 at 6 26 26 PM](https://github.com/Rajdip019/mamavault-backend-go/assets/91758830/0b73a0af-78e5-4d34-a1ed-6c3312398c84)

### Step 1
Mobile App makes a req to secure AWS API Gateway which triggers a cloud function.

### Step 2 :
The cloud function then gets all the `panic_info` of the user from the database and processes it.

### Step 3:
After the data processing it sends emergency message to the emergency contacts and calls first emergency contact with a automated voice and briefs about the emergency.

### Step 4 : 
Then it takes the use latitte and longitude and user Azure Maps api to get the nearby hospital data and send back to the user.
