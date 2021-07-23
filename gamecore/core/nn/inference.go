package nn

import (
	"fmt"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

type NeuralNet struct {
	Type             string
	inf_pb_path      string
	saved_model      *tf.SavedModel
	input_s_op         tf.Output
	input_d_op			tf.Output
	input_c_op			tf.Output
	output_policy_0_op tf.Output
	output_policy_1_op tf.Output
	output_policy_3_op tf.Output
	output_policy_4_op tf.Output
	output_value_op  tf.Output
}

func (nn *NeuralNet) Test() {
	s := op.NewScope()
	c := op.Const(s, "Hello0 from tensorflow version "+tf.Version())
	graph, err := s.Finalize()
	if err != nil {
		panic(err)
	}

	// Execute the graph in a session
	sess, err := tf.NewSession(graph, nil)
	if err != nil {
		panic(err)
	}

	output, err := sess.Run(nil, []tf.Output{c}, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(output[0].Value())
}

func (nn *NeuralNet) Load(pb_path string) {
	
	var err error
	var option tf.SessionOptions
	option.Config = []byte{0x32, 0x2, 0x20, 0x1}
	nn.saved_model, err = tf.LoadSavedModel(pb_path, []string{"serve"}, &option)
	
	if err != nil {
		panic(err)
	}

	nn.input_s_op = nn.saved_model.Graph.Operation("input/multi_s").Output(0)
	nn.input_d_op = nn.saved_model.Graph.Operation("input/multi_d").Output(0)
	nn.input_c_op = nn.saved_model.Graph.Operation("input/multi_c").Output(0)
	nn.output_policy_0_op = nn.saved_model.Graph.Operation("policy_net_new/policy_net_new/policy_head_0_head").Output(0)
	nn.output_policy_1_op = nn.saved_model.Graph.Operation("policy_net_new/policy_net_new/policy_head_1_head").Output(0)
	nn.output_policy_3_op = nn.saved_model.Graph.Operation("policy_net_new/policy_net_new/policy_head_3_head").Output(0)
	nn.output_policy_4_op = nn.saved_model.Graph.Operation("policy_net_new/policy_net_new/policy_head_4_head").Output(0)
	nn.output_value_op = nn.saved_model.Graph.Operation("policy_net_new/policy_net_new/value_output").Output(0)
}

func (nn *NeuralNet) Release() {
	nn.saved_model.Session.Close()
}

func (nn *NeuralNet) Ref(input_data_s [][][][]float32, input_data_d [][][][]float32, input_data_c [][][][]int32) []*tf.Tensor {
	input_tensor_s, err := tf.NewTensor(input_data_s)
	if err != nil {
		panic(err)
	}

	input_tensor_d, err := tf.NewTensor(input_data_d)
	if err != nil {
		panic(err)
	}

	input_tensor_c, err := tf.NewTensor(input_data_c)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("shape: %v, %v", input_tensor_s.Shape(), input_tensor_d.Shape())

	result, run_err := nn.saved_model.Session.Run(
		map[tf.Output]*tf.Tensor{
			nn.input_s_op: input_tensor_s,
			nn.input_d_op: input_tensor_d,
			nn.input_c_op: input_tensor_c,
		},
		[]tf.Output{nn.output_policy_0_op, nn.output_policy_1_op, nn.output_policy_3_op, nn.output_policy_4_op},
		nil)

	if run_err != nil {
		fmt.Printf("Error: %v", run_err)
		panic(run_err)
	}


	return result
}
