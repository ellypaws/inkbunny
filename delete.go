package inkbunny

import (
	"net/url"
	"strconv"
)

type DeleteFileResponse struct {
	SubmissionID IntString `json:"submission_id"`
	FileID       IntString `json:"file_id"`
}

type ReorderFileResponse struct {
	SubmissionID IntString `json:"submission_id"`
	FileID       IntString `json:"file_id"`
	NewPosition  IntString `json:"newpos"`
}

func (u *User) DeleteFile(id int) (DeleteFileResponse, error) {
	if u.SID == "" {
		return DeleteFileResponse{FileID: IntString(id)}, ErrNotLoggedIn
	}
	return PostDecode[DeleteFileResponse](u.Client(), ApiUrl("delfile"), url.Values{"sid": {u.SID}, "file_id": {strconv.Itoa(id)}})
}

func (u *User) ReorderFile(id int, position int) (ReorderFileResponse, error) {
	if u.SID == "" {
		return ReorderFileResponse{FileID: IntString(id), NewPosition: IntString(position)}, ErrNotLoggedIn
	}
	values := url.Values{"sid": {u.SID}, "file_id": {strconv.Itoa(id)}, "newpos": {strconv.Itoa(position)}}
	return PostDecode[ReorderFileResponse](u.Client(), ApiUrl("reorderfile"), values)
}
