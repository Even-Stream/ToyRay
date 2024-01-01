package main 

import (
    "image/color"
    "math"

//"fmt"
)

type Triangle struct {
    Vertices [3][3]float64
    Color color.NRGBA
    Specular float64
    Reflective float64
}


var Triangles []Triangle

/*
var Triangles = [6]Triangle {
    Triangle{Vertices: [3][3]float64{{-1, 1, 5}, {-1, -1, 5}, {1, 1, 5}}, Color: color.NRGBA{163, 73, 164, 0xff},
        Specular: 500, Reflective: 0},
    Triangle{Vertices: [3][3]float64{{1, -1, 5}, {-1, -1, 5}, {1, 1, 5}}, Color: color.NRGBA{255, 201, 14, 0xff},
        Specular: 500, Reflective: 0},

    Triangle{Vertices: [3][3]float64{{1, -1, 5}, {1, 1, 7}, {1, 1, 5}}, Color: color.NRGBA{35, 169, 203, 0xff},
        Specular: 500, Reflective: 0},
    Triangle{Vertices: [3][3]float64{{1, -1, 5}, {1, 1, 7}, {1, -1, 7}}, Color: color.NRGBA{255, 155, 89, 0xff},
        Specular: 500, Reflective: 0},

    Triangle{Vertices: [3][3]float64{{-1, 1, 5}, {1, 1, 7}, {1, 1, 5}}, Color: color.NRGBA{240, 96, 99, 0xff},
        Specular: 500, Reflective: 0},
    Triangle{Vertices: [3][3]float64{{-1, 1, 5}, {-1, 1, 7}, {1, 1, 7}}, Color: color.NRGBA{81, 245, 10, 0xff},
        Specular: 500, Reflective: 0},
}
*/

func FindTriNormal(vertices [3][3]float64) [3]float64 {
    v1 := [3]float64{
        vertices[1][0] -  vertices[0][0], vertices[1][1] -  vertices[0][1], vertices[1][2] -  vertices[0][2]}

    v2 := [3]float64{
        vertices[2][0] -  vertices[0][0], vertices[2][1] -  vertices[0][1], vertices[2][2] -  vertices[0][2]}

    cross := FindCross(v1, v2)

    return cross
}

func FindPlane(vertices [3][3]float64) ([3]float64, [3]float64) {
    cross := FindTriNormal(vertices)

    point := vertices[0]
    //distance from origin(0,0,0) to plane
    pdist := [3]float64 {(cross[0] * -point[0]), (cross[1] * -point[1]), (cross[2] * -point[2])}

    return cross, pdist
}

//finds intersection of ray(given origin and direction) and plane, and whether it intersects with the triangle
func IntersectRayTriangle(O, D [3]float64, tri Triangle) (float64, [3]float64) {
    normal, _ := FindPlane(tri.Vertices)

    dndot := Dot(normal, D)

    if math.Abs(dndot) < 0.0001 {
        //fmt.Println("does not intersect")
        return math.Inf(0), normal
    }

    point0 := tri.Vertices[1]
    t := Dot([3]float64{(point0[0] - O[0]), (point0[1] - O[1]), (point0[2] - O[2])}, normal)
    t = t / dndot 
    
    intpoint := [3]float64{O[0] + (D[0] * t), O[1] + (D[1] * t), O[2] + (D[2] * t)}

    a := tri.Vertices[0]
    b := tri.Vertices[1]
    c := tri.Vertices[2]

    pa := [3]float64{a[0] - intpoint[0], a[1] - intpoint[1], a[2] - intpoint[2]}
    pb := [3]float64{b[0] - intpoint[0], b[1] - intpoint[1], b[2] - intpoint[2]}
    pc := [3]float64{c[0] - intpoint[0], c[1] - intpoint[1], c[2] - intpoint[2]} 

    pbxpc := FindCross(pb, pc) 
    pcxpa := FindCross(pc, pa)

    triarea := Vector_length(normal) / 2

    N_len := Vector_length(normal)
    N := [3]float64{normal[0] / N_len, normal[1] / N_len, normal[2] / N_len}

    alpha := Dot(pbxpc, N) / (triarea * 2)
    beta := Dot(pcxpa, N) / (triarea * 2)
    gamma := 1 - alpha - beta

    if (0.0 <= alpha && alpha <= 1.0) && (0.0 <= beta && beta <= 1.0) && (0.0 <= gamma && gamma <= 1.0) {
        return t, normal
    }

    return math.Inf(0), normal
}