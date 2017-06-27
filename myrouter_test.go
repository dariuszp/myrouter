package myrouter

import (
	"fmt"
	"testing"
)

func createRouter() *MyRouter {
	var router = NewMyRouter("http", "example.com", 0)

	router.AddRoute("dashboard", []string{}, "/dashboard", make(map[string]string))
	router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	router.AddRoute("message", []string{"POST", "PUT"}, "/message/{channel}/{type}", map[string]string{"type": "error|success"})

	return router
}

func TestAddInvalidRoute(t *testing.T) {
	var router = createRouter()
	var success, err = router.AddRoute("error", []string{}, "", make(map[string]string))
	if success {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

func TestRemoveRoute(t *testing.T) {
	var router = createRouter()
	router.RemoveRoute("profile")

	var route, ok = router.routes["profile"]
	if ok {
		t.Fail()
	}

	if route != nil {
		t.Fail()
	}
}

func TestMatchPath(t *testing.T) {
	var router = createRouter()
	var route, params, err = router.MatchPath("/user/poltorak-dariusz")
	if route.name != "profile" {
		t.Fail()
	}

	if len(params) > 1 {
		t.Fail()
	}

	if params["slug"] != "poltorak-dariusz" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestMatchInvalidPath(t *testing.T) {
	var router = createRouter()
	var route, params, err = router.MatchPath("/user/5")
	if route != nil {
		t.Fail()
	}

	if len(params) > 0 {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

func TestMatchPathByMethod(t *testing.T) {
	var router = createRouter()
	var route, params, err = router.MatchPathByMethod("GET", "/user/poltorak-dariusz")
	if route.name != "profile" {
		t.Fail()
	}

	if len(params) > 1 {
		t.Fail()
	}

	if params["slug"] != "poltorak-dariusz" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestMatchPathByWrongMethod(t *testing.T) {
	var router = createRouter()
	var route, params, err = router.MatchPathByMethod("POST", "/user/poltorak-dariusz")

	if route != nil {
		t.Fail()
	}

	if len(params) > 0 {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

func TestGetRouteByName(t *testing.T) {
	var router = createRouter()
	var route = router.GetRouteByName("profile")

	if route.name != "profile" {
		t.Fail()
	}
}

func TestGetRouteByInvalidName(t *testing.T) {
	var router = createRouter()
	var route = router.GetRouteByName("mispelled-profile")

	if route != nil {
		t.Fail()
	}
}

func TestGeneratePath(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.GeneratePath("profile", map[string]string{"slug": "poltorak-dariusz"})

	if path != "/user/poltorak-dariusz" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestGeneratePathWithExtraParam(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.GeneratePath("profile", map[string]string{"slug": "poltorak-dariusz", "id": "5"})

	if path != "/user/poltorak-dariusz?id=5" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestGeneratePathWithExtraParamAndFragment(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.GeneratePathWithFragment("profile", map[string]string{"slug": "poltorak-dariusz", "id": "5"}, "test")

	if path != "/user/poltorak-dariusz?id=5#test" {
		fmt.Println(path)
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestGeneratePathWithMissingParams(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.GeneratePath("profile", make(map[string]string))

	if path != "" {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

func TestGenerateRouterURL(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.GenerateURL("profile", map[string]string{"slug": "poltorak-dariusz"})

	if path != "http://example.com/user/poltorak-dariusz" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestGenerateRouterURLWithFragment(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.GenerateURLWithFragment("profile", map[string]string{"slug": "poltorak-dariusz"}, "test")

	if path != "http://example.com/user/poltorak-dariusz#test" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestGenerateUnsecureRouterURL(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.GenerateUnsecureURL("profile", "test", "test2", map[string]string{"slug": "poltorak-dariusz"})

	if path != "http://test:test2@example.com/user/poltorak-dariusz" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestGenerateUnsecureRouterURLWithFragment(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.GenerateUnsecureURLWithFragment("profile", "test", "test2", map[string]string{"slug": "poltorak-dariusz"}, "test")

	if path != "http://test:test2@example.com/user/poltorak-dariusz#test" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}
