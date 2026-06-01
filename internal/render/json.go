package render

import (
	"encoding/json"
	"os"

	"wii/internal/model"
)

type jsonRenderer struct {
	cfg *model.Config
}

func (r *jsonRenderer) Render(devEnvs []model.DevEnv, sections []model.Section) {
	type jsonDevEnv struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	type jsonTool struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Path    string `json:"path"`
		Source  string `json:"source"`
	}
	type jsonOutput struct {
		DevEnvs []jsonDevEnv `json:"devEnvs"`
		Tools   []jsonTool   `json:"tools"`
	}

	out := jsonOutput{
		DevEnvs: make([]jsonDevEnv, 0),
		Tools:   make([]jsonTool, 0),
	}

	for _, env := range devEnvs {
		if env.Installed {
			out.DevEnvs = append(out.DevEnvs, jsonDevEnv{
				Name:    env.Name,
				Version: env.Version,
			})
		}
	}

	for _, sec := range sections {
		for _, t := range sec.Tools {
			out.Tools = append(out.Tools, jsonTool{
				Name:    t.Name,
				Version: t.Version,
				Path:    displayPath(t.Path, r.cfg.FullPath),
				Source:  sec.Label,
			})
		}
	}

	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(out)
}
