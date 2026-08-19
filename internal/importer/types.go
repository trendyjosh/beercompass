package importer

import "github.com/paulmach/osm"

// Pub holds the data extracted from an OSM pub node.
type Pub struct {
	OsmID osm.NodeID
	Name  string
	Lat   float64
	Lon   float64
}
