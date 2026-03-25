package smoke

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestBookingFlow(t *testing.T) {
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		t.Skip("APP_BASE_URL is not set")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	waitForReady(t, client, baseURL)
	targetDate := time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")

	adminToken := dummyLogin(t, client, baseURL, "admin")
	userToken := dummyLogin(t, client, baseURL, "user")

	roomID := createRoom(t, client, baseURL, adminToken)
	createSchedule(t, client, baseURL, adminToken, roomID)
	slotID := getFirstSlot(t, client, baseURL, userToken, roomID, targetDate)
	bookingID := createBooking(t, client, baseURL, userToken, slotID)
	cancelBooking(t, client, baseURL, userToken, bookingID)
}

func waitForReady(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()

	deadline := time.Now().Add(40 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/_info")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	t.Fatal("application is not ready")
}

func dummyLogin(t *testing.T, client *http.Client, baseURL, role string) string {
	t.Helper()

	var response struct {
		Token string `json:"token"`
	}

	doJSON(t, client, request{
		method: http.MethodPost,
		url:    baseURL + "/dummyLogin",
		body:   map[string]any{"role": role},
		token:  "",
		want:   http.StatusOK,
		out:    &response,
	})

	return response.Token
}

func createRoom(t *testing.T, client *http.Client, baseURL, token string) string {
	t.Helper()

	var response struct {
		Room struct {
			ID string `json:"id"`
		} `json:"room"`
	}

	doJSON(t, client, request{
		method: http.MethodPost,
		url:    baseURL + "/rooms/create",
		body: map[string]any{
			"name":     "Blue",
			"capacity": 8,
		},
		token: token,
		want:  http.StatusCreated,
		out:   &response,
	})

	return response.Room.ID
}

func createSchedule(t *testing.T, client *http.Client, baseURL, token, roomID string) {
	t.Helper()

	doJSON(t, client, request{
		method: http.MethodPost,
		url:    fmt.Sprintf("%s/rooms/%s/schedule/create", baseURL, roomID),
		body: map[string]any{
			"daysOfWeek": []int{1, 2, 3, 4, 5, 6, 7},
			"startTime":  "09:00",
			"endTime":    "18:00",
		},
		token: token,
		want:  http.StatusCreated,
	})
}

func getFirstSlot(t *testing.T, client *http.Client, baseURL, token, roomID, date string) string {
	t.Helper()

	var response struct {
		Slots []struct {
			ID string `json:"id"`
		} `json:"slots"`
	}

	doJSON(t, client, request{
		method: http.MethodGet,
		url:    fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, roomID, date),
		token:  token,
		want:   http.StatusOK,
		out:    &response,
	})

	if len(response.Slots) == 0 {
		t.Fatal("expected at least one slot")
	}

	return response.Slots[0].ID
}

func createBooking(t *testing.T, client *http.Client, baseURL, token, slotID string) string {
	t.Helper()

	var response struct {
		Booking struct {
			ID string `json:"id"`
		} `json:"booking"`
	}

	doJSON(t, client, request{
		method: http.MethodPost,
		url:    baseURL + "/bookings/create",
		body: map[string]any{
			"slotId":               slotID,
			"createConferenceLink": true,
		},
		token: token,
		want:  http.StatusCreated,
		out:   &response,
	})

	return response.Booking.ID
}

func cancelBooking(t *testing.T, client *http.Client, baseURL, token, bookingID string) {
	t.Helper()

	doJSON(t, client, request{
		method: http.MethodPost,
		url:    fmt.Sprintf("%s/bookings/%s/cancel", baseURL, bookingID),
		token:  token,
		want:   http.StatusOK,
	})
}

type request struct {
	method string
	url    string
	body   any
	token  string
	want   int
	out    any
}

func doJSON(t *testing.T, client *http.Client, req request) {
	t.Helper()

	var bodyBytes []byte
	var err error
	if req.body != nil {
		bodyBytes, err = json.Marshal(req.body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
	}

	httpReq, err := http.NewRequest(req.method, req.url, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	if req.body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	if req.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+req.token)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != req.want {
		t.Fatalf("unexpected status for %s %s: got %d want %d", req.method, req.url, resp.StatusCode, req.want)
	}

	if req.out != nil {
		if err := json.NewDecoder(resp.Body).Decode(req.out); err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}
}
