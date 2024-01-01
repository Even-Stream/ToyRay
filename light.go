package main 

import (
    "math"
)

type Light_type int
const (
    Ambient     Light_type = iota
    Point
    Directional
)

type Light struct {
    Ltype Light_type
    Intensity float64
    Position [3]float64
    Direction [3]float64
}

var Lights = [3]Light{
    Light{Ltype: Ambient, Intensity: .1},
    Light{Ltype: Directional, Intensity: .2, Direction: [3]float64{1, 2, 1}},
    Light{Ltype: Point, Intensity: .9, Position: [3]float64{2, 4, 20}},
}

func ReflectRay(R, N [3]float64) [3]float64 {
    dotnr := Dot(N,R)
    return [3]float64{(2 * N[0] * dotnr) - R[0], (2 * N[1] * dotnr) - R[1], (2 * N[2] * dotnr) - R[2]}
}

func ComputeLighting(P, N, V [3]float64, s float64) float64 {
    i := 0.0
    var t_max float64    

    for _, light := range Lights {
        if light.Ltype == Ambient {
            i += light.Intensity
        } else {
            var L [3]float64 //direction from point on sphere to light
            if light.Ltype == Point {
                t_max = 1
                L = [3]float64{light.Position[0] - P[0], light.Position[1] - P[1], light.Position[2] - P[2]}
            } else {
                t_max = math.Inf(0)
                L = light.Direction
            }

            //shadow check
            _, _, found := ClosestIntersection(P, L, 0.001, t_max)

            if found == true {
                continue
            }


            //diffusion reflection
            n_dot_l := Dot(N, L)
            if n_dot_l > 0 {
                i += light.Intensity * n_dot_l / (Vector_length(N) * Vector_length(L))
            }

            //specular reflection
            if s != -1 {
                R := ReflectRay(L, N)

                rvdot := Dot(R, V)
 
                if rvdot > 0 {
                    i += light.Intensity * math.Pow(rvdot / (Vector_length(R) * Vector_length(V)), s)
                }
            }
    }}

    return i
}
