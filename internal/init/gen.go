/*
Copyright 2023 Enzo Nocera <enzo@nocera.eu>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package init

import (
	ttemplate "github.com/Raskyld/go-tektasker/internal/init/internal/template"
	ttmarkers "github.com/Raskyld/go-tektasker/pkg/markers"
	"log/slog"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"
	"text/template"
)

type Generator struct {
	Logger   *slog.Logger
	Template *template.Template
}

func New(logger *slog.Logger) *Generator {
	gen := &Generator{
		Logger:   logger,
		Template: &template.Template{},
	}

	gen.RegisterTemplate(ttemplate.TaskfileName, ttemplate.TaskfileTpl)

	return gen
}

func (*Generator) RegisterMarkers(into *markers.Registry) error {
	return ttmarkers.Register(into)
}

func (g *Generator) RegisterTemplate(name string, tpl string) *Generator {
	nt, err := g.Template.New(name).Parse(tpl)
	if err != nil {
		panic(err)
	}

	g.Template = nt
	return g
}

func (g *Generator) Generate(output genall.OutputRule) error {
	wr, err := output.Open(nil, "Taskfile.yaml")
	if err != nil {
		return err
	}

	err = g.Template.ExecuteTemplate(wr, ttemplate.TaskfileName, nil)
	if err != nil {
		return err
	}

	return nil
}
