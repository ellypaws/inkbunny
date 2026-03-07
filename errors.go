package inkbunny

type ErrorResponse struct {
	Code    *int   `json:"error_code,omitempty"`
	Message string `json:"error_message"`
}

func (error ErrorResponse) Error() string {
	return error.Message
}

// Error code constants matching the Inkbunny API error codes.
// See https://wiki.inkbunny.net/wiki/API#Error_Codes for reference.
const (
	ErrInvalidLogin                      = 0   // Invalid login. Username and password incorrect or account does not have API Access enabled in account Settings.
	ErrEmptySessionID                    = 1   // No Session ID sent as variable 'sid'. If this error appears then a valid session ID is required as part of the query, but it was not received by the script. Session Ids are obtained by logging in using the api_login.php Login script.
	ErrInvalidSessionID                  = 2   // Invalid Session ID sent as variable 'sid'. This error will appear if you send a Session ID (sid) that is not valid, has been logged out or has expired.
	ErrInvalidResultsID                  = 3   // Invalid Results ID sent as variable 'rid'. It contains invalid characters. Results Ids can only contain Hexadecimal values.
	ErrNoResultsFound                    = 4   // No results found for Results ID sent as variable 'rid'. The Results ID (rid) sent has either expired or is not valid. Results sets will automatically be removed after not being accessed for a period of time, or when a user has created too many results sets (in which case the oldest results sets will be removed first).
	ErrNoPermissionToUpload              = 5   // Current user does not have permission to upload files.
	ErrHourlyLimitReached                = 6   // Submission Hourly Limit Reached. This policy exists to prevent spamming and flooding. We apologise for the inconvenience. Please wait a while and try again to add more uploads.
	ErrDatabaseError                     = 7   // Database error. Unable to create a new submission.
	ErrNoValidSubmissionID               = 8   // No valid submission id given.
	ErrNoPermissionToEditSubmission      = 9   // Current user does not have permission to edit this submission.
	ErrNoPermissionToEditFile            = 10  // Current user does not have permission to edit this file.
	ErrCouldNotCreateEntry               = 11  // Could not create an entry for the file (filename) in our database. Please try again. If the problem persists, contact an administrator.
	ErrNoPermissionToBulkUpload          = 12  // User does not have permission to BULK upload multiple pages/files at once.
	ErrZIPFileTooBig                     = 13  // ZIP file is too big.
	ErrInvalidFileName                   = 14  // Incoming file names cannot contain a double-dot '..'. Please rename your file and try again.
	ErrInvalidCharactersInFilenames      = 15  // Invalid characters detected in filenames inside your zip file. The server said (error message).
	ErrCouldNotExtractFiles              = 16  // Could not extract any files from that ZIP. ZIP files cannot have subdirectories in them for Bulk Upload. Please check the zip file is not damaged and that it has files in it. If you are sure your zip file is fine, please contact an administrator and tell them about this message.
	ErrNotZIPFile                        = 17  // The file you uploaded was not a ZIP file. It was of a non-ZIP file type (type). Please try again with a valid ZIP file. If the problem persists, contact an administrator.
	ErrZIPUploadFailed                   = 18  // ZIP upload failed. You might not have provided a file, or the file may have been too big for the size restrictions. If you are sure the file is fine, then the server may be out of space or the tmp uploads directory is not writable. Please contact a system administrator and tell them about this message if you are sure the problem is on our end.
	ErrFileCouldNotBeRead                = 19  // File could not be read. File name was (file name).
	ErrFileTooLarge                      = 20  // The file you uploaded (file name) was too large in file size. Please try again with a smaller file size. If the problem persists, contact an administrator.
	ErrUnsupportedFileType               = 21  // The file you uploaded (file name) was of an unsupported file type (type). Please try again with a supported file type as listed. If the problem persists, contact an administrator.
	ErrFileTooLargeInPixelSize           = 22  // The file you uploaded (file name) was too large in pixel size (width and/or height). Please try again with dimensions not exceeding those listed. If the problem persists, contact an administrator.
	ErrFileNotInRGBOrGreyscale           = 23  // The file you uploaded (file name) was not in RGB or Greyscale color mode (most likely it was in CMYK mode). Please try again with an RGB or Greyscale image. If the problem persists, contact an administrator.
	ErrUnknownErrorCheckingFile          = 24  // There was an unknown error when trying to check your uploaded file (file name). Please check your file meets all the listed requirements and try again. If the problem persists, contact an administrator.
	ErrUnsupportedThumbnailType          = 25  // The thumbnail you uploaded (file name) was of an unsupported file type (type). Please try again with a supported thumbnail type as listed. If the problem persists, contact an administrator.
	ErrThumbnailTooLargeInPixelSize      = 26  // The thumbnail you uploaded (file name) was too large in pixel size (width and/or height). Please try again with dimensions not exceeding those listed. If the problem persists, contact an administrator.
	ErrThumbnailNotInRGBOrGreyscale      = 27  // The thumbnail you uploaded (file name) was not in RGB or Greyscale color mode (most likely it was in CMYK mode). Please try again with an RGB or Greyscale thumbnail. If the problem persists, contact an administrator.
	ErrUnknownErrorCheckingThumbnail     = 28  // There was an unknown error when trying to check your uploaded thumbnail (file name). Please check your thumbnail meets all the listed requirements and try again. If the problem persists, contact an administrator.
	ErrTooManySubmissionIDs              = 29  // If you upload all the files in this ZIP file, you will exceed the maximum limit of files/pages per submission. None of the files from your zip were added. Please upload less pages in the one zip file or start a new submission for the remaining pages.
	ErrCancellationRequest               = 30  // Received cancellation request or didn't receive response from browser within timeout limit. Some files may have been uploaded successfully. Stopped uploading on file (file name). That file and any after it in the zip were not added.
	ErrMaxAllowedNumberOfFiles           = 31  // You have reached the maximum allowed number of files/pages for this submission. Stopped uploading on file (file name). That file and any after it in the zip were not added.
	ErrCouldNotUploadThumbnail           = 32  // Could not upload the thumbnail (file name). Please try again. If the problem persists, contact an administrator.
	ErrCouldNotCreateCopyOfFile          = 33  // Could not create copy of that file in our system. Please try again. If the file is a PNG, make sure it is in RGB color mode and not Indexed color. If the problem persists, contact an administrator.
	ErrCouldNotCreateThumbnail           = 34  // Could not create a thumbnail for the file (file name). Thumbnail file was called (thumbnail file name). If the file is a PNG, make sure it is in RGB color mode and not Indexed color. Please try again. If the problem persists, contact an administrator.
	ErrUserCanceled                      = 35  // User canceled.
	ErrInvalidProgressKey                = 36  // Invalid Progress Key.
	ErrNoPermissionToDeleteSubmission    = 37  // Current user does not have permission to delete this submission.
	ErrNoPermissionToRemoveFile          = 38  // Current user does not have permission to remove this file.
	ErrNoPermissionToChangeOrderOfFile   = 39  // Current user does not have permission to change the order of this file.
	ErrSubmissionDeleted                 = 40  // That submission has been deleted.
	ErrTooManySubmissionIDsToQuery       = 41  // Too many submission ids to query. Limit exceeded.
	ErrNoPermissionToGetFavlist          = 42  // Current user does not have permission to get the favlist of this submission.
	ErrCouldNotCreateZIPtmpExtractionDir = 43  // Couldn't create zip tmp extraction dir.
	ErrCouldNotRenameFileInUnzipProcess  = 44  // Could not rename file in unzip process.
	ErrInvalidKeywordID                  = 45  // Invalid Keyword ID.
	ErrCouldNotReplaceFileOrThumbnail    = 46  // Could not replace that file or thumbnail.
	ErrRequestNotSentInHTTPSMode         = 999 // Request was not sent in HTTPS mode.
)
