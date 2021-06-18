package render

import (
	"fmt"
	"math"

	"../common"
	"../core"
	"github.com/ungerik/go3d/vec3"
)

func UpdatePosDir_test(vertice []float32, dir vec3.T) []float32 {
	var output_vert [3]float32
	dir.Normalize()

	// Calculate the rotation quaternion b = dir, b*a^-1
	theta := math.Acos(float64(dir[0]))

	var n2_n3_normal vec3.T
	n2_n3_normal[0] = float32(0)
	n2_n3_normal[1] = -float32(dir[2])
	n2_n3_normal[2] = float32(dir[1])
	n2_n3_normal.Normalize()

	var quatRotation common.Quaternion

	half_cos_theta := float32(math.Cos(theta / 2.0))
	half_sin_theta := float32(math.Sin(theta / 2.0))
	quatRotation.W = half_cos_theta
	quatRotation.X = 0
	quatRotation.Y = n2_n3_normal[1] * half_sin_theta
	quatRotation.Z = n2_n3_normal[2] * half_sin_theta

	var quatRotationConjugate common.Quaternion
	quatRotationConjugate.Copy(&quatRotation)
	quatRotationConjugate.Conjugate()

	PointSize := 5
	for _idx := 0; _idx < 1; _idx++ {
		var quatVert common.Quaternion
		var quatTmp common.Quaternion
		quatVert.SetVert(vertice[_idx*PointSize : _idx*PointSize+3])
		quatTmp.Copy(&quatRotation)
		quatTmp.Multi(&quatVert)
		quatTmp.Multi(&quatRotationConjugate)
		quatTmp.GetVert(output_vert[:])
	}

	return output_vert[:]
}

func UpdatePosDir_test2(vertice []float32, dir vec3.T) []float32 {
	var output_vert [3]float32
	dir.Normalize()

	var quatRotation common.Quaternion
	core.LogStr(fmt.Sprintf("View Dir: %v", dir))
	core.LogStr(fmt.Sprintf("Calculation: %v", math.Atan2(float64(dir[2]), float64(dir[0]))))
	X := 0.0                                                                                     //math.Atan2(float64(dir[2]), math.Abs(float64(dir[1])))
	Y := -math.Atan2(float64(dir[2]), math.Abs(math.Sqrt(float64(dir[0]*dir[0]+dir[1]*dir[1])))) // float64(dir[0])
	Z := math.Atan2(float64(dir[1]), float64(dir[0]))
	quatRotation.W = float32(math.Cos(Y/2)*math.Cos(Z/2)*math.Cos(X/2) + math.Sin(Y/2)*math.Sin(Z/2)*math.Sin(X/2))
	quatRotation.X = float32(math.Cos(Y/2)*math.Cos(Z/2)*math.Sin(X/2) - math.Sin(Y/2)*math.Sin(Z/2)*math.Cos(X/2))
	quatRotation.Y = float32(math.Sin(Y/2)*math.Cos(Z/2)*math.Cos(X/2) + math.Cos(Y/2)*math.Sin(Z/2)*math.Sin(X/2))
	quatRotation.Z = float32(math.Cos(Y/2)*math.Sin(Z/2)*math.Cos(X/2) - math.Sin(Y/2)*math.Cos(Z/2)*math.Sin(X/2))

	var quatRotationConjugate common.Quaternion
	quatRotationConjugate.Copy(&quatRotation)
	quatRotationConjugate.Conjugate()
	PointSize := 5
	for _idx := 0; _idx < 1; _idx++ {
		var quatVert common.Quaternion
		var quatTmp common.Quaternion
		quatVert.SetVert(vertice[_idx*PointSize : _idx*PointSize+3])
		quatTmp.Copy(&quatRotation)
		quatTmp.Multi(&quatVert)
		quatTmp.Multi(&quatRotationConjugate)
		quatTmp.GetVert(output_vert[:])
	}

	return output_vert[:]
}
