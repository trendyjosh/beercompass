package importer

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

const osmDownloadURL = "https://download.geofabrik.de/europe/great-britain-latest.osm.pbf"
const tmpPath string = "downloads/great-britain-latest.osm.pbf.tmp"
const osmPath string = "downloads/great-britain-latest.osm.pbf"
const downloadTimeout time.Duration = 45 * time.Minute
const osmScannerWorkers int = 4

// Download OSM data from geofrabik (updated daily).
func DownloadOsm() error {
	// Make sure downloads directory exists
	if err := os.MkdirAll("downloads", 0755); err != nil {
		return fmt.Errorf("downloadOsm: creating downloads directory: %w", err)
	}

	// Use a client with a timeout appropriate for a large PBF file
	client := &http.Client{Timeout: downloadTimeout}

	// Retrieve feed from endpoint
	resp, err := client.Get(osmDownloadURL)
	if err != nil {
		return fmt.Errorf("downloadOsm: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle unsuccessful responses
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloadOsm: unexpected status: %s", resp.Status)
	}

	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("downloadOsm: creating temp file: %w", err)
	}

	// Clean up the temp file if anything goes wrong after creation
	defer func() {
		if err != nil {
			os.Remove(tmpPath)
		}
	}()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		f.Close()
		return fmt.Errorf("downloadOsm: writing temp file: %w", err)
	}

	// Validate the download size against Content-Length if provided
	if resp.ContentLength > 0 && n != resp.ContentLength {
		f.Close()
		return fmt.Errorf("downloadOsm: incomplete download: got %d bytes, expected %d", n, resp.ContentLength)
	}

	// Flush to disk before renaming to guard against data loss on crash
	if err = f.Sync(); err != nil {
		f.Close()
		return fmt.Errorf("downloadOsm: syncing temp file: %w", err)
	}
	f.Close()

	if err = os.Rename(tmpPath, osmPath); err != nil {
		return fmt.Errorf("downloadOsm: renaming temp file: %w", err)
	}

	log.Printf("downloadOsm: downloaded %d bytes to %s", n, osmPath)
	return nil
}

// Parse the OSM data to extract only pub nodes.
func ParseOsm() ([]Pub, error) {
	f, err := os.Open(osmPath)
	if err != nil {
		return nil, fmt.Errorf("parseOsm: opening OSM file: %w", err)
	}
	defer f.Close()

	scanner := osmpbf.New(context.Background(), f, osmScannerWorkers)
	defer scanner.Close()

	var pubs []Pub

	for scanner.Scan() {
		obj := scanner.Object()

		// Limit to only nodes
		node, ok := obj.(*osm.Node)
		if !ok {
			continue
		}

		// Limit to only pubs
		if node.Tags.Find("amenity") != "pub" {
			continue
		}

		pubs = append(pubs, Pub{
			OsmID: node.ID,
			Name:  node.Tags.Find("name"),
			Lat:   node.Lat,
			Lon:   node.Lon,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parseOsm: scanning OSM data: %w", err)
	}

	return pubs, nil
}
