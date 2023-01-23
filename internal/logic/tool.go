package logic

//nolint:goimports //i like it
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/exec"
	"gopkg.in/yaml.v3"
)

func setSaveMode(b bool) {
	k2 := "build.save"
	cmdr.Set(k2, b)
}

func isSaveMode() (b bool) {
	return cmdr.GetBoolR("build.save")
}

func isDryRunMode() (b bool) {
	return cmdr.GetBoolR("dry-run")
}

func setBuildScope(scope string) {
	k2 := "bgo.Scope"
	cmdr.Set(k2, scope)
}

func buildScopeFromCmdr(cmd *cmdr.Command) string {
	// ToggleGroup value of building scope is: appName + '.' + "Scope"
	// k2 := cmd.GetDottedNamePath() + ".Scope"
	k2 := "bgo.Scope"
	buildScope := cmdr.GetStringR(k2)
	// set the choice from command-line option to option store
	// so that we can retrieve it in extracting BgoSettings
	cmdr.Set("bgo.build.scope", buildScope)
	return buildScope
}

func findStringInFile(where, what string) (has bool, ver string) {
	file, err := os.Open(where)
	if err != nil {
		logx.Error("%v", err)
		return
	}

	defer func() {
		if err = file.Close(); err != nil {
			logx.Error("%v", err)
		}
	}()

	var b []byte
	if b, err = io.ReadAll(file); err != nil {
		logx.Error("%v", err)
		return
	}

	if strings.Contains(string(b), what) {
		has = true
		ver = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(b)), what))
	}
	return
}

func uniAdd(target []string, item string) []string {
	for _, s := range target {
		if item == s {
			return target
		}
	}
	return append(target, item)
}

//nolint:gocognit //needs split
func ifLdflags(bc *build.Context) {
	pairs := make(map[string]string)

	ifLdflagsCmdrSpecials(bc, pairs)
	ifLdflagsExtends(bc, pairs)
	ifLdflagsReduce(bc, pairs)

	for k, v := range pairs {
		bc.Ldflags = append(bc.Ldflags, k+v)
	}
}

func cleanupBs(bs *BgoSettings) {
	rootCommon := bs.Common
	bs.Scope = ""
	for _, g := range bs.Projects {
		for pn, p := range g.Items {
			logx.Trace("filter project %q", pn)
			cleanupCommon(p.Common, g.Common)
			cleanupCommon(p.Common, rootCommon)
		}
		cleanupCommon(g.Common, rootCommon)
	}
}

func cleanupCommon(c, ref *build.Common) {
	if ref == nil || c == nil {
		return
	}

	// cp:= *cmdr.StandardCopier
	// cp.ZeroIfEqualsFrom=true
	// cp.KeepIfFromIsNil=true
	// cp.KeepIfFromIsZero=true
	// cp.EachFieldAlways=true
	cp := *cmdr.GormDefaultCopier
	cp.IgnoreIfNotEqual = true

	// clear target field if equals to source
	_ = cp.Copy(c, ref)
}

func uniappend(a []string, s string) []string {
	for _, t := range a {
		if t == s {
			return a
		}
	}
	a = append(a, s)
	return a
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func tplExpand(tpl, name string, bc interface{}) (output string, err error) {
	var sb strings.Builder
	t := template.Must(template.New(name).Parse(tpl))
	if err = t.Execute(&sb, bc); err == nil {
		output = sb.String()
	}
	return
}

func yamlText(obj interface{}) string {
	var sb strings.Builder
	e := yaml.NewEncoder(&sb)
	e.SetIndent(2) //nolint:gomnd //i like it
	err := e.Encode(obj)
	if err != nil {
		return ""
	}
	err = e.Close()
	if err != nil {
		return ""
	}

	return sb.String()
}

func leftPad(s string, pad int) string {
	if pad <= 0 {
		return s
	}

	var sb strings.Builder
	padstr := strings.Repeat(" ", pad)
	scanner := bufio.NewScanner(bufio.NewReader(strings.NewReader(s)))
	for scanner.Scan() {
		sb.WriteString(padstr)
		sb.WriteString(scanner.Text())
		sb.WriteRune('\n')
	}
	return sb.String()
}
