package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	fileHeader = "// Code Generated by asm_gen.go DO NOT EDIT.\npackage instructions\nvar\n(\n%s\n"
	funcString = "func (%c %s) %s() string {\n%s\n}\n\n"
)

var (
	targetLang           string
	instructionMatcher   = regexp.MustCompile("([a-zA-Z]*)->\n(((\t|    ).*\n?)*)")
	specialStructMatcher = regexp.MustCompile("([a-zA-Z]*)=>\n(((\t|    ).*\n?)*)")
	commentMatcher       = regexp.MustCompile(`\/\/[^\n\r]+?(?:\*\)|[\n\r])`)
)

//Formats the body of the function for generating instruction "name"
func formatBody(name string) string {
	return fmt.Sprintf(
		`var sb strings.Builder
		err := %sTmpl.Execute(&sb, %c)
		if err != nil {
			fmt.Println("%s instruction has invalid templating -", err)
		}
		return sb.String()`, name, name[0], name)
}

//Formats a list into a list of strings, with the appropriate indexing
func formatSpecialBody(name string, fields []string) string {
	fieldsStr := make([]string, len(fields))

	for i, field := range fields {
		fieldsStr[i] = fmt.Sprintf(`"%s"`, strings.TrimSpace(field))
	}

	return fmt.Sprintf(
		`return []string{%s}[%c]`, strings.Join(fieldsStr, ","), name[0])
}

//Creates the template used to expand an instruction
func formatTmpl(name, tmpl string) string {
	return fmt.Sprintf("%sTmpl, _ = template.New(\"\").Delims(\"<\", \">\").Parse(`%s`)", name, tmpl)
}

func parseFile(targetLang, suffix string) (map[string]string, []string, error) {
	filename := targetLang + suffix
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	//Preprocess string
	withoutComments := commentMatcher.ReplaceAllString(string(data), "")
	operandsExpanded := operand.apply(withoutComments)

	labelToBody := make(map[string]string)

	cmds := instructionMatcher.FindAllStringSubmatch(operandsExpanded, -1)
	tmpls := make([]string, len(cmds))

	for i, cmd := range cmds {
		label := cmd[1]
		body := strings.TrimRight(cmd[2], " \n")
		body = oneLine.apply(body)
		tmpls[i] = formatTmpl(label, body)
		labelToBody[label] = formatBody(label)
	}

	spCmds := specialStructMatcher.FindAllStringSubmatch(withoutComments, -1)

	for _, spCmd := range spCmds {
		label := spCmd[1]
		fields := strings.Split(strings.TrimSpace(spCmd[2]), "\n")
		labelToBody[label] = formatSpecialBody(label, fields)
	}
	return labelToBody, tmpls, nil
}

func main() {
	targetLang = os.Args[1]
	funcs, tmpls, err := parseFile(targetLang, "_instructions.tmpl")
	if err != nil {
		panic(err)
	}

	opFuncs, opTmpls, err := parseFile(targetLang, "_operands.tmpl")

	dest := fmt.Sprintf(fileHeader, strings.Join(tmpls, "\n"))
	dest += fmt.Sprintf("\n%s\n)\n", strings.Join(opTmpls, "\n"))

	for k, v := range funcs {
		dest += fmt.Sprintf(funcString, k[0], k, strings.Title(targetLang), v)
	}

	for k, v := range opFuncs {
		dest += fmt.Sprintf(funcString, k[0], k, strings.Title(targetLang)+"Operand", v)
	}

	goFile := "../" + targetLang + ".go"
	err = ioutil.WriteFile(goFile, []byte(dest), 0644)
	if err != nil {
		panic(err)
	}

	//Apply formatting to generated file
	err = exec.Command("../../../../wacc-goimports", "-w", goFile).Run()
	if err != nil {
		panic(err)
	}
	err = exec.Command("gofmt", "-w", "-s", "-r", `fmt.Sprintf("") -> ""`, goFile).Run()
	if err != nil {
		panic(err)
	}
}
