package inkbunny

import (
	"net/url"
	"strconv"
)

// DeleteFileResponse is the response from api_delfile.php.
type DeleteFileResponse struct {
	// SubmissionID is the submission the file was removed from.
	SubmissionID IntString `json:"submission_id"`
	// FileID is the deleted file ID.
	FileID IntString `json:"file_id"`
}

// ReorderFileResponse is the response from api_reorderfile.php.
type ReorderFileResponse struct {
	// SubmissionID is the submission that owns the file.
	SubmissionID IntString `json:"submission_id"`
	// FileID is the reordered file ID.
	FileID IntString `json:"file_id"`
	// NewPosition is the file's new zero-based position.
	NewPosition IntString `json:"newpos"`
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
