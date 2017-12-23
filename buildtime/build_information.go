package buildtime

import "github.com/qb0C80aE/clay/extensions"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime: "unknown",
		claySubModuleInformationList: []*claySubModuleInformation{
			{
				name:     "clay",
				revision: "d9d014dc4bbc169869bf05e6ac0b74ee8d63b3a7",
			},
		},
	}
	extensions.RegisterProgramInformation(programInformation)
}
