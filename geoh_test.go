package geoh

import (
	"os"
	"testing"
)

func Test_Geohashes_FeatureCollection(t *testing.T) {
	if file, err := os.ReadFile("./testdata/geojson-sf.json"); err != nil {
		t.Fatalf("Failed to read file: %v", err)
	} else {
		geojson := string(file)
		geohashes := Geohashes(geojson, 6, 2)
		if len(geohashes) == 0 {
			t.Fatalf("Expected geohashes to be generated, but got none.")
		}
	}
}

func Test_Geohashes_Feature(t *testing.T) {
	if file, err := os.ReadFile("./testdata/geojson-feature.json"); err != nil {
		t.Fatalf("Failed to read file: %v", err)
	} else {
		geojson := string(file)
		geohashes := Geohashes(geojson, 6, 2)
		if len(geohashes) == 0 {
			t.Fatalf("Expected geohashes to be generated, but got none.")
		}
	}
}
