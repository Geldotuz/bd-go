package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var client *http.Client

type CatFact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

type RandomUser struct {
	Results []UserResult
}

type UserResult struct {
	Name    UserName
	Email   string
	Picture UserPicture
}

type UserName struct {
	Title string
	First string
	Last  string
}

type UserPicture struct {
	Large     string
	Medium    string
	Thumbnail string
}

func GetCatFact() {
	url := "https://catFact.ninja/fact"

	var catFact CatFact

	err := GetJson(url, &catFact)
	if err != nil {
		fmt.Printf("error getting cat fact: %s\n", err.Error())
		return
	}
	fmt.Printf("a super interesting Cat Fact: %s\n", catFact.Fact)

}

func GetRandomUser() {
	url := "https://randomUser.me/api/?inc=name,email,picture"

	var user RandomUser

	err := GetJson(url, &user)

	if err != nil {
		fmt.Printf("error getting json: %s\n", err.Error())
	} else {
		fmt.Printf("User: %s %s %s\nEmail: %s\nThumbnail: %s",
			user.Results[0].Name.Title,
			user.Results[0].Name.First,
			user.Results[0].Name.Last,
			user.Results[0].Email,
			user.Results[0].Picture.Thumbnail)
	}
}

func GetJson(url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)

}

func main() {
	client = &http.Client{}

	GetCatFact()
	GetRandomUser()
}
