package inkbunny

import (
	"net/url"
	"strconv"

	"github.com/ellypaws/inkbunny/types"
)

type DeleteFileResponse struct {
	SubmissionID types.IntString `json:"submission_id"`
	FileID       types.IntString `json:"file_id"`
}

type ReorderFileResponse struct {
	SubmissionID types.IntString `json:"submission_id"`
	FileID       types.IntString `json:"file_id"`
	NewPosition  types.IntString `json:"new_position"`
}

func (u *User) DeleteFile(id int) (DeleteFileResponse, error) {
	if u.SID == "" {
		return DeleteFileResponse{FileID: types.IntString(id)}, ErrNotLoggedIn
	}
	return PostDecode[DeleteFileResponse](u.Client(), ApiUrl("delfile"), url.Values{"sid": {u.SID}, "file_id": {strconv.Itoa(id)}})
}

func (u *User) ReorderFile(id int, position int) (ReorderFileResponse, error) {
	if u.SID == "" {
		return ReorderFileResponse{FileID: types.IntString(id), NewPosition: types.IntString(position)}, ErrNotLoggedIn
	}
	values := url.Values{"sid": {u.SID}, "file_id": {strconv.Itoa(id)}, "newpos": {strconv.Itoa(position)}}
	return PostDecode[ReorderFileResponse](u.Client(), ApiUrl("reorderfile"), values)
}
