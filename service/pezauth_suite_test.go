package pezauth_test

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezauth/service"
	"github.com/pivotal-pez/pezdispenser/service"
	"github.com/xchapter7x/cloudcontroller-client"
	"github.com/xchapter7x/goutil"
	"gopkg.in/mgo.v2"

	"html/template"
	"testing"
	"time"
)

func TestPezAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pez Auth Suite")
}

func setVcapApp(uri string) {
	os.Setenv("VCAP_APPLICATION", fmt.Sprintf(`{  "application_name": "pezauthdev_73b90a93043eb59ee9b3d202dd525f762e865130",  "application_uris": [   "%s"  ],  "application_version": "d744bf29-1465-4634-905d-4fd8a1c19777",  "limits": {   "disk": 1024,   "fds": 16384,   "mem": 1024  },  "name": "pezauthdev_73b90a93043eb59ee9b3d202dd525f762e865130",  "space_id": "49b3e004-702a-4f2c-835c-f25d022882c9",  "space_name": "pez-test",  "uris": [   "%s"  ],  "users": null,  "version": "d744bf29-1465-4634-905d-4fd8a1c19777" }`, uri, uri))
}

func setVcapServ() {
	os.Setenv("VCAP_SERVICES", `{ }`)
}

type mockTokens struct{}

func (s *mockTokens) Access() (r string) {
	return
}

func (s *mockTokens) Refresh() (r string) {
	return
}
func (s *mockTokens) Expired() (r bool) {
	return
}
func (s *mockTokens) ExpiryTime() (r time.Time) {
	return
}

type mockResponseWriter struct {
	StatusCode int
	Body       []byte
}

func (s *mockResponseWriter) WriteHeader(i int) {
	s.StatusCode = i
}

func (s *mockResponseWriter) Header() (r http.Header) {
	return
}

func (s *mockResponseWriter) Write(x []byte) (a int, b error) {
	s.Body = x
	return
}

type (
	mockDoer struct {
		nilResponse bool
		guid        string
		fail        bool
	}
	mockGUIDMaker struct {
		guid string
	}
)

var (
	errDoerCallFailure = errors.New("Failure calling doer")
)

func (s *mockGUIDMaker) Create() string {
	return strings.Split(s.guid, ":")[1]
}

func (s *mockDoer) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if commandName == "KEYS" || commandName == "HMGET" {
		reply = []interface{}{
			[]byte(s.guid),
			[]byte(`"active":  true,"details": "put somethings here"`),
		}
	}

	if s.fail {
		err = errDoerCallFailure
		reply = []interface{}{
			0,
			[]interface{}{},
		}
	}

	if s.nilResponse {
		reply = nil
	}
	return
}

func getKeygen(fail bool, guid string, nilResponse bool) KeyGenerator {
	d := &mockDoer{fail: fail, guid: guid, nilResponse: nilResponse}
	g := &mockGUIDMaker{guid: guid}
	return NewKeyGen(d, g)
}

type mockRenderer struct {
	StatusCode     int
	ResponseObject interface{}
}

func (r *mockRenderer) JSON(status int, v interface{}) {
	r.StatusCode = status
	r.ResponseObject = v
}

func (r *mockRenderer) HTML(status int, name string, binding interface{}, htmlOpt ...render.HTMLOptions) {
}

func (r *mockRenderer) XML(status int, v interface{}) {
}

func (r *mockRenderer) Data(status int, v []byte) {
}

func (r *mockRenderer) Error(status int) {
}

func (r *mockRenderer) Status(status int) {
}

func (r *mockRenderer) Redirect(location string, status ...int) {
}

func (r *mockRenderer) Template() (t *template.Template) {
	return
}

func (r *mockRenderer) Header() (h http.Header) {
	return
}

func setGetUserInfo(domain string, username string) {
	var oldGetUserInfo func(tokens oauth2.Tokens) map[string]interface{}

	BeforeEach(func() {
		oldGetUserInfo = GetUserInfo
		GetUserInfo = func(tokens oauth2.Tokens) map[string]interface{} {
			return map[string]interface{}{
				"domain": domain,
				"emails": []interface{}{
					map[string]interface{}{
						"value": "garbage",
					},
					map[string]interface{}{
						"value": username,
					},
				},
			}
		}
	})

	AfterEach(func() {
		GetUserInfo = oldGetUserInfo
	})
}

type mockRedisCreds struct {
	pass string
	uri  string
}

func (s *mockRedisCreds) Pass() string {
	return s.pass
}

func (s *mockRedisCreds) Uri() string {
	return s.uri
}

type mockMongo struct {
	err    error
	result interface{}
}

func (s *mockMongo) Collection() pezdispenser.Persistence {
	return &mockPersistence{
		err:    s.err,
		result: s.result,
	}
}

func (s *mockMongo) Find(query interface{}) *mgo.Query {
	return new(mgo.Query)
}

func (s *mockMongo) Remove(selector interface{}) (err error) {
	err = s.err
	return
}

func (s *mockMongo) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	err = s.err
	return
}

type mockPersistence struct {
	result interface{}
	err    error
}

func (s *mockPersistence) Remove(selector interface{}) (err error) {
	return
}

func (s *mockPersistence) FindOne(query interface{}, result interface{}) (err error) {
	goutil.Unpack([]interface{}{s.result}, result)
	err = s.err
	return
}

func (s *mockPersistence) Upsert(selector interface{}, update interface{}) (err error) {
	return
}

type mockHeritageClient struct {
	*ccclient.Client
	res *http.Response
}

func (s *mockHeritageClient) CreateAuthRequest(verb, requestURL, path string, args interface{}) (*http.Request, error) {
	return &http.Request{}, nil
}

func (s *mockHeritageClient) CCTarget() string {
	return ccclient.URLPWSLogin
}

func (s *mockHeritageClient) HttpClient() ccclient.ClientDoer {
	return &mockClientDoer{
		res: s.res,
	}
}

func (s *mockHeritageClient) Login() (c *ccclient.Client, err error) {
	return
}

type mockClientDoer struct {
	req *http.Request
	res *http.Response
	err error
}

func (s *mockClientDoer) Do(rq *http.Request) (rs *http.Response, e error) {
	s.req = rq
	rs = s.res
	e = s.err
	return
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func getMockNewOrg(showP, createP, safeCreateP *PivotOrg, showErr, createErr, safeCreateErr error) func(username string, log *log.Logger, tokens oauth2.Tokens, store pezdispenser.Persistence, authClient AuthRequestCreator) OrgManager {
	return func(username string, log *log.Logger, tokens oauth2.Tokens, store pezdispenser.Persistence, authClient AuthRequestCreator) OrgManager {
		s := &mockNewOrg{
			showP:         showP,
			createP:       createP,
			safeCreateP:   safeCreateP,
			showErr:       showErr,
			createErr:     createErr,
			safeCreateErr: safeCreateErr,
		}
		return s
	}
}

type mockNewOrg struct {
	showErr       error
	createErr     error
	safeCreateErr error
	showP         *PivotOrg
	createP       *PivotOrg
	safeCreateP   *PivotOrg
}

func (s *mockNewOrg) Show() (result *PivotOrg, err error) {
	return s.showP, s.showErr
}

func (s *mockNewOrg) SafeCreate() (record *PivotOrg, err error) {
	return s.safeCreateP, s.safeCreateErr
}
