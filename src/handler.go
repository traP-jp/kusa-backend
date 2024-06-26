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

func getMeHandler(c echo.Context) error {
	var payload AuthHeader
	(&echo.DefaultBinder{}).BindHeaders(c, &payload)

	user := User{
		Name:    payload.UserId,
		IconUri: "https://q.trap.jp/api/v3/public/icon/" + payload.UserId,
	}
	return c.JSON(http.StatusOK, user)
}

func getRankingsHandler(c echo.Context) error {
	countStr := c.QueryParam("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	levelStr := c.QueryParam("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	rankingsFromDb := []RankingDb{}
	err = db.Select(&rankingsFromDb, "SELECT * FROM rankings WHERE level = ? ORDER BY score DESC LIMIT ?", level, count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	rankings := []Ranking{}
	for i, ranking := range rankingsFromDb {
		ranking := Ranking{
			Rank:      i + 1,
			UserName:  ranking.UserName,
			Score:     ranking.Score,
			Level:     ranking.Level,
			TimeStamp: ranking.TimeStamp.Format("2006/01/02 15:04"),				
			IconUri: "https://q.trap.jp/api/v3/public/icon/" + ranking.UserName,
		}
		rankings = append(rankings, ranking)
	}

	return c.JSON(http.StatusOK, rankings)
}

func postRankingsHandler(c echo.Context) error {
	fmt.Println("postRankingsHandler")
	rankingRequest := RankingsRequest{}
	err := c.Bind(&rankingRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	timeStamp, err := time.Parse("2006/01/02 15:04Z07:00", rankingRequest.TimeStamp+"+09:00")
	if err != nil {
		fmt.Println("error in parsing timeStamp", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	_, err = db.Exec("INSERT INTO rankings (userName, score, level, timeStamp) VALUES (?, ?, ?, ?)", rankingRequest.UserName, rankingRequest.Score, rankingRequest.Level, timeStamp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
