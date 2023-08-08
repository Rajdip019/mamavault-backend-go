## Get timeline formatted

**Path** : `/get-timeline`

```json
{
	"documents" : [{}, {}] //array of all documents full documents not id
}
```

### Expected Response

```jsx
  "documents" : [{time : number, documents : []}, {time : number, documents : []}]
```
formats the docuemnts according to their timeline .
