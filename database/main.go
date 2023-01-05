package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type User struct {
	ChatID int64 `json:"id"`
	Hero   Hero  `json:"hero"`
}

type Hero struct {
	Url   string   `json:"url"`
	Roles []string `json:"roles"`
	Name  string   `json:"name"`
}

var path string

func init() {
	path = "./database.json"
}

func Write(chatid int64, url, name string, roles []string) {
	var users []User

	file, _ := os.ReadFile(path)
	if len(file) != 0 {
		err := json.Unmarshal(file, &users)
		if err != nil {
			log.Fatal(err)
		}
		for index, user := range users {
			if user.ChatID == chatid {
				users = remove(users, index)
			}
		}
	}
	chosenHero := Hero{url, roles, name}
	users = append(users, User{chatid, chosenHero})
	data, err := json.Marshal(users)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(path, data, 1)

}

func Get(chatid int64) (Hero, error) {
	var users []User
	file, _ := os.ReadFile(path)
	err := json.Unmarshal(file, &users)
	if err != nil {
		return Hero{}, err
	}
	for _, user := range users {
		if user.ChatID == chatid {
			return user.Hero, nil

		}
	}
	return Hero{}, errors.New("HERO NOT CHOSEN")
}

// ебучий голанг заставляет писать меня ФУНКЦИИ КОТОРЫЕ ДОЛЖНЫ БЫТЬ ВСТРОЕННЫМИ
func remove(slice []User, s int) []User {
	return append(slice[:s], slice[s+1:]...)
}
