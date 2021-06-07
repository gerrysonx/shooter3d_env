package render

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"

	"../core"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ungerik/go3d/vec3"
)

const vert_point_size = 5

type Actor struct {
	// Render related
	_tex         uint32
	_vao         uint32
	_vbo         uint32
	_vert_buffer []float32

	_vert_buffer_size uint32
	_vert_buffer_len  uint32
	_vert_count       uint32

	_size       vec3.T
	_renderer   *Renderer
	_sub_actors []*Actor
}

func (actor *Actor) LoadModel(model_name string) {
	root_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Load the texture
	texture := uint32(0)
	if _, ok := actor._renderer.textures[model_name]; ok == false {
		full_path := fmt.Sprintf("%s/%s", root_dir, model_name)
		texture, err = newTexture(full_path)
		if err != nil {
			log.Fatalln(err)
		}
		actor._renderer.textures[model_name] = texture
	} else {
		texture = actor._renderer.textures[model_name]
	}

	actor._tex = texture

	actor._vert_buffer = []float32{
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

	actor._vert_buffer_len = uint32(len(actor._vert_buffer))
	var tmp = reflect.TypeOf(actor._vert_buffer[0])
	var tmp2 = uint32(tmp.Size())

	actor._vert_buffer_size = actor._vert_buffer_len * tmp2
	actor._vert_count = actor._vert_buffer_len / vert_point_size

	var vao_hero uint32
	gl.GenVertexArrays(1, &vao_hero)
	gl.BindVertexArray(vao_hero)
	actor._vao = vao_hero

	var vbo_hero uint32
	gl.GenBuffers(1, &vbo_hero)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_hero)
	gl.BufferData(gl.ARRAY_BUFFER, int(actor._vert_buffer_size), gl.Ptr(actor._vert_buffer), gl.STATIC_DRAW)
	actor._vbo = vbo_hero

	vertAttrib := uint32(gl.GetAttribLocation(actor._renderer.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(actor._renderer.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

}

func (actor *Actor) DrawStatic(f0 *core.StaticUnit) {
	pos := f0.Pos
	actor._size[0] = f0.Extent[0] * f0.Scale[0]
	actor._size[1] = f0.Extent[1] * f0.Scale[1]
	actor._size[2] = f0.Extent[2] * f0.Scale[2]

	// 向缓冲中写入数据
	colorUniform := gl.GetUniformLocation(actor._renderer.program, gl.Str("camp_color\x00"))
	change_color := float32(0.4)
	gl.Uniform3f(colorUniform, change_color, change_color, change_color)

	// Draw actor
	if !f0.Cached {
		UpdatePosWithRotation(actor._vert_buffer, pos, actor._size, f0.Rotation)

		// Copy
		f0.Vertice = make([]float32, len(actor._vert_buffer))
		for i := 0; i < len(actor._vert_buffer); i++ {
			f0.Vertice[i] = actor._vert_buffer[i]

			switch i % 5 {
			case 0:
				if f0.Vertice[i] < f0.BB.Xmin {
					f0.BB.Xmin = f0.Vertice[i]
				}
				if f0.Vertice[i] > f0.BB.Xmax {
					f0.BB.Xmax = f0.Vertice[i]
				}
			case 1:
				if f0.Vertice[i] < f0.BB.Ymin {
					f0.BB.Ymin = f0.Vertice[i]
				}
				if f0.Vertice[i] > f0.BB.Ymax {
					f0.BB.Ymax = f0.Vertice[i]
				}
			case 2:
				if f0.Vertice[i] < f0.BB.Zmin {
					f0.BB.Zmin = f0.Vertice[i]
				}
				if f0.Vertice[i] > f0.BB.Zmax {
					f0.BB.Zmax = f0.Vertice[i]
				}
			}

		}

		f0.Cached = true

	}

	gl.BindBuffer(gl.ARRAY_BUFFER, actor._vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(actor._vert_buffer_size), gl.Ptr(f0.Vertice), gl.DYNAMIC_DRAW) //

	gl.BindVertexArray(actor._vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, actor._tex)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(actor._vert_count))

	for _, sub_actor := range actor._sub_actors {
		sub_actor.DrawStatic(f0)
	}
}

func (actor *Actor) Draw(f0 core.BaseFunc) {
	dir := f0.Viewdir()
	dir.Normalize()
	camp := f0.Camp()
	// How to bind value?
	// hp := f0.Health()
	pos := f0.Position()

	// 向缓冲中写入数据
	colorUniform := gl.GetUniformLocation(actor._renderer.program, gl.Str("camp_color\x00"))
	change_color := float32(0.0)
	if camp == 0 {
		gl.Uniform3f(colorUniform, change_color, 0.8, change_color)
	} else {
		gl.Uniform3f(colorUniform, 0.8, change_color, change_color)
	}

	// Draw actor
	UpdatePos(actor._vert_buffer, pos, actor._size)

	gl.BindBuffer(gl.ARRAY_BUFFER, actor._vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(actor._vert_buffer_size), gl.Ptr(actor._vert_buffer), gl.DYNAMIC_DRAW)

	gl.BindVertexArray(actor._vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, actor._tex)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(actor._vert_count))

	if core.GameInst.ShowFrustum {
		// Draw view pyramid
		if f0.GetId() > 0 {
			if math.Abs(float64(dir[0])) > 1e-3 || math.Abs(float64(dir[1])) > 1e-3 {

				change_color = float32(0.1)

				gl.Uniform3f(colorUniform, change_color, change_color, 0.6)

				transparencyUniform := gl.GetUniformLocation(actor._renderer.program, gl.Str("transparency\x00"))
				change_color := float32(0.3)
				gl.Uniform1f(transparencyUniform, change_color)

				UpdatePosDir(actor._vert_buffer, pos, dir, f0)

				gl.BindBuffer(gl.ARRAY_BUFFER, actor._vbo)
				gl.BufferData(gl.ARRAY_BUFFER, int(actor._vert_buffer_size), gl.Ptr(actor._vert_buffer), gl.DYNAMIC_DRAW)

				gl.BindVertexArray(actor._vao)

				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, actor._tex)

				gl.DrawArrays(gl.TRIANGLES, 0, int32(actor._vert_count))
				change_color = float32(1.0)
				gl.Uniform1f(transparencyUniform, change_color)
			}
		}
	}

	for _, sub_actor := range actor._sub_actors {
		sub_actor.Draw(f0)
	}
	// Draw hp bar
	// actor.renderer.DrawHealthBar(colorUniform, f0)

	// Draw hero direction
	/*
		update_dir_vert(actor.renderer.vert_hero_dir, dir[0], dir[1], pos[0], pos[1])

		gl.BindBuffer(gl.ARRAY_BUFFER, actor.renderer.vbo_hero_dir)
		gl.BufferData(gl.ARRAY_BUFFER, len(actor.renderer.vert_hero_dir)*4, gl.Ptr(actor.renderer.vert_hero_dir), gl.STATIC_DRAW)

		gl.BindVertexArray(actor.renderer.vao_hero_dir)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, actor._tex)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)
	*/
}

func (actor *Actor) DrawDepthMap(tex uint32) {
	// Set camera to ortho mode
	projection := mgl32.Ortho(0, 1000, 0, 1000, -1000, 1000) //mgl32.Ident4() //
	projectionUniform := gl.GetUniformLocation(actor._renderer.program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.Ident4()
	cameraUniform := gl.GetUniformLocation(actor._renderer.program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	colorUniform := gl.GetUniformLocation(actor._renderer.program, gl.Str("camp_color\x00"))
	change_color := float32(1.0)
	gl.Uniform3f(colorUniform, change_color, change_color, change_color)

	scale := float32(200)
	center_x := float32(800)
	center_y := float32(800)

	actor._vert_buffer[0] = center_x - scale
	actor._vert_buffer[1] = center_y - scale
	actor._vert_buffer[3] = 0
	actor._vert_buffer[4] = 0

	actor._vert_buffer[5] = center_x + scale
	actor._vert_buffer[6] = center_y + scale
	actor._vert_buffer[8] = 1
	actor._vert_buffer[9] = 1

	actor._vert_buffer[10] = center_x + scale
	actor._vert_buffer[11] = center_y - scale
	actor._vert_buffer[13] = 1
	actor._vert_buffer[14] = 0

	actor._vert_buffer[15] = center_x - scale
	actor._vert_buffer[16] = center_y - scale
	actor._vert_buffer[18] = 0
	actor._vert_buffer[19] = 0

	actor._vert_buffer[20] = center_x - scale
	actor._vert_buffer[21] = center_y + scale
	actor._vert_buffer[23] = 0
	actor._vert_buffer[24] = 1

	actor._vert_buffer[25] = center_x + scale
	actor._vert_buffer[26] = center_y + scale
	actor._vert_buffer[28] = 1
	actor._vert_buffer[29] = 1

	gl.BindBuffer(gl.ARRAY_BUFFER, actor._vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(actor._vert_buffer_size), gl.Ptr(actor._vert_buffer), gl.DYNAMIC_DRAW)

	gl.BindVertexArray(actor._vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(actor._vert_count))
}
