package inkbunny

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
)

// FileUpload is one entry in UploadRequest.Files.
//
// Use MainFile for the file itself, Thumbnail for an optional custom thumbnail,
// and Replace when updating an existing file on a submission.
type FileUpload struct {
	// Replace is the existing file ID to replace, when modifying a submission.
	Replace string
	// MainFile is the primary uploaded file.
	MainFile *FileContent
	// Thumbnail is the optional uploaded thumbnail.
	Thumbnail *FileContent
}

// FileContent is the package upload wrapper around a filename and io.Reader.
type FileContent struct {
	// Name is the filename sent to Inkbunny.
	Name string
	// File is the file content stream.
	File io.Reader
}

// UploadRequest configures Client.Upload.
//
// This type intentionally uses plain Go values instead of mirroring the raw
// form expected by Inkbunny
// form exactly. File payloads are carried separately from JSON
// tags, and Context can override the client's default request context.
type UploadRequest struct {
	Context context.Context `json:"-"` // Override the context.Context used instead of the one in Client.

	// SID overrides the client's current session for this upload.
	SID string `json:"sid"`
	// SubmissionID appends to an existing submission. Leave it empty to create a new one.
	SubmissionID string `json:"submission_id,omitempty"`
	// ProgressKey enables UploadProgress when Inkbunny supports it.
	ProgressKey string `json:"progress_key,omitempty"` // Deprecated: currently broken in the API
	// Notify follows the upload endpoint's watcher-notification behavior.
	Notify BooleanYN `json:"notify"`
	// Files is the ordered list of files to upload.
	Files []FileUpload `json:"-"`
	// ZipFile uploads one ZIP archive instead of explicit Files.
	ZipFile *FileContent `json:"-"`
}

// UploadResponse is returned by Client.Upload.
type UploadResponse struct {
	// SID is the current session ID.
	SID string `json:"sid"`
	// SubmissionID is the submission created or modified by the upload.
	SubmissionID string `json:"submission_id"`

	ProgressKey *string `json:"-"`
	cancel      func() error
	client      *Client
}

// UploadProgressResponse is returned by UploadProgress.
//
// Numeric fields use IntString because Inkbunny returns string-encoded numbers.
type UploadProgressResponse struct {
	// Status describes the transfer portion of the upload.
	Status UploadStatus `json:"status"`
	// FilesCount is the total number of files being processed after transfer.
	// FilesCount is the total number of files being processed.
	FilesCount IntString `json:"filescount"`
	// FilesCompleteCount is the number of completed files.
	FilesCompleteCount IntString `json:"filescompletecount"`
	// CurFilename is the file currently being processed.
	CurFilename string `json:"curfilename"`
	// LastUserResponseEpochSecs is the last poll time, in Unix epoch seconds.
	LastUserResponseEpochSecs IntString `json:"lastuserresponse_epoch_secs"`
	// UserCancelled reports whether the user cancelled the upload.
	UserCancelled string `json:"usercancelled"`
}

// UploadStatus contains the transfer progress reported by UploadProgress.
type UploadStatus struct {
	// Total is the total upload size in bytes.
	Total IntString `json:"total"`
	// Current is the number of uploaded bytes so far.
	Current IntString `json:"current"`
	// Rate is the transfer rate in bytes per second.
	Rate IntString `json:"rate"`
	// Filename is the current file name.
	Filename FalsyString `json:"filename"`
	// Name is the multipart field name.
	Name FalsyString `json:"name"`
	// CancelUpload reports whether cancellation has been requested.
	CancelUpload IntString `json:"cancel_upload"`
	// Done reports whether the transfer portion is complete.
	Done IntString `json:"done"`
}

// DeleteSubmissionResponse is the deletion result used by UploadResponse.Delete.
type DeleteSubmissionResponse struct {
	// SubmissionID is the deleted submission ID.
	SubmissionID string `json:"submission_id"`
}

var (
	ErrAlreadyCancelled       = errors.New("upload already cancelled")
	ErrEmptySubmissionID      = errors.New("empty submission id")
	ErrUnexpectedSubmissionID = errors.New("unexpected submission id")
	ErrResponseNoSubmissionID = errors.New("no submission id after first upload")
)

func (u *UploadStatus) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '{' {
		type uploadStatus UploadStatus
		return json.Unmarshal(data, (*uploadStatus)(u))
	}

	var aux []json.RawMessage
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	fields := []any{
		&u.Total,
		&u.Current,
		&u.Rate,
		&u.Filename,
		&u.Name,
		&u.CancelUpload,
		&u.Done,
	}
	for i, v := range aux {
		if err := json.Unmarshal(v, fields[i]); err != nil {
			return err
		}
	}
	return nil
}

// Upload uploads one or more files (and optional thumbnails) to an Inkbunny submission.
// If multiple FileUpload entries are provided, each is sent in sequence, using the returned submission_id from the previous upload for subsequent calls.
// URL: https://inkbunny.net/api_upload.php
func (u *User) Upload(req UploadRequest) (UploadResponse, error) {
	if req.SID == "" {
		if u.SID == "" {
			return UploadResponse{}, ErrNotLoggedIn
		}
		req.SID = u.SID
	}
	return u.Client().Upload(req)
}

func (c *Client) Upload(req UploadRequest) (UploadResponse, error) {
	if req.SID == "" {
		return UploadResponse{}, ErrEmptySID
	}
	if len(req.Files) == 0 && req.ZipFile == nil {
		return UploadResponse{}, errors.New("no files to upload")
	}

	lastResp := UploadResponse{SID: req.SID}
	shouldSingle := req.ZipFile != nil || slices.ContainsFunc(req.Files, func(f FileUpload) bool {
		return f.Thumbnail != nil || f.Replace != ""
	})
	if shouldSingle {
		if req.ZipFile != nil {
			resp, err := uploadZip(c, req)
			if err != nil {
				return resp, fmt.Errorf("could not upload zip: %w", err)
			}
			lastResp = resp
			req.SubmissionID = resp.SubmissionID
		}
		for i := range req.Files {
			if i > 0 {
				if lastResp.SubmissionID == "" {
					return lastResp, ErrResponseNoSubmissionID
				}
				req.SubmissionID = lastResp.SubmissionID
			}
			resp, err := uploadSingle(c, req, i)
			if err != nil {
				return resp, err
			}
			lastResp = resp
		}
	} else {
		resp, err := uploadMultiple(c, req)
		if err != nil {
			return resp, err
		}
		lastResp = resp
	}
	if req.ProgressKey != "" {
		lastResp.ProgressKey = &req.ProgressKey
		lastResp.cancel = func() error {
			_, err := UploadProgress(c, req.ProgressKey, true)
			if err != nil {
				return err
			}
			return nil
		}
	}
	lastResp.client = c
	return lastResp, nil
}

func Upload(req UploadRequest) (UploadResponse, error) {
	return DefaultClient.Upload(req)
}

func (u *UploadResponse) Delete() error {
	if u.SID == "" {
		return ErrEmptySID
	}
	if u.SubmissionID == "" {
		return ErrEmptySubmissionID
	}
	response, err := PostDecode[DeleteSubmissionResponse](u.client.Get(), ApiUrl("delsubmission"), url.Values{"sid": {u.SID}, "submission_id": {u.SubmissionID}})
	if err != nil {
		return err
	}
	if u.SubmissionID != response.SubmissionID {
		return ErrUnexpectedSubmissionID
	}
	return nil
}

// Deprecated: currently the api does not work, as it uses UploadProgress, which is also broken.
// Instead, use UploadRequest.Context along with context.WithCancel to cancel the upload.
func (u *UploadResponse) Cancel() error {
	if u.cancel != nil {
		err := u.cancel()
		if err != nil {
			return err
		}
		u.cancel = nil
	}
	return ErrAlreadyCancelled
}

// uploadMultiple performs the actual multipart/form-data POST for multiple FileUploads.
// Assumes that there are no thumbnails, otherwise use uploadSingle.
// URL: https://inkbunny.net/api_upload.php
func uploadMultiple(c *Client, r UploadRequest) (UploadResponse, error) {
	endpoint := ApiUrl("upload")

	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()

	w := multipart.NewWriter(pipeWriter)
	go func() {
		var lastErr error
		defer func() { pipeWriter.CloseWithError(lastErr) }()
		err := structToMultipartWriter(w, r)
		if err != nil {
			lastErr = err
			return
		}
		for i, pair := range r.Files {
			file := pair.MainFile
			fw, err := w.CreateFormFile(fmt.Sprintf("uploadedfile[%d]", i), file.Name)
			if err != nil {
				lastErr = err
				return
			}
			if _, err := io.Copy(fw, file.File); err != nil {
				lastErr = err
				return
			}
		}
		lastErr = w.Close()
	}()

	req, err := http.NewRequestWithContext(cmp.Or(r.Context, c.ctx), http.MethodPost, endpoint.String(), pipeReader)
	if err != nil {
		return UploadResponse{}, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	httpResp, err := c.client.Do(req)
	if err != nil {
		return UploadResponse{}, err
	}

	uploadResp, err := parseResponse[UploadResponse](httpResp)
	if err != nil {
		return UploadResponse{}, err
	}
	return uploadResp, nil
}

// uploadZip performs the actual multipart/form-data POST for a UploadRequest.ZipFile.
func uploadZip(c *Client, r UploadRequest) (UploadResponse, error) {
	endpoint := ApiUrl("upload")

	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()

	w := multipart.NewWriter(pipeWriter)
	go func() {
		var lastErr error
		defer func() { pipeWriter.CloseWithError(lastErr) }()
		err := structToMultipartWriter(w, r)
		if err != nil {
			lastErr = err
			return
		}

		fw, err := w.CreateFormFile("zipfile", r.ZipFile.Name)
		if err != nil {
			lastErr = err
			return
		}
		if _, err := io.Copy(fw, r.ZipFile.File); err != nil {
			lastErr = err
			return
		}
		lastErr = w.Close()
	}()

	req, err := http.NewRequestWithContext(cmp.Or(r.Context, c.ctx), http.MethodPost, endpoint.String(), pipeReader)
	if err != nil {
		return UploadResponse{}, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	httpResp, err := c.client.Do(req)
	if err != nil {
		return UploadResponse{}, err
	}

	uploadResp, err := parseResponse[UploadResponse](httpResp)
	if err != nil {
		return UploadResponse{}, err
	}

	return uploadResp, nil
}

// uploadSingle performs the actual multipart/form-data POST for a single FileUpload.
func uploadSingle(c *Client, r UploadRequest, index int) (UploadResponse, error) {
	endpoint := ApiUrl("upload")

	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()

	w := multipart.NewWriter(pipeWriter)
	go func() {
		var lastErr error
		defer func() { pipeWriter.CloseWithError(lastErr) }()
		err := structToMultipartWriter(w, r)
		if err != nil {
			lastErr = err
			return
		}
		file := r.Files[index].MainFile
		if r.Files[index].Replace != "" {
			err := w.WriteField("replace", r.Files[index].Replace)
			if err != nil {
				lastErr = err
				return
			}
		}

		fw, err := w.CreateFormFile(fmt.Sprintf("uploadedfile[%d]", index), file.Name)
		if err != nil {
			lastErr = err
			return
		}
		if _, err := io.Copy(fw, file.File); err != nil {
			lastErr = err
			return
		}

		if thumb := r.Files[index].Thumbnail; thumb != nil {
			tw, _ := w.CreateFormFile(fmt.Sprintf("uploadedthumbnail[%d]", index), thumb.Name)
			if _, err := io.Copy(tw, thumb.File); err != nil {
				lastErr = err
				return
			}
		}
		lastErr = w.Close()
	}()

	req, err := http.NewRequestWithContext(cmp.Or(r.Context, c.ctx), http.MethodPost, endpoint.String(), pipeReader)
	if err != nil {
		return UploadResponse{}, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	httpResp, err := c.client.Do(req)
	if err != nil {
		return UploadResponse{}, err
	}

	uploadResp, err := parseResponse[UploadResponse](httpResp)
	if err != nil {
		return UploadResponse{}, err
	}

	return uploadResp, nil
}

// UploadProgress retrieves or cancels upload progress based on a unique progress key.
// Set cancel to true to instruct the server to cancel the upload in progress.
// Deprecated: currently broken in the API
func UploadProgress(c *Client, progressKey string, cancel bool) (UploadProgressResponse, error) {
	values := url.Values{"progress_key": {progressKey}}
	if cancel {
		values.Set("cancel", "yes")
	}
	httpResp, err := c.PostForm(ApiUrl("progress"), values)
	if err != nil {
		return UploadProgressResponse{}, err
	}
	return parseResponse[UploadProgressResponse](httpResp)
}
