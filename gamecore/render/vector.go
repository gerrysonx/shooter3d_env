package render

import (
	"math"

	"../common"
	"github.com/ungerik/go3d/vec3"
	"github.com/ungerik/go3d/vec4"
)

/*
func RotateVertWithQuat(vert []float32, rotation core.Vector4) {
	var vertQuat core.Vector4
	vertQuat.W = 0
	vertQuat.X = vert[0]
	vertQuat.Y = vert[1]
	vertQuat.Z = vert[2]

}
*/
func SetConeOffset(vertice []float32, x_new float32, y_new float32, z_new float32,
	unit_width float32, unit_height float32, unit_depth float32) {
	far_face_ratio := float32(40.0)
	// Bottom
	offset := 0
	vertice[offset+0] = x_new - unit_width
	vertice[offset+1] = y_new + unit_height
	vertice[offset+2] = z_new - unit_depth

	vertice[offset+5] = x_new - unit_width
	vertice[offset+6] = y_new - unit_height
	vertice[offset+7] = z_new - unit_depth

	vertice[offset+10] = x_new + unit_width
	vertice[offset+11] = y_new - unit_height*far_face_ratio
	vertice[offset+12] = z_new - unit_depth*far_face_ratio

	vertice[offset+15] = x_new - unit_width
	vertice[offset+16] = y_new + unit_height
	vertice[offset+17] = z_new - unit_depth

	vertice[offset+20] = x_new + unit_width
	vertice[offset+21] = y_new - unit_height*far_face_ratio
	vertice[offset+22] = z_new - unit_depth*far_face_ratio

	vertice[offset+25] = x_new + unit_width
	vertice[offset+26] = y_new + unit_height*far_face_ratio
	vertice[offset+27] = z_new - unit_depth*far_face_ratio

	// Top
	offset = 30
	vertice[offset+0] = x_new - unit_width
	vertice[offset+1] = y_new + unit_height
	vertice[offset+2] = z_new + unit_depth

	vertice[offset+5] = x_new + unit_width
	vertice[offset+6] = y_new + unit_height*far_face_ratio
	vertice[offset+7] = z_new + unit_depth*far_face_ratio

	vertice[offset+10] = x_new + unit_width
	vertice[offset+11] = y_new - unit_height*far_face_ratio
	vertice[offset+12] = z_new + unit_depth*far_face_ratio

	vertice[offset+15] = x_new - unit_width
	vertice[offset+16] = y_new + unit_height
	vertice[offset+17] = z_new + unit_depth

	vertice[offset+20] = x_new + unit_width
	vertice[offset+21] = y_new - unit_height*far_face_ratio
	vertice[offset+22] = z_new + unit_depth*far_face_ratio

	vertice[offset+25] = x_new - unit_width
	vertice[offset+26] = y_new - unit_height
	vertice[offset+27] = z_new + unit_depth

	// Left
	offset = 60
	vertice[offset+0] = x_new - unit_width
	vertice[offset+1] = y_new - unit_height
	vertice[offset+2] = z_new + unit_depth

	vertice[offset+5] = x_new - unit_width
	vertice[offset+6] = y_new - unit_height
	vertice[offset+7] = z_new - unit_depth

	vertice[offset+10] = x_new - unit_width
	vertice[offset+11] = y_new + unit_height
	vertice[offset+12] = z_new - unit_depth

	vertice[offset+15] = x_new - unit_width
	vertice[offset+16] = y_new - unit_height
	vertice[offset+17] = z_new + unit_depth

	vertice[offset+20] = x_new - unit_width
	vertice[offset+21] = y_new + unit_height
	vertice[offset+22] = z_new - unit_depth

	vertice[offset+25] = x_new - unit_width
	vertice[offset+26] = y_new + unit_height
	vertice[offset+27] = z_new + unit_depth

	// Right
	offset = 90
	vertice[offset+0] = x_new + unit_width
	vertice[offset+1] = y_new - unit_height*far_face_ratio
	vertice[offset+2] = z_new + unit_depth*far_face_ratio

	vertice[offset+5] = x_new + unit_width
	vertice[offset+6] = y_new + unit_height*far_face_ratio
	vertice[offset+7] = z_new + unit_depth*far_face_ratio

	vertice[offset+10] = x_new + unit_width
	vertice[offset+11] = y_new + unit_height*far_face_ratio
	vertice[offset+12] = z_new - unit_depth*far_face_ratio

	vertice[offset+15] = x_new + unit_width
	vertice[offset+16] = y_new - unit_height*far_face_ratio
	vertice[offset+17] = z_new + unit_depth*far_face_ratio

	vertice[offset+20] = x_new + unit_width
	vertice[offset+21] = y_new + unit_height*far_face_ratio
	vertice[offset+22] = z_new - unit_depth*far_face_ratio

	vertice[offset+25] = x_new + unit_width
	vertice[offset+26] = y_new - unit_height*far_face_ratio
	vertice[offset+27] = z_new - unit_depth*far_face_ratio

	// Back
	offset = 120
	vertice[offset+0] = x_new - unit_width
	vertice[offset+1] = y_new - unit_height
	vertice[offset+2] = z_new + unit_depth

	vertice[offset+5] = x_new - unit_width
	vertice[offset+6] = y_new - unit_height
	vertice[offset+7] = z_new - unit_depth

	vertice[offset+10] = x_new + unit_width
	vertice[offset+11] = y_new - unit_height*far_face_ratio
	vertice[offset+12] = z_new - unit_depth*far_face_ratio

	vertice[offset+15] = x_new - unit_width
	vertice[offset+16] = y_new - unit_height
	vertice[offset+17] = z_new + unit_depth

	vertice[offset+20] = x_new + unit_width
	vertice[offset+21] = y_new - unit_height*far_face_ratio
	vertice[offset+22] = z_new - unit_depth*far_face_ratio

	vertice[offset+25] = x_new + unit_width
	vertice[offset+26] = y_new - unit_height*far_face_ratio
	vertice[offset+27] = z_new + unit_depth*far_face_ratio

	// Front

	offset = 150
	vertice[offset+0] = x_new - unit_width
	vertice[offset+1] = y_new + unit_height
	vertice[offset+2] = z_new + unit_depth

	vertice[offset+5] = x_new + unit_width
	vertice[offset+6] = y_new + unit_height*far_face_ratio
	vertice[offset+7] = z_new + unit_depth*far_face_ratio

	vertice[offset+10] = x_new + unit_width
	vertice[offset+11] = y_new + unit_height*far_face_ratio
	vertice[offset+12] = z_new - unit_depth*far_face_ratio

	vertice[offset+15] = x_new - unit_width
	vertice[offset+16] = y_new + unit_height
	vertice[offset+17] = z_new + unit_depth

	vertice[offset+20] = x_new + unit_width
	vertice[offset+21] = y_new + unit_height*far_face_ratio
	vertice[offset+22] = z_new - unit_depth*far_face_ratio

	vertice[offset+25] = x_new - unit_width
	vertice[offset+26] = y_new + unit_height
	vertice[offset+27] = z_new - unit_depth
}

func SetOffset(vertice []float32, x_new float32, y_new float32, z_new float32,
	unit_width float32, unit_height float32, unit_depth float32) {

	// Bottom
	offset := 0
	vertice[offset+0] = x_new - unit_width
	vertice[offset+1] = y_new + unit_height
	vertice[offset+2] = z_new - unit_depth

	vertice[offset+5] = x_new - unit_width
	vertice[offset+6] = y_new - unit_height
	vertice[offset+7] = z_new - unit_depth

	vertice[offset+10] = x_new + unit_width
	vertice[offset+11] = y_new - unit_height
	vertice[offset+12] = z_new - unit_depth

	vertice[offset+15] = x_new - unit_width
	vertice[offset+16] = y_new + unit_height
	vertice[offset+17] = z_new - unit_depth

	vertice[offset+20] = x_new + unit_width
	vertice[offset+21] = y_new - unit_height
	vertice[offset+22] = z_new - unit_depth

	vertice[offset+25] = x_new + unit_width
	vertice[offset+26] = y_new + unit_height
	vertice[offset+27] = z_new - unit_depth

	// Top
	offset = 30
	vertice[offset+0] = x_new - unit_width
	vertice[offset+1] = y_new + unit_height
	vertice[offset+2] = z_new + unit_depth

	vertice[offset+5] = x_new + unit_width
	vertice[offset+6] = y_new + unit_height
	vertice[offset+7] = z_new + unit_depth

	vertice[offset+10] = x_new + unit_width
	vertice[offset+11] = y_new - unit_height
	vertice[offset+12] = z_new + unit_depth

	vertice[offset+15] = x_new - unit_width
	vertice[offset+16] = y_new + unit_height
	vertice[offset+17] = z_new + unit_depth

	vertice[offset+20] = x_new + unit_width
	vertice[offset+21] = y_new - unit_height
	vertice[offset+22] = z_new + unit_depth

	vertice[offset+25] = x_new - unit_width
	vertice[offset+26] = y_new - unit_height
	vertice[offset+27] = z_new + unit_depth

	// Left
	offset = 60
	vertice[offset+0] = x_new - unit_width
	vertice[offset+1] = y_new - unit_height
	vertice[offset+2] = z_new + unit_depth

	vertice[offset+5] = x_new - unit_width
	vertice[offset+6] = y_new - unit_height
	vertice[offset+7] = z_new - unit_depth

	vertice[offset+10] = x_new - unit_width
	vertice[offset+11] = y_new + unit_height
	vertice[offset+12] = z_new - unit_depth

	vertice[offset+15] = x_new - unit_width
	vertice[offset+16] = y_new - unit_height
	vertice[offset+17] = z_new + unit_depth

	vertice[offset+20] = x_new - unit_width
	vertice[offset+21] = y_new + unit_height
	vertice[offset+22] = z_new - unit_depth

	vertice[offset+25] = x_new - unit_width
	vertice[offset+26] = y_new + unit_height
	vertice[offset+27] = z_new + unit_depth

	// Right
	offset = 90
	vertice[offset+0] = x_new + unit_width
	vertice[offset+1] = y_new - unit_height
	vertice[offset+2] = z_new + unit_depth

	vertice[offset+5] = x_new + unit_width
	vertice[offset+6] = y_new + unit_height
	vertice[offset+7] = z_new + unit_depth

	vertice[offset+10] = x_new + unit_width
	vertice[offset+11] = y_new + unit_height
	vertice[offset+12] = z_new - unit_depth

	vertice[offset+15] = x_new + unit_width
	vertice[offset+16] = y_new - unit_height
	vertice[offset+17] = z_new + unit_depth

	vertice[offset+20] = x_new + unit_width
	vertice[offset+21] = y_new + unit_height
	vertice[offset+22] = z_new - unit_depth

	vertice[offset+25] = x_new + unit_width
	vertice[offset+26] = y_new - unit_height
	vertice[offset+27] = z_new - unit_depth

	// Back
	offset = 120
	vertice[offset+0] = x_new - unit_width
	vertice[offset+1] = y_new - unit_height
	vertice[offset+2] = z_new + unit_depth

	vertice[offset+5] = x_new - unit_width
	vertice[offset+6] = y_new - unit_height
	vertice[offset+7] = z_new - unit_depth

	vertice[offset+10] = x_new + unit_width
	vertice[offset+11] = y_new - unit_height
	vertice[offset+12] = z_new - unit_depth

	vertice[offset+15] = x_new - unit_width
	vertice[offset+16] = y_new - unit_height
	vertice[offset+17] = z_new + unit_depth

	vertice[offset+20] = x_new + unit_width
	vertice[offset+21] = y_new - unit_height
	vertice[offset+22] = z_new - unit_depth

	vertice[offset+25] = x_new + unit_width
	vertice[offset+26] = y_new - unit_height
	vertice[offset+27] = z_new + unit_depth

	// Front

	offset = 150
	vertice[offset+0] = x_new - unit_width
	vertice[offset+1] = y_new + unit_height
	vertice[offset+2] = z_new + unit_depth

	vertice[offset+5] = x_new + unit_width
	vertice[offset+6] = y_new + unit_height
	vertice[offset+7] = z_new + unit_depth

	vertice[offset+10] = x_new + unit_width
	vertice[offset+11] = y_new + unit_height
	vertice[offset+12] = z_new - unit_depth

	vertice[offset+15] = x_new - unit_width
	vertice[offset+16] = y_new + unit_height
	vertice[offset+17] = z_new + unit_depth

	vertice[offset+20] = x_new + unit_width
	vertice[offset+21] = y_new + unit_height
	vertice[offset+22] = z_new - unit_depth

	vertice[offset+25] = x_new - unit_width
	vertice[offset+26] = y_new + unit_height
	vertice[offset+27] = z_new - unit_depth
}

func OffsetObject(vertice []float32, x_new float32, y_new float32, z_new float32,
	unit_width float32, unit_height float32, unit_depth float32) {

	// Bottom
	offset := 0
	vertice[offset+0] += x_new - unit_width
	vertice[offset+1] += y_new + unit_height
	vertice[offset+2] += z_new - unit_depth

	vertice[offset+5] += x_new - unit_width
	vertice[offset+6] += y_new - unit_height
	vertice[offset+7] += z_new - unit_depth

	vertice[offset+10] += x_new + unit_width
	vertice[offset+11] += y_new - unit_height
	vertice[offset+12] += z_new - unit_depth

	vertice[offset+15] += x_new - unit_width
	vertice[offset+16] += y_new + unit_height
	vertice[offset+17] += z_new - unit_depth

	vertice[offset+20] += x_new + unit_width
	vertice[offset+21] += y_new - unit_height
	vertice[offset+22] += z_new - unit_depth

	vertice[offset+25] += x_new + unit_width
	vertice[offset+26] += y_new + unit_height
	vertice[offset+27] += z_new - unit_depth

	// Top
	offset = 30
	vertice[offset+0] += x_new - unit_width
	vertice[offset+1] += y_new + unit_height
	vertice[offset+2] += z_new + unit_depth

	vertice[offset+5] += x_new + unit_width
	vertice[offset+6] += y_new + unit_height
	vertice[offset+7] += z_new + unit_depth

	vertice[offset+10] += x_new + unit_width
	vertice[offset+11] += y_new - unit_height
	vertice[offset+12] += z_new + unit_depth

	vertice[offset+15] += x_new - unit_width
	vertice[offset+16] += y_new + unit_height
	vertice[offset+17] += z_new + unit_depth

	vertice[offset+20] += x_new + unit_width
	vertice[offset+21] += y_new - unit_height
	vertice[offset+22] += z_new + unit_depth

	vertice[offset+25] += x_new - unit_width
	vertice[offset+26] += y_new - unit_height
	vertice[offset+27] += z_new + unit_depth

	// Left
	offset = 60
	vertice[offset+0] += x_new - unit_width
	vertice[offset+1] += y_new - unit_height
	vertice[offset+2] += z_new + unit_depth

	vertice[offset+5] += x_new - unit_width
	vertice[offset+6] += y_new - unit_height
	vertice[offset+7] += z_new - unit_depth

	vertice[offset+10] += x_new - unit_width
	vertice[offset+11] += y_new + unit_height
	vertice[offset+12] += z_new - unit_depth

	vertice[offset+15] += x_new - unit_width
	vertice[offset+16] += y_new - unit_height
	vertice[offset+17] += z_new + unit_depth

	vertice[offset+20] += x_new - unit_width
	vertice[offset+21] += y_new + unit_height
	vertice[offset+22] += z_new - unit_depth

	vertice[offset+25] += x_new - unit_width
	vertice[offset+26] += y_new + unit_height
	vertice[offset+27] += z_new + unit_depth

	// Right
	offset = 90
	vertice[offset+0] += x_new + unit_width
	vertice[offset+1] += y_new - unit_height
	vertice[offset+2] += z_new + unit_depth

	vertice[offset+5] += x_new + unit_width
	vertice[offset+6] += y_new + unit_height
	vertice[offset+7] += z_new + unit_depth

	vertice[offset+10] += x_new + unit_width
	vertice[offset+11] += y_new + unit_height
	vertice[offset+12] += z_new - unit_depth

	vertice[offset+15] += x_new + unit_width
	vertice[offset+16] += y_new - unit_height
	vertice[offset+17] += z_new + unit_depth

	vertice[offset+20] += x_new + unit_width
	vertice[offset+21] += y_new + unit_height
	vertice[offset+22] += z_new - unit_depth

	vertice[offset+25] += x_new + unit_width
	vertice[offset+26] += y_new - unit_height
	vertice[offset+27] += z_new - unit_depth

	// Back
	offset = 120
	vertice[offset+0] += x_new - unit_width
	vertice[offset+1] += y_new - unit_height
	vertice[offset+2] += z_new + unit_depth

	vertice[offset+5] += x_new - unit_width
	vertice[offset+6] += y_new - unit_height
	vertice[offset+7] += z_new - unit_depth

	vertice[offset+10] += x_new + unit_width
	vertice[offset+11] += y_new - unit_height
	vertice[offset+12] += z_new - unit_depth

	vertice[offset+15] += x_new - unit_width
	vertice[offset+16] += y_new - unit_height
	vertice[offset+17] += z_new + unit_depth

	vertice[offset+20] += x_new + unit_width
	vertice[offset+21] += y_new - unit_height
	vertice[offset+22] += z_new - unit_depth

	vertice[offset+25] += x_new + unit_width
	vertice[offset+26] += y_new - unit_height
	vertice[offset+27] += z_new + unit_depth

	// Front

	offset = 150
	vertice[offset+0] += x_new - unit_width
	vertice[offset+1] += y_new + unit_height
	vertice[offset+2] += z_new + unit_depth

	vertice[offset+5] += x_new + unit_width
	vertice[offset+6] += y_new + unit_height
	vertice[offset+7] += z_new + unit_depth

	vertice[offset+10] += x_new + unit_width
	vertice[offset+11] += y_new + unit_height
	vertice[offset+12] += z_new - unit_depth

	vertice[offset+15] += x_new - unit_width
	vertice[offset+16] += y_new + unit_height
	vertice[offset+17] += z_new + unit_depth

	vertice[offset+20] += x_new + unit_width
	vertice[offset+21] += y_new + unit_height
	vertice[offset+22] += z_new - unit_depth

	vertice[offset+25] += x_new - unit_width
	vertice[offset+26] += y_new + unit_height
	vertice[offset+27] += z_new - unit_depth
}

func UpdatePosWithRotation(vertice []float32, pos vec3.T, extent vec3.T, rotation vec4.T) {
	unit_width := float32(extent[0])  // x
	unit_height := float32(extent[1]) // z
	unit_depth := float32(extent[2])  // y
	SetOffset(vertice, 0, 0, 0, unit_width, unit_height, unit_depth)

	var quatRotation common.Quaternion
	var quatRotationConjugate common.Quaternion
	quatRotation.W = rotation[0]
	quatRotation.X = rotation[1]
	quatRotation.Y = rotation[2]
	quatRotation.Z = rotation[3]
	quatRotationConjugate.Copy(&quatRotation)
	quatRotationConjugate.Conjugate()

	// Rotate every vertex using quaternion
	PointSize := 5
	for _idx := 0; _idx < 36; _idx++ {
		var quatVert common.Quaternion
		var quatTmp common.Quaternion
		quatVert.SetVert(vertice[_idx*PointSize : _idx*PointSize+3])
		quatTmp.Copy(&quatRotation)
		quatTmp.Multi(&quatVert)
		quatTmp.Multi(&quatRotationConjugate)
		quatTmp.GetVert(vertice[_idx*PointSize : _idx*PointSize+3])
	}

	OffsetObject(vertice, pos[0], pos[1], pos[2], 0, 0, 0)

}

func UpdatePos(vertice []float32, pos vec3.T, extent vec3.T) {
	unit_width := float32(extent[0])  // x
	unit_height := float32(extent[1]) // z
	unit_depth := float32(extent[2])  // y
	x_new := pos[0]
	y_new := pos[1]
	z_new := pos[2]

	SetOffset(vertice, x_new, y_new, z_new, unit_width, unit_height, unit_depth)
}

func UpdatePosDir(vertice []float32, pos vec3.T, dir vec3.T) {
	dir.Normalize()
	unit_width := float32(500) // x
	unit_height := float32(10) // z
	unit_depth := float32(10)  // y

	ahead_dist := unit_width + 50
	x_new := pos[0] + dir[0]*ahead_dist
	y_new := pos[1] + dir[1]*ahead_dist
	z_new := pos[2] + dir[2]*ahead_dist

	SetConeOffset(vertice, 0, 0, 0, unit_width, unit_height, unit_depth)
	// Calculate the rotation quaternion b = dir, b*a^-1
	theta := math.Acos(float64(dir[0]))
	if dir[1] < 0 {
		theta = math.Pi*2 - theta
	}

	var quatRotation common.Quaternion

	half_cos_theta := float32(math.Cos(theta / 2.0))
	half_sin_theta := float32(math.Sin(theta / 2.0))
	quatRotation.W = half_cos_theta
	quatRotation.X = 0
	quatRotation.Y = 0
	quatRotation.Z = half_sin_theta

	var quatRotationConjugate common.Quaternion
	quatRotationConjugate.Copy(&quatRotation)
	quatRotationConjugate.Conjugate()

	PointSize := 5
	for _idx := 0; _idx < 36; _idx++ {
		var quatVert common.Quaternion
		var quatTmp common.Quaternion
		quatVert.SetVert(vertice[_idx*PointSize : _idx*PointSize+3])
		quatTmp.Copy(&quatRotation)
		quatTmp.Multi(&quatVert)
		quatTmp.Multi(&quatRotationConjugate)
		quatTmp.GetVert(vertice[_idx*PointSize : _idx*PointSize+3])
	}
	OffsetObject(vertice, x_new, y_new, z_new, 0, 0, 0)
}
