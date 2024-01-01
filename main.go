package main 

import (
    "image"
    "image/color"
    "image/png"
    "os"
    "math"
    //"sync"
//"fmt"
)

const (
    Vw = 1 
    Vh = 1 
    //distance from camera to viewport
    d = 4.2
    Cw = 400
    Ch = 400
)

type Axis int
const (
    Xax     Axis = iota
    Yax
    Zax
)

var inf = math.Inf(0)

var Background_color = color.NRGBA{255, 255, 255, 0xff}

//0, 0 is center of canvas
func PutPixel(canvas *image.NRGBA, x, y int, col color.NRGBA) {
    x += canvas.Rect.Max.X / 2
    y = (canvas.Rect.Max.Y / 2) - y

    canvas.Set(x, y, col)
}

func Vector_length(V [3]float64) float64 {
    a := V[0] * V[0]
    b := V[1] * V[1]
    c := V[2] * V[2]

    return math.Sqrt(a + b + c)
}

//gets  the dot products of 3d vectors 
func Dot(V1, V2 [3]float64) float64{
    return (V1[0] * V2[0]) + (V1[1] * V2[1]) + (V1[2] * V2[2])
}

func FindCross(v1, v2 [3]float64) [3]float64 {
    cross := [3]float64{
        (v1[1] * v2[2]) - (v1[2] * v2[1]),
        (v1[2] * v2[0]) - (v1[0] * v2[2]),
        (v1[0] * v2[1]) - (v1[1] * v2[0]),
    }

    return cross
}

func cam_rotation(vector [3]float64, angle float64, ax Axis) [3]float64 {
    rads := angle * (math.Pi / 180)

    var matrix [3][3]float64
    cos_res := math.Cos(rads)
    sin_res := math.Sin(rads)
    switch {
        case ax == Xax:
            matrix = [3][3]float64 {
                {1, 0, 0},
                {0, cos_res, -sin_res},
                {0, sin_res, cos_res},}
        case ax == Yax:
            matrix = [3][3]float64 {
                {cos_res, 0, sin_res},
                {0, 1, 0},
                {-sin_res, 0, cos_res},}
        case ax == Zax:
            matrix = [3][3]float64 {
                {cos_res, -sin_res, 0},
                {sin_res, cos_res, 0},
                {0, 0, 1},}
    }

    vector_product := [3]float64{
        matrix[0][0] * vector[0] + matrix[0][1] * vector[1] + matrix[0][2] * vector[2],
        matrix[1][0] * vector[0] + matrix[1][1] * vector[1] + matrix[1][2] * vector[2],
        matrix[2][0] * vector[0] + matrix[2][1] * vector[1] + matrix[2][2] * vector[2],}
    return vector_product
}


func ClosestIntersection(O, D [3]float64, t_min, t_max float64) (Sphere, float64, bool) {
    var closest_t float64
    closest_t = inf
    var closest_sphere Sphere
    found := false

    for _, sph := range Spheres {
        t1, t2 := IntersectRaySphere(O, D, sph)
        if t1 >= t_min && t1 <= t_max && t1 < closest_t {
            closest_t = t1
            closest_sphere = sph
            found = true
        }
        if t2 >= t_min && t2 <= t_max && t2 < closest_t {
            closest_t = t2
            closest_sphere = sph
            found = true
        }
    }

    return closest_sphere, closest_t, found
}

func TraceRay(O, D [3]float64, t_min, t_max float64, recursion_depth int) color.NRGBA {
    closest_sphere, closest_t, found := ClosestIntersection(O, D, t_min, t_max)

    closest_color := closest_sphere.Color
    closest_specular := closest_sphere.Specular
    closest_reflective := closest_sphere.Reflective

    tri_check := false
    var tri_normal [3]float64
    var tri_t float64

    for _, tri := range Triangles {
        tri_t, tri_normal = IntersectRayTriangle(O, D, tri)
        if tri_t >= t_min && tri_t <= t_max && tri_t < closest_t {
            tri_check = true
            closest_t = tri_t
            closest_color = tri.Color
            closest_specular = tri.Specular
            closest_reflective = tri.Reflective
            found = true
        }
    }

    if found == false {
        return Background_color     
    }

    P := [3]float64{O[0] + (D[0] * closest_t), O[1] + (D[1] * closest_t), O[2] + (D[2] * closest_t)}

    var N [3]float64
    if tri_check == false {
        N = [3]float64{P[0] - closest_sphere.Center[0], P[1] - closest_sphere.Center[1], P[2] - closest_sphere.Center[2]}
    } else {
        N = tri_normal
    }

    N_len := Vector_length(N)
    N = [3]float64{N[0] / N_len, N[1] / N_len, N[2] / N_len}

    neg_D := [3]float64{-D[0], -D[1], -D[2]}
    
    col_lum := ComputeLighting(P, N, neg_D, closest_specular)
    local_color := [3]float64{float64(closest_color.R) * col_lum,
        float64(closest_color.G) * col_lum,
        float64(closest_color.B) * col_lum,}

    r := closest_reflective
    if recursion_depth > 0 && r > 0 {
        R := ReflectRay(neg_D, N)
        reflected_color := TraceRay(P, R, 0.001, inf, recursion_depth - 1)
        local_color = [3]float64 {local_color[0] * (1 - r) + (float64(reflected_color.R) * r),
            local_color[1] * (1 - r) + (float64(reflected_color.G) * r),
            local_color[2] * (1 - r) + (float64(reflected_color.B) * r),}
    }

    for index, val := range local_color {
        if val > 255 {
            local_color[index] = 255
        }
    }
    post_color := color.NRGBA{uint8(local_color[0]), uint8(local_color[1]), uint8(local_color[2]), closest_color.A}
    return post_color
}

//define Vw and Vh
func CanvasToViewport(x, y float64) [3]float64{
    return [3]float64{x*Vw/Cw, y*Vh/Ch, d}
}

func main() {
    canvas := image.NewNRGBA(image.Rect(0,0, Cw, Ch))
    f, _ := os.OpenFile("test.png", os.O_WRONLY|os.O_CREATE, 0666)
    defer f.Close()
    
    ori := [3]float64{5, 70, 30}

    Triangles = Parse("teapot.obj")
 
    for x := -Cw / 2; x <= Cw/2; x++ {
        for y := -Ch/2; y <= Ch/2; y++ {     
            ctvp := CanvasToViewport(float64(x), float64(y))
            dir := ctvp
            //dir = cam_rotation(dir, -40, Yax)
            dir = cam_rotation(dir, 110, Xax)
            //dir = cam_rotation(dir, 10, Zax)
            col := TraceRay(ori, dir, 1, inf, 3)
            PutPixel(canvas, x, y, col)
    }}
   
    png.Encode(f, canvas.SubImage(canvas.Rect))
}
