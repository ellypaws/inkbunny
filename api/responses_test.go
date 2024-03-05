package api

import (
	"fmt"
	"github.com/ellypaws/inkbunny/api/utils"
	"strings"
	"testing"
)

func TestRatings_String(t *testing.T) {
	r := Ratings{
		General:        true,
		Nudity:         true,
		MildViolence:   false,
		Sexual:         true,
		StrongViolence: false,
	}

	t.Logf("Testing Ratings.String: %#v", r)

	if r.String() != "1101" {
		t.Fatalf("Expected 1101, got %s", r.String())
	}

	t.Logf("Output: %s", r.String())
}

func TestRatings(t *testing.T) {
	r := "1101"

	t.Logf("Testing binary bitmask: %s", r)

	expected := Ratings{
		General:        true,
		Nudity:         true,
		MildViolence:   false,
		Sexual:         true,
		StrongViolence: false,
	}
	if parseMask(r) != expected {
		t.Errorf("Expected %+v, got %+v", expected, parseMask(r))
	}
}

func TestRatingsUrlValues(t *testing.T) {
	if BooleanYN(true).String() != "yes" {
		t.Errorf("Expected yes, got %s", BooleanYN(true).String())
	}

	r := Ratings{
		General:        true,
		Nudity:         true,
		MildViolence:   true,
		Sexual:         true,
		StrongViolence: true,
	}

	values := utils.StructToUrlValues(r)
	for i := 2; i <= 5; i++ {
		if values.Get(fmt.Sprintf("tag[%d]", i)) != "yes" {
			t.Errorf("tag[%d] expected yes, got %s", i, values.Get(fmt.Sprintf("tag[%d]", i)))
		}
	}

	user := &Credentials{Sid: "sid", Username: "username", Ratings: r}
	values = utils.StructToUrlValues(user)
	for i := 2; i <= 5; i++ {
		if values.Get(fmt.Sprintf("tag[%d]", i)) != "yes" {
			t.Errorf("tag[%d] expected yes, got %s", i, values.Get(fmt.Sprintf("tag[%d]", i)))
		}
	}

	if values.Get("sid") != "sid" {
		t.Errorf("sid expected sid, got %s", values.Get("sid"))
	}

	if values.Get("username") != "username" {
		t.Errorf("username expected username, got %s", values.Get("username"))
	}

	if values.Get("password") != "" {
		t.Errorf("password expected empty, got %s", values.Get("password"))
	}

	v := strings.Replace(values.Encode(), "%5B", "[", -1)
	v = strings.Replace(v, "%5D", "]", -1)
	if v != "sid=sid&tag[1]=yes&tag[2]=yes&tag[3]=yes&tag[4]=yes&tag[5]=yes&username=username" {
		t.Errorf("Expected values to be sid=sid&tag[1]=yes&tag[2]=yes&tag[3]=yes&tag[4]=yes&tag[5]=yes&username=username, got %s", v)
	}
}
