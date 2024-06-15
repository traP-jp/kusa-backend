package main

import (
	"time"
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
	MessageId         string  `json:"messageId"`
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
	MessageId         string    `db:"messageId"`
}
type StampDb struct {
	TaskId  string `db:"taskId"`
	StampId string `db:"stampId"`
	Count   int    `db:"count"`
}

type User struct {
	Name    string `json:"name"`
	IconUri string `json:"iconUri"`
}

type AuthHeader struct {
	UserId string `header:"X-Showcase-User"`
}
