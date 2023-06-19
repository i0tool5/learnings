# How to deploy your own *OpenStreetMap*-like service

The most difficult thing in this journey is to find proper way. There is plethora of obsolete, low-level, incomplete tools, obsolete docs and other things that can lead astray. 

### Steps
1. Download **osm.pbf** files from [Geofabric](https://download.geofabrik.de/) or other sources
2. Create tiles using [tilemaker](https://github.com/systemed/tilemaker) or other tools
3. Host tiles, using one of the tile servers. One of the most feature-rich is [tileserver-gl](https://github.com/maptiler/tileserver-gl)
4. Create front-end for map rendering. Most popular libraries are [MapLibreGL](https://github.com/maplibre/maplibre-gl-js) and [Leaflet](https://github.com/Leaflet/Leaflet)

#### Note about styles
To display beautiful and informative maps you need styles. You can create your own styles, or take them from [Open Map Tiles](https://openmaptiles.org/styles/). You need to configure properly your service to use this styles

## Tilemaker

```sh
docker run \ 
    -v .:/map --rm --name tailmaking-test tailmaker:latest \ 
    --input /map/your-file.osm.pbf --output /map/out-name.mbtiles \ 
    --config /map/config-openmaptiles.json  \ 
    --process /map/process-openmaptiles.lua
```

> `--config` and `--process` flags can be omitted

# Links

- [OpenStreetMap](https://wiki.openstreetmap.org/wiki/Main_Page) Wiki
- [OpenMapTiles](https://openmaptiles.org/)
- [Switch to OSM](https://switch2osm.org/)
- More about styles on [MapBox](https://docs.mapbox.com/mapbox-gl-js/style-spec/)
