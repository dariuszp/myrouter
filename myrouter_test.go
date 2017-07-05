package myrouter

import (
	"fmt"
	"testing"
)

func createRouter() *MyRouter {
	var router = NewMyRouter("http", "example.com", 0)

	router.Add("dashboard", []string{}, "/dashboard", make(map[string]string))
	router.Add("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	router.Add("message", []string{"POST", "PUT"}, "/message/{channel}/{type}", map[string]string{"type": "error|success"})

	return router
}

func createRouterWithPort() *MyRouter {
	var router = NewMyRouter("http", "example.com", 3000)

	router.Add("dashboard", []string{}, "/dashboard", make(map[string]string))
	router.Add("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	router.Add("message", []string{"POST", "PUT"}, "/message/{channel}/{type}", map[string]string{"type": "error|success"})

	return router
}

func TestAddInvalidRoute(t *testing.T) {
	var router = createRouter()
	var success, err = router.Add("error", []string{}, "", make(map[string]string))
	if success {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

func TestRemoveRoute(t *testing.T) {
	var router = createRouter()
	router.Remove("profile")

	var route, ok = router.routes["profile"]
	if ok {
		t.Fail()
	}

	if route != nil {
		t.Fail()
	}
}

func TestMatchPathByWrongMethod(t *testing.T) {
	var router = createRouter()
	var myURL, err = router.MatchPathByMethod("POST", "/user/poltorak-dariusz")
	var route = myURL.Route
	var params = myURL.Parameters

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
	var route = router.Get("profile")

	if route.name != "profile" {
		t.Fail()
	}
}

func TestGetRouteByInvalidName(t *testing.T) {
	var router = createRouter()
	var route = router.Get("mispelled-profile")

	if route != nil {
		t.Fail()
	}
}

func TestPath(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.Path("profile", map[string][]string{"slug": []string{"poltorak-dariusz"}})

	if path != "/user/poltorak-dariusz" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestPathWithExtraParam(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.Path("profile", map[string][]string{"slug": []string{"poltorak-dariusz"}, "id": []string{"5"}})

	if path != "/user/poltorak-dariusz?id=5" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestPathWithExtraParamAndFragment(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.PathWithFragment("profile", map[string][]string{"slug": []string{"poltorak-dariusz"}, "id": []string{"5"}}, "test")

	if path != "/user/poltorak-dariusz?id=5#test" {
		fmt.Println(path)
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestPathWithMissingParams(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.Path("profile", make(map[string][]string))

	if path != "" {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

func TestRouterURL(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.URL("profile", map[string][]string{"slug": []string{"poltorak-dariusz"}})

	if path != "http://example.com/user/poltorak-dariusz" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestRouterURLWithFragment(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.URLWithFragment("profile", map[string][]string{"slug": []string{"poltorak-dariusz"}}, "test")

	if path != "http://example.com/user/poltorak-dariusz#test" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestUnsecureRouterURL(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.UnsecureURL("profile", "test:test3", map[string][]string{"slug": []string{"poltorak-dariusz"}})

	if path != "http://test:test3@example.com/user/poltorak-dariusz" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestUnsecureRouterURLWithFragment(t *testing.T) {
	var router = createRouter()

	//router.AddRoute("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
	var path, err = router.UnsecureURLWithFragment("profile", "test:test2", map[string][]string{"slug": []string{"poltorak-dariusz"}}, "test")

	if path != "http://test:test2@example.com/user/poltorak-dariusz#test" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestMatchPath(t *testing.T) {
	var router = createRouter()
	var myURL, err = router.MatchPath("/user/poltorak-dariusz")
	var route = myURL.Route
	var params = myURL.Parameters

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
	var myURL, err = router.MatchPath("/user/5")
	var route = myURL.Route
	var params = myURL.Parameters

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
	var myURL, err = router.MatchPathByMethod("GET", "/user/poltorak-dariusz")
	var route = myURL.Route
	var params = myURL.Parameters

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

func TestMatchURL(t *testing.T) {
	var router = createRouter()
	var myURL, err = router.MatchURL("http://example.com/user/poltorak-dariusz?x=1#test")
	var route = myURL.Route
	var params = myURL.Parameters

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

func TestMatchURLWithPort(t *testing.T) {
	var router = createRouterWithPort()
	var myURL, err = router.MatchURL("http://example.com:3000/user/poltorak-dariusz?x=1#test")
	var route = myURL.Route
	var params = myURL.Parameters

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
