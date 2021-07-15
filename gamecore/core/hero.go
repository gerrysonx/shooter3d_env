package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ungerik/go3d/vec3"
)

const (
	SkillCount = 5
)

type HeroFunc interface {
	SetTargetPos(x float32, y float32)
	TargetPos() vec3.T
	SetSkillTargetPos(x float32, y float32)
	SkillTargetPos() vec3.T
	UseSkill(skill_idx uint8, a ...interface{})
	GetSkills() []Skill
	SetSkills([]Skill)
	CopyHero(HeroFunc)
}

type Hero struct {
	BaseInfo
	targetpos      vec3.T
	skilltargetpos vec3.T
	skills         [SkillCount]Skill
}

func (baseinfo *Hero) SetTargetPos(x float32, y float32) {
	/*
		player_camp := baseinfo.Camp()
		if player_camp == 0 {
			LogStr(fmt.Sprintf("SetTargetPos camp:%v move is called, target pos from (%v, %v) to (%v, %v)",
				player_camp, baseinfo.targetpos[0], baseinfo.targetpos[1], x, y))
		}
	*/
	baseinfo.targetpos[0] = x
	baseinfo.targetpos[1] = y
	baseinfo.targetpos[2] = baseinfo.position[2]
}

func (baseinfo *Hero) TargetPos() vec3.T {
	return baseinfo.targetpos
}

func (baseinfo *Hero) SetSkillTargetPos(x float32, y float32) {
	baseinfo.skilltargetpos[0] = x
	baseinfo.skilltargetpos[1] = y
	baseinfo.skilltargetpos[2] = baseinfo.position[2]
}

func (baseinfo *Hero) SkillTargetPos() vec3.T {
	return baseinfo.skilltargetpos
}

func (hero *Hero) GetSkills() []Skill {
	return hero.skills[0:]
}

func (hero *Hero) SetSkills(skills []Skill) {
	for i := 0; i < SkillCount; i += 1 {
		hero.skills[i] = skills[i]
	}
}

func (hero *Hero) Tick(gap_time float64) {
	now := time.Now()
	fmt.Printf("Hero is ticking, %v gap_time is:%v\n", now, gap_time)
}

func (hero *Hero) CopyHero(src HeroFunc) {
	hero.SetSkills(src.GetSkills())
}

var hero_template Hero

func (hero *Hero) ManualCtrl(gap_time float64) {
	pos := hero.Position()
	// Check milestone distance
	targetPos := hero.TargetPos()
	// Todo: remain z value
	targetPos[2] = pos[2]
	dist := vec3.Distance(&pos, &targetPos)
	pos_ready := false
	if dist < 5 {
		// Already the last milestone
		pos_ready = true
	}

	isEnemyNearby, enemy := CheckEnemyInFrustum(hero.Camp(), hero)
	if isEnemyNearby && pos_ready {
		// Check if time to make hurt
		NormalAttackEnemy(hero, enemy)
	} else {
		if pos_ready {
			// Do nothing
		} else {
			Chase(hero, targetPos, gap_time)
		}
	}
}

func (hero *Hero) InitFromJson(full_path string, id int32) bool {
	file_handle, err := os.Open(full_path)
	if err != nil {
		return false
	}

	defer file_handle.Close()

	buffer := make([]byte, 400)
	read_count, err := file_handle.Read(buffer)
	if err != nil {
		return false
	}
	buffer = buffer[:read_count]
	var jsoninfo JsonInfo
	if err = json.Unmarshal(buffer, &jsoninfo); err == nil {
		hero.attack_range = jsoninfo.AttackRange
		hero.attack_freq = jsoninfo.AttackFreq
		hero.health = jsoninfo.Health
		hero.max_health = jsoninfo.Health
		hero.damage = jsoninfo.Damage
		hero.speed = jsoninfo.Speed
		hero.view_range = jsoninfo.ViewRange
		hero.type_id = id
		hero.name = jsoninfo.Name
		hero.fov = jsoninfo.Fov
		hero.mag_size = jsoninfo.MagSize
		hero.mag_remain = jsoninfo.MagSize
		hero.bullet_num = jsoninfo.BulletNum
		hero.reload_time = jsoninfo.ReloadTime

		hero.extent[0] = jsoninfo.Extent[0]
		hero.extent[1] = jsoninfo.Extent[1]
		hero.extent[2] = jsoninfo.Extent[2]

		for i := 0; i < len(jsoninfo.Skills); i += 1 {
			if jsoninfo.Skills[i] == -1 {
				hero.skills[i].Name = ""
			} else {
				hero.skills[i] = SkillMgrInst.skills[jsoninfo.Skills[i]]
			}
		}

	} else {
		return false
	}

	return true
}

func (hero *Hero) DumpState() {
	LogStr(fmt.Sprintf("hero name is:%v, hero id:%v", hero.name, hero.GetId()))
	for i := 0; i < 4; i += 1 {
		skill_idx := i
		LogStr(fmt.Sprintf("skill_idx:%v, name is:%v, call back is:%v",
			skill_idx,
			hero.skills[skill_idx].Name,
			hero.skills[skill_idx].SkillFunc))
	}

}

func (hero *Hero) UseSkill(skill_idx uint8, a ...interface{}) {
	game := &GameInst
	// Check CD
	now_seconds := game.LogicTime
	old_skill_use_time := hero.LastSkillUseTime(skill_idx)
	hero.DumpState()

	if (old_skill_use_time + hero.skills[skill_idx].Life) > now_seconds {
		LogStr(fmt.Sprintf("%v CD time not come, skill_idx:%v, old_skill_use_time:%v, from:%v at time:%v",
			hero.skills[skill_idx].Name,
			skill_idx,
			old_skill_use_time,
			hero.GetId(), game.LogicTime))
		return
	} else {
		LogStr(fmt.Sprintf("%v CD time OK, skill_idx:%v, old_skill_use_time:%v, from:%v at time:%v",
			hero.skills[skill_idx].Name,
			skill_idx,
			old_skill_use_time,
			hero.GetId(), game.LogicTime))
	}

	hero.SetLastSkillUseTime(skill_idx, now_seconds)

	switch hero.skills[skill_idx].Type {
	case SkillTypeDir:
		SkillMgrInst.skills[hero.skills[skill_idx].Id].SkillFunc(hero, a...)
	case SkillTypeRadius:
		SkillMgrInst.skills[hero.skills[skill_idx].Id].SkillFunc(hero)
	case SkillTypeSpot:
	}
}

type HeroMgr struct {
	heroes map[int32]*Hero
}

var HeroMgrInst HeroMgr

func (heromgr *HeroMgr) LoadCfg(id int32, config_file_name string) {
	hero := new(Hero)
	hero.InitFromJson(config_file_name, id)
	heromgr.heroes[id] = hero
}

func GetType(obj interface{}) reflect.Type {
	return reflect.TypeOf(obj)
}

func (heromgr *HeroMgr) LoadCfgFolder(config_file_folder string) {
	// Load all skill configs under folder
	heromgr.heroes = make(map[int32]*Hero)
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
		heromgr.LoadCfg(id32, cfg_full_file_name)
	}
}

func (heromgr *HeroMgr) Spawn(a ...interface{}) BaseFunc {
	hero_id := a[0].(int32)
	wanted_camp := a[1].(int32)
	pos_x := a[2].(float32)
	pos_y := a[3].(float32)
	hero_template := heromgr.heroes[hero_id]
	var pos_z float32
	if len(a) > 4 {
		pos_z = a[4].(float32)
	} else {
		pos_z = 80.0
	}
	new_hero := GetHeroByName(hero_template.name, pos_x, pos_y, pos_z)

	new_hero.Copy(hero_template)
	hero_unit, ok := new_hero.(HeroFunc)
	if ok {
		hero_unit.CopyHero(hero_template)
	}

	new_hero.SetLastAttackTime(GameInst.LogicTime)
	if len(a) > 4 {
		InitHeroWithCamp(new_hero, wanted_camp, pos_x, pos_y, a[4].(float32))
	} else {
		InitHeroWithCamp(new_hero, wanted_camp, pos_x, pos_y)
	}

	return (new_hero)
}
