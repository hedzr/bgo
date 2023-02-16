package logic

import (
	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/evendeep"
	"github.com/hedzr/evendeep/flags/cms"
)

func (s *BgoSettings) PullDownCommonSettings() {
	c := evendeep.New(
		evendeep.WithCopyStrategyOpt,
		evendeep.WithStrategies(cms.OmitIfEmpty),
	) // , evendeep.WithWipeTargetSliceFirstOpt)
	// c := cmdr.StandardCopier // make a copy
	// c.KeepIfFromIsNil = true
	// c.KeepIfFromIsZero = true
	// c.EachFieldAlways = true
	// c.IgnoreIfNotEqual = true

	// log.Debugf("  .  PROJECTS.Common = %+v", s.Common)

	// Copying the common settings from the up level if necessary.
	// It means:
	//   the top level (app.bgo.build.xxx) -> project.common
	//   the project-group level -> project.common
	for gn, grp := range s.Projects {
		// log.Debugf("  .  PROJECTS[%q].Common = %+v\n", gn, s.Projects[gn].Common)
		for pn, prj := range grp.Items {
			var common = build.NewCommon()
			_ = c.CopyTo(s.Common, &common)
			_ = c.CopyTo(grp.Common, &common)
			_ = c.CopyTo(prj.Common, &common)
			s.Projects[gn].Items[pn].Common = common
			if s.Projects[gn].Items[pn].Disabled {
				continue
			}
			// log.Debugf("  .  PROJECTS[%q][%q]: %+v | Install = %+v\n",
			// gn, pn, s.Projects[gn].Items[pn], s.Projects[gn].Items[pn].Common.Install)
		}
	}
}
