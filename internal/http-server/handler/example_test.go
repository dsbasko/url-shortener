package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/http-server/middlewares"
	"github.com/dsbasko/yandex-go-shortener/internal/jwt"
	"github.com/dsbasko/yandex-go-shortener/internal/storage"
	"github.com/dsbasko/yandex-go-shortener/internal/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
)

func ExampleHandler_Ping_success() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/ping", h.Ping)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().Ping(gomock.Any()).Return(nil)

	// Send a request.
	resp, body := test.Request(t, ts, &test.RequestArgs{
		Method: "GET",
		Path:   "/ping",
	})
	defer resp.Body.Close()

	// Print the result.
	fmt.Printf("Status code: %v\n", resp.StatusCode)
	fmt.Printf("Body: %v\n", body)

	// Output:
	// Status code: 200
	// Body: pong
}

func ExampleHandler_Ping_failed() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/ping", h.Ping)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().Ping(gomock.Any()).Return(errors.New("some error"))

	// Send a request.
	resp, _ := test.Request(t, ts, &test.RequestArgs{
		Method: "GET",
		Path:   "/ping",
	})
	defer resp.Body.Close()

	// Print the result.
	fmt.Printf("Status code: %v\n", resp.StatusCode)

	// Output:
	// Status code: 400
}

func ExampleHandler_Redirect_found() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/{short_url}", h.Redirect)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().
		GetURLByShortURL(gomock.Any(), gomock.Any()).
		Return(entities.URL{OriginalURL: "https://ya.ru/"}, nil)

	// Send a request.
	resp, _ := test.Request(t, ts, &test.RequestArgs{
		Method: "GET",
		Path:   "/HkNj",
	})
	defer resp.Body.Close()

	// Read the response.
	location := resp.Header.Get("Location")

	// Print the result.
	fmt.Printf("Status code: %v\n", resp.StatusCode)
	fmt.Printf("Location: %v\n", location)

	// Output:
	// Status code: 307
	// Location: https://ya.ru/
}

func ExampleHandler_Redirect_notFound() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/{short_url}", h.Redirect)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().
		GetURLByShortURL(gomock.Any(), gomock.Any()).
		Return(entities.URL{}, errors.New("some error"))

	// Send a request.
	resp, _ := test.Request(t, ts, &test.RequestArgs{
		Method: "GET",
		Path:   "/HkNj",
	})
	defer resp.Body.Close()

	// Read the response.
	fmt.Printf("Status code: %v\n", resp.StatusCode)

	// Output:
	// Status code: 400
}

func ExampleHandler_Redirect_isDeleted() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/{short_url}", h.Redirect)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().
		GetURLByShortURL(gomock.Any(), gomock.Any()).
		Return(entities.URL{DeletedFlag: true}, nil)

	// Send a request.
	resp, _ := test.Request(t, ts, &test.RequestArgs{
		Method: "GET",
		Path:   "/HkNj",
	})
	defer resp.Body.Close()

	// Read the response.
	fmt.Printf("Status code: %v\n", resp.StatusCode)

	// Output:
	// Status code: 410
}

func ExampleHandler_GetURLsByUserID_found() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/api/user/urls", h.GetURLsByUserID)
	ts := httptest.NewServer(router)
	defer ts.Close()
	token, _ := jwt.GenerateToken()

	// Mock the store.
	store.EXPECT().
		GetURLsByUserID(gomock.Any(), gomock.Any()).
		Return([]entities.URL{
			{
				ID:          "1",
				UserID:      "1",
				ShortURL:    "HkNj",
				OriginalURL: "https://ya.ru/",
			},
		}, nil)

	// Send a request.
	resp, body := test.Request(t, ts, &test.RequestArgs{
		Method: "GET",
		Path:   "/api/user/urls",
		Cookie: &http.Cookie{Name: jwt.CookieKey, Value: token},
	})
	defer resp.Body.Close()

	// Print the result.
	fmt.Printf("Status code: %v\n", resp.StatusCode)
	fmt.Printf("Body: %v\n", body)

	// Output:
	// Status code: 200
	// Body: [{"id":"1","user_id":"1","short_url":"http://localhost:8080/HkNj","original_url":"https://ya.ru/","is_deleted":false}]
}

func ExampleHandler_GetURLsByUserID_notFound() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/api/user/urls", h.GetURLsByUserID)
	ts := httptest.NewServer(router)
	defer ts.Close()
	token, _ := jwt.GenerateToken()

	// Mock the store.
	store.EXPECT().
		GetURLsByUserID(gomock.Any(), gomock.Any()).
		Return([]entities.URL{}, nil)

	// Send a request.
	resp, _ := test.Request(t, ts, &test.RequestArgs{
		Method: "GET",
		Path:   "/api/user/urls",
		Cookie: &http.Cookie{Name: jwt.CookieKey, Value: token},
	})
	defer resp.Body.Close()

	// Read the response.
	fmt.Printf("Status code: %v\n", resp.StatusCode)

	// Output:
	// Status code: 204
}

func ExampleHandler_GetURLsByUserID_unauthorized() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/api/user/urls", h.GetURLsByUserID)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().
		GetURLsByUserID(gomock.Any(), gomock.Any()).
		Return([]entities.URL{}, nil)

	// Send a request.
	resp, _ := test.Request(t, ts, &test.RequestArgs{
		Method: "GET",
		Path:   "/api/user/urls",
	})
	defer resp.Body.Close()

	// Read the response.
	fmt.Printf("Status code: %v\n", resp.StatusCode)

	// Output:
	// Status code: 401
}

func ExampleHandler_GetURLsByUserID_someError() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/api/user/urls", h.GetURLsByUserID)
	ts := httptest.NewServer(router)
	defer ts.Close()
	token, _ := jwt.GenerateToken()

	// Mock the store.
	store.EXPECT().
		GetURLsByUserID(gomock.Any(), gomock.Any()).
		Return([]entities.URL{}, errors.New("some error"))

	// Send a request.
	resp, _ := test.Request(t, ts, &test.RequestArgs{
		Method: "GET",
		Path:   "/api/user/urls",
		Cookie: &http.Cookie{Name: jwt.CookieKey, Value: token},
	})
	defer resp.Body.Close()

	// Read the response.
	fmt.Printf("Status code: %v\n", resp.StatusCode)

	// Output:
	// Status code: 400
}

func ExampleHandler_CreateURLTextPlain_created() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	mw := middlewares.New(log)
	h := New(log, store, urlsService)
	router.With(mw.JWT).Post("/", h.CreateURLTextPlain)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().
		CreateURL(gomock.Any(), gomock.Any()).
		Return(entities.URL{
			ShortURL: "HkNj",
		}, true, nil)

	// Send a request.
	resp, body := test.Request(t, ts, &test.RequestArgs{
		Method: "POST",
		Path:   "/",
		Body:   []byte("https://ya.ru/"),
	})
	defer resp.Body.Close()

	fmt.Printf("Status code: %v\n", resp.StatusCode)
	fmt.Printf("Body: %v\n", body)

	// Output:
	// Status code: 201
	// Body: http://localhost:8080/HkNj
}

func ExampleHandler_CreateURLTextPlain_duplicate() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	mw := middlewares.New(log)
	h := New(log, store, urlsService)
	router.With(mw.JWT).Post("/", h.CreateURLTextPlain)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().
		CreateURL(gomock.Any(), gomock.Any()).
		Return(entities.URL{
			ShortURL: "HkNj",
		}, false, nil)

	// Send a request.
	resp, body := test.Request(t, ts, &test.RequestArgs{
		Method: "POST",
		Path:   "/",
		Body:   []byte("https://ya.ru/"),
	})
	defer resp.Body.Close()

	fmt.Printf("Status code: %v\n", resp.StatusCode)
	fmt.Printf("Body: %v\n", body)

	// Output:
	// Status code: 409
	// Body: http://localhost:8080/HkNj
}

func ExampleHandler_CreateURLTextPlain_someError() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	mw := middlewares.New(log)
	h := New(log, store, urlsService)
	router.With(mw.JWT).Post("/", h.CreateURLTextPlain)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().
		CreateURL(gomock.Any(), gomock.Any()).
		Return(entities.URL{}, false, errors.New("some error"))

	// Send a request.
	resp, _ := test.Request(t, ts, &test.RequestArgs{
		Method: "POST",
		Path:   "/",
		Body:   []byte("https://ya.ru/"),
	})
	defer resp.Body.Close()

	fmt.Printf("Status code: %v\n", resp.StatusCode)

	// Output:
	// Status code: 400
}

func ExampleHandler_CreateURLJSON_created() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	mw := middlewares.New(log)
	h := New(log, store, urlsService)
	router.With(mw.JWT).Post("/api/shorten", h.CreateURLJSON)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().
		CreateURL(gomock.Any(), gomock.Any()).
		Return(entities.URL{
			ShortURL: "HkNj",
		}, true, nil)

	// Send a request.
	resp, body := test.Request(t, ts, &test.RequestArgs{
		Method:      "POST",
		Path:        "/api/shorten",
		Body:        []byte(`{"url":"https://ya.ru/"}`),
		ContentType: "application/json",
	})
	defer resp.Body.Close()

	fmt.Printf("Status code: %v\n", resp.StatusCode)
	fmt.Printf("Body: %v\n", body)

	// Output:
	// Status code: 201
	// Body: {"result":"http://localhost:8080/HkNj"}
}

func ExampleHandler_CreateURLJSON_conflict() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	mw := middlewares.New(log)
	h := New(log, store, urlsService)
	router.With(mw.JWT).Post("/api/shorten", h.CreateURLJSON)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().
		CreateURL(gomock.Any(), gomock.Any()).
		Return(entities.URL{
			ShortURL: "HkNj",
		}, false, nil)

	// Send a request.
	resp, body := test.Request(t, ts, &test.RequestArgs{
		Method:      "POST",
		Path:        "/api/shorten",
		Body:        []byte(`{"url":"https://ya.ru/"}`),
		ContentType: "application/json",
	})
	defer resp.Body.Close()

	fmt.Printf("Status code: %v\n", resp.StatusCode)
	fmt.Printf("Body: %v\n", body)

	// Output:
	// Status code: 409
	// Body: {"result":"http://localhost:8080/HkNj"}
}

func ExampleHandler_CreateURLJSON_someError() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	mw := middlewares.New(log)
	h := New(log, store, urlsService)
	router.With(mw.JWT).Post("/api/shorten", h.CreateURLJSON)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Mock the store.
	store.EXPECT().
		CreateURL(gomock.Any(), gomock.Any()).
		Return(entities.URL{}, false, errors.New("some error"))

	// Send a request.
	resp, _ := test.Request(t, ts, &test.RequestArgs{
		Method:      "POST",
		Path:        "/api/shorten",
		Body:        []byte(`{"url":"https://ya.ru/"}`),
		ContentType: "application/json",
	})
	defer resp.Body.Close()

	// Read the response.
	fmt.Printf("Status code: %v\n", resp.StatusCode)

	// Output:
	// Status code: 400
}

func ExampleHandler_DeleteURLs() {
	// Initialize the server.
	t := &testing.T{}
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Delete("/api/user/urls", h.DeleteURLs)
	ts := httptest.NewServer(router)
	defer ts.Close()
	token, _ := jwt.GenerateToken()

	// Send a request.
	resp, _ := test.Request(t, ts, &test.RequestArgs{
		Method: "DELETE",
		Path:   "/api/user/urls",
		Cookie: &http.Cookie{Name: jwt.CookieKey, Value: token},
		Body:   []byte(`["HkNj","HmNj"]`),
	})
	defer resp.Body.Close()

	// Read the response.
	fmt.Printf("Status code: %v\n", resp.StatusCode)

	// Output:
	// Status code: 202
}
