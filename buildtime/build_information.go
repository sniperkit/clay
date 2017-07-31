package buildtime

import "github.com/qb0C80aE/clay/extensions"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime: "20170731203045",
		claySubModuleInformationList: []*claySubModuleInformation{
			{
				name:     "clay",
				revision: "5b9738ee004cdec0594306a7aea14fdb57ec8eb2",
			},
			{
				name:     "github.com/qb0C80aE/loam",
				revision: "408ece1087c973979487cf47d8b8deb2671e14a1",
			},
			{
				name:     "github.com/qb0C80aE/pottery",
				revision: "7972ce0e6b854cc9faab5ab17d30fc6dfbf9b45d",
			},
		},
	}
	extensions.RegisterProgramInformation(programInformation)
}
