package autostruct

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type Animal struct {
	Name           string `json:"name" jsonschema_description:"the name of it"`
	Biome          string `json:"biome" jsonschema_description:"where biome does it live in" jsonschema:"enum=mountains,enum=plains,enum=desert,enum=forrest,enum=lake,enum=ocean"`
	LifeExpectancy int    `json:"lifeExpectancy" jsonschema_description:"how many years do they live"`
	FAQs           []struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	} `json:"faqs" jsonschema_description:"what are 2 common questions and answerts people have about it"`
}

func TestFillObject(t *testing.T) {
	Key = getKeyFromEnv()
	animal := []Animal{}
	err := Fill("salmon", &animal)
	if err != nil {
		panic(err)
	}

	res, _ := json.Marshal(animal)
	fmt.Println(string(res))
}

func TestFillArray(t *testing.T) {
	Key = getKeyFromEnv()
	animals := []Animal{}
	err := Fill("3 popular animals", &animals)
	if err != nil {
		panic(err)
	}

	if len(animals) != 3 {
		panic("there should be 3 animals")
	}

	res, _ := json.Marshal(animals)
	fmt.Println(string(res))
}

func getKeyFromEnv() string {
	return os.Getenv("OPENAI_KEY")
}
