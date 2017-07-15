package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/antonholmquist/jason"
)

func main() {
	// receiver

	// Get Code

	// PullReq Diff Files

	// exec phpcs
	err := os.Mkdir("tmp", 0777)
	if err != nil {
		fmt.Println(err)
	}

	os.Chdir("./tmp/201707152116")
	wDir, _ := os.Getwd()
	fmt.Println("WokingDir: " + wDir)

	cmd := exec.Command("./vendor/bin/phpcs", "--report=json", "--standard=PSR2", "hello.php")
	fmt.Println("ExecCmd: " + strings.Join(cmd.Args, " "))
	out, err := cmd.Output()
	fmt.Println(string(out))

	fmt.Println("- Result phpcs ------------")
	v, err := jason.NewObjectFromBytes(out)
	files, err := v.GetObject("files")
	for key := range files.Map() {
		fileName := strings.Replace(key, wDir+"/", "", 1)
		fmt.Println("-- " + fileName + ": ")
		messages, _ := v.GetObjectArray("files", key, "messages")

		for _, message := range messages {
			line, _ := message.GetNumber("line")
			typ, _ := message.GetString("type")
			msg, _ := message.GetString("message")
			fmt.Println(string(line) + ": " + typ + ": " + msg)
		}
	}
	var json string = ""

	// JSON format
	var fmt_json string = ConvertJSONToLintEchoesFormat(json)
	fmt.Println(fmt_json)
}

func ConvertJSONToLintEchoesFormat(json string) string {
	// var cnv_json string

	return json
}
