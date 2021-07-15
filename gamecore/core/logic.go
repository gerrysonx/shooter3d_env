package core

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"./nn"

	"github.com/ungerik/go3d/vec3"
)

type GameTrainState struct {
	SelfHeroPosX   float32
	SelfHeroPosY   float32
	SelfHeroHealth float32
	OppoHeroPosX   float32
	OppoHeroPosY   float32
	OppoHeroHealth float32

	SlowBuffState      float32
	SlowBuffRemainTime float32

	SelfWin int32
}

type GameMultiPlayerTrainState struct {
	SelfHero0PosX       float32
	SelfHero0PosY       float32
	SelfHero0Health     float32
	SelfHero0HealthFull float32

	SelfHero1PosX       float32
	SelfHero1PosY       float32
	SelfHero1Health     float32
	SelfHero1HealthFull float32

	OppoHeroPosX       float32
	OppoHeroPosY       float32
	OppoHeroHealth     float32
	OppoHeroHealthFull float32

	SlowBuffState      float32
	SlowBuffRemainTime float32

	SelfWin int32
}

type GameVarPlayerTrainState struct {
	SelfHeroCount       int32
	SelfHeroPosX        []float32
	SelfHeroPosY        []float32
	SelfHeroPosZ        []float32
	SelfHeroHealth      []float32
	SelfHeroHealthFull  []float32
	SelfHeroDepthMap    [][]float32
	SelfHeroDirX        []float32
	SelfHeroDirY        []float32
	SelfHeroDirZ        []float32
	SelfHeroAttkRmnTime []float32

	OppoHeroCount      int32
	OppoHeroPosX       []float32
	OppoHeroPosY       []float32
	OppoHeroPosZ       []float32
	OppoHeroHealth     []float32
	OppoHeroHealthFull []float32
	OppoHeroDirX       []float32
	OppoHeroDirY       []float32
	OppoHeroDirZ       []float32

	ModelIndex []string

	SlowBuffState      float32
	SlowBuffRemainTime float32

	SelfWin int32
}

type Game struct {
	SceneId     int
	CurrentTime float64
	LogicTime   float64
	BattleUnits []BaseFunc
	AddUnits    []BaseFunc
	BattleField *BattleField

	SelfHeroes      []HeroFunc
	OppoHeroes      []HeroFunc
	ManualCtrlEnemy bool
	MultiPlayer     bool

	roundCnt        int
	RoundAlterModel int

	train_state              GameTrainState
	multi_player_train_state GameMultiPlayerTrainState
	var_player_train_state   GameVarPlayerTrainState
	skill_targets            []SkillTarget
	skill_targets_add        []SkillTarget

	DepthMapSize int
	Paused       bool
	ShowMiniMap  bool
	ShowDepthMap bool
	ShowFrustum  bool
	LogoutStack  bool

	FlagPos          vec3.T
	FlagRadius       float32
	isSecuring       uint8
	TimeInitSecuring float64
	Secured          bool
	Time2Secure      float64
}

type TestConfig struct {
	Restricted_x   float32
	Restricted_y   float32
	Restricted_w   float32
	Restricted_h   float32
	OppoHeroes     []int32
	SelfHeroes     []int32
	SpawnAreaWidth float32
	SelfTowers     []float32
	OppoTowers     []float32
	SelfCreeps     []float32
	OppoCreeps     []float32
	Width          float32
	Height         float32
	DepthMapSize   int
}

func get_rand_pos(dist float32) []float32 {
	var rand_pos [2]float32
	_seed := rand.Float32() * math.Pi * 2
	rand_num_1 := math.Cos(float64(_seed)) * float64(dist)
	rand_num_2 := math.Sin(float64(_seed)) * float64(dist)
	rand_pos[0] = float32(rand_num_1)
	rand_pos[1] = float32(rand_num_2)
	return rand_pos[:]
}

func get_rand_pos2(dist float32) []float32 {
	var rand_pos [2]float32
	rand_section_idx := rand.Int31n(4)
	var rand_radiant float32
	rand_radiant = rand.Float32() * math.Pi / 12.0
	switch rand_section_idx {
	case 0:
		rand_pos[0] = float32(rand.Int31n(70) + 350)
		rand_pos[1] = float32(rand.Int31n(300) + 100)
		return rand_pos[:]
		rand_radiant += 0.3 * math.Pi
	case 1:
		rand_pos[0] = -float32(rand.Int31n(70) + 350)
		rand_pos[1] = float32(rand.Int31n(300) + 100)
		return rand_pos[:]
		rand_radiant += 0.8 * math.Pi
	case 2:
		rand_pos[0] = float32(rand.Int31n(70) + 350)
		rand_pos[1] = -float32(rand.Int31n(300) + 100)
		return rand_pos[:]
		rand_radiant += 1.2 * math.Pi
	case 3:
		rand_pos[0] = -float32(rand.Int31n(70) + 350)
		rand_pos[1] = -float32(rand.Int31n(300) + 100)
		return rand_pos[:]
		rand_radiant += 1.8 * math.Pi
	}

	rand_num_1 := math.Cos(float64(rand_radiant)) * float64(dist)
	rand_num_2 := math.Sin(float64(rand_radiant)) * float64(dist)
	rand_pos[0] = float32(rand_num_1)
	rand_pos[1] = float32(rand_num_2)
	return rand_pos[:]
}

func (game *Game) LoadTestCase(test_cfg_name string) {
	root_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	full_cfg_name := fmt.Sprintf("%s/%s", root_dir, test_cfg_name)

	file_handle, err := os.Open(full_cfg_name)
	if err != nil {
		return
	}

	buffer := make([]byte, 10000)
	read_count, err := file_handle.Read(buffer)
	if err != nil {
		return
	}
	buffer = buffer[:read_count]
	var testconfig TestConfig
	now := time.Now()
	rand.Seed(now.UnixNano())

	if err = json.Unmarshal(buffer, &testconfig); err == nil {
		const battle_field_width = 1000
		const battle_field_height = 1000

		game.BattleField = &BattleField{Restricted_x: testconfig.Restricted_x,
			Restricted_y: testconfig.Restricted_y,
			Restricted_w: testconfig.Restricted_w,
			Restricted_h: testconfig.Restricted_h,
			Width:        int32(testconfig.Width),
			Height:       int32(testconfig.Height)}

		game.DepthMapSize = testconfig.DepthMapSize
		full_map_cfg_name := fmt.Sprintf("%s/cfg/maps/%s", root_dir, "a.txt")
		game.BattleField.LoadProps(full_map_cfg_name)

		born_area_side_width := float64(testconfig.SpawnAreaWidth)

		/*for _, v := range game.BattleField.Props {
			pos := vec3.T{237.38586, 257.19244, 10}
			within := v.CheckWithin(pos)
			LogStr(fmt.Sprintf("In Wall: %v", within))
		}*/

		// Self heroes count
		self_heroes_count := len(testconfig.SelfHeroes)
		oppo_heroes_count := len(testconfig.OppoHeroes)

		rand_pos := get_rand_pos2(float32(born_area_side_width))
		rand_num_1, rand_num_2 := rand_pos[0], rand_pos[1]

		for idx := 0; idx < self_heroes_count; idx += 1 {
			// rand_pos := get_rand_pos2(float32(born_area_side_width))
			// rand_num_1, rand_num_2 := rand_pos[0], rand_pos[1]
			self_hero := HeroMgrInst.Spawn(testconfig.SelfHeroes[idx], int32(0),
				float32(rand_num_1),
				float32(rand_num_2),
				float32(80))
			/*self_hero := HeroMgrInst.Spawn(testconfig.SelfHeroes[idx], int32(0),
			float32(650),
			float32(450))*/
			self_hero.SetDirection(vec3.T{0, -1, 0})
			self_hero.SetViewdir(self_hero.Direction())
			LogStr(fmt.Sprintf("Spawn Locations: %v", self_hero.Position()))

			game.BattleUnits = append(game.BattleUnits, self_hero)
			game.SelfHeroes = append(game.SelfHeroes, self_hero.(HeroFunc))
		}

		// Oppo heroes count
		for idx := 0; idx < oppo_heroes_count; idx += 1 {
			// rand_pos := get_rand_pos2(float32(born_area_side_width))
			// rand_num_1, rand_num_2 := rand_pos[0], rand_pos[1]
			oppo_hero := HeroMgrInst.Spawn(testconfig.OppoHeroes[idx], int32(1),
				float32(-rand_num_1),
				float32(rand_num_2),
				float32(80))
			oppo_hero.SetDirection(vec3.T{0, -1, 0})
			oppo_hero.SetViewdir(oppo_hero.Direction())
			game.BattleUnits = append(game.BattleUnits, oppo_hero)
			game.OppoHeroes = append(game.OppoHeroes, oppo_hero.(HeroFunc))
		}

		const tower_attributes_count = 3
		self_tower_count := len(testconfig.SelfTowers) / tower_attributes_count
		oppo_tower_count := len(testconfig.OppoTowers) / tower_attributes_count
		for idx := 0; idx < self_tower_count; idx += 1 {
			tower_x := testconfig.SelfTowers[idx*tower_attributes_count+0] * battle_field_width
			tower_y := testconfig.SelfTowers[idx*tower_attributes_count+1] * battle_field_height
			tower_id := int32(testconfig.SelfTowers[idx*tower_attributes_count+2])

			self_hero := HeroMgrInst.Spawn(tower_id, int32(0), float32(tower_x), float32(tower_y))
			game.BattleUnits = append(game.BattleUnits, self_hero)
		}

		for idx := 0; idx < oppo_tower_count; idx += 1 {
			tower_x := testconfig.OppoTowers[idx*tower_attributes_count+0] * battle_field_width
			tower_y := testconfig.OppoTowers[idx*tower_attributes_count+1] * battle_field_height
			tower_id := int32(testconfig.OppoTowers[idx*tower_attributes_count+2])

			self_hero := HeroMgrInst.Spawn(tower_id, int32(1), float32(tower_x), float32(tower_y))
			game.BattleUnits = append(game.BattleUnits, self_hero)
		}

		const creep_attributes_count = 7
		self_creep_count := len(testconfig.SelfCreeps) / creep_attributes_count
		for idx := 0; idx < self_creep_count; idx += 1 {
			creep_spawn_x := testconfig.SelfCreeps[idx*creep_attributes_count+0] * battle_field_width
			creep_spawn_y := testconfig.SelfCreeps[idx*creep_attributes_count+1] * battle_field_height
			creep_target_x := testconfig.SelfCreeps[idx*creep_attributes_count+2] * battle_field_width
			creep_target_y := testconfig.SelfCreeps[idx*creep_attributes_count+3] * battle_field_height
			creep_spawn_freq := float64(testconfig.SelfCreeps[idx*creep_attributes_count+4])
			creep_id := int32(testconfig.SelfCreeps[idx*creep_attributes_count+5])
			creep_count := int32(testconfig.SelfCreeps[idx*creep_attributes_count+6])

			creep_spawn_mgr := new(CreepMgr)
			creep_spawn_mgr.SetAttackFreq(creep_spawn_freq)
			creep_spawn_mgr.SetPosition(vec3.T{creep_spawn_x, creep_spawn_y})
			creep_spawn_mgr.SetDirection(vec3.T{creep_target_x, creep_target_y})
			creep_spawn_mgr.SetCamp(int32(0))
			creep_spawn_mgr.SetHealth(1)
			creep_spawn_mgr.SetId(0)
			creep_spawn_mgr.CreepId = creep_id
			creep_spawn_mgr.CreepCount = creep_count
			game.BattleUnits = append(game.BattleUnits, creep_spawn_mgr)
		}

		oppo_creep_count := len(testconfig.OppoCreeps) / creep_attributes_count
		for idx := 0; idx < oppo_creep_count; idx += 1 {
			creep_spawn_x := testconfig.OppoCreeps[idx*creep_attributes_count+0] * battle_field_width
			creep_spawn_y := testconfig.OppoCreeps[idx*creep_attributes_count+1] * battle_field_height
			creep_target_x := testconfig.OppoCreeps[idx*creep_attributes_count+2] * battle_field_width
			creep_target_y := testconfig.OppoCreeps[idx*creep_attributes_count+3] * battle_field_height
			creep_spawn_freq := float64(testconfig.OppoCreeps[idx*creep_attributes_count+4])
			creep_id := int32(testconfig.OppoCreeps[idx*creep_attributes_count+5])
			creep_count := int32(testconfig.OppoCreeps[idx*creep_attributes_count+6])

			creep_spawn_mgr := new(CreepMgr)
			creep_spawn_mgr.SetAttackFreq(creep_spawn_freq)
			creep_spawn_mgr.SetPosition(vec3.T{creep_spawn_x, creep_spawn_y})
			creep_spawn_mgr.SetDirection(vec3.T{creep_target_x, creep_target_y})
			creep_spawn_mgr.SetCamp(int32(1))
			creep_spawn_mgr.SetHealth(1)
			creep_spawn_mgr.SetId(0)
			creep_spawn_mgr.CreepId = creep_id
			creep_spawn_mgr.CreepCount = creep_count
			game.BattleUnits = append(game.BattleUnits, creep_spawn_mgr)
		}

	} else {
		fmt.Println("Error is:", err)
	}
}

func (game *Game) Init() {
	now := time.Now()
	game.LogicTime = float64(now.UnixNano()) / 1e9

	// map_units := game.BattleField.LoadMap("./map/3_corridors.png")

	game.BattleUnits = []BaseFunc{}
	game.AddUnits = []BaseFunc{}
	game.skill_targets = []SkillTarget{}
	game.skill_targets_add = []SkillTarget{}
	game.SelfHeroes = []HeroFunc{}
	game.OppoHeroes = []HeroFunc{}
	game.multi_player_train_state.SelfHero0Health = 0
	game.multi_player_train_state.SelfHero1Health = 0
	game.multi_player_train_state.OppoHeroHealth = 0

	game.FlagPos = vec3.T{0, 0, 80}
	game.FlagRadius = 50
	game.isSecuring = 0
	game.TimeInitSecuring = -1
	game.Secured = false
	game.Time2Secure = 1

	if game.roundCnt%game.RoundAlterModel == 0 {
		for _, v := range HeroMgrInst.heroes {
			if v.name == "lusian" {
				index := game.SelectEnemyModel("weighted")
				//fmt.Printf("Selected Model: %v \n\n\n\n", index)
				v.model.InitModel(index)
			}
		}
		game.roundCnt = 0
	}
	game.roundCnt += 1

	game.LoadTestCase(fmt.Sprintf("./cfg/maps/%d.json", game.SceneId))

	LogStr("Game Inited.")
}

func (game *Game) SelectEnemyModel(mode string) string {
	root_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	file_handle, err := os.Open(fmt.Sprintf("%s/../model/model_list.json", root_dir))
	if err != nil {
		panic("Open Model File Error")
	}

	defer file_handle.Close()

	buffer := make([]byte, 100000)
	read_count, err := file_handle.Read(buffer)
	if err != nil {
		panic("Open Model File Error")
	}
	buffer = buffer[:read_count]
	var jsoninfo nn.ModelJsonInfo
	var model_index []string
	var model_score []int
	if err = json.Unmarshal(buffer, &jsoninfo); err == nil {
		length := len(jsoninfo.Model)
		for i := 0; i < length; i += 1 {
			model_index = append(model_index, jsoninfo.Model[i])
			model_score = append(model_score, jsoninfo.Score[i])
		}
	}
	// fmt.Printf("List: %v\n", err)
	// fmt.Printf("List: %v\n", model_index)
	rand.Seed(time.Now().UTC().UnixNano())
	switch mode {
	case "random":
		idx := rand.Intn(len(model_index))
		return model_index[idx]
	case "weighted":
		var score []float32
		var total_score float32
		for i := range model_score {
			model_score[i] += -1000+rand.Intn(100)
			total_score += float32(model_score[i])
		}
		for i := range model_score {
			score = append(score, float32(model_score[i]) / total_score)
		}
		//fmt.Printf("score: %v", score)
		prob := rand.Float32()
		for i, o := range score {
			prob -= o
			if prob < 0 {
				return model_index[i]
				break
			}
		}
		return "0001"
	default:
		return "0001"
	}
}

func (game *Game) AddTarget(target SkillTarget) {
	game.skill_targets_add = append(game.skill_targets_add, target)
}

func (game *Game) HandleInput() {

}

func (game *Game) HandleCallback() {

}

func (game *Game) GetGameState(reverse bool) [][]float32 {

	oppo_count := int(game.var_player_train_state.OppoHeroCount)
	self_count := int(game.var_player_train_state.SelfHeroCount)
	width := float32(game.BattleField.Width)

	var game_state [][]float32

	for i := 0; i < oppo_count; i += 1 {
		game_state = append(game_state, []float32{game.var_player_train_state.OppoHeroPosX[i] / width})
		game_state = append(game_state, []float32{game.var_player_train_state.OppoHeroPosY[i] / width})
		game_state = append(game_state, []float32{game.var_player_train_state.OppoHeroPosZ[i]/250 - 1})
		game_state = append(game_state, []float32{game.var_player_train_state.OppoHeroHealth[i] / game.var_player_train_state.OppoHeroHealthFull[i]})
		game_state = append(game_state, []float32{game.var_player_train_state.OppoHeroDirX[i]})
		game_state = append(game_state, []float32{game.var_player_train_state.OppoHeroDirY[i]})
		game_state = append(game_state, []float32{game.var_player_train_state.OppoHeroDirZ[i]})
	}

	for i := 0; i < self_count; i += 1 {
		game_state = append(game_state, []float32{game.var_player_train_state.SelfHeroPosX[i] / width})
		game_state = append(game_state, []float32{game.var_player_train_state.SelfHeroPosY[i] / width})
		game_state = append(game_state, []float32{game.var_player_train_state.SelfHeroPosZ[i]/250 - 1})
		game_state = append(game_state, []float32{game.var_player_train_state.SelfHeroHealth[i] / game.var_player_train_state.SelfHeroHealthFull[i]})
		game_state = append(game_state, []float32{game.var_player_train_state.SelfHeroDirX[i]})
		game_state = append(game_state, []float32{game.var_player_train_state.SelfHeroDirY[i]})
		game_state = append(game_state, []float32{game.var_player_train_state.SelfHeroDirZ[i]})
	}
	return game_state
}

func (game *Game) ClearGameStateData() {
	game.var_player_train_state.SelfHeroCount = 0
	game.var_player_train_state.SelfHeroPosX = game.var_player_train_state.SelfHeroPosX[:0]
	game.var_player_train_state.SelfHeroPosY = game.var_player_train_state.SelfHeroPosY[:0]
	game.var_player_train_state.SelfHeroPosZ = game.var_player_train_state.SelfHeroPosZ[:0]
	game.var_player_train_state.SelfHeroDirX = game.var_player_train_state.SelfHeroDirX[:0]
	game.var_player_train_state.SelfHeroDirY = game.var_player_train_state.SelfHeroDirY[:0]
	game.var_player_train_state.SelfHeroDirZ = game.var_player_train_state.SelfHeroDirZ[:0]
	game.var_player_train_state.SelfHeroAttkRmnTime = game.var_player_train_state.SelfHeroAttkRmnTime[:0]
	game.var_player_train_state.SelfHeroHealth = game.var_player_train_state.SelfHeroHealth[:0]
	game.var_player_train_state.SelfHeroHealthFull = game.var_player_train_state.SelfHeroHealthFull[:0]
	game.var_player_train_state.SelfHeroDepthMap = game.var_player_train_state.SelfHeroDepthMap[:0]

	game.var_player_train_state.OppoHeroCount = 0
	game.var_player_train_state.OppoHeroPosX = game.var_player_train_state.OppoHeroPosX[:0]
	game.var_player_train_state.OppoHeroPosY = game.var_player_train_state.OppoHeroPosY[:0]
	game.var_player_train_state.OppoHeroPosZ = game.var_player_train_state.OppoHeroPosZ[:0]
	game.var_player_train_state.OppoHeroDirX = game.var_player_train_state.OppoHeroDirX[:0]
	game.var_player_train_state.OppoHeroDirY = game.var_player_train_state.OppoHeroDirY[:0]
	game.var_player_train_state.OppoHeroDirZ = game.var_player_train_state.OppoHeroDirZ[:0]
	game.var_player_train_state.OppoHeroHealth = game.var_player_train_state.OppoHeroHealth[:0]
	game.var_player_train_state.OppoHeroHealthFull = game.var_player_train_state.OppoHeroHealthFull[:0]
	game.var_player_train_state.ModelIndex = game.var_player_train_state.ModelIndex[:0]
}

func (game *Game) DumpVarPlayerGameState() []byte {
	game.ClearGameStateData()
	all_oppo_heroes_dead := true
	all_self_heroes_dead := true

	game.var_player_train_state.SelfWin = 0
	game.var_player_train_state.SelfHeroCount = int32(len(game.SelfHeroes))

	game.var_player_train_state.SelfHeroDepthMap = make([][]float32, game.var_player_train_state.SelfHeroCount)
	for i := 0; i < int(game.var_player_train_state.SelfHeroCount); i += 1 {
		self_hero_unit := game.SelfHeroes[i].(BaseFunc)
		game.var_player_train_state.SelfHeroPosX = append(game.var_player_train_state.SelfHeroPosX, self_hero_unit.Position()[0])
		game.var_player_train_state.SelfHeroPosY = append(game.var_player_train_state.SelfHeroPosY, self_hero_unit.Position()[1])
		game.var_player_train_state.SelfHeroPosZ = append(game.var_player_train_state.SelfHeroPosZ, self_hero_unit.Position()[2])
		viewdir := self_hero_unit.Viewdir()
		viewdir.Normalize()
		game.var_player_train_state.SelfHeroDirX = append(game.var_player_train_state.SelfHeroDirX, viewdir[0])
		game.var_player_train_state.SelfHeroDirY = append(game.var_player_train_state.SelfHeroDirY, viewdir[1])
		game.var_player_train_state.SelfHeroDirZ = append(game.var_player_train_state.SelfHeroDirZ, viewdir[2])
		attackremaintime := 1 - float32((game.LogicTime-self_hero_unit.LastAttackTime())/self_hero_unit.AttackFreq())
		if attackremaintime < 0 {
			attackremaintime = 0
		}
		game.var_player_train_state.SelfHeroAttkRmnTime = append(game.var_player_train_state.SelfHeroAttkRmnTime, attackremaintime)
		game.var_player_train_state.SelfHeroHealth = append(game.var_player_train_state.SelfHeroHealth, self_hero_unit.Health())
		game.var_player_train_state.SelfHeroHealthFull = append(game.var_player_train_state.SelfHeroHealthFull, self_hero_unit.MaxHealth())

		depth_map := self_hero_unit.ViewDepth()
		game.var_player_train_state.SelfHeroDepthMap[i] = depth_map

		if self_hero_unit.Health() > float32(0.0) {
			all_self_heroes_dead = false
		}
	}

	game.var_player_train_state.OppoHeroCount = int32(len(game.OppoHeroes))
	for i := 0; i < int(game.var_player_train_state.OppoHeroCount); i += 1 {
		oppo_hero_unit := game.OppoHeroes[i].(BaseFunc)
		game.var_player_train_state.OppoHeroPosX = append(game.var_player_train_state.OppoHeroPosX, oppo_hero_unit.Position()[0])
		game.var_player_train_state.OppoHeroPosY = append(game.var_player_train_state.OppoHeroPosY, oppo_hero_unit.Position()[1])
		game.var_player_train_state.OppoHeroPosZ = append(game.var_player_train_state.OppoHeroPosZ, oppo_hero_unit.Position()[2])
		viewdir := oppo_hero_unit.Viewdir()
		viewdir.Normalize()
		game.var_player_train_state.OppoHeroDirX = append(game.var_player_train_state.OppoHeroDirX, viewdir[0])
		game.var_player_train_state.OppoHeroDirY = append(game.var_player_train_state.OppoHeroDirY, viewdir[1])
		game.var_player_train_state.OppoHeroDirZ = append(game.var_player_train_state.OppoHeroDirZ, viewdir[2])
		game.var_player_train_state.OppoHeroHealth = append(game.var_player_train_state.OppoHeroHealth, oppo_hero_unit.Health())
		game.var_player_train_state.OppoHeroHealthFull = append(game.var_player_train_state.OppoHeroHealthFull, oppo_hero_unit.MaxHealth())
		if oppo_hero_unit.Health() > float32(0.0) {
			all_oppo_heroes_dead = false
		}
		game.var_player_train_state.ModelIndex = append(game.var_player_train_state.ModelIndex, oppo_hero_unit.GetModelIndex())
	}

	// if game.Secured {
	// 	game.var_player_train_state.SelfWin = 0
	// } else
	if all_oppo_heroes_dead {
		game.var_player_train_state.SelfWin = 1
	} else if all_self_heroes_dead {
		game.var_player_train_state.SelfWin = -1
	}

	e, err := json.Marshal(game.var_player_train_state)
	/*
		var unmarshaled_gamestate GameVarPlayerTrainState
		json.Unmarshal(e, &unmarshaled_gamestate)
	*/
	if err != nil {
		return []byte(fmt.Sprintf("Marshal train_state failed.%v", game.multi_player_train_state))
	}

	return e
}

func (game *Game) Tick(gap_time float64) {
	//fmt.Printf("->gameTick, len(game.BattleUnits) is:%d, logictime:%f\n", len(game.BattleUnits), game.LogicTime)
	game.LogicTime += gap_time
	now := time.Now()
	game.CurrentTime = float64(now.UnixNano()) / 1e9
	// fmt.Printf("Game is ticking, %v-%v, gap time is:%v\n", now.Second(), game.CurrentTime, gap_time)
	var temp_arr []BaseFunc
	for _, v := range game.BattleUnits {
		if v.Health() > 0 {
			v.Tick(gap_time)
			temp_arr = append(temp_arr, v)
		} else {
			// Remove v from array, or just leave it there
		}
	}

	for _, v := range game.AddUnits {
		temp_arr = append(temp_arr, v)
	}

	game.AddUnits = []BaseFunc{}
	game.BattleUnits = temp_arr

	// Handle skill targets callbacks
	var temp_arr2 []SkillTarget
	for _, v := range game.skill_targets {
		if v.trigger_time < game.LogicTime {
			v.callback(&v)
		} else {
			temp_arr2 = append(temp_arr2, v)
		}
	}

	if game.isSecuring != 0 {
		if game.TimeInitSecuring == -1 {
			game.TimeInitSecuring = game.LogicTime
			LogStr(fmt.Sprintf("Init Time: %v", game.TimeInitSecuring))
		} else if game.TimeInitSecuring+game.Time2Secure < game.LogicTime {
			game.Secured = true
			LogStr("Secured!")
		}
	} else {
		game.TimeInitSecuring = -1
	}
	temp_arr2 = append(temp_arr2, game.skill_targets_add...)
	game.skill_targets_add = []SkillTarget{}
	game.skill_targets = temp_arr2
	game.isSecuring = 0
	LogStr(fmt.Sprintf("Time++++: %v", game.TimeInitSecuring+game.Time2Secure))
	LogStr(fmt.Sprintf("Time: %v", game.LogicTime))

}

func (game *Game) HandleMultiPlayerAction(player_idx int, camp bool, action_code_0 int, action_code_1 int, action_code_2 int) {
	if action_code_0 == 9 {
		game.Init()
		return
	}

	var battle_unit BaseFunc
	if camp == true {
		battle_unit = game.SelfHeroes[player_idx].(BaseFunc)
	} else {
		battle_unit = game.OppoHeroes[player_idx].(BaseFunc)
	}
	cur_pos := battle_unit.Position()
	if battle_unit.Health() <= 0 {
		return
	}

	offset_x := float32(0)
	offset_y := float32(0)

	CalculateViewDepth(battle_unit)

	switch action_code_0 {
	case 0: // do nothing
		// Remain the same position
		battle_unit.(HeroFunc).SetTargetPos(cur_pos[0], cur_pos[1])
	case 1:
		// move
		dir := ConvertNum2Dir(action_code_1)
		if camp == true {
			LogStr(fmt.Sprintf("Input Direction: %v", dir))
		} else {
			LogStr(fmt.Sprintf("Enemy AI Direction: %v", dir))
		}
		offset_x = dir[0]
		offset_y = dir[1]
		target_pos_x := float32(cur_pos[0] + offset_x)
		target_pos_y := float32(cur_pos[1] + offset_y)
		battle_unit.(HeroFunc).SetTargetPos(target_pos_x, target_pos_y)
		dir.Normalize()
		battle_unit.SetDirection(dir)

	case 2:
		// normal attack
		// Remain the same position
		battle_unit.(HeroFunc).SetTargetPos(cur_pos[0], cur_pos[1])
	case 3:
		// skill 1
		dir := ConvertNum2Dir(action_code_2)
		offset_x = dir[0]
		offset_y = dir[1]
		battle_unit.(HeroFunc).UseSkill(0, offset_x, offset_y)
		// Set skill target
	case 4:
		// skill 2
		dir := ConvertNum2Dir(action_code_2)
		offset_x = dir[0]
		offset_y = dir[1]
		battle_unit.(HeroFunc).UseSkill(1, offset_x, offset_y)
	case 5:
		// skill 3
		dir := ConvertNum2Dir(action_code_2)
		offset_x = dir[0]
		offset_y = dir[1]
		battle_unit.(HeroFunc).UseSkill(2, offset_x, offset_y)
	case 6:
		// extra skill
		dir := ConvertNum2Dir(action_code_2)
		offset_x = dir[0]
		offset_y = dir[1]
		battle_unit.(HeroFunc).UseSkill(3, offset_x, offset_y)

	}
}

func (game *Game) HandleMultiPlayerVision(player_idx int, camp bool, action_code_0 int, action_code_1 int) {

	var battle_unit BaseFunc
	if camp == true {
		battle_unit = game.SelfHeroes[player_idx].(BaseFunc)
	} else {
		battle_unit = game.OppoHeroes[player_idx].(BaseFunc)
	}
	view_dir := battle_unit.Viewdir()
	if battle_unit.Health() <= 0 {
		return
	}
	switch action_code_0 {
	case 0: // do nothing
		// Remain the same position

	case 1:
		// rotate left
		//LogStr(fmt.Sprintf("Input Direction: %v", dir))
		length := math.Sqrt(float64(view_dir[0]*view_dir[0] + view_dir[1]*view_dir[1]))
		theta := math.Atan2(float64(view_dir[1]), float64(view_dir[0]))
		view_dir[0] = float32(math.Cos(theta-0.2) * length)
		view_dir[1] = float32(math.Sin(theta-0.2) * length)
		battle_unit.SetViewdir(view_dir.Normalized())
	case 2:
		// rotate right
		// Remain the same position
		length := math.Sqrt(float64(view_dir[0]*view_dir[0] + view_dir[1]*view_dir[1]))
		theta := math.Atan2(float64(view_dir[1]), float64(view_dir[0]))
		view_dir[0] = float32(math.Cos(theta+0.2) * length)
		view_dir[1] = float32(math.Sin(theta+0.2) * length)
		battle_unit.SetViewdir(view_dir.Normalized())
	}
	switch action_code_1 {
	case 0: // do nothing
		// Remain the same position

	case 1:
		// rotate up

		length := math.Sqrt(float64(view_dir[0]*view_dir[0] + view_dir[1]*view_dir[1]))
		theta := math.Atan2(float64(view_dir[2]), length)
		theta += 0.01
		if theta > math.Pi/2-1.5 {
			theta = math.Pi/2 - 1.5
		} else if theta < -math.Pi/2+1.5 {
			theta = -math.Pi/2 + 1.5
		}
		view_dir[2] = float32(math.Tan(theta) * length)
		battle_unit.SetViewdir(view_dir.Normalized())
	case 2:
		// rotate down
		// Remain the same position
		length := math.Sqrt(float64(view_dir[0]*view_dir[0] + view_dir[1]*view_dir[1]))
		theta := math.Atan2(float64(view_dir[2]), length)
		theta -= 0.01
		if theta > math.Pi/2-1.5 {
			theta = math.Pi/2 - 1.5
		} else if theta < -math.Pi/2+1.5 {
			theta = -math.Pi/2 + 1.5
		}
		view_dir[2] = float32(math.Tan(theta) * length)
		battle_unit.SetViewdir(view_dir.Normalized())
	}
}

// Global game object
var GameInst Game

func (game *Game) HandleEnemyAICommand() {
	hero := game.OppoHeroes[0]
	game_state := game.GetGameState(true)
	depth_map := hero.(BaseFunc).ViewDepth()

	action_type := hero.(*Lusian).model.SampleAction(game_state, depth_map)
	game.HandleMultiPlayerAction(0, false, action_type[0], action_type[1], 0) // bug
	game.HandleMultiPlayerVision(0, false, action_type[2], action_type[3])    // bug
}
