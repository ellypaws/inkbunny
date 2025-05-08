<p align="center">
  <img src="https://inkbunny.net/images81/elephant/logo/bunny.png" width="100" />
  <img src="https://inkbunny.net/images81/elephant/logo/text.png" width="300" />
  <br>
  <h1 align="center">Inkbunny Go API</h1>
</p>

<p align="center">
  <a href="https://inkbunny.net/">
    <img alt="Inkbunny" src="https://img.shields.io/badge/website-inkbunny.net-blue">
  </a>
  <a href="https://wiki.inkbunny.net/wiki/API">
    <img alt="API" src="https://img.shields.io/badge/api-inkbunny.net-blue">
  </a>
  <a href="https://pkg.go.dev/github.com/ellypaws/inkbunny">
    <img alt="api reference" src="https://img.shields.io/badge/api-inkbunny/api-007d9c?logo=go&logoColor=white">
  </a>
  <a href="https://github.com/ellypaws/inkbunny">
    <img alt="api github" src="https://img.shields.io/badge/github-inkbunny/api-007d9c?logo=github&logoColor=white">
  </a>
  <a href="https://goreportcard.com/report/github.com/ellypaws/inkbunny">
    <img src="https://goreportcard.com/badge/github.com/ellypaws/inkbunny" alt="Go Report Card" />
  </a>
  <br>
  <a href="https://github.com/ellypaws/inkbunny/graphs/contributors">
    <img alt="Inkbunny ML contributors" src="https://img.shields.io/github/contributors/ellypaws/inkbunny">
  </a>
  <a href="https://github.com/ellypaws/inkbunny/commits/main">
    <img alt="Commit Activity" src="https://img.shields.io/github/commit-activity/m/ellypaws/inkbunny">
  </a>
  <a href="https://github.com/ellypaws/inkbunny">
    <img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/ellypaws/inkbunny?style=social">
  </a>
</p>

--------------

<p align="right"><i>Disclaimer: This project is not affiliated or endorsed by Inkbunny.</i></p>

<img src="https://go.dev/images/gophers/ladder.svg" width="48" alt="Go Gopher climbing a ladder." align="right">

## Installation

Install [Go 1.24.2](https://go.dev/dl/) or later and set up your Go environment.

```bash
go get github.com/ellypaws/inkbunny
```

## Authentication

### Login

To use the Inkbunny API, you need to authenticate and obtain a Session ID (SID). This library provides a simple way to
do this:

```go
// Create a client and login
user, err := inkbunny.Login("username", "password")
if err != nil {
    log.Fatalf("Failed to login: %v", err)
}

// The SID is stored in the user object
fmt.Println("Logged in with SID:", user.SID)

// You can also login as a guest
guestUser, err := inkbunny.Login("guest", "")
if err != nil {
    log.Fatalf("Failed to login as guest: %v", err)
}
```

> [!IMPORTANT]  
> User accounts are only accessible for login via the API if "Enable API Access" is enabled in the user's [Account
> Settings](https://inkbunny.net/account.php)

### Setting Content Ratings

Inkbunny uses a rating system to filter content. You can set which ratings you want to see:

> [!TIP]
> For guest users, rating changes only affect the current session. For registered users, rating changes in the API
> only affect the current session and are not saved to their account.

```go
// Change ratings to see General and Nudity content only
err := user.ChangeRatings(types.Ratings{
    General: &types.Yes,
    Nudity:  &types.Yes,
})
if err != nil {
    log.Fatalf("Failed to change ratings: %v", err)
}

// Alternatively, you can use the ParseMaskU function with constants
ratings := types.ParseMaskU(types.General | types.Nudity)
err = user.ChangeRatings(ratings)
if err != nil {
    log.Fatalf("Failed to change ratings: %v", err)
}
```

The available ratings are:

- `General` - Suitable for all ages
- `Nudity` - Nonsexual nudity exposing breasts or genitals
- `MildViolence` - Mild violence
- `Sexual` - Erotic imagery, sexual activity or arousal
- `StrongViolence` - Strong violence, blood, serious injury or death

## API Usage

### Searching Submissions

You can search for submissions using various parameters:

```go
// Create a search request
searchReq := inkbunny.SubmissionSearchRequest{
    SID:                "SID", // overrides the SID, otherwise uses user.SID if blank
    Text:               "fox",
    Type:               inkbunny.SubmissionTypes{inkbunny.SubmissionTypePicturePinup},
    Page:               types.IntString(1),
    SubmissionsPerPage: types.IntString(10),
    OrderBy:            types.OrderByViews,
    Random:             types.No,
    DaysLimit:          types.IntString(30),
    Keywords:           &types.Yes,  // Search in keywords (boolean)
    Title:              &types.Yes,  // Also search in titles (boolean)
    UserID:             types.IntString(12345), // Limit results to those uploaded/owned by user with this User ID.
    Username:           "artist_name", // Limit results to those uploaded/owned by user with this Username only.
    RID:                "abc123",   // Results ID for paging through results (Mode 2)
    GetRID:             types.Yes,  // Get a Results ID for this search (for Mode 2)
    PoolID:             types.IntString(789),
    CountLimit:         types.IntString(100), // Limit results to 100 submissions
    Scraps:             inkbunny.ScrapsBoth,  // Show submissions from both main and scraps galleries
}

// Perform the search
results, err := user.SearchSubmissions(searchReq)
if err != nil {
    log.Fatalf("Search failed: %v", err)
}

// Process the results
fmt.Printf("Found %d submissions total, %d on this page\n", 
    results.ResultsCountAll.Int(), 
    results.ResultsCountThisPage.Int())
for _, submission := range results.Submissions {
    fmt.Printf("Submission ID: %s, Title: %s\n", submission.SubmissionID, submission.Title)
}

// If you set GetRID: types.Yes, you can paginate through all results
// without running the search again (Mode 2)
if results.RID != "" {
    fmt.Printf("Results ID: %s (expires in %s)\n", 
        results.RID, 
        results.RIDTTLDuration)

    // Iterate through all pages
    for resp, err := range results.AllPages() {
        if err != nil {
            log.Fatalf("Failed to get page: %v", err)
        }
        fmt.Printf("Page %d of %d\n", resp.Page.Int(), resp.PagesCount.Int())
    }

    // Or iterate through all submissions across all pages
    for submissions, err := range results.AllSubmissions() {
        if err != nil {
            log.Fatalf("Failed to get submissions: %v", err)
        }
        for _, sub := range submissions {
            fmt.Printf("Submission ID: %s\n", sub.SubmissionID)
        }
    }
}
```

### Getting Submission Details

You can get detailed information about specific submissions:

```go
// Create a submission details request
detailsReq := inkbunny.SubmissionDetailsRequest{
    SID:           user.SID,
    SubmissionIDs: "123456,789012", // Comma-separated list of submission IDs
    // Alternatively, you can use SubmissionIDSlice
    SubmissionIDSlice: []string{"123456", "789012"},
    ShowDescription:             types.Yes,
    ShowDescriptionBbcodeParsed: types.Yes,
    ShowWriting:                 types.Yes,
    ShowWritingBbcodeParsed:     types.Yes,
    ShowPools:                   types.Yes,
}

// Get the details
details, err := user.SubmissionDetails(detailsReq)
if err != nil {
    log.Fatalf("Failed to get submission details: %v", err)
}

// Process the details
for _, submission := range details.Submissions {
    fmt.Printf("Title: %s\n", submission.Title)
    fmt.Printf("Description: %s\n", submission.Description)
    fmt.Printf("File URL: %s\n", submission.FileURLFull)

    // Access files
    for _, file := range submission.Files {
        fmt.Printf("File ID: %s, Name: %s\n", file.FileID, file.FileName)
        fmt.Printf("URL: %s\n", file.FileURLFull)
    }
}
```

> [!NOTE]  
> The difference between `SubmissionIDs` and `SubmissionIDSlice`:
> - `SubmissionIDs` is a comma-separated string of submission IDs
> - `SubmissionIDSlice` is a slice of strings that will be joined into `SubmissionIDs`

### Editing Submissions

You can edit submissions using the `EditSubmission` method:

```go
// Create a title and description
title := "My Updated Submission"
description := "This is an updated description for my submission."

// Create an edit request
editReq := inkbunny.SubmissionEditRequest{
    SID:          user.SID,
    SubmissionID: types.IntString("123456"),
    Title:        &title,
    Description:  &description,
    Public:       &types.Yes,
    Scraps:       &types.No,
    Keywords:     []string{"updated", "edited", "new"},
}

// Edit the submission
err := user.EditSubmission(editReq)
if err != nil {
    log.Fatalf("Failed to edit submission: %v", err)
}
```

#### Understanding Pointer Values

In the `SubmissionEditRequest` struct, many fields are pointers. This is important because it allows you to control
whether a field should be updated, cleared, or left unchanged:

1. **Not setting a field** (nil pointer): The field's current value will be preserved
2. **Setting a field to empty** (pointer to empty string): The field will be cleared
3. **Setting a field to a new value** (pointer to string): The field will be updated

Example:

```go
// Example 1: Update the title, preserve the description
newTitle := "Updated Title"
editReq := inkbunny.SubmissionEditRequest{
    SID:          user.SID,
    SubmissionID: types.IntString("123456"),
    Title:        &newTitle,     // Will update the title
    Description:  nil,           // Will preserve the current description
}

// Example 2: Update the title, clear the description
newTitle := "Updated Title"
emptyDesc := ""
editReq := inkbunny.SubmissionEditRequest{
    SID:          user.SID,
    SubmissionID: types.IntString("123456"),
    Title:        &newTitle,     // Will update the title
    Description:  &emptyDesc,    // Will clear the description
}

// Example 3: Update both title and description
newTitle := "Updated Title"
newDesc := "New description"
editReq := inkbunny.SubmissionEditRequest{
    SID:          user.SID,
    SubmissionID: types.IntString("123456"),
    Title:        &newTitle,     // Will update the title
    Description:  &newDesc,      // Will update the description
}
```

### Uploading Files

You can upload files to Inkbunny:

```go
// Open a file
file, err := os.Open("image.jpg")
if err != nil {
    log.Fatalf("Failed to open file: %v", err)
}
defer file.Close()

// Create an upload request
uploadReq := inkbunny.UploadRequest{
    SID:    user.SID,
    Notify: true,
    Files: []inkbunny.FileUpload{
        {
            MainFile: &inkbunny.FileContent{
                Name: "image.jpg",
                File: file,
            },
        },
    },
}

// Upload the file
resp, err := user.Upload(uploadReq)
if err != nil {
    log.Fatalf("Failed to upload file: %v", err)
}

fmt.Printf("Uploaded file with submission ID: %s\n", resp.SubmissionID)
```

### Deleting Submissions

You can delete submissions:

```go
// Delete a submission
err := user.DeleteSubmission("123456")
if err != nil {
    log.Fatalf("Failed to delete submission: %v", err)
}
```

---

### Broken API Methods

As of May 9, 2025, the following API methods are broken or deprecated:

- **ZIP File Upload**: The `ZipFile` field in `UploadRequest` is currently broken in the API.
- **Upload Progress**: The `ProgressKey` field in `UploadRequest` and the `UploadProgress` function are currently broken
  in the API.
- **Cancel Upload**: The `Cancel` method of `UploadResponse` is deprecated as it uses the broken `UploadProgress`
  function.

> [!WARNING]  
> Instead of using the deprecated `Cancel` method, use `UploadRequest.Context` along with `context.WithCancel` to cancel
> uploads.

### API Access Requirements

> [!IMPORTANT]  
> User accounts are only accessible for login via the API if "Enable API Access" is enabled in the user's Account
> Settings at https://inkbunny.net/account.php
