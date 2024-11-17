package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/saarzur123/task-management/backend/mocks/serviceMock"
	"github.com/saarzur123/task-management/backend/models"
	"github.com/saarzur123/task-management/backend/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manager Suite")
}

var _ = Describe("TaskHandler", func() {
	var (
		mockDB           *serviceMock.MockTaskRepository
		handler          *TaskHandler
		responseRecorder *httptest.ResponseRecorder
		request          *http.Request
		testErr          error
		multipleTasks    = []models.Task{{Title: "Task 1"}, {Title: "Task 2"}}
		errMock          = errors.New("mock error")
	)

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		mockDB = serviceMock.NewMockTaskRepository(mockCtrl)
		handler = &TaskHandler{DB: mockDB}
		responseRecorder = httptest.NewRecorder()
	})

	Describe("CreateTask", func() {
		var (
			task = models.Task{Title: "Test Task"}
		)

		BeforeEach(func() {
			body, err := json.Marshal(task)
			Expect(err).To(Succeed())
			request, err = http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
			Expect(err).To(Succeed())
		})

		It("should create a task successfully", func() {
			mockDB.EXPECT().Create(&task).Return(nil)

			handler.CreateTask(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusCreated))
			var responseTask models.Task
			json.NewDecoder(responseRecorder.Body).Decode(&responseTask)
			Expect(responseTask).To(Equal(task))
		})

		It("returns 400 when failed to decode request body", func() {
			reqCreateInalid, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte("{invalid-json")))
			Expect(err).To(Succeed())

			handler.CreateTask(responseRecorder, reqCreateInalid)
			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("Invalid input"))
		})

		It("returns 500 when database error occurred", func() {
			mockDB.EXPECT().Create(&task).Return(errMock)

			handler.CreateTask(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			Expect(responseRecorder.Body.String()).To(ContainSubstring(errMock.Error()))
		})
	})

	Describe("GetAllTasks", func() {

		BeforeEach(func() {
			request, testErr = http.NewRequest("GET", "/tasks", nil)
			Expect(testErr).To(Succeed())
		})

		It("succeeds to return all tasks", func() {
			mockDB.EXPECT().GetAll().Return(multipleTasks, nil)

			handler.GetAllTasks(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			var responseTasks []models.Task
			json.NewDecoder(responseRecorder.Body).Decode(&responseTasks)
			Expect(responseTasks).To(HaveLen(2))
			Expect(responseTasks).To(ConsistOf(multipleTasks))
		})

		It("returns 500 when database error occurred", func() {
			mockDB.EXPECT().GetAll().Return(nil, errMock)
			handler.GetAllTasks(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			Expect(responseRecorder.Body.String()).To(ContainSubstring(errMock.Error()))
		})

		It("doesn't return error when no tasks were found", func() {
			mockDB.EXPECT().GetAll().Return([]models.Task{}, nil)
			handler.GetAllTasks(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			var responseTask []models.Task
			json.NewDecoder(responseRecorder.Body).Decode(&responseTask)
			Expect(responseTask).To(BeEmpty())
		})
	})

	Describe("GetTask", func() {

		BeforeEach(func() {
			request, testErr = http.NewRequest("GET", "/tasks/1", nil)
			Expect(testErr).To(Succeed())
			request = mux.SetURLVars(request, map[string]string{"id": "1"})
		})

		It("returns the task by ID", func() {
			task := models.Task{Id: "1", Title: "Task 1"}
			mockDB.EXPECT().GetByID("1").Return(&task, nil)

			handler.GetTask(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			var responseTask models.Task
			json.NewDecoder(responseRecorder.Body).Decode(&responseTask)
			Expect(responseTask).To(Equal(task))
		})

		It("returns 404 if task not found", func() {
			mockDB.EXPECT().GetByID("1").Return(nil, sql.ErrNoRows)

			handler.GetTask(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("Task not found"))
		})

		It("returns 500 when database error occurred", func() {
			mockDB.EXPECT().GetByID("1").Return(nil, errMock)
			handler.GetTask(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			Expect(responseRecorder.Body.String()).To(ContainSubstring(errMock.Error()))
		})
	})

	Describe("UpdateTask", func() {
		var (
			task = models.Task{Id: "1", Title: "Updated Task"}
		)

		BeforeEach(func() {
			body, err := json.Marshal(task)
			Expect(err).To(Succeed())
			request, err = http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
			Expect(err).To(Succeed())
			request = mux.SetURLVars(request, map[string]string{"id": "1"})
		})

		It("succeeds to update task", func() {
			mockDB.EXPECT().Update(&task).Return(nil)
			handler.UpdateTask(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			var responseTask models.Task
			json.NewDecoder(responseRecorder.Body).Decode(&responseTask)
			Expect(responseTask).To(Equal(task))
		})

		It("returns 400 failed to decode request body", func() {
			request, err := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer([]byte("{invalid-json")))
			Expect(err).To(Succeed())
			handler.UpdateTask(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("Invalid input"))
		})

		It("should return 500 when database error occurred", func() {
			mockDB.EXPECT().Update(&task).Return(errMock)

			handler.UpdateTask(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			Expect(responseRecorder.Body.String()).To(ContainSubstring(errMock.Error()))
		})

		It("should return 404 when didn't find row to update", func() {
			mockDB.EXPECT().Update(&task).Return(service.ErrNotFound)

			handler.UpdateTask(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("Task not found"))
		})
	})

	Describe("DeleteTask", func() {

		BeforeEach(func() {
			request, testErr = http.NewRequest("DELETE", "/tasks/1", nil)
			Expect(testErr).To(Succeed())
			request = mux.SetURLVars(request, map[string]string{"id": "1"})
		})

		It("succeeds to delete the task", func() {
			mockDB.EXPECT().Delete("1").Return(nil)
			handler.DeleteTask(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNoContent))
		})

		It("returns 500 when database error occurred", func() {
			mockDB.EXPECT().Delete("1").Return(errMock)
			handler.DeleteTask(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			Expect(responseRecorder.Body.String()).To(ContainSubstring(errMock.Error()))
		})

		It("should return 404 when didn't find row to delete", func() {
			mockDB.EXPECT().Delete("1").Return(service.ErrNotFound)

			handler.DeleteTask(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("Task not found"))
		})
	})
})
