# autostruct

A simple library to generate data using json schema struct tags with OpenAI functions.

## Install
Fetch the package with
```bash
go get github.com/carterharrison/autostruct
```
And import it into your programs with
```golang
import "github.com/carterharrison/autostruct"
```

## Examples
Assume we have the following struct,
```golang
type President struct {
    Name  string `json:"name"`
    Start int    `json:"start" jsonschema_description:"the year the president started office"`
}
```

Then we can fill it in,
```golang
autostruct.Key = "sk_idneuwbfibeiwb" // open ai api key
president := President{}
err := autostruct.Fill("first president of the united states", &president)
if err != nil {
    panic(err)
}

res, _ := json.Marshal(president)
fmt.Println(string(res))
```
```json
{
    "name": "George Washington",
    "startOfPresidency": 1789
}
```
We can also fill in a slice,
```golang
autostruct.Key = "sk_idneuwbfibeiwb" // open ai api key
presidents := []President{}
err := autostruct.Fill("last 10 presidents of the united states", &presidents)
if err != nil {
    panic(err)
}

res, _ := json.Marshal(presidents)
fmt.Println(string(res))
```
```json
[
    {"name": "Joe Biden", "startOfPresidency": 2021},
    {"name": "Donald Trump", "startOfPresidency": 2017},
    {"name": "Barack Obama", "startOfPresidency": 2009},
    {"name": "George W. Bush", "startOfPresidency": 2001},
    {"name": "Bill Clinton", "startOfPresidency": 1993},
    {"name": "George H.W. Bush", "startOfPresidency": 1989},
    {"name": "Ronald Reagan", "startOfPresidency": 1981},
    {"name": "Jimmy Carter", "startOfPresidency": 1977},
    {"name": "Gerald Ford", "startOfPresidency": 1974},
    {"name": "Richard Nixon" ,"startOfPresidency": 1969}
]
```
