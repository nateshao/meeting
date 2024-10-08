package service

import "time"

type MeetingListRequest struct {
	Page int    `json:"page" form:"page"`
	Size int    `json:"size" form:"size"`
	Name string `json:"name" form:"name"`
}

type MeetingListReply struct {
	Identity string    `json:"identity"`
	Name     string    `json:"name,omitempty"`
	BeginAt  time.Time `json:"begin_at"`
	EndAt    time.Time `json:"end_at"`
}

type MeetingCreateRequest struct {
	Name    string `json:"name,omitempty"`
	BeginAt int64  `json:"begin_at"`
	EndAt   int64  `json:"end_at"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MeetingEditRequest struct {
	Identity string `json:"identity"`
	*MeetingCreateRequest
}

type WsP2PConnectionRequest struct {
	RoomIdentity string `json:"room_identity" uri:"room_identity"`
	UserIdentity string `json:"user_identity" uri:"user_identity"`
}

type WsP2PConnectionMessage struct {
	RoomIdentity string `json:"room_identity"`
	UserIdentity string `json:"user_identity"`
	Key          string `json:"key"`
	Value        any    `json:"value"`
}
