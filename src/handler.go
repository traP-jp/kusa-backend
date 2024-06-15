package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TasksRequest struct {
	Level       int  `json:"level"`
	Count       int  `json:"count"`
	IsSensitive bool `json:"isSensitive"`
}

type Task struct {
	Content           string `json:"content"`
	Yomi              string `json:"yomi"`
	IconUri           string `json:"iconUri"`
	AutherDisplayName string `json:"authorDisplayName"`
	Grade             string `json:"grade"`
	AutherName        string `json:"authorName"`
	UpdatedAt         string `json:"updatedAt"`
	KusaCount         int    `json:"kusaCount"`
}

type TaskDb struct {
	Id                int    `db:"id"`
	Content           string `db:"content"`
	Yomi              string `db:"yomi"`
	IconUri           string `db:"iconUri"`
	AutherDisplayName string `db:"authorDisplayName"`
	Grade             string `db:"grade"`
	AutherName        string `db:"authorName"`
	UpdatedAt         string `db:"updatedAt"`
	KusaCount         int    `db:"kusaCount"`
	Level             int    `db:"level"`
	IsSensitive       bool   `db:"isSensitive"`
}

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func tasksHandler(c echo.Context) error {
	tasksRequest := &TasksRequest{}
	err := c.Bind(tasksRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	tasks, err := getTaskFromDb(tasksRequest.Level, tasksRequest.Count, tasksRequest.IsSensitive)
	if err != nil {
		fmt.Println("error in getTaskFromDb", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, tasks)
}

func getTaskFromDb(level int, count int, isSensitive bool) ([]Task, error) {
	var dbContentsCount int
	err := db.Get(&dbContentsCount, "SELECT COUNT(*) FROM tasks WHERE level = ? AND isSensitive = ?", level, isSensitive)
	if err != nil {
		return []Task{}, err
	}

	returnTasks := []Task{}

	// countがDBのレコード数より多い場合は、すべてのレコードを返す
	tasksFromDb := []TaskDb{}
	err = db.Select(&tasksFromDb, "SELECT * FROM tasks WHERE level = ? AND isSensitive = ? ORDER BY RAND() LIMIT ?", level, isSensitive, count)
	if err != nil {
		return []Task{}, err
	}
	for _, task := range tasksFromDb {
		returnTasks = append(returnTasks, Task{
			Content:           task.Content,
			Yomi:              task.Yomi,
			IconUri:           task.IconUri,
			AutherDisplayName: task.AutherDisplayName,
			Grade:             task.Grade,
			AutherName:        task.AutherName,
			UpdatedAt:         task.UpdatedAt,
			KusaCount:         task.KusaCount,
		})
	}

	return returnTasks, nil
}
