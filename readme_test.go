package myrouter

import (
	"fmt"
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
	// README.md code
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

func TestReadmeMatch(t *testing.T) {
	// README.md code
	var router = NewMyRouter("http", "madmanlabs.com", 0)
	var _, err = router.Add("user_profile", []string{}, "/user/{slug}", map[string]string{"slug": "[a-z]+[a-z\\-]*"})
	var result, err2 = router.Match("GET", "http://madmanlabs.com:40/user/poltorak-dariusz?tab=contacts")

	// Tests
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if err2 != nil {
		fmt.Println(err2)
		t.Fail()
	}

	if result == nil {
		t.Fail()
		return
	}

	if false == (result.Route != nil && result.Route.Name == "user_profile") {
		fmt.Println(result.Route)
		t.Fail()
	}

	if result.Fragment != "" {
		t.Fail()
	}

	if result.Host != "madmanlabs.com" {
		t.Fail()
	}

	if result.Port != 40 {
		t.Fail()
	}

	if result.Path != "/user/poltorak-dariusz" {
		t.Fail()
	}

	if result.Parameters["slug"] != "poltorak-dariusz" {
		t.Fail()
	}

	if result.Query["tab"][0] != "contacts" {
		t.Fail()
	}
}
