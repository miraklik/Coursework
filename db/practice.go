package db

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Students struct {
	gorm.Model
	Full_name  string `gorm:"size:100;not null" json:"full_name"`
	Group_name string `gorm:"size:50;not null" json:"group_name"`
}

type Organizations struct {
	gorm.Model
	Name    string `gorm:"size:100;not null" json:"name"`
	Address string `gorm:"size:200" json:"address"`
}

func GetStudentByName(db *gorm.DB, full_name string) ([]Students, error) {
	var student []Students

	if err := db.Where("full_name = ?", full_name).First(&student).Error; err != nil {
		log.Printf("Failed to get book by name: %v", err)
		return nil, fmt.Errorf("failed to get book by name: %v", err)
	}

	return student, nil
}

func GetOrganizationsByName(db *gorm.DB, name string) ([]Organizations, error) {
	var organization []Organizations

	if err := db.Where("name = ?", name).First(&organization).Error; err != nil {
		log.Printf("Failed to get organization by name: %v", err)
		return nil, fmt.Errorf("failed to get organization by name: %v", err)
	}

	return organization, nil
}

func GetOrganizationsByAddress(db *gorm.DB, address string) ([]Organizations, error) {
	var organization []Organizations

	if err := db.Where("address = ?", address).First(&organization).Error; err != nil {
		log.Printf("Failed to get organization by address: %v", err)
		return nil, fmt.Errorf("failed to get organization by address: %v", err)
	}

	return organization, nil
}

func GetAllStudents(db *gorm.DB) ([]Students, error) {
	var students []Students

	if err := db.Find(&students).Error; err != nil {
		log.Printf("Failed to get all books: %v", err)
		return nil, fmt.Errorf("failed to get all books: %v", err)
	}

	return students, nil
}

func GetAllOrganizations(db *gorm.DB) ([]Organizations, error) {
	var organizations []Organizations

	if err := db.Find(&organizations).Error; err != nil {
		log.Printf("Failed to get all organizations: %v", err)
		return nil, fmt.Errorf("failed to get all organizations: %v", err)
	}

	return organizations, nil
}

func DeleteStudens(db *gorm.DB, full_name string) error {

	if err := db.Unscoped().Where("full_name = ?", full_name).Delete(&Students{}).Error; err != nil {
		log.Printf("Failed to delete students: %v", err)
		return fmt.Errorf("failed to delete students: %v", err)
	}

	return nil
}

func DeleteOrganization(db *gorm.DB, name string) error {

	if err := db.Unscoped().Where("name = ?", name).Delete(&Organizations{}).Error; err != nil {
		log.Printf("Failed to delete book: %v", err)
		return fmt.Errorf("failed to delete book: %v", err)
	}

	return nil
}
