package render

import (
	"encoding/csv"
	"os"

	"github.com/oopsunix/wii/internal/model"
)

type csvRenderer struct {
	cfg *model.Config
}

func (r *csvRenderer) Render(_ []model.DevEnv, sections []model.Section) {
	w := csv.NewWriter(os.Stdout)
	_ = w.Write([]string{"name", "version", "path", "source"})
	for _, sec := range sections {
		for _, t := range sec.Tools {
			_ = w.Write([]string{
				t.Name,
				t.Version,
				displayPath(t.Path, r.cfg.FullPath),
				sec.Label,
			})
		}
	}
	w.Flush()
}
