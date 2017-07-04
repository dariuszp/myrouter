package myrouter

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func testRegExpFromPathVerification(t *testing.T, path string, url string, expect string, requirements map[string]*regexp.Regexp) {
	var regexp, err = generateRegexpFromPath(path, requirements)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if regexp.String() != expect {
		fmt.Print(regexp.String())
		fmt.Print("\n")
		t.Fail()
	}

	if !regexp.MatchString(url) {
		t.Fail()
	}
}

func testRegExpFromPathValueVerification(t *testing.T, path string, url string, expect []string, requirements map[string]*regexp.Regexp) {
	var regexp, err = generateRegexpFromPath(path, requirements)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	var result = regexp.FindAllStringSubmatch(url, -1)
	var parameters []string
	for i := 1; i < len(result[0]); i++ {
		parameters = append(parameters, result[0][i])
	}

	if !regexp.MatchString(url) {
		t.Fail()
	}

	if !arrayCompareString(expect, parameters) {
		fmt.Print(strings.Join(parameters, ", "))
		fmt.Print("\n")
		t.Fail()
	}
}

func TestEmptyPath(t *testing.T) {
	var path, err = generatePath("", make(map[string][]string), make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}
	if path != "" {
		t.Fail()
	}
}

func TestPathWithMissingParameter(t *testing.T) {
	var path, err = generatePath("/{id}", make(map[string][]string), make(map[string]*regexp.Regexp))
	if err == nil {
		t.Fail()
	}

	if path != "" {
		t.Fail()
	}
}

func TestPathWithParameter(t *testing.T) {
	var path, err = generatePath("/{id}", map[string][]string{"id": []string{"test"}}, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}

	if path != "/test" {
		t.Fail()
	}
}

func TestPathWithExtraParameter(t *testing.T) {
	var path, err = generatePath("/{id}", map[string][]string{"id": []string{"test"}, "slug": []string{"poltorak-dariusz"}}, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}

	if path != "/test?slug=poltorak-dariusz" {
		t.Fail()
	}
}

func TestPathWithTwoParameters(t *testing.T) {
	var path, err = generatePath("/{id}/{slug}", map[string][]string{"id": []string{"test"}, "slug": []string{"poltorak-dariusz"}}, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}

	if path != "/test/poltorak-dariusz" {
		t.Fail()
	}
}

func TestRegExpFromPathNoArgs(t *testing.T) {
	var path = "/api/user"
	var url = "/api/user"
	var expect = "/api/user"

	testRegExpFromPathVerification(t, path, url, expect, map[string]*regexp.Regexp{})
}

func TestRegExpFromPathSingleArg(t *testing.T) {
	var path = "/api/user/{id}"
	var url = "/api/user/5"
	var expect = "/api/user/([^/]+)"

	testRegExpFromPathVerification(t, path, url, expect, map[string]*regexp.Regexp{})
}

func TestRegExpFromPathSingleCopiedArg(t *testing.T) {
	var path = "/api/user/{id}/user-{id}"
	var url = "/api/user/5/user-5"
	var expect = "/api/user/([^/]+)/user-([^/]+)"

	testRegExpFromPathVerification(t, path, url, expect, map[string]*regexp.Regexp{})
}

func TestRegExpFromPathTwoArgs(t *testing.T) {
	var path = "/api/user/{id}/client-{slug}"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = "/api/user/([^/]+)/client-([^/]+)"

	testRegExpFromPathVerification(t, path, url, expect, map[string]*regexp.Regexp{})
}

func TestRegExpFromPathNoValues(t *testing.T) {
	var path = "/api/user/5/client-poltorak-dariusz"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = []string{}

	testRegExpFromPathValueVerification(t, path, url, expect, map[string]*regexp.Regexp{})
}

func TestRegExpFromPathOneValue(t *testing.T) {
	var path = "/api/user/{id}"
	var url = "/api/user/5"
	var expect = []string{"5"}

	testRegExpFromPathValueVerification(t, path, url, expect, map[string]*regexp.Regexp{})
}

func TestRegExpFromPathTwoValues(t *testing.T) {
	var path = "/api/user/{id}/client-{slug}"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = []string{"5", "poltorak-dariusz"}

	testRegExpFromPathValueVerification(t, path, url, expect, map[string]*regexp.Regexp{})
}

func TestRegExpFromPathWithRequirement(t *testing.T) {
	var path = "/api/user/{id}/client-{slug}"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = "/api/user/([^/]+)/client-([a-z\\-]+)"
	var requirements = map[string]*regexp.Regexp{"slug": regexp.MustCompile("[a-z\\-]+")}

	testRegExpFromPathVerification(t, path, url, expect, requirements)
}

func TestRegExpFromPathWithTwoRequirements(t *testing.T) {
	var path = "/api/user/{id}/client-{slug}"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = "/api/user/([1-9]+[0-9]*)/client-([a-z\\-]+)"
	var requirements = map[string]*regexp.Regexp{"id": regexp.MustCompile("[1-9]+[0-9]*"), "slug": regexp.MustCompile("[a-z\\-]+")}

	testRegExpFromPathVerification(t, path, url, expect, requirements)
}

func TestUrl(t *testing.T) {
	var expect = "https://example.com/api/user/5"
	var scheme = "https"
	var host = "example.com"
	var port int
	var path = "/api/user/{id}"
	var parameters = map[string][]string{"id": []string{"5"}}

	var url, err = generateURL(scheme, "", host, port, path, parameters, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}
	if url != expect {
		t.Fail()
	}
}

func TestUrlTwoParameters(t *testing.T) {
	var expect = "https://example.com/api/user/5/dariusz-poltorak"
	var scheme = "https"
	var host = "example.com"
	var port int
	var path = "/api/user/{id}/{slug}"
	var parameters = map[string][]string{"slug": []string{"dariusz-poltorak"}, "id": []string{"5"}}

	var url, err = generateURL(scheme, "", host, port, path, parameters, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}
	if url != expect {
		t.Fail()
	}
}

func TestUrlMissingParameters(t *testing.T) {
	var expect = "https://example.com/api/user/5/dariusz-poltorak"
	var scheme = "https"
	var host = "example.com"
	var port int
	var path = "/api/user/{id}/{slug}"
	var parameters = map[string][]string{"id": []string{"5"}}

	var url, err = generateURL(scheme, "", host, port, path, parameters, make(map[string]*regexp.Regexp))
	if err == nil {
		t.Fail()
	}
	if url == expect {
		t.Fail()
	}
}

func TestUrlTwoParametersWithPort(t *testing.T) {
	var expect = "https://example.com:80/api/user/5/dariusz-poltorak"
	var scheme = "https"
	var host = "example.com"
	var port = 80
	var path = "/api/user/{id}/{slug}"
	var parameters = map[string][]string{"slug": []string{"dariusz-poltorak"}, "id": []string{"5"}}

	var url, err = generateURL(scheme, "", host, port, path, parameters, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}
	if url != expect {
		t.Fail()
	}
}

func TestPathWithTwoParametersWithInvalidRequirement(t *testing.T) {
	var requirements = map[string]*regexp.Regexp{"id": regexp.MustCompile("[0-9]+")}
	var path, err = generatePath("/{id}/{slug}", map[string][]string{"id": []string{"test"}, "slug": []string{"poltorak-dariusz"}}, requirements)
	if err == nil {
		t.Fail()
	}

	if path != "" {
		t.Fail()
	}
}

func TestPathWithTwoParametersWithTwoInvalidRequirements(t *testing.T) {
	var requirements = map[string]*regexp.Regexp{"id": regexp.MustCompile("[0-9]+"), "slug": regexp.MustCompile("[A-Z]+")}
	var path, err = generatePath("/{id}/{slug}", map[string][]string{"id": []string{"test"}, "slug": []string{"poltorak-dariusz"}}, requirements)
	if err == nil {
		t.Fail()
	}

	if path != "" {
		t.Fail()
	}
}

func TestPathWithTwoParametersWithValidRequirement(t *testing.T) {
	var requirements = map[string]*regexp.Regexp{"id": regexp.MustCompile("[0-9]+")}
	var path, err = generatePath("/{id}/{slug}", map[string][]string{"id": []string{"5"}, "slug": []string{"poltorak-dariusz"}}, requirements)
	if err != nil {
		t.Fail()
	}

	if path != "/5/poltorak-dariusz" {
		t.Fail()
	}
}

func TestPathWithTwoParametersWithTwoValidRequirement(t *testing.T) {
	var requirements = map[string]*regexp.Regexp{"id": regexp.MustCompile("[0-9]+"), "slug": regexp.MustCompile("[a-z\\-]+")}
	var path, err = generatePath("/{id}/{slug}", map[string][]string{"id": []string{"5"}, "slug": []string{"poltorak-dariusz"}}, requirements)
	if err != nil {
		t.Fail()
	}

	if path != "/5/poltorak-dariusz" {
		t.Fail()
	}
}

func TestPathWithTwoParametersWithTwoValidRequirementAndExtra(t *testing.T) {
	var requirements = map[string]*regexp.Regexp{"id": regexp.MustCompile("[0-9]+"), "slug": regexp.MustCompile("[a-z\\-]+")}
	var path, err = generatePath("/{id}/{slug}", map[string][]string{"id": []string{"5"}, "slug": []string{"poltorak-dariusz"}, "test": []string{"4"}}, requirements)
	if err != nil {
		t.Fail()
	}

	if path != "/5/poltorak-dariusz?test=4" {
		t.Fail()
	}
}
