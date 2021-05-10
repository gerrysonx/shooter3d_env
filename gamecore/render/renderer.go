package render

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"../common"
	"../core"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ungerik/go3d/vec3"
)

type Renderer struct {
	// Render related
	window   *glfw.Window
	program  uint32
	textures map[string]uint32
	actors   map[int32]*Actor

	_view_angle  float32
	_view_dist   float32
	_view_radius float64
	_fov         float32
	_wh_ratio    float32

	game *core.Game
}

const windowWidth = 1000
const windowHeight = 1000

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

var vertexShader = `
#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 330
uniform sampler2D tex;
uniform vec3 camp_color;

in vec2 fragTexCoord;
out vec4 outputColor;
void main() {
	outputColor = texture(tex, fragTexCoord);
	outputColor.rgb = outputColor.rgb * camp_color;
}
` + "\x00"

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	0, 1000.0, 0, 0.0, 1.0,
	1000.0, 1000.0, 0, 1.0, 1.0,
	1000.0, 0, 0, 1.0, 0.0,
	0, 1000.0, 0, 0.0, 1.0,
	1000.0, 0, 0, 1.0, 0.0,
	0, 0, 0, 0.0, 0.0}

func key_call_back(w *glfw.Window, char rune) {
	ch_input := string(char)
	switch ch_input {
	case "1":
		fmt.Println("Key 1 is pressed.")
		core.GameInst.SelfHeroes[0].UseSkill(0)

	case "2":
		fmt.Println("Key 2 is pressed.")
		core.GameInst.SelfHeroes[0].UseSkill(1)

	case "3":
		fmt.Println("Key 3 is pressed.")
		core.GameInst.SelfHeroes[0].UseSkill(2)

	case "4":
		fmt.Println("Key 4 is pressed.")
		core.GameInst.SelfHeroes[0].UseSkill(3)

	case "5":
		core.GameInst.Init()
		game_state_str := core.GameInst.DumpVarPlayerGameState()
		fmt.Printf("%d@%s\n", 9999999, game_state_str)
	case "7":
		if RendererInst._view_radius > 0 {
			RendererInst._view_radius -= 10
		}

	case "8":
		RendererInst._view_radius += 10

	case "9":
		RendererInst._view_angle += 0.5

	case "0":
		RendererInst._view_angle -= 0.5

	case "-":
		RendererInst._view_dist -= 100

	case "=":
		RendererInst._view_dist += 100
	case "q":
		RendererInst._fov -= 0.1

	case "w":
		RendererInst._fov += 0.1
	}

}

func (renderer *Renderer) get_view_pos() vec3.T {
	view_angle := renderer._view_angle / 20.0
	view_angle = view_angle * math.Pi
	inner_circle_radius := renderer._view_radius
	pos := vec3.T{
		float32(math.Cos(float64(view_angle)) * inner_circle_radius),
		float32(math.Sin(float64(view_angle)) * inner_circle_radius),
		renderer._view_dist}
	return pos
}

func mouse_button_call_back(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	x, y := w.GetCursorPos()

	// Calculate the target pos
	// 1. calculate the rotation, from vector 0 0 -1, to  new view dir

	new_view_dir := RendererInst.get_view_pos()
	_view_pos := new_view_dir
	new_view_dir.Invert()
	new_view_dir.Normalize()
	old_view_dir := vec3.T{0, 0, -1}
	rotate_axis := vec3.Cross(&old_view_dir, &new_view_dir)
	cos_angle := vec3.Dot(&old_view_dir, &new_view_dir)
	rotate_angle := float32(math.Acos(float64(cos_angle)))
	rotate_axis.Normalize()
	// 2. apply the rotation to vector
	var view_ray_vec vec3.T
	view_ray_vec_z := 100.0
	half_height := math.Tan(float64(RendererInst._fov/2)) * view_ray_vec_z
	half_width := half_height * float64(RendererInst._wh_ratio)
	view_ray_vec_x := (float64(x)/windowWidth - 0.5) * half_width * 2
	view_ray_vec_y := (float64(y)/windowHeight - 0.5) * half_height * 2

	view_ray_vec[0] = float32(-view_ray_vec_x)
	view_ray_vec[1] = float32(view_ray_vec_y)
	view_ray_vec[2] = float32(-view_ray_vec_z)
	view_ray_vec.Normalize()
	rotated_view_ray_vec := common.RotateVector(view_ray_vec, rotate_axis, rotate_angle)

	target_pos := rotated_view_ray_vec
	target_pos.Scale(10000).Add(&_view_pos)

	// 3. check cross point
	collide_pos, _ret := core.GetCollidePosT(_view_pos, target_pos, []vec3.T{vec3.T{0, -30000, 20}, vec3.T{-30000, 30000, 20}, {30000, 30000, 20}})
	var logic_pos_x, logic_pos_y float32
	if _ret {
		logic_pos_x, logic_pos_y = collide_pos[0], collide_pos[1]
	}

	switch {
	case action == glfw.Release && button == glfw.MouseButtonLeft:

		//	fmt.Println("Left mouse button is released.", button, action, mod, x, y)
		if core.GameInst.ManualCtrlEnemy {
			core.GameInst.OppoHeroes[0].SetTargetPos(logic_pos_x, logic_pos_y)
		} else {
			if core.GameInst.SelfHeroes[0] != nil {
				core.GameInst.SelfHeroes[0].SetTargetPos(logic_pos_x, logic_pos_y)
			}

		}

	case action == glfw.Release && button == glfw.MouseButtonRight:
		if core.GameInst.SelfHeroes[0] != nil {
			core.GameInst.SelfHeroes[0].SetSkillTargetPos(logic_pos_x, logic_pos_y)
		}

	}
}

func (renderer *Renderer) Release() {
	glfw.Terminate()
}

func (renderer *Renderer) InitRenderEnv(game *core.Game) {
	renderer.game = game
	renderer._view_angle = 0.0
	renderer._view_dist = 3000.0
	renderer._view_radius = 0.0
	renderer._fov = mgl32.DegToRad(45.0)
	renderer._wh_ratio = float32(windowWidth) / windowHeight

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	renderer.window = window
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)
	renderer.program = program

	// Init camera
	projection := mgl32.Perspective(renderer._fov, renderer._wh_ratio, 0.1, 10000.0) // mgl32.Ortho2D(-1, 1, -1, 1) //mgl32.Ident4() //
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// camera := mgl32.LookAtV(mgl32.Vec3{windowWidth / 20, windowHeight / 20, -50}, mgl32.Vec3{windowWidth / 2, windowHeight / 2, 0}, mgl32.Vec3{0, -1, 0}) // mgl32.Ortho(0, 1000, 0, 1000, -1, 1) //
	// cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	// gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	renderer.LoadCfgFolder()

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	// ---------------------------------------------------------------
	// core.GameInst.SetRenderParam(window, program, vao, texture, tex_footman, vao_footman, vbo_footman, vert_footman,
	//	tex_bullet, vao_bullet, vbo_bullet, vert_bullet,
	//	tex_hero, vao_hero, vbo_hero, vert_hero)

	window.SetMouseButtonCallback(mouse_button_call_back)
	window.SetCharCallback(key_call_back)
}

type JsonInfo struct {
	AttackRange float32
	AttackFreq  float64
	Health      float32
	Damage      float32
	Stub1       float32
	Speed       float32
	ViewRange   float32
	Id          int32
	Skills      []int32
	Name        string
	Texture     string
	Size        float32
	Extent      []float32
}

func (renderer *Renderer) LoadActor(id int32, full_path string) *Actor {
	file_handle, err := os.Open(full_path)
	if err != nil {
		return nil
	}

	defer file_handle.Close()

	buffer := make([]byte, 1000)
	read_count, err := file_handle.Read(buffer)
	if err != nil {
		return nil
	}

	buffer = buffer[:read_count]
	var jsoninfo JsonInfo

	if err = json.Unmarshal(buffer, &jsoninfo); err == nil {
		actor := new(Actor)
		actor._renderer = renderer

		actor._size[0] = jsoninfo.Extent[0]
		actor._size[1] = jsoninfo.Extent[1]
		actor._size[2] = jsoninfo.Extent[2]

		actor.LoadModel(jsoninfo.Texture)

		return actor
	}

	return nil
}

func (renderer *Renderer) LoadCfgFolder() {

	root_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	config_file_folder := fmt.Sprintf("%s/cfg/heroes", root_dir)

	// Load all skill configs under folder
	renderer.textures = make(map[string]uint32)
	renderer.actors = make(map[int32]*Actor)

	files, err := ioutil.ReadDir(config_file_folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		cfg_file_name := f.Name()
		segs := strings.Split(cfg_file_name, ".")
		id, _ := strconv.Atoi(segs[0])
		id32 := int32(id)
		// fmt.Println(cfg_file_name)
		cfg_full_file_name := fmt.Sprintf("%s/%s", config_file_folder, cfg_file_name)
		if _, ok := renderer.actors[id32]; ok == false {
			renderer.actors[id32] = renderer.LoadActor(id32, cfg_full_file_name)
		}
	}
}

func update_health_bar_pos(vertice []float32, x_new float32, y_new float32, unit_width float32) {
	y_scale := float32(3)
	vertice[0] = x_new - unit_width
	vertice[1] = y_new + y_scale
	vertice[5] = x_new + unit_width
	vertice[6] = y_new + y_scale
	vertice[10] = x_new + unit_width
	vertice[11] = y_new - y_scale

	vertice[15] = x_new - unit_width
	vertice[16] = y_new + y_scale
	vertice[20] = x_new + unit_width
	vertice[21] = y_new - y_scale
	vertice[25] = x_new - unit_width
	vertice[26] = y_new - y_scale
}

func update_dir_vert(vertice []float32, x_dir float32, y_dir float32, x_src float32, y_src float32) {
	var scale_val float32
	scale_val = 20.0
	vertice[0] = x_src + x_dir*(scale_val+10)
	vertice[1] = y_src + y_dir*(scale_val+10)

	vertice[5] = x_src + y_dir*scale_val
	vertice[6] = y_src - x_dir*scale_val

	vertice[10] = x_src - y_dir*scale_val
	vertice[11] = y_src + x_dir*scale_val
}

func (renderer *Renderer) Render() {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	projection := mgl32.Perspective(renderer._fov, renderer._wh_ratio, 0.1, 10000.0) // mgl32.Ortho2D(-1, 1, -1, 1) //mgl32.Ident4() //
	projectionUniform := gl.GetUniformLocation(renderer.program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	pos := renderer.get_view_pos()

	camera := mgl32.LookAtV(mgl32.Vec3{
		pos[0],
		pos[1],
		pos[2]},
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, -1, 0}) // mgl32.Ortho(0, 1000, 0, 1000, -1, 1) //
	cameraUniform := gl.GetUniformLocation(renderer.program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	// Render
	gl.UseProgram(renderer.program)
	// gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// We should draw Static Objects and the Actors
	// 1. Actors

	for _, v := range renderer.game.BattleField.Props {
		actor := renderer.actors[9999]
		actor.DrawStatic(v)
	}

	// 1. Actors
	for _, v := range renderer.game.BattleUnits {
		if v.Health() > 0 {
			actor := renderer.actors[v.GetId()]
			actor.Draw(v)
		}
	}

	// Maintenance
	renderer.window.SwapBuffers()
	glfw.PollEvents()
}

var RendererInst Renderer
