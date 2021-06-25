package core

import "github.com/ungerik/go3d/vec3"

func GetHeroByName(name string, a ...float32) BaseFunc {
	switch name {
	case "lusian":
		lusian := new(Lusian)
		lusian.enemyLastPosition = vec3.T{a[0], a[1], a[2]} //initialize the last position
		return lusian
	case "vi":
		return new(Vi)
	case "vayne":
		return new(Vayne)
	case "bullet":
		return new(Bullet)
	case "tower":
		return new(Tower)
	case "meleecreep":
		return new(MeleeCreep)
	case "rangecreep":
		return new(RangeCreep)
	case "siegecreep":
		return new(SiegeCreep)
	}

	return nil
}

func GetSkillFuncByName(name string) SkillFuncDef {
	switch name {
	case "GrabEnemyAtNose":
		return GrabEnemyAtNose
	case "PushEnemyAway":
		return PushEnemyAway
	case "DoStompHarm":
		return DoStompHarm
	case "BloodSucker":
		return BloodSucker
	case "DoDirHarm":
		return DoDirHarm
	case "SlowDirEnemy":
		return SlowDirEnemy
	case "JumpTowardsEnemy":
		return JumpTowardsEnemy
	}

	return nil
}
