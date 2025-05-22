package handlers

import (
	"bookstore/db"
	"fmt"

	"gorm.io/gorm"
)

func CreateStudents(database *gorm.DB, full_name, group_name string) (*db.Students, error) {
	student := db.Students{Full_name: full_name, Group_name: group_name}
	if err := database.Create(&student).Error; err != nil {
		return nil, fmt.Errorf("failed to create book: %v", err)
	}

	return &student, nil
}

func CreateOrganization(database *gorm.DB, name, address string) (*db.Organizations, error) {
	organizations := db.Organizations{Name: name, Address: address}
	if err := database.Create(&organizations).Error; err != nil {
		return nil, fmt.Errorf("failed to create book")
	}

	return &organizations, nil
}

func GetAllStudents(database *gorm.DB) (*[]db.Students, error) {
	students, err := db.GetAllStudents(database)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	return &students, nil
}

func GetStudentByName(database *gorm.DB, full_name string) (*[]db.Students, error) {
	students, err := db.GetStudentByName(database, full_name)
	if err != nil {
		return nil, fmt.Errorf("failed to get student by name")
	}

	return &students, nil
}

func DeleteStudens(database *gorm.DB, full_name string) error {
	if err := db.DeleteStudens(database, full_name); err != nil {
		return fmt.Errorf("failed to delete students: %v", err)
	}

	return nil
}

func GetAllOrganizations(database *gorm.DB) (*[]db.Organizations, error) {
	organizations, err := db.GetAllOrganizations(database)
	if err != nil {
		return nil, fmt.Errorf("failed to get all organizations: %v", err)
	}

	return &organizations, nil
}

func DeleteOrganization(database *gorm.DB, name string) error {
	if err := db.DeleteOrganization(database, name); err != nil {
		return fmt.Errorf("failed to delete organization: %v", err)
	}

	return nil
}

func GetOrganizationsByName(database *gorm.DB, name string) (*[]db.Organizations, error) {
	organization, err := db.GetOrganizationsByName(database, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get organozations by name: %v", err)
	}

	return &organization, nil
}

func GetOrganizationsByAddress(database *gorm.DB, address string) (*[]db.Organizations, error) {
	organization, err := db.GetOrganizationsByAddress(database, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get organozations by address: %v", err)
	}

	return &organization, nil
}
