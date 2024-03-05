package api

import (
	"github.com/ellypaws/inkbunny/api/utils"
	"net/url"
	"strings"
	"testing"
)

func TestInkbunnyURL(t *testing.T) {
	url := InkbunnyUrl("path", url.Values{"key": {"value"}})
	if url.String() != "https://inkbunny.net/path?key=value" {
		t.Errorf("Expected https://inkbunny.net/path?key=value, got %s", url.String())
	}

	t.Logf("Inkbunny URL: %s", url.String())
}

func TestApiURL(t *testing.T) {
	url := ApiUrl("path", url.Values{"key": {"value"}})
	if url.String() != "https://inkbunny.net/api_path.php?key=value" {
		t.Errorf("Expected https://inkbunny.net/api_path.php?key=value, got %s", url.String())
	}

	t.Logf("API URL: %s", url.String())
}

func TestApiWithStruct(t *testing.T) {
	user := Credentials{Sid: "sid", Username: "username", Ratings: Ratings{
		General:        false,
		Nudity:         false,
		MildViolence:   false,
		Sexual:         false,
		StrongViolence: false,
	}}

	url := ApiUrl("path", utils.StructToUrlValues(user))
	if url.String() != "https://inkbunny.net/api_path.php?sid=sid&username=username" {
		t.Errorf("Expected https://inkbunny.net/api_path.php?sid=sid&username=username, got %s", url.String())
	}

	t.Logf("API URL with struct: %s", url.String())
}

func TestCredentials_ChangeRating(t *testing.T) {
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

	testUser := user
	testUser.Ratings = Ratings{
		General:        true,
		Nudity:         true,
		MildViolence:   false,
		Sexual:         true,
		StrongViolence: false,
	}

	testVals := utils.StructToUrlValues(testUser)

	v := strings.Replace(testVals.Encode(), "%5B", "[", -1)
	v = strings.Replace(v, "%5D", "]", -1)
	if !strings.Contains(v, "tag[1]=yes&tag[2]=yes&tag[4]=yes&username=guest") {
		t.Errorf("Expected values to contain tag[1]=yes&tag[2]=yes&tag[4]=yes&username=guest, got %s", v)
	}

	ratings := Ratings{
		General:        true,
		Nudity:         true,
		MildViolence:   false,
		Sexual:         true,
		StrongViolence: false,
	}
	t.Logf("Changing ratings: %v", ApiUrl("userrating", utils.StructToUrlValues(ratings), url.Values{"sid": {user.Sid}}))
	err = user.ChangeRating(ratings)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if user.Ratings != ratings {
		t.Errorf("Expected ratings to be %+v, got %+v", ratings, user.Ratings)
	}

	if user.Ratings.String() != "1101" {
		t.Errorf("Expected 1101, got %s", user.Ratings.String())
	}

	t.Logf("Ratings mask correctly set to %s", user.Ratings.String())
}
