package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"tasksREST/pkg/service"
)

// NewHandler return a router for handling http request
func NewHandler(service *service.Service) http.Handler {
	// create router
	m := http.NewServeMux()

	// setup router path
	m.HandleFunc("/task", func(w http.ResponseWriter, req *http.Request) {
		handleNewTask(w, req, service)
	})
	m.HandleFunc("/task/", func(w http.ResponseWriter, req *http.Request) {
		handleGetTask(w, req, service)
	})
	m.HandleFunc("/tasks", func(w http.ResponseWriter, req *http.Request) {
		handleTasks(w, req, service)
	})
	return m
}

func handleNewTask(w http.ResponseWriter, req *http.Request, service *service.Service) {
	switch req.Method {
	case "POST":
		// create new task and return UUID
		UUID, err := service.NewTask()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			w.Write([]byte("\n"))
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(UUID))
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 bad request"))
	}

}

func handleGetTask(w http.ResponseWriter, req *http.Request, srv *service.Service) {
	switch req.Method {
	case "GET":
		// parse parameters
		params := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
		if len(params) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 bad request"))
			return
		}

		// get tasks
		task, err := srv.Task(params[1])
		if err != nil {
			switch {
			case err.Error() == "invalid UUID format":
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 bad request" + "\n" + err.Error()))
			case strings.Contains(err.Error(), "invalid UUID length:"):
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 bad request" + "\n" + err.Error()))
			case err == service.ErrTaskNotFound:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 page not found" + "\n" + err.Error()))
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 Internal Server Error" + "\n" + err.Error()))
			}
			return
		}

		// convert to json and write a responce
		b, err := json.Marshal(task)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error" + "\n" + err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 bad request"))
	}

}

func handleTasks(w http.ResponseWriter, req *http.Request, service *service.Service) {
	switch req.Method {
	case "GET":
		ts, err := service.AllTasks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error" + "\n" + err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		for _, t := range ts {
			w.Write([]byte(t.UUID.String()))
			w.Write([]byte(" | "))
			w.Write([]byte(t.Status))
			w.Write([]byte(" | "))
			w.Write([]byte(t.Timestamp.String()))
			w.Write([]byte("\n"))

		}
	case "DELETE":
		err := service.DeleteAllTasks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			w.Write([]byte("\n"))
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("All tasks were deleted"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 bad request"))
	}

}
