package logic

//nolint:goimports
import (
	"bufio"
	"fmt"
	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/exec"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
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
	// ToggleGroup value of building scope is: "Scope"
	k2 := cmd.GetDottedNamePath() + ".Scope" // nolint
	k2 = "bgo.Scope"
	buildScope := cmdr.GetStringR(k2)
	// set the choice from command-line option to option store
	// so that we can retrieve it in extracting BgoSettings
	cmdr.Set("bgo.build.scope", buildScope)
	return buildScope
}

func findStringInFile(where, what string) (has bool) {
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
	if b, err = ioutil.ReadAll(file); err != nil {
		logx.Error("%v", err)
		return
	}

	if strings.Contains(string(b), what) {
		has = true
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

func ifLdflags(bc *build.Context) {
	pairs := make(map[string]string)

	if bc.HasGoMod {
		where, what := bc.GoModFile, "github.com/hedzr/cmdr"
		bc.CmdrSpecials = findStringInFile(where, what)
	} else {
		bc.CmdrSpecials = true
	}

	if bc.CmdrSpecials {
		const W = "github.com/hedzr/cmdr/conf"
		var str string
		str = fmt.Sprintf("-X %s.AppName=", W)
		pairs[str] = bc.AppName
		str = fmt.Sprintf("-X %s.Version=", W)
		ver := bc.CalcVersion()
		pairs[str] = strings.TrimPrefix(ver, "v")
		str = fmt.Sprintf("-X %s.Buildstamp=", W)
		pairs[str] = bc.BuildTime
		str = fmt.Sprintf("-X %s.Githash=", W)
		pairs[str] = bc.GitRevision
		str = fmt.Sprintf("-X %s.GoVersion=", W)
		pairs[str] = strings.ReplaceAll(bc.GoVersion, " ", "_")
		// fmt.Sprintf("-X '%s.AppName=%s'", W,bc.AppName),

		str = fmt.Sprintf("-X %s.Serial=", W)
		pairs[str] = strconv.FormatInt(bc.Serial, 10) //nolint:gomnd
		str = fmt.Sprintf("-X %s.SerialString=", W)
		pairs[str] = bc.RandomString
	}

	for _, pnv := range bc.Common.Extends {
		if pnv.Package == "" {
			continue
		}
		for n, v := range pnv.Values {
			if n == "" || v == "" {
				continue
			}
			if v[0] == '`' && v[len(v)-1] == '`' {
				// shell scripts
				script := v[1 : len(v)-1]
				if re, err := tplExpand(script, "set-name-and-value-in-package", bc); err == nil {
					script = re
				}
				if err := exec.New().
					WithCommand("bash", "-c", script).
					WithOnOK(func(retCode int, stdoutText string) {
						v = strings.ReplaceAll(strings.TrimSuffix(stdoutText, "\n"), " ", "_")
					}).RunAndCheckError(); err != nil {
					continue
				}
			} else {
				if re, err := tplExpand(v, "set-name-and-value-in-package", bc); err == nil {
					v = re
				}
			}
			str := fmt.Sprintf("-X %s.%s=", pnv.Package, n)
			pairs[str] = v
		}
	}

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
	e.SetIndent(2) //nolint:gomnd
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
