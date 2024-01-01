package main 

import (
    "image/color"
    "math"
)

type Sphere struct {
    Center [3]float64
    Radius float64
    Color color.NRGBA
    Specular float64
    Reflective float64
}


var Spheres = [1]Sphere{

}

/*
var Spheres = [4]Sphere{
    Sphere{Center: [3]float64{0, -1, 3}, Radius: 1, Color: color.NRGBA{237, 54, 90, 0xff}, 
        Specular: 500, Reflective: .2},
    Sphere{Center: [3]float64{2, 0, 4}, Radius: 1, Color: color.NRGBA{37, 190, 235, 0xff}, 
        Specular: 500, Reflective: .1},
    Sphere{Center: [3]float64{-2, 0, 4}, Radius: 1, Color: color.NRGBA{210, 210, 210, 0xff}, 
        Specular: 10, Reflective: .7},
    Sphere{Center: [3]float64{0, -5001, 0}, Radius: 5000, Color: color.NRGBA{232, 198, 164, 0xff}, 
        Specular: 500, Reflective: 0},
}
*/

func IntersectRaySphere(O, D [3]float64, sph Sphere) (float64, float64) {
    r := sph.Radius
    CO := [3]float64{O[0] - sph.Center[0], O[1] - sph.Center[1], O[2] - sph.Center[2]}

    a := Dot(D, D)
    b := 2 * Dot(CO, D)
    c := Dot(CO, CO) - r*r

    discriminant := b*b - 4*a*c
    if discriminant < 0 {
        return inf, inf
    }

    t1 := (-b + math.Sqrt(discriminant)) / (2*a)
    t2 := (-b - math.Sqrt(discriminant)) / (2*a)

    return t1, t2
}