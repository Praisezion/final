package api

import (
	"net/http"

	"final/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 50

	search := r.FormValue("search")

	var tasks []*db.Task
	var err error
	// poisk ***
	if search != "" {
		tasks, err = db.TasksWithSearch(limit, search)
	} else {
		tasks, err = db.Tasks(limit)
	}

	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
