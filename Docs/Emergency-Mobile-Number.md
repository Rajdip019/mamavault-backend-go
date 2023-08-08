## Verify Emergency Mobile number

**Path** : `/add-emergency-mobile-number`

```json
{
	"uid" : "SGQyh1JTc8SK4MiNY4JnAeWO5Un2",
	"name":"Roshni",
  "number":"+91XXXX",
}
```

### Expected Response

```jsx
{
  "status" : 200,
  "message" : "Rajdeep has request you to be an emergency contact. Click this link below to verify your number \n " + verification_link + "sent by MamaVault",
  "sentTo" : "number"
}
```

## Confirm Emergency mobile number

**This action is performed by the web-app only**

**Path** : `/update-verification-status`

```json
{
  "uid" : "JiIYn6W4EkR5tsS34xKjeCmnTrp1", // user id
  "verification_id" : "KL9wJsFxu4pJ3r28GETf", // mobile number id
  "action" : "Confirm"
}
```

### Expected Response

```jsx
{
  "status" : 200,
  "message" : "Emergency Mobile Number verified",
}
```

## Delete Emergency mobile number ( This action occurs when emergency number is rejected )

**Path** : `/update-verification-status`

```json
{
  "uid" : "JiIYn6W4EkR5tsS34xKjeCmnTrp1", // user id
  "verification_id" : "2IFnJ6KRe577H6qneNJ2", // mobile number id
  "action" : "Delete"
}
```

### Expected Response

```jsx
{
  "status" : 200,
  "message" : "Emergency Mobile Number deleted",
}
```

## How it works ⏬
![Screenshot 2023-08-08 at 5 41 01 PM](https://github.com/Rajdip019/mamavault-backend-go/assets/91758830/a83fd205-f4cb-4fa4-934c-f251c3346a99)

### Step 1
Mobile app sends a req to secure end-point AWS api gateway which triggers a cloud function.

### Step 2 
The cloud function adds the emergency contact to `users/{uid}/panic_info/{info_id}` with a "waiting" status.

### Step 3
The cloud function also generates link for validation which links to the `mamavault.vercel.app` webapp. And sends a the verification link with a small message to the number using twillo.

### Step 4 :
When the person clicks the link on the message it takes the user to the MamaVault Webapp where his number is verified and now ready for sending emegenncy messages and calls.
