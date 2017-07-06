package myrouter

import (
	"testing"
)

func TestCreateRouter(t *testing.T) {
	// README.md code
	var router = NewMyRouter("http", "madmanlabs.com", 0)
	var router2 = NewMyRouter("http", "madmanlabs.com", 443)

	// Tests
	if router == nil {
		t.Fail()
	}

	if router2 == nil {
		t.Fail()
	}
}

func TestAddUserProfile(t *testing.T) {
	// README.md code
	var router = NewMyRouter("http", "madmanlabs.com", 0)
	router.Add("user_profile", []string{}, "/user/{id}", map[string]string{})

	// tests
	if router.Get("user_profile") == nil {
		t.Fail()
	}

	if router.Get("test") != nil {
		t.Fail()
	}
}

func TestAddUserProfileWithRequirements(t *testing.T) {
	// README.md code
	var router = NewMyRouter("http", "madmanlabs.com", 0)
	var status, err = router.Add("user_profile", []string{}, "/user/{id}", map[string]string{"id": "[1-9]+[0-9]*"})

	// Tests
	if status == nil {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestAddCustomUserProfile(t *testing.T) {
	// README.md code
	var router = NewMyRouter("http", "madmanlabs.com", 0)
	var route, err = router.AddCustom("user_profile", []string{}, "https", "mylogin:mypassword", "secure.example.com", 440, "/user/{id}", map[string]string{"id": "[1-9]+[0-9]*"})

	// Tests
	if err != nil {
		t.Fail()
	}

	if route.Name != "user_profile" {
		t.Fail()
	}
}

func TestAddRouteWithMetods(t *testing.T) {
	// README.md code
	var router = NewMyRouter("http", "madmanlabs.com", 0)
	var route, err = router.Add("user_profile", []string{"get", "post"}, "/user/{id}", map[string]string{"id": "[1-9]+[0-9]*"})

	// Tests
	if err != nil {
		t.Fail()
	}

	if route.Name != "user_profile" {
		t.Fail()
	}
}

func TestRemovingRoute(t *testing.T) {
	var router = NewMyRouter("http", "madmanlabs.com", 0)
	// Adding route
	var _, err = router.Add("user_profile", []string{"get", "post"}, "/user/{id}", map[string]string{"id": "[1-9]+[0-9]*"})
	// Removing aded route
	router.Remove("user_profile")

	// Tests
	var route = router.Get("user_profile")

	// Tests
	if err != nil {
		t.Fail()
	}

	if route != nil {
		t.Fail()
	}
}
