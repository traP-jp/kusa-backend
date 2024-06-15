package main

import (
	"context"
	"fmt"
	"io"
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

type Stamp struct {
	StampId string `json:"stampId"`
	Count   int    `json:"count"`
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
}
type StampDb struct {
	TaskId  int    `db:"taskId"`
	StampId string `db:"stampId"`
	Count   int    `db:"count"`
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
	err := db.Select(&tasksFromDb, "SELECT id,content,yomi,iconUri,authorDisplayName, grade,authorName,updatedAt, citated,image FROM tasks WHERE level = ? AND isSensitive = ? ORDER BY RAND() LIMIT ?", level, isSensitive, count)
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
			UpdatedAt:         task.UpdatedAt.Format("2006/01/02 15:04"),
			Stamps:            []Stamp{},
			Citated:           task.Citated,
			Image:             task.Image,
		})
	}

	return returnTasks, nil
}
func getStampHandler(c echo.Context) error {
	stampId := c.Param("id")
	fmt.Println(stampId)
	_, r, _ := bot.API().StampApi.GetStampImage(context.Background(), stampId).Execute()

	// レスポンスヘッダ
	response := c.Response()
	response.Header().Set("Cache-Control", "no-store")
	response.Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	response.Header().Set(echo.HeaderAccessControlExposeHeaders, "Content-Disposition")
	response.Header().Set(echo.HeaderContentDisposition, "attachment; filename="+stampId)
	// ステータスコード
	response.WriteHeader(200)
	// レスポンスのライターに対して、バイナリデータをコピーする
	io.Copy(response.Writer, r.Body)
	return c.NoContent(http.StatusOK)
}
