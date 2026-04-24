package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"crypto/rand"
	"encoding/hex"
)

type coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type trip struct {
	ID            string     `json:"id"`
	RiderID       string     `json:"rider_id"`
	Status        string     `json:"status"`
	Pickup        coordinate `json:"pickup"`
	Dropoff       coordinate `json:"dropoff"`
	EstimatedFare float64    `json:"estimated_fare"`
	FinalFare     float64    `json:"final_fare,omitempty"`
	Currency      string     `json:"currency"`
	RequestedAt   time.Time  `json:"requested_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}

type createTripRequest struct {
	RiderID string     `json:"rider_id"`
	Pickup  coordinate `json:"pickup"`
	Dropoff coordinate `json:"dropoff"`
}

type completeTripRequest struct {
	FinalFare float64 `json:"final_fare"`
}

var (
	trips   = map[string]trip{}
	tripsMu sync.RWMutex
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"service": "trip-service", "status": "ok"})
	})

	mux.HandleFunc("/v1/trips/request", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req createTripRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
			return
		}

		if req.RiderID == "" || !validCoord(req.Pickup) || !validCoord(req.Dropoff) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "rider_id, pickup, and dropoff are required"})
			return
		}

		estimatedFare := estimateFare(req.Pickup, req.Dropoff)
		newTrip := trip{
			ID:            newID(),
			RiderID:       req.RiderID,
			Status:        "REQUESTED",
			Pickup:        req.Pickup,
			Dropoff:       req.Dropoff,
			EstimatedFare: estimatedFare,
			Currency:      "GHS",
			RequestedAt:   time.Now().UTC(),
		}

		tripsMu.Lock()
		trips[newTrip.ID] = newTrip
		tripsMu.Unlock()

		writeJSON(w, http.StatusAccepted, newTrip)
	})

	mux.HandleFunc("/v1/trips/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/trips/")
		parts := strings.Split(strings.Trim(path, "/"), "/")

		if len(parts) == 1 && r.Method == http.MethodGet {
			handleGetTrip(w, parts[0])
			return
		}

		if len(parts) == 2 && parts[1] == "complete" && r.Method == http.MethodPost {
			handleCompleteTrip(w, r, parts[0])
			return
		}

		writeJSON(w, http.StatusNotFound, map[string]string{"error": "route not found"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	log.Printf("trip-service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func handleGetTrip(w http.ResponseWriter, tripID string) {
	tripsMu.RLock()
	t, ok := trips[tripID]
	tripsMu.RUnlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "trip not found"})
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func handleCompleteTrip(w http.ResponseWriter, r *http.Request, tripID string) {
	var req completeTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
		return
	}

	if req.FinalFare <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "final_fare must be greater than 0"})
		return
	}

	tripsMu.Lock()
	t, ok := trips[tripID]
	if !ok {
		tripsMu.Unlock()
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "trip not found"})
		return
	}

	now := time.Now().UTC()
	t.Status = "COMPLETED"
	t.FinalFare = req.FinalFare
	t.CompletedAt = &now
	trips[tripID] = t
	tripsMu.Unlock()

	writeJSON(w, http.StatusOK, t)
}

func validCoord(c coordinate) bool {
	return c.Lat >= -90 && c.Lat <= 90 && c.Lng >= -180 && c.Lng <= 180
}

func estimateFare(pickup, dropoff coordinate) float64 {
	const baseFare = 8.0
	const perKm = 2.2
	const metersPerDegree = 111_000.0

	deltaLat := pickup.Lat - dropoff.Lat
	deltaLng := pickup.Lng - dropoff.Lng
	distanceMeters := math.Sqrt(deltaLat*deltaLat+deltaLng*deltaLng) * metersPerDegree
	distanceKm := distanceMeters / 1000.0

	total := baseFare + (distanceKm * perKm)
	return math.Round(total*100) / 100
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b)
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
