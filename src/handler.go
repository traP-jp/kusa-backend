package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type TasksRequest struct {
	Level       int  `json:"level"`
	Count       int  `json:"count"`
	IsSensitive bool `json:"isSensitive"`
}

type Task struct {
	Content           string  `json:"content"`
	Yomi              string  `json:"yomi"`
	IconUri           string  `json:"iconUri"`
	AutherDisplayName string  `json:"authorDisplayName"`
	Grade             string  `json:"grade"`
	AutherName        string  `json:"authorName"`
	UpdatedAt         string  `json:"updatedAt"`
	Citated           string  `json:"citated"`
	Image             string  `json:"image"`
	Stamps            []Stamp `json:"stamps"`
}

type StampDb struct {
	TaskId  string `json:"taskId" db:"taskId"`
	StampId string `json:"stampId" db:"stampId"`
	Count   int    `json:"count" db:"count"`
}

type Stamp struct {
	StampId string `json:"stampId" db:"stampId"`
	Count   int    `json:"count" db:"count"`
}

type TaskDb struct {
	Id                int       `db:"id"`
	Content           string    `db:"content"`
	Yomi              string    `db:"yomi"`
	IconUri           string    `db:"iconUri"`
	AutherDisplayName string    `db:"authorDisplayName"`
	Grade             string    `db:"grade"`
	AutherName        string    `db:"authorName"`
	UpdatedAt         time.Time `db:"updatedAt"`
	Citated           string    `db:"citated"`
	Image             string    `db:"image"`
	MessageId         string    `db:"messageId"`
}

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func tasksHandler(c echo.Context) error {
	level, err := strconv.Atoi(c.QueryParam("level"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	count, err := strconv.Atoi(c.QueryParam("count"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	isSensitive, err := strconv.ParseBool(c.QueryParam("isSensitive"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	tasks, err := getTaskFromDb(level, count, isSensitive)
	if err != nil {
		fmt.Println("error in getTaskFromDb", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, tasks)
}

func getTaskFromDb(level int, count int, isSensitive bool) ([]Task, error) {

	returnTasks := []Task{}

	// countがDBのレコード数より多い場合は、すべてのレコードを返す
	tasksFromDb := []TaskDb{}
	err := db.Select(&tasksFromDb, "SELECT id,content,yomi,iconUri,authorDisplayName, grade,authorName,updatedAt, citated,image, messageId FROM tasks WHERE level = ? AND isSensitive = ? ORDER BY RAND() LIMIT ?", level, isSensitive, count)
	if err != nil {
		fmt.Println("error in getting tasks", err)
		return []Task{}, err
	}

	for _, task := range tasksFromDb {
		stampsFromDb := []StampDb{}
		err := db.Select(&stampsFromDb, "SELECT * FROM stamps WHERE taskId = ?", task.MessageId)
		if err != nil {
			fmt.Println("error in getting stamps", err)
			return []Task{}, err
		}

		stamps := []Stamp{}
		for _, stamp := range stampsFromDb {
			stamps = append(stamps, Stamp{
				StampId: stamp.StampId,
				Count:   stamp.Count,
			})
		}

		returnTasks = append(returnTasks, Task{
			Content:           task.Content,
			Yomi:              task.Yomi,
			IconUri:           task.IconUri,
			AutherDisplayName: task.AutherDisplayName,
			Grade:             task.Grade,
			AutherName:        task.AutherName,
			UpdatedAt:         task.UpdatedAt.Format("2006/01/02 15:04"),
			Stamps:            stamps,
			Citated:           task.Citated,
			Image:             task.Image,
		})
	}

	return returnTasks, nil
}
