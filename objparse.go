package main

import (
    "os"
    "bufio"
    "strings"
    "strconv"
    "image/color"
    "math/rand"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func Parse(name string) []Triangle {
    triangles := make([]Triangle, 0)
    vertices := make([][3]float64, 0)

    f, err := os.Open(name)
    check(err)
    defer f.Close()

    fileScanner := bufio.NewScanner(f)
    fileScanner.Split(bufio.ScanLines)

    //every vertex has three coordinates, every triangle has 3 vertices
    for fileScanner.Scan() {
        text := fileScanner.Text()
        if strings.HasPrefix(text, "v ") {
            coord_array := strings.Split(text, " ")
            x, _ := strconv.ParseFloat(coord_array[1], 32)
            y, _ := strconv.ParseFloat(coord_array[2], 32)
            z, _ := strconv.ParseFloat(coord_array[3], 32)
            vertices = append(vertices, [3]float64{x * 10,y * 10,z * 10})
        } else if strings.HasPrefix(text, "f ") {
            varray := strings.Split(text, " ")
            i1, _ := strconv.Atoi(strings.Split(varray[1], "//")[0])
            i2, _ := strconv.Atoi(strings.Split(varray[2], "//")[0])
            i3, _ := strconv.Atoi(strings.Split(varray[3], "//")[0])

            v1 := vertices[i1 - 1]
            v2 := vertices[i2 - 1]
            v3 := vertices[i3 - 1]

            ctriangle := Triangle{Vertices: [3][3]float64{v1, v2, v3}, Color: color.NRGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 0xff},
                Specular: 500, Reflective: 0}
            //ctriangle.Reflective += 1
            triangles = append(triangles, ctriangle)
            //fmt.Println(ctriangle)
        }
    }

   // fmt.Println("done")
    return triangles
}