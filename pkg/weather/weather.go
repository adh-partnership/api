/*
 * Copyright Daniel Hawton
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package weather

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/adh-partnership/api/pkg/logger"
)

const (
	addsMetarsCache = "https://aviationweather.gov/data/cache/metars.cache.csv.gz"
	addsTafsCache   = "https://aviationweather.gov/data/cache/tafs.cache.xml.gz"
)

var (
	ErrorNoWeather = errors.New("no weather available")

	weatherCache = make(map[string]*Weather)
	cacheLocked  sync.RWMutex
	log          = logger.Logger.WithField("component", "weather")
)

func GetWeather(icao string) (*Weather, error) {
	cacheLocked.RLock()
	defer cacheLocked.RUnlock()
	if _, ok := weatherCache[icao]; !ok {
		return nil, ErrorNoWeather
	}

	return weatherCache[icao], nil
}

// We should, in theory, be putting this in REDIS or similar... but since all of us right now are using single instance APIs, using an in-memory cache
// is sufficient.
func UpdateWeatherCache() error {
	newCache := make(map[string]*Weather)

	metars, err := processMetars()
	if err != nil {
		return fmt.Errorf("failed to process metars: %s", err)
	}
	for _, metar := range metars.METARs {
		if _, ok := newCache[metar.StationID]; !ok {
			newCache[metar.StationID] = &Weather{}
		}
		newCache[metar.StationID].METAR = metar.RawText
	}

	tafs, err := processTafs()
	if err != nil {
		return fmt.Errorf("failed to process tafs: %s", err)
	}
	for _, taf := range tafs.TAFs {
		if _, ok := newCache[taf.StationID]; !ok {
			newCache[taf.StationID] = &Weather{}
		}
		newCache[taf.StationID].TAF = taf.RawText
	}

	log.Infof("Updating weather cache with, %d weather stations, %d METARs and %d TAFs", len(newCache), len(metars.METARs), len(tafs.TAFs))

	cacheLocked.Lock()
	weatherCache = newCache
	cacheLocked.Unlock()

	return nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func ungzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to ungzip: %s", err)
	}
	defer r.Close()

	return io.ReadAll(r)
}

func processMetars() (*response, error) {
	metars, err := downloadFile(addsMetarsCache)
	if err != nil {
		return nil, fmt.Errorf("failed to download metars: %s", err)
	}

	metarsRaw, err := ungzip(metars)
	if err != nil {
		return nil, fmt.Errorf("failed to ungzip metars: %s", err)
	}

	resp := response{}
	if err := xml.NewDecoder(bytes.NewReader(metarsRaw)).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode metars: %s", err)
	}

	return &resp, nil
}

func processTafs() (*response, error) {
	tafs, err := downloadFile(addsTafsCache)
	if err != nil {
		return nil, fmt.Errorf("failed to download tafs: %s", err)
	}

	tafsRaw, err := ungzip(tafs)
	if err != nil {
		return nil, fmt.Errorf("failed to ungzip tafs: %s", err)
	}

	resp := response{}
	if err := xml.NewDecoder(bytes.NewReader(tafsRaw)).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode tafs: %s", err)
	}

	return &resp, nil
}
