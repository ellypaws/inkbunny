package api

import "testing"

func TestCredentials_SubmissionDetails(t *testing.T) {
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

	details, err := user.SubmissionDetails(SubmissionDetailsRequest{
		SubmissionIDs:   "14576",
		ShowDescription: Yes,
	})
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if len(details.Submissions) == 0 {
		t.Fatalf("Expected at least one submission, got none")
	}

	if details.Submissions[0].Title == "" {
		t.Errorf("Expected title to be non-empty, got empty")
	}

	if details.Submissions[0].Title != "Inkbunny Logo (Mascot Only)" {
		t.Errorf("Expected title to be Inkbunny Logo (Mascot Only), got %s", details.Submissions[0].Title)
	}

	expectedDescription := `This image is for use by people creating authorised promotional material for Inkbunny.

Make sure you view it in "Large" mode to see it with transparency ans the original PNG file. Or use this link to download it: https://inkbunny.net/files/screen/16/16930_inkbunny_inkbunnylogo_trans_rev.png`

	if details.Submissions[0].Description != expectedDescription {
		t.Errorf("Expected description to be %s, got %s", expectedDescription, details.Submissions[0].Description)
	}

	if len(details.Submissions[0].Files) == 0 {
		t.Fatalf("Expected at least one file, got none")
	}

	if details.Submissions[0].Files[0].FullFileMD5 != "146e44def7e7e325d0be6b9bad6fb27c" {
		t.Errorf("Expected md5 to be 146e44def7e7e325d0be6b9bad6fb27c, got %s", details.Submissions[0].Files[0].FullFileMD5)
	}
}
