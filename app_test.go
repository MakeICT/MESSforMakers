// package main

// import (
// 	"fmt"
// 	. "github.com/onsi/gomega"
// 	"github.com/sclevine/agouti"
// 	// 	"github.com/sclevine/agouti/matchers"
// 	"net/http"
// 	"os"
// 	"runtime/debug"
// 	"testing"
// )

// var (
// 	driver *agouti.WebDriver
// 	page   *agouti.Page
// )

// func TestMain(m *testing.M) {
// 	var t *testing.T
// 	var err error

// 	driver = agouti.ChromeDriver()
// 	driver.Start()

// 	go startWebsite()

// 	page, err = agouti.NewPage(
// 		driver.URL(),
// 		agouti.Desired(agouti.Capabilities{
// 			"chromeOptions": map[string][]string{
// 				"args": []string{
// 					"disable-gpu",
// 					"no-sandbox",
// 					"headless",
// 				},
// 			},
// 		}),
// 	)

// 	if err != nil {
// 		t.Error("Failed to open page.")
// 	}

// 	RegisterTestingT(t)
// 	test := m.Run()

// 	driver.Stop()
// 	os.Exit(test)

// }
// func startWebsite() {
// 	config, err := InitConfig("config.json")
// 	if err != nil {
// 		fmt.Print("Cannot parse the configuration file")
// 		panic(1)
// 	}

// 	// create the app with user-defined settings
// 	app := newApplication(config)

// 	// make sure the logger releases it's resources if the server shuts down.
// 	defer app.logger.Close()

// 	app.logger.Println("Starting Application")
// 	app.logger.Fatal(http.ListenAndServe(":8080", app.Router))
// }

// func StopDriverOnPanic() {
// 	//var t *testing.T
// 	if r := recover(); r != nil {
// 		debug.PrintStack()
// 		fmt.Println("Recovered in StopDriverOnPanic", r)
// 		//driver.Stop()
// 		//t.Fail()
// 	}
// }

// func TestPage(t *testing.T) {
// 	defer StopDriverOnPanic()
// 	Expect(page.Navigate("http://localhost:8080")).To(Succeed())
// }

// func TestForm(t *testing.T) {
// 	defer StopDriverOnPanic()

// 	Expect(page.Navigate("http://localhost:8080/")).To(Succeed()) //fmt.Sprintf("%v/user", baseUrl)
// 	err := page.Find("#submit-butt").Click()
// 	fmt.Println(page.Find("#submit-butt").String())
// 	fmt.Println(err)
// 	fmt.Println(Succeed().Match(err))
// 	Expect(page.Find("#submit-butt").Click()).To(Succeed())

// }

//******************************************************************************************

// package main

// import (
// 	"github.com/onsi/ginkgo"
// 	"github.com/onsi/gomega"

// 	"github.com/onsi/gomega/gexec"
// 	"github.com/sclevine/agouti"
// 	"os"
// 	"os/exec"
// 	"testing"
// )

// var (
// 	agoutiDriver   *agouti.WebDriver
// 	websiteSession *gexec.Session
// )

// func TestWebsite(t *testing.T) {
// 	gomega.RegisterFailHandler(ginkgo.Fail)
// 	ginkgo.RunSpecs(t, "Website Suite")
// }

// var _ = ginkgo.BeforeSuite(func() {
// 	agoutiDriver = agouti.ChromeDriver()
// 	gomega.Expect(agoutiDriver.Start()).To(gomega.Succeed())

// 	startWebsite()
// })

// var _ = ginkgo.AfterSuite(func() {
// 	gomega.Expect(agoutiDriver.Stop()).To(gomega.Succeed())
// 	websiteSession.Kill()
// })

// func getPage() *agouti.Page {
// 	var page *agouti.Page
// 	var err error

// 	if os.Getenv("TEST_ENV") == "CI" {
// 		page, err = agouti.NewPage(agoutiDriver.URL(),
// 			agouti.Desired(agouti.Capabilities{
// 				"chromeOptions": map[string][]string{
// 					"args": []string{
// 						"disable-gpu",
// 						"no-sandbox",
// 					},
// 				},
// 			}),
// 		)
// 	} else {
// 		//page, err = agoutiDriver.NewPage(agouti.Browser("chrome"))
// 		page, err = agouti.NewPage(agoutiDriver.URL(),
// 			agouti.Desired(agouti.Capabilities{
// 				"chromeOptions": map[string][]string{
// 					"args": []string{
// 						"disable-gpu",
// 						"no-sandbox",
// 						"headless",
// 					},
// 				},
// 			}),
// 		)
// 	}
// 	gomega.Expect(err).NotTo(gomega.HaveOccurred())

// 	return page
// }

// func startWebsite() {
// 	command1 := exec.Command("go", "build")
// 	command2 := exec.Command("./MESSforMakers")
// 	gomega.Eventually(func() error {
// 		var err error
// 		websiteSession, err = gexec.Start(command1, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
// 		websiteSession, err = gexec.Start(command2, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
// 		return err
// 	}).Should(gomega.Succeed())
// }

// *********************************************************************************

package main

import (
	"bytes"
	"flag"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// set up anything that should be run once before all tests here

	//parse flags so that "go test" respects command line flags
	flag.Parse()

	// run the test suite and store the code
	exitCode := m.Run()

	// do any teardown needed once after all tests

	// Exit and return the code
	os.Exit(exitCode)
}

func TestNewApplication(t *testing.T) {

	// Test that the app initializer panics if there is a bad config supplied
	cfg := &Config{}
	cfg.Database.Username = "none"
	t.Run("bad config should panic", testNewAppFunc(cfg, true))

	// check that the app initializer rerturn OK if a good config is supplied
	// TODO: set up a testing database so that connection is possible.
	cfg = &Config{}
	cfg.Database.Username = "postgres_test"
	t.Run("good config should not panic", testNewAppFunc(cfg, false))

}

//pass a Config struct in, along with whether the test is expected to panic
func testNewAppFunc(cfg *Config, expectToPanic bool) func(*testing.T) {
	return func(t *testing.T) {
		defer func() {
			if expectToPanic {
				if r := recover(); r == nil {
					t.Error("app did not panic with bad config")
				}
			} else {
				if r := recover(); r != nil {
					t.Error("app panicked with good config")
				}
			}
		}()

		_ = newApplication(cfg)

	}
}

type AppTestServer struct {
	client *http.Client
	app    *application
	t      *testing.T
	server *httptest.Server
}

func TestRoutes(t *testing.T) {

	//set up a functional configuration
	//TODO make this a test config, not a real config.  Needs test database set up first
	cfg, err := InitConfig("config.json")

	// create a new app
	app := newApplication(cfg)

	// start a test server running that app
	server := httptest.NewServer(app.Router)

	//make sure the server gets shut down after testing
	defer server.Close()

	// request the root route
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Error(err)
	}

	// verify the response is correct
	buf := &bytes.Buffer{}
	buf.ReadFrom(resp.Body)
	if strings.Index(buf.String(), "root handler") == -1 {
		t.Error("Root should say  root handler")
	}
}

func TestCookies(t *testing.T) {
	cfg, err := InitConfig("config.json")
	app := newApplication(cfg)
	server := httptest.NewServer(app.Router)
	defer server.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{Jar: jar}
	resp, err := client.Get(server.URL + "/")
	if err != nil {
		t.Error(err)
	}
	buf := &bytes.Buffer{}
	buf.ReadFrom(resp.Body)
	if strings.Index(buf.String(), "Who are you") == -1 {
		t.Error("Root should ask who on first visit")
	}

	resp, err = client.PostForm(
		server.URL+"/",
		url.Values{"name": {"somebody"}},
	)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	buf.ReadFrom(resp.Body)
	if strings.Index(buf.String(), "Hi somebody") == -1 {
		t.Error("root should say hi after form is posted")
	}
}
