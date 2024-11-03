package cclLoader

import (
	"github.com/ALiwoto/ccl/src/cclGenerators"
	"github.com/ALiwoto/ccl/src/cclGenerators/csGenerator"
	"github.com/ALiwoto/ccl/src/cclGenerators/gdGenerator"
	"github.com/ALiwoto/ccl/src/cclGenerators/goGenerator"
	"github.com/ALiwoto/ccl/src/cclGenerators/pyGenerator"
)

func LoadGenerators() {
	// c# generator
	for _, currentAlias := range csGenerator.LanguageAliases {
		cclGenerators.CodeGenerators[currentAlias] = csGenerator.GenerateCode
	}

	// gd generator
	for _, currentAlias := range gdGenerator.LanguageAliases {
		cclGenerators.CodeGenerators[currentAlias] = gdGenerator.GenerateCode
	}

	// go generator
	for _, currentAlias := range goGenerator.LanguageAliases {
		cclGenerators.CodeGenerators[currentAlias] = goGenerator.GenerateCode
	}

	// py generator
	for _, currentAlias := range pyGenerator.LanguageAliases {
		cclGenerators.CodeGenerators[currentAlias] = pyGenerator.GenerateCode
	}
}
