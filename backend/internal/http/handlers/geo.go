package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type GeoHandler struct {
	apiKey string
}

func NewGeoHandler(apiKey string) *GeoHandler {
	return &GeoHandler{apiKey: apiKey}
}

type geocodeResult struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	MapURL string `json:"map_url"`
}

func (h *GeoHandler) Geocode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address is required", http.StatusBadRequest)
		return
	}

	apiURL := fmt.Sprintf(
		"https://geocode-maps.yandex.ru/1.x/?apikey=%s&geocode=%s&format=json",
		h.apiKey,
		url.QueryEscape(address),
	)

	resp, err := http.Get(apiURL) //nolint:noctx
	if err != nil {
		http.Error(w, "geocoding failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "geocoding failed", http.StatusInternalServerError)
		return
	}

	var yandexResp struct {
		Response struct {
			GeoObjectCollection struct {
				FeatureMember []struct {
					GeoObject struct {
						Name        string `json:"name"`
						Description string `json:"description"`
						Point       struct {
							Pos string `json:"pos"`
						} `json:"Point"`
					} `json:"GeoObject"`
				} `json:"featureMember"`
			} `json:"GeoObjectCollection"`
		} `json:"response"`
	}

	if err := json.Unmarshal(body, &yandexResp); err != nil {
		http.Error(w, "geocoding failed", http.StatusInternalServerError)
		return
	}

	members := yandexResp.Response.GeoObjectCollection.FeatureMember
	if len(members) == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	geo := members[0].GeoObject
	parts := strings.Fields(geo.Point.Pos) // "lon lat"
	if len(parts) != 2 {
		http.Error(w, "invalid coordinates", http.StatusInternalServerError)
		return
	}

	lon, err1 := strconv.ParseFloat(parts[0], 64)
	lat, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		http.Error(w, "invalid coordinates", http.StatusInternalServerError)
		return
	}

	result := geocodeResult{
		Name:        geo.Name,
		Description: geo.Description,
		MapURL:      fmt.Sprintf("https://yandex.ru/maps/?ll=%s,%s&z=15&pt=%s,%s", parts[0], parts[1], parts[0], parts[1]),
	}
	result.Coordinates.Latitude = lat
	result.Coordinates.Longitude = lon

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result) //nolint:errcheck
}
