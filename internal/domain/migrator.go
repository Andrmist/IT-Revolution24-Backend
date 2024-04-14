package domain

import (
	"itrevolution-backend/internal/types"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Pet{}, &Food{}, &Message{}, &PetShop{}, &FoodShop{}); err != nil {
		return err
	}

	if err := createPetShop(db); err != nil {
		return err
	}

	if err := createFoodShop(db); err != nil {
		return err
	}

	return nil
}

func createPetShop(db *gorm.DB) error {
	fish := PetShop{
		Type: types.TYPE_FISH,
		Cost: types.FISH_COST,
	}

	if err := db.Create(&fish).Error; err != nil {
		return err
	}

	snails := PetShop{
		Type: types.TYPE_SNAIL,
		Cost: types.SNAIL_COST,
	}

	if err := db.Create(&snails).Error; err != nil {
		return err
	}

	seahorse := PetShop{
		Type: types.TYPE_SEAHORSE,
		Cost: types.SEAHORSE_COST,
	}

	if err := db.Create(&seahorse).Error; err != nil {
		return err
	}

	return nil
}

func createFoodShop(db *gorm.DB) error {
	algae := FoodShop{
		Type: types.TYPE_ALGAE,
		Cost: types.ALGAE_COST,
	}

	if err := db.Create(&algae).Error; err != nil {
		return err
	}

	return nil
}

func DropDB(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&PetShop{}, &FoodShop{}); err != nil {
		return err
	}

	return nil
}
