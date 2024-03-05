package api

import (
	"testing"
)

func TestCredentials_SearchSubmissions(t *testing.T) {
	user, err := Guest().Login()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if user == nil {
		t.Fatalf("Expected user to not be nil, got nil")
	}

	if user.Sid == "" {
		t.Errorf("Expected sid to not be empty, got empty")
	}

	searchResponse, err := user.SearchSubmissions(SubmissionSearchRequest{
		Text:  "Inkbunny Logo (Mascot Only)",
		Title: Yes,
	})
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
