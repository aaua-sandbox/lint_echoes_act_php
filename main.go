package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/antonholmquist/jason"
)

type LintEchoesFormat struct {
	File     string                    `json:"file"`
	Messages []LintEchoesFormatMessage `json:"messages"`
}

type LintEchoesFormatMessage struct {
	Type    string `json:"type"`
	Line    int64  `json:"line"`
	Message string `json:"message"`
}

func main() {
	// TODO: receiver

	// TODO: send status to slack

	// make tmp dir
	if err := os.Mkdir("./tmp", 0777); err != nil {
		fmt.Println(err)
	}
	if err := os.RemoveAll("./tmp/src"); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir("./tmp/src", 0777); err != nil {
		fmt.Println(err)
	}

	// Get Code
	// TODO: URL and private repo
	gitCloneURL := "https://github.com/aaua-sandbox/php_sandbox.git"
	tmpDir := "./tmp/src/" + strings.Replace(time.Now().Format("20060102150405.000"), ".", "", 1)
	cmdGetCode := exec.Command("git", "clone", gitCloneURL, tmpDir)
	fmt.Println("ExecCmd: " + strings.Join(cmdGetCode.Args, " "))
	_, err := cmdGetCode.Output()
	if err != nil {
		fmt.Println(err)
	}

	os.Chdir(tmpDir)
	wDir, _ := os.Getwd()
	fmt.Println("WokingDir: " + wDir)

	// TODO: PullReq Diff Files

	// exec phpcs
	var fileNames []string
	fileNames = append(fileNames, "hello.php")
	cmd := exec.Command("phpcs", "--report=json", "--standard=PSR2", strings.Join(fileNames, " ")) // TODO: file name
	fmt.Println("ExecCmd: " + strings.Join(cmd.Args, " "))
	json, err := cmd.Output()
	fmt.Println(string(json))

	// format
	fmtJSON, err := ConvertJSONToLintEchoesFormat(json, wDir)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmtJSON)

	// TODO: filter PullReq Change Lines

	// TODO: filter Lint Level

	// TODO: Comment to PullReq

	// TODO: ResultComment to PullReq

	// TODO: send status to slack
}

func ConvertJSONToLintEchoesFormat(json []byte, wDir string) (LintEchoesFormat, error) {
	fmt.Println("- ConvertJSONToLintEchoesFormat ------------")
	result := LintEchoesFormat{}

	v, err := jason.NewObjectFromBytes(json)
	if err != nil {
		return result, err
	}

	files, err := v.GetObject("files")
	if err != nil {
		return result, err
	}

	for key := range files.Map() {
		fileName := strings.Replace(key, wDir+"/", "", 1)
		jMessages, _ := v.GetObjectArray("files", key, "messages")

		lintEchoesFormatMessages := []LintEchoesFormatMessage{}
		for _, jMessage := range jMessages {
			jType, _ := jMessage.GetString("type")
			jLine, _ := jMessage.GetInt64("line")
			jMessageText, _ := jMessage.GetString("message")
			lintEchoesFormatMessages = append(lintEchoesFormatMessages, LintEchoesFormatMessage{
				Type:    jType,
				Line:    jLine,
				Message: jMessageText,
			})
		}

		result.File = fileName
		result.Messages = lintEchoesFormatMessages
	}

	return result, nil
}
