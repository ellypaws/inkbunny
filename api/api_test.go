package api

import (
	"github.com/ellypaws/inkbunny/api/utils"
	"net/url"
	"strings"
	"testing"
)

func TestInkbunnyURL(t *testing.T) {
	url := inkbunnyURL("path", url.Values{"key": {"value"}})
	if url.String() != "https://inkbunny.net/path?key=value" {
		t.Errorf("Expected https://inkbunny.net/path?key=value, got %s", url.String())
	}
}

func TestApiURL(t *testing.T) {
	url := apiURL("path", url.Values{"key": {"value"}})
	if url.String() != "https://inkbunny.net/api_path.php?key=value" {
		t.Errorf("Expected https://inkbunny.net/api_path.php?key=value, got %s", url.String())
	}
}

func TestApiWithStruct(t *testing.T) {
	user := Credentials{Sid: "sid", Username: "username", Ratings: Ratings{
		RatingsMask:    "",
		General:        false,
		Nudity:         false,
		MildViolence:   false,
		Sexual:         false,
		StrongViolence: false,
	}}

	url := apiURL("path", utils.StructToUrlValues(user))
	if url.String() != "https://inkbunny.net/api_path.php?sid=sid&username=username" {
		t.Errorf("Expected https://inkbunny.net/api_path.php?sid=sid&username=username, got %s", url.String())
	}
}

func TestCredentials_ChangeRating(t *testing.T) {
	user, err := Guest().Login()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
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

	err = user.ChangeRating(Ratings{
		General:        true,
		Nudity:         true,
		MildViolence:   false,
		Sexual:         true,
		StrongViolence: false,
	})
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if user.Ratings.General != true {
		t.Errorf("Expected true, got %t", user.Ratings.General)
	}

	if user.Ratings.String() != "1101" {
		t.Errorf("Expected 1101, got %s", user.Ratings.String())
	}

	t.Logf("Ratings mask: %s", user.Ratings.String())
}
