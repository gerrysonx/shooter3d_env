package common

import (
	"github.com/ungerik/go3d/vec3"
	"github.com/ungerik/go3d/vec4"
)

func UpdatePosWithRotation(vertice_new *[]float32, pos vec3.T, extent vec3.T, rotation vec4.T) {

	vertice := SetVertice()

	unit_width := float32(extent[0])  // x
	unit_height := float32(extent[1]) // z
	unit_depth := float32(extent[2])  // y
	SetOffset(vertice, 0, 0, 0, unit_width, unit_height, unit_depth)

	var quatRotation Quaternion
	var quatRotationConjugate Quaternion
	quatRotation.W = rotation[0]
	quatRotation.X = rotation[1]
	quatRotation.Y = rotation[2]
	quatRotation.Z = rotation[3]
	quatRotationConjugate.Copy(&quatRotation)
	quatRotationConjugate.Conjugate()

	// Rotate every vertex using quaternion
	PointSize := 5
	for _idx := 0; _idx < 36; _idx++ {
		var quatVert Quaternion
		var quatTmp Quaternion
		quatVert.SetVert(vertice[_idx*PointSize : _idx*PointSize+3])
		quatTmp.Copy(&quatRotation)
		quatTmp.Multi(&quatVert)
		quatTmp.Multi(&quatRotationConjugate)
		quatTmp.GetVert(vertice[_idx*PointSize : _idx*PointSize+3])
	}

	OffsetObject(vertice, pos[0], pos[1], pos[2], 0, 0, 0)

	*vertice_new = vertice

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

func SetVertice() []float32 {
	vertice := []float32{
		//  X, Y, Z, U, V
		// Bottom
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 0.0,

		// Top
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 0.0,

		// Left
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 0.0,

		// Right
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 0.0,

		// Front
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 0.0,

		// Back
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 1.0,
		0, 0, 0, 1.0, 0.0,
		0, 0, 0, 0.0, 0.0,
	}
	return vertice
}

/*
actor._vert_buffer_len = uint32(len(actor._vert_buffer))
var tmp = reflect.TypeOf(actor._vert_buffer[0])
var tmp2 = uint32(tmp.Size())

actor._vert_buffer_size = actor._vert_buffer_len * tmp2
actor._vert_count = actor._vert_buffer_len / vert_point_size*/
