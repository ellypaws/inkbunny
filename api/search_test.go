package api

import (
	"github.com/ellypaws/inkbunny/api/utils"
	"testing"
)

func TestCredentials_SearchSubmissions(t *testing.T) {
	guest := Guest()
	t.Logf("Logging in as guest: %v", ApiUrl("login", utils.StructToUrlValues(guest)))
	user, err := guest.Login()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if user == nil {
		t.Fatalf("Expected user to not be nil, got nil")
	}

	if user.Sid == "" {
		t.Errorf("Expected sid to not be empty, got empty")
	}

	request := SubmissionSearchRequest{
		SID:   user.Sid,
		Text:  "Inkbunny Logo (Mascot Only)",
		Title: Yes,
	}
	t.Logf("Searching for submissions: %v", ApiUrl("search", utils.StructToUrlValues(request)))
	searchResponse, err := user.SearchSubmissions(request)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if len(searchResponse.Submissions) == 0 {
		t.Fatalf("Expected at least one submission, got none")
	}

	if searchResponse.Submissions[0].Title == "" {
		t.Errorf("Expected title to be non-empty, got empty")
	}

	if searchResponse.Submissions[0].Title != "Inkbunny Logo (Mascot Only)" {
		t.Errorf("Expected title to be Inkbunny Logo (Mascot Only), got %s", searchResponse.Submissions[0].Title)
	}

	if searchResponse.Submissions[0].SubmissionID != "14576" {
		t.Errorf("Expected submission id to be 14576, got %s", searchResponse.Submissions[0].SubmissionID)
	}
}

func TestCredentials_SearchSubmissionsRandom(t *testing.T) {
	guest := Guest()
	t.Logf("Logging in as guest: %v", ApiUrl("login", utils.StructToUrlValues(guest)))
	user, err := guest.Login()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = user.ChangeRating(Ratings{
		General:        true,
		Nudity:         true,
		MildViolence:   true,
		Sexual:         true,
		StrongViolence: true,
	})

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	t.Logf("Logged in as %s, sid: %s\n", user.Username, user.Sid)

	t.Logf("Searching for submissions")
	searchResponse, err := user.SearchSubmissions(SubmissionSearchRequest{
		SubmissionIDsOnly:  true,
		SubmissionsPerPage: 5,
		Page:               1,
		Text:               "inkbunny",
		Type:               SubmissionTypes{SubmissionTypePicturePinup},
		OrderBy:            "views",
		Random:             true,
		Scraps:             "both",
	})
	if err != nil {
		t.Errorf("Error searching submissions: %v", err)
	}

	if len(searchResponse.Submissions) == 0 {
		t.Fatal("No submissions found")
	}

	var submissionIDs string
	const maxSubmissions = 5
	for i := 0; i < min(maxSubmissions, len(searchResponse.Submissions)); i++ {
		submissionIDs += searchResponse.Submissions[i].SubmissionID
		if i != min(maxSubmissions-1, len(searchResponse.Submissions)-1) {
			submissionIDs += ","
		}
	}

	if submissionIDs == "" {
		t.Fatal("No submission IDs found")
	}

	t.Logf("Getting submission details for IDs: %s", submissionIDs)
	details, err := user.SubmissionDetails(
		SubmissionDetailsRequest{
			SubmissionIDs:   submissionIDs,
			ShowDescription: Yes,
		})
	if err != nil {
		t.Fatalf("Error getting submission details: %v", err)
	}

	if len(details.Submissions) == 0 {
		t.Fatalf("Expected at least one submission, got none")
	}

	if details.Submissions[0].Title == "" {
		t.Errorf("Expected title to be non-empty, got empty")
	}
}
