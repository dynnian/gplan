package repository

import (
	"codeberg.org/dynnian/gplan/models"
)

// CreateProject inserts a new project into the database.
func (r *Repository) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

// GetProjectByID retrieves a project from the database by its ID.
func (r *Repository) GetProjectByID(id int) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetProjectByName retrieves a project from the database by its name.
func (r *Repository) GetProjectByName(name string) (*models.Project, error) {
	var project models.Project
	err := r.db.Where("name = ?", name).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetAllProjects retrieves all projects from the database.
func (r *Repository) GetAllProjects() ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.Find(&projects).Error
	return projects, err
}

// GetSubprojects retrieves all subprojects that have the given parent project ID.
func (r *Repository) GetSubprojects(parentProjectID int) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.Where("parent_project_id = ?", parentProjectID).Find(&projects).Error
	return projects, err
}

// UpdateProject updates an existing project in the database.
func (r *Repository) UpdateProject(project *models.Project) error {
	return r.db.Save(project).Error
}

// DeleteProject removes a project from the database by its ID.
func (r *Repository) DeleteProject(id int) error {
	return r.db.Delete(&models.Project{}, id).Error
}

// GetNextProjectID retrieves the next available project ID in the database.
// It selects the maximum project ID and adds 1 to determine the next available ID.
func (r *Repository) GetNextProjectID() (int, error) {
	var maxID int
	err := r.db.Model(&models.Project{}).Select("COALESCE(MAX(id), 0) + 1").Scan(&maxID).Error
	return maxID, err
}
