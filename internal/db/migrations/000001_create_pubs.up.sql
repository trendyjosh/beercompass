CREATE TABLE pubs (
    id       BIGSERIAL PRIMARY KEY,
    osm_id   BIGINT UNIQUE NOT NULL,
    name     TEXT,
    location GEOMETRY(Point, 4326)
);

CREATE INDEX pubs_location_idx ON pubs USING GIST (location);