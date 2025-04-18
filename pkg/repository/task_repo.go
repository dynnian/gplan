package repository

import (
	"time"

	"codeberg.org/dynnian/gplan/pkg/models"
)

func (r *Repository) CreateTask(task *models.Task) error {
	task.CreationDate = time.Now()
	task.LastUpdatedDate = time.Now()

	// Find the lowest unused ID
	var id int
	err := r.db.QueryRow(`
		SELECT COALESCE(MIN(t1.ID + 1), 1)
		FROM Tasks t1
		LEFT JOIN Tasks t2 ON t1.ID + 1 = t2.ID
		WHERE t2.ID IS NULL`).Scan(&id)
	if err != nil {
		return err
	}

	// Insert the task with the found ID
	_, err = r.db.Exec(
		`INSERT INTO Tasks (ID,
	                        Name,
							Description,
							ProjectID,
							TaskCompleted,
							DueDate,
							CompletionDate,
							CreationDate,
							LastUpdatedDate,
							Priority)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		task.Name,
		task.Description,
		task.ProjectID,
		task.TaskCompleted,
		task.DueDate,
		task.CompletionDate,
		task.CreationDate,
		task.LastUpdatedDate,
		task.Priority,
	)
	if err != nil {
		return err
	}

	task.ID = id
	return nil
}

func (r *Repository) GetTaskByID(id int) (*models.Task, error) {
	task := &models.Task{}

	err := r.db.QueryRow(`SELECT ID,
	                             Name,
								 Description,
								 ProjectID,
								 TaskCompleted,
								 DueDate,
								 CompletionDate,
								 CreationDate,
								 LastUpdatedDate,
								 Priority
						FROM Tasks
						WHERE ID = ?`, id).
		Scan(&task.ID,
			&task.Name,
			&task.Description,
			&task.ProjectID,
			&task.TaskCompleted,
			&task.DueDate,
			&task.CompletionDate,
			&task.CreationDate,
			&task.LastUpdatedDate,
			&task.Priority)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *Repository) GetAllTasks() ([]*models.Task, error) {
	rows, err := r.db.Query(
		`SELECT ID,
		        Name,
				Description,
				ProjectID,
				TaskCompleted,
				DueDate,
				CompletionDate,
				CreationDate,
				LastUpdatedDate,
				Priority
		 FROM Tasks`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		scanErr := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.ProjectID,
			&task.TaskCompleted,
			&task.DueDate,
			&task.CompletionDate,
			&task.CreationDate,
			&task.LastUpdatedDate,
			&task.Priority,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		tasks = append(tasks, task)
	}

	// Check for any error that occurred during row iteration
	if rowErr := rows.Err(); rowErr != nil {
		return nil, rowErr
	}

	return tasks, nil
}

func (r *Repository) GetTasksByProjectID(projectID int) ([]*models.Task, error) {
	rows, err := r.db.Query(
		`SELECT ID,
		        Name,
				Description,
				ProjectID,
				TaskCompleted,
				DueDate,
				CompletionDate,
				CreationDate,
				LastUpdatedDate,
				Priority
		 FROM Tasks
		 WHERE ProjectID = ?`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		scanErr := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.ProjectID,
			&task.TaskCompleted,
			&task.DueDate,
			&task.CompletionDate,
			&task.CreationDate,
			&task.LastUpdatedDate,
			&task.Priority,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		tasks = append(tasks, task)
	}

	// Check for any error that occurred during row iteration
	if rowErr := rows.Err(); rowErr != nil {
		return nil, rowErr
	}

	return tasks, nil
}

func (r *Repository) UpdateTask(task *models.Task) error {
	task.LastUpdatedDate = time.Now()
	_, err := r.db.Exec(
		`UPDATE Tasks SET
		    Name = ?,
			Description = ?,
			ProjectID = ?,
			TaskCompleted = ?,
			DueDate = ?,
			CompletionDate = ?,
			LastUpdatedDate = ?,
			Priority = ?
		 WHERE ID = ?`,
		task.Name,
		task.Description,
		task.ProjectID,
		task.TaskCompleted,
		task.DueDate,
		task.CompletionDate,
		task.LastUpdatedDate,
		task.Priority,
		task.ID,
	)
	return err
}

func (r *Repository) DeleteTask(id int) error {
	_, err := r.db.Exec(`DELETE FROM Tasks WHERE ID = ?`, id)
	return err
}
