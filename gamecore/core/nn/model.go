package nn

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

var LogHandle *os.File

func LogStr(log string) {
	if LogHandle == nil {
		return
	}
	LogHandle.WriteString(log)

	LogHandle.WriteString("\n")
	LogHandle.Sync()
}

type Model struct {
	inference_pb    string
	bouquet_state   [][]float32
	formatted_state [][]float32
	model           NeuralNet
	r               *rand.Rand
}

type ModelFunc interface {
	InitModel(index string)
	SampleAction(game_state [][]float32, depth_map []float32) []int
	GetIndex() string
}

type ModelJsonInfo struct {
	Model []string
	Score []int
}

func (hero *Model) InitModel(index string) {
	root_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	full_path := fmt.Sprintf("%s/../model/model_%s", root_dir, index)
	hero.inference_pb = index
	if err != nil {
		log.Fatal(err)
	}

	hero.model = NeuralNet{}
	if hero.model.saved_model == nil {
		hero.model.Load(full_path)
	}
	hero.r = rand.New(rand.NewSource(0))
}

func (model *Model) SampleAction(game_state [][]float32, depth_map []float32) []int {

	// Inference from nn
	var temp_depth_map [][]float32
	for _, item := range depth_map {
		temp_depth_map = append(temp_depth_map, []float32{item})
	}
	formatted_game_state := [][][][]float32{[][][]float32{game_state}}
	formatted_depth_map := [][][][]float32{[][][]float32{temp_depth_map}}

	predict := model.model.Ref(formatted_game_state, formatted_depth_map)

	var action []int
	for _, p := range predict {
		prob := model.r.Float32()
		for i, o := range p.Value().([][]float32)[0] {
			prob -= o
			if prob < 0 {
				action = append(action, i)
				break
			}
		}
	}

	return action

}

func (model *Model) GetIndex() string {
	return model.inference_pb
}

//Deprecated in the newest version. Archived
/*
func (hero *Model) GetInput(input_state [][]float32) [][][]float32 {
	// from bouquet states to formatted states
	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			hero.formatted_state[i][j] = hero.bouquet_state[j][i]
		}
	}

	formatted_input := [][][]float32{input_state}
	return formatted_input
}

func (hero *Model) PrepareInput() {
	hero.bouquet_state = make([][]float32, 4)
	for i := 0; i < 4; i++ {
		hero.bouquet_state[i] = make([]float32, 6)
	}

	hero.formatted_state = make([][]float32, 6)
	for i := 0; i < 6; i++ {
		hero.formatted_state[i] = make([]float32, 4)
	}

	root_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	full_path := fmt.Sprintf("%s/../model/model_%s", root_dir, "0001")
	if err != nil {
		log.Fatal(err)
	}

	hero.model = NeuralNet{}
	hero.model.Load(full_path)
}

func (hero *Model) UpdateInput(one_state []float32) {
	hero.bouquet_state = append(hero.bouquet_state, one_state)
	hero.bouquet_state = hero.bouquet_state[1:5]
}*/
