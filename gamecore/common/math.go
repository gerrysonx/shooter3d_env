package common

import (
	"math"

	"github.com/ungerik/go3d/vec3"
)

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

type Quaternion struct {
	X float32
	Y float32
	Z float32
	W float32
}

func (self *Quaternion) Copy(other *Quaternion) {
	self.W = other.W
	self.X = other.X
	self.Y = other.Y
	self.Z = other.Z
}

func (self *Quaternion) Multi(other *Quaternion) {
	w := self.W*other.W - self.X*other.X - self.Y*other.Y - self.Z*other.Z
	x := self.W*other.X + self.X*other.W + self.Y*other.Z - self.Z*other.Y
	y := self.W*other.Y + self.Y*other.W + self.Z*other.X - self.X*other.Z
	z := self.W*other.Z + self.Z*other.W + self.X*other.Y - self.Y*other.X

	self.W, self.X, self.Y, self.Z = w, x, y, z

}

func (self *Quaternion) Conjugate() {
	self.X = -self.X
	self.Y = -self.Y
	self.Z = -self.Z
}

func (self *Quaternion) SetVert(vert []float32) {
	self.W = 0
	self.X = vert[0]
	self.Y = vert[1]
	self.Z = vert[2]
}

func (self *Quaternion) SetVertT(vert vec3.T) {
	self.W = 0
	self.X = vert[0]
	self.Y = vert[1]
	self.Z = vert[2]
}

func (self *Quaternion) GetVert(vert []float32) {
	vert[0] = self.X
	vert[1] = self.Y
	vert[2] = self.Z
}
func (self *Quaternion) GetVertT(vert *vec3.T) {
	vert[0] = self.X
	vert[1] = self.Y
	vert[2] = self.Z
}

// Judge if a ray is within a certain radius of sphere
func CheckRayShot(start_pos vec3.T, dir vec3.T, target_pos vec3.T) {

}

func RotateVector(a vec3.T, rotate_axis vec3.T, rotate_angle float32) vec3.T {
	half_angle := float64(rotate_angle / 2.0)
	cos_half_angle := float32(math.Cos(half_angle))
	sin_half_angle := float32(math.Sin(half_angle))

	var quat_rotate Quaternion
	quat_rotate.W = cos_half_angle
	quat_rotate.X = sin_half_angle * rotate_axis[0]
	quat_rotate.Y = sin_half_angle * rotate_axis[1]
	quat_rotate.Z = sin_half_angle * rotate_axis[2]

	var quatRotationConjugate Quaternion
	quatRotationConjugate.Copy(&quat_rotate)
	quatRotationConjugate.Conjugate()

	var quatVert Quaternion
	var quatTmp Quaternion
	quatVert.SetVertT(a)
	quatTmp.Copy(&quat_rotate)
	quatTmp.Multi(&quatVert)
	quatTmp.Multi(&quatRotationConjugate)

	var vec_after_rotation vec3.T
	quatTmp.GetVertT(&vec_after_rotation)
	return vec_after_rotation

}
