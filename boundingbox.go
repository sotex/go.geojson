package geojson

import (
	"fmt"
	"math"
)

func decodeBoundingBox(bb interface{}) ([]float64, error) {
	if bb == nil {
		return nil, nil
	}

	switch f := bb.(type) {
	case []float64:
		return f, nil
	case []interface{}:
		bb := make([]float64, 0, 4)
		for _, v := range f {
			switch c := v.(type) {
			case float64:
				bb = append(bb, c)
			default:
				return nil, fmt.Errorf("bounding box coordinate not usable, got %T", v)
			}

		}
		return bb, nil
	default:
		return nil, fmt.Errorf("bounding box property not usable, got %T", bb)
	}
}

func checkBoundingBox(min, max [2]float64) bool {
	return min[0] <= max[0] && min[1] <= max[1]
}

// GetBoundingBox  Get Geometry geo outsourcing box
func (g *Geometry) GetBoundingBox() (min, max [2]float64, ok bool) {
	// x0,y0,x1,y1
	if g.BoundingBox != nil {
		min[0] = g.BoundingBox[0]
		min[1] = g.BoundingBox[1]
		max[0] = g.BoundingBox[2]
		max[1] = g.BoundingBox[3]
		return min, max, checkBoundingBox(min, max)
	}

	min[0] = math.MaxFloat64
	min[1] = math.MaxFloat64
	max[0] = -math.MaxFloat64
	max[1] = -math.MaxFloat64

	switch g.Type {
	case GeometryPoint:
		min[0] = g.Point.X
		min[1] = g.Point.Y
		max[0] = min[0]
		max[1] = min[1]
	case GeometryMultiPoint:
		for i := 0; i < len(g.MultiPoint); i++ {
			min[0] = math.Min(min[0], g.MultiPoint[i].X)
			min[1] = math.Min(min[1], g.MultiPoint[i].Y)
			max[0] = math.Max(max[0], g.MultiPoint[i].X)
			max[1] = math.Max(max[1], g.MultiPoint[i].Y)
		}
	case GeometryLineString:
		for i := 0; i < len(g.LineString); i++ {
			min[0] = math.Min(min[0], g.LineString[i].X)
			min[1] = math.Min(min[1], g.LineString[i].Y)
			max[0] = math.Max(max[0], g.LineString[i].X)
			max[1] = math.Max(max[1], g.LineString[i].Y)
		}
	case GeometryMultiLineString:
		for line := 0; line < len(g.MultiLineString); line++ {
			linestring := g.MultiLineString[line]
			for i := 0; i < len(linestring); i++ {
				min[0] = math.Min(min[0], linestring[i].X)
				min[1] = math.Min(min[1], linestring[i].Y)
				max[0] = math.Max(max[0], linestring[i].X)
				max[1] = math.Max(max[1], linestring[i].Y)
			}
		}
	case GeometryPolygon:
		for line := 0; line < len(g.Polygon); line++ {
			linestring := g.Polygon[line]
			for i := 0; i < len(linestring); i++ {
				min[0] = math.Min(min[0], linestring[i].X)
				min[1] = math.Min(min[1], linestring[i].Y)
				max[0] = math.Max(max[0], linestring[i].X)
				max[1] = math.Max(max[1], linestring[i].Y)
			}
		}
	case GeometryMultiPolygon:
		for poly := 0; poly < len(g.MultiPolygon); poly++ {
			for line := 0; line < len(g.MultiPolygon[poly]); line++ {
				linestring := g.MultiPolygon[poly][line]
				for i := 0; i < len(linestring); i++ {
					min[0] = math.Min(min[0], linestring[i].X)
					min[1] = math.Min(min[1], linestring[i].Y)
					max[0] = math.Max(max[0], linestring[i].X)
					max[1] = math.Max(max[1], linestring[i].Y)
				}
			}
		}
	case GeometryCollection:
		for i := 0; i < len(g.Geometries); i++ {
			tmin, tmax, ok := g.Geometries[i].GetBoundingBox()
			if ok {
				min[0] = math.Min(min[0], tmin[0])
				min[1] = math.Min(min[1], tmin[1])
				max[0] = math.Max(max[0], tmax[0])
				max[1] = math.Max(max[1], tmax[1])
			}
		}
	}
	return min, max, checkBoundingBox(min, max)
}

// GetBoundingBox  Get Feature geo outsourcing box
func (f *Feature) GetBoundingBox() (min, max [2]float64, ok bool) {
	// x0,y0,x1,y1
	if f.BoundingBox != nil {
		min[0] = f.BoundingBox[0]
		min[1] = f.BoundingBox[1]
		max[0] = f.BoundingBox[2]
		max[1] = f.BoundingBox[3]
		return min, max, checkBoundingBox(min, max)
	}
	return f.Geometry.GetBoundingBox()
}

// GetBoundingBox  Get FeatureCollection geo outsourcing box
func (fc *FeatureCollection) GetBoundingBox() (min, max [2]float64, ok bool) {
	// x0,y0,x1,y1
	if fc.BoundingBox != nil {
		min[0] = fc.BoundingBox[0]
		min[1] = fc.BoundingBox[1]
		max[0] = fc.BoundingBox[2]
		max[1] = fc.BoundingBox[3]
		return min, max, checkBoundingBox(min, max)
	}
	min[0] = math.MaxFloat64
	min[1] = math.MaxFloat64
	max[0] = -math.MaxFloat64
	max[1] = -math.MaxFloat64
	for i := 0; i < len(fc.Features); i++ {
		tmin, tmax, ok := fc.Features[i].GetBoundingBox()
		if ok {
			min[0] = math.Min(min[0], tmin[0])
			min[1] = math.Min(min[1], tmin[1])
			max[0] = math.Max(max[0], tmax[0])
			max[1] = math.Max(max[1], tmax[1])
		}
	}
	return min, max, checkBoundingBox(min, max)
}
