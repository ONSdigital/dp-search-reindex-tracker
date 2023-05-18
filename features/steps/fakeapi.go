package steps

import (
	"context"
	"net/http"

	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/maxcnunes/httpfake"
)

// FakeAPI contains all the information for a fake component API
type FakeAPI struct {
	fakeHTTP      *httpfake.HTTPFake
	healthRequest *httpfake.Request
}

// NewFakeAPI creates a new fake component API
func NewFakeAPI() *FakeAPI {
	return &FakeAPI{
		fakeHTTP: httpfake.New(),
	}
}

// Close closes the fake API
func (f *FakeAPI) Close() {
	f.fakeHTTP.Close()
}

// getMockAPIHTTPClient mocks HTTP Client
func (f *FakeAPI) getMockAPIHTTPClient() *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(paths []string) {},
		GetPathsWithNoRetriesFunc: func() []string { return []string{} },
		DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return f.fakeHTTP.Server.Client().Do(req)
		},
	}
}
