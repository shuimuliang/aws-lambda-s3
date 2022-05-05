package apiserver

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
    . "github.com/smartystreets/goconvey/convey"
)

var a App

func TestMain(m *testing.M) {
	var err error
	a = App{}

	cfg := &Config{
		Debug:       true,
		ListenAddr:  "0.0.0.0:1230",
		LogInfoPath: "/tmp/app_info.log",
		LogErrPath:  "/tmp/app_error.log",
	}
	err = a.Initialize(cfg)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

// ------------------------------------------------------------------------

func TestPing(t *testing.T) {
	var err error
	req, _ := http.NewRequest("GET", "/ping", nil)
	response := executeRequest(req)

	Convey("TestPing", t, func() {
		So(response.Code, ShouldEqual, http.StatusOK)
		res := &Response{}
		err = DecodeResponse(response.Body, res)
		So(err, ShouldBeNil)
	})
}
