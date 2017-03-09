package spinnaker

type orcaRefResponse struct {
	Ref string `json:"ref"`
}

type StartPipelineRequest struct {
	Application   string
	Name          string
	Template      DCDPipelineTemplate
	Configuration DCDPipelineConfig
}

type StartPipeline struct {
	Type   string      `json:"type"`
	Config interface{} `json:"config"`
}

type Pipeline struct {
	Pipeline PipelineConfig `json:"pipeline"`
}

type PipelineConfig struct {
	Application   string              `json:"application"`
	Name          string              `json:"name"`
	Template      DCDPipelineTemplate `json:"template"`
	Configuration DCDPipelineConfig   `json:"configuration,omitempty"`
}

type DCDPipelineTemplate struct {
	Source string `json:"source"`
}

type DCDPipelineConfig struct{}
