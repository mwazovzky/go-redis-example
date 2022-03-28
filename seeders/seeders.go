package seeders

import (
	"math/rand"
)

func getRandomArrayItem(items []string) string {
	index := rand.Intn(len(items))
	return items[index]
}
