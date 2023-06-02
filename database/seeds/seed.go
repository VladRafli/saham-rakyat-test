package seeds

import (
	"fmt"
	"sahamrakyat_test/database"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// Down
	db.Exec("TRUNCATE TABLE histories RESTART IDENTITY CASCADE;")
	db.Exec("TRUNCATE TABLE orders RESTART IDENTITY CASCADE;")
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE;")

	// Up
	// Histories
	var histories = []database.Histories{}

	for i := 0; i < 10; i++ {
		historyInterface := database.Histories{}
		err := faker.FakeData(&historyInterface)

		if err != nil {
			fmt.Println(err)
		}

		histories = append(histories, historyInterface)
	}

	db.CreateInBatches(histories, len(histories))

	// Users
	var users = []database.Users{}
	for i := 0; i < 10; i++ {
		idx := uint(i) + 1
		userInterface := database.Users{HistoriesID: &idx}
		err := faker.FakeData(&userInterface)

		if err != nil {
			fmt.Println(err)
		}

		users = append(users, userInterface)
	}

	db.CreateInBatches(users, len(users))

	// Orders
	// Note: Seeding data more than 10 data in batches is slower than seeding 10 data in 2 or more batches
	for j := 0; j < 2; j++ {
		var orders = []database.Orders{}
		for i := 0; i < 10; i++ {
			idx := uint(i) + 1
			orderInterface := database.Orders{HistoriesID: &idx}
			err := faker.FakeData(&orderInterface)

			if err != nil {
				fmt.Println(err)
			}

			orders = append(orders, orderInterface)
		}

		db.CreateInBatches(orders, len(orders))
	}
}
