package test

import (
	"fmt"
	"math"

	"../render"
	"github.com/ungerik/go3d/vec3"
)

func RotateVert(vertice []float32) {
	var dir vec3.T
	for _x := -1.0; _x < 1.0; _x += 0.2 {
		for _y := -1.0; _y < 1.0; _y += 0.2 {
			for _z := -1.0; _z < 1.0; _z += 0.2 {
				dir[0] = float32(_x)
				dir[1] = float32(_y)
				dir[2] = float32(_z)
				dir.Normalize()
				output_vert1 := render.UpdatePosDir_test(vertice[:], dir)
				output_vert2 := render.UpdatePosDir_test2(vertice[:], dir)
				if math.Abs(float64(output_vert1[0]-output_vert2[0])) > 1e-3 || math.Abs(float64(output_vert1[1]-output_vert2[1])) > 1e-3 || math.Abs(float64(output_vert1[2]-output_vert2[2])) > 1e-3 { //
					fmt.Printf("Inconsistent value:%v:%v \n", output_vert1, output_vert2)
				} else {
					fmt.Printf("Consistent value:%v:%v \n", output_vert1, output_vert2)
				}

				/*
					if math.Abs(float64(dir[0]-vertice[0])) > 1e-3 || math.Abs(float64(dir[1]-vertice[1])) > 1e-3 || math.Abs(float64(dir[2]-vertice[2])) > 1e-3 { //
						fmt.Printf("Inconsistent value:%v:(%v, %v, %v) \n", dir, vertice[0], vertice[1], vertice[2])
					} else {
						fmt.Printf("Consistent value:%v:(%v, %v, %v) \n", dir, vertice[0], vertice[1], vertice[2])
					}
				*/
			}
		}
	}
}

func Main() {
	var vertice [180]float32

	for _x := -1.0; _x < 1.0; _x += 0.2 {
		for _y := -1.0; _y < 1.0; _y += 0.2 {
			for _z := -1.0; _z < 1.0; _z += 0.2 {
				vertice[0] = float32(_x)
				vertice[1] = float32(_y)
				vertice[2] = float32(_z)
				RotateVert(vertice[:])
			}
		}
	}
	// test code end

	//	render.UpdatePosDir_test(vertice[:], pos, dir, &f0)
}

// Test code
/*
	view_ray_vec := vec3.T{1.0, 0, 0}
	rotate_axis := vec3.T{0, 1.0, 0}
	cross := vec3.Cross(&rotate_axis, &view_ray_vec)
	rotate_angle := float32(math.Pi / 4.0)
	rotated_view_ray_vec := common.RotateVector(view_ray_vec, rotate_axis, rotate_angle)
	fmt.Println("rotated_view_ray_vec val is:%v, %v", rotated_view_ray_vec, cross)
*/
