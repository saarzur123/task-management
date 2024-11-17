package service

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/saarzur123/task-management/backend/models"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "db Suite")
}

var (
	errMock = errors.New("mock error")
)

var _ = Describe("TaskManager", func() {
	const (
		taskID1 = "1"
	)
	var (
		manager  *TaskManager
		database *sql.DB
		mockSQL  sqlmock.Sqlmock
		oldTask  = models.Task{
			Title:       "old task",
			Description: "desc",
			Status:      "pending",
		}
		task = models.Task{
			Title:       "Test Task",
			Description: "This is a test task",
			Status:      "pending",
		}
		columns = []string{"id", "title", "description", "status", "created_at"}
		err     error
	)

	BeforeEach(func() {
		database, mockSQL, err = sqlmock.New()
		Expect(err).To(Succeed())
		manager = &TaskManager{DB: database}
	})

	AfterEach(func() {
		database.Close()
	})

	Describe("Create", func() {
		It("succeeds to create new task when database is empty", func() {
			mockSQL.ExpectExec("INSERT INTO tasks").WithArgs(task.Title, task.Description, task.Status, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

			err := manager.Create(&task)
			Expect(err).To(Succeed())
			Expect(task.Id).To(Equal("1"), "defined by the database")
			Expect(task.CreatedAt).ToNot(BeZero())
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})

		It("succeeds to create new task when database is not empty", func() {
			mockSQL.ExpectExec("INSERT INTO tasks").WithArgs(oldTask.Title, oldTask.Description, oldTask.Status, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			err := manager.Create(&oldTask)
			Expect(err).To(Succeed())
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())

			mockSQL.ExpectExec("INSERT INTO tasks").WithArgs(task.Title, task.Description, task.Status, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(2, 1))
			err = manager.Create(&task)
			Expect(err).To(Succeed())
			Expect(task.Id).To(Equal("2"), "defined by the database")
			Expect(task.CreatedAt).ToNot(BeZero())
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})

		It("returns error and doesn't create new task when failed on exec", func() {
			mockSQL.ExpectExec("INSERT INTO tasks").WithArgs(task.Title, task.Description, task.Status, sqlmock.AnyArg()).WillReturnError(errMock)

			err := manager.Create(&task)
			Expect(err).To(MatchError(errMock))
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})

		It("returns error and doesn't create new task when failed on getting LastInsertId", func() {
			mockSQL.ExpectExec("INSERT INTO tasks").WithArgs(task.Title, task.Description, task.Status, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewErrorResult(errMock))

			err := manager.Create(&task)
			Expect(err).To(MatchError(errMock))
			Expect(task.Id).To(Equal("2"), "defined by the database - last inserted")
			Expect(task.CreatedAt).ToNot(BeZero())
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})
	})

	Describe("GetByID", func() {
		It("succeeds to get task by ID", func() {
			mockSQL.ExpectQuery("SELECT id, title, description, status, created_at FROM tasks").
				WithArgs(taskID1).
				WillReturnRows(sqlmock.NewRows(columns).
					AddRow(taskID1, task.Title, task.Description, task.Status, task.CreatedAt))

			resultTask, err := manager.GetByID(taskID1)
			Expect(err).To(Succeed())
			Expect(resultTask.Id).To(Equal(taskID1), "defined by the database")
			Expect(resultTask.CreatedAt).NotTo(BeZero())
			Expect(resultTask.Title).To(Equal(task.Title))
			Expect(resultTask.Description).To(Equal(task.Description))
			Expect(resultTask.Status).To(Equal(task.Status))
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})

		It("should return an error if the task is not found", func() {
			mockSQL.ExpectQuery("SELECT id, title, description, status, created_at FROM tasks").
				WithArgs(taskID1).
				WillReturnError(sql.ErrNoRows)

			task, err := manager.GetByID(taskID1)
			Expect(err).To(MatchError(sql.ErrNoRows))
			Expect(task).To(Equal(&models.Task{}))
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})
	})

	Describe("Update", func() {
		var (
			updatedTask = &models.Task{Title: task.Title, Description: task.Description, Status: task.Status, CreatedAt: oldTask.CreatedAt, Id: taskID1}
		)

		It("succeeds to update task", func() {
			// fill data
			mockSQL.ExpectExec("INSERT INTO tasks").WithArgs(oldTask.Title, oldTask.Description, oldTask.Status, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			err := manager.Create(&oldTask)
			Expect(err).To(Succeed())
			Expect(oldTask.Id).To(Equal("1"), "defined by the database")
			Expect(oldTask.CreatedAt).ToNot(BeZero())
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())

			updatedTask.CreatedAt = oldTask.CreatedAt

			// update
			mockSQL.ExpectExec(`UPDATE tasks SET title = \?, description = \?, status = \? WHERE id = \?`).
				WithArgs(updatedTask.Title, updatedTask.Description, updatedTask.Status, updatedTask.Id).
				WillReturnResult(sqlmock.NewResult(0, 1))

			err = manager.Update(updatedTask)
			Expect(err).To(Succeed())
			Expect(updatedTask.Id).To(Equal(oldTask.Id), "shouldn't be changed")
			Expect(updatedTask.CreatedAt).To(Equal(oldTask.CreatedAt), "shouldn't be changed")
			Expect(updatedTask.Title).To(Not(Equal(oldTask.Title)), "should have changed")
			Expect(updatedTask.Description).To(Not(Equal(oldTask.Title)), "should have changed")
			Expect(updatedTask.Status).To(Equal(oldTask.Status), "wasn't changed")
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})

		It("returns an error if the update query fails", func() {
			mockSQL.ExpectExec(`UPDATE tasks SET title = \?, description = \?, status = \? WHERE id = \?`).
				WithArgs(updatedTask.Title, updatedTask.Description, updatedTask.Status, updatedTask.Id).
				WillReturnError(errMock)

			err := manager.Update(updatedTask)
			Expect(err).To(MatchError(errMock))
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})

		It("returns an error when failed on getting rows affected", func() {
			mockSQL.ExpectExec(`UPDATE tasks SET title = \?, description = \?, status = \? WHERE id = \?`).
				WithArgs(updatedTask.Title, updatedTask.Description, updatedTask.Status, updatedTask.Id).
				WillReturnResult(sqlmock.NewErrorResult(errMock))

			err := manager.Update(updatedTask)
			Expect(err).To(MatchError(errMock))
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})

		It("returns an error when no rows were updated", func() {
			mockSQL.ExpectExec(`UPDATE tasks SET title = \?, description = \?, status = \? WHERE id = \?`).
				WithArgs(updatedTask.Title, updatedTask.Description, updatedTask.Status, updatedTask.Id).
				WillReturnResult(sqlmock.NewResult(0, 0))

			err := manager.Update(updatedTask)
			Expect(err).To(MatchError(ErrNotFound))
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})
	})

	Describe("Delete", func() {
		It("succeeds to delete task", func() {
			// fill data
			mockSQL.ExpectExec("INSERT INTO tasks").WithArgs(oldTask.Title, oldTask.Description, oldTask.Status, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			err := manager.Create(&oldTask)
			Expect(err).To(Succeed())
			Expect(oldTask.Id).To(Equal("1"), "defined by the database")
			Expect(oldTask.CreatedAt).ToNot(BeZero())
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())

			// delete
			mockSQL.ExpectExec(`DELETE FROM tasks WHERE id = \?`).
				WithArgs(taskID1).
				WillReturnResult(sqlmock.NewResult(0, 1))

			err = manager.Delete(taskID1)
			Expect(err).To(Succeed())
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})

		It("returns an error when fails on exec", func() {
			mockSQL.ExpectExec(`DELETE FROM tasks WHERE id = \?`).
				WithArgs(taskID1).
				WillReturnError(errMock)

			err := manager.Delete(taskID1)
			Expect(err).To(MatchError(errMock))
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})

		It("returns an error when failed on getting rows affected", func() {
			mockSQL.ExpectExec(`DELETE FROM tasks WHERE id = \?`).
				WithArgs(taskID1).
				WillReturnResult(sqlmock.NewErrorResult(errMock))

			err := manager.Delete(taskID1)
			Expect(err).To(MatchError(errMock))
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})

		It("returns an error when no rows were deleted", func() {
			mockSQL.ExpectExec(`DELETE FROM tasks WHERE id = \?`).
				WithArgs(taskID1).
				WillReturnResult(sqlmock.NewResult(0, 0))

			err := manager.Delete(taskID1)
			Expect(err).To(MatchError(ErrNotFound))
			Expect(mockSQL.ExpectationsWereMet()).To(Succeed())
		})
	})

	Describe("GetAll", func() {
		var (
			time1 = time.Now()
			time2 = time.Now()
			task1 = models.Task{"1", "Task 1", "Description 1", "pending", time1}
			task2 = models.Task{"2", "Task 2", "Description 2", "completed", time2}
			row1  = []driver.Value{"1", "Task 1", "Description 1", "pending", time1}
		)

		It("succeeds to get all tasks", func() {
			taskRows := sqlmock.NewRows(columns).
				AddRow(row1...).
				AddRow("2", "Task 2", "Description 2", "completed", time2)
			mockSQL.ExpectQuery(`SELECT id, title, description, status, created_at FROM tasks`).
				WillReturnRows(taskRows)

			tasks, err := manager.GetAll()
			Expect(err).To(Succeed())
			Expect(tasks).To(HaveLen(2))
			Expect(tasks[0]).To(Equal(task1))
			Expect(tasks[1]).To(Equal(task2))
		})

		It("returns an empty slice when no tasks exist", func() {
			mockSQL.ExpectQuery(`SELECT id, title, description, status, created_at FROM tasks`).
				WillReturnRows(sqlmock.NewRows(columns))

			tasks, err := manager.GetAll()
			Expect(err).To(Succeed())
			Expect(tasks).To(BeEmpty())
		})

		It("returns an error when fails on exec query", func() {
			mockSQL.ExpectQuery(`SELECT id, title, description, status, created_at FROM tasks`).
				WillReturnError(errMock)

			tasks, err := manager.GetAll()
			Expect(err).To(MatchError(errMock))
			Expect(tasks).To(BeNil())
		})

		It("returns an error when row scanning fails", func() {
			taskRowsFail := sqlmock.NewRows(columns).
				AddRow(row1...).
				AddRow(nil, "Task 2", "Description 2", "completed", time.Now())
			mockSQL.ExpectQuery(`SELECT id, title, description, status, created_at FROM tasks`).
				WillReturnRows(taskRowsFail)

			tasks, err := manager.GetAll()
			Expect(err).To(HaveOccurred())
			Expect(tasks).To(BeNil())
		})

		It("returns an error when rows.Err() returns an error", func() {
			taskRows := sqlmock.NewRows(columns).
				AddRow(row1...)
			mockSQL.ExpectQuery(`SELECT id, title, description, status, created_at FROM tasks`).
				WillReturnRows(taskRows).
				WillReturnError(errMock)

			tasks, err := manager.GetAll()
			Expect(err).To(MatchError(errMock))
			Expect(tasks).To(BeNil())
		})
	})

})
