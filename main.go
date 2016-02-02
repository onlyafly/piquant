package main

import (
	"flag"
	"fmt"

	"log"
	"os"
	"strings"

	"piquant/util"

	"piquant/lang/ast"
	"piquant/lang/interpreter"

	"github.com/peterh/liner"
)

const (
	version         = `0.1.0-alpha`
	versionDate     = `2015-02-20`
	historyFilename = "/tmp/.piquant_liner_history"
)

var (
	// TODO add functionality for these missing commands
	commandCompletions = []string{":quit" /*":load ", ":reset", ":help",*/, ":inspect "}
	// TODO wordCompletions    = []string{"def", "update!"}
)

func configureLiner(linerState *liner.State) {
	linerState.SetCtrlCAborts(true)

	linerState.SetCompleter(func(line string) (c []string) {
		for _, n := range commandCompletions {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	/* TODO
	// WordCompleter takes the currently edited line with the cursor position and
	// returns the completion candidates for the partial word to be completed. If
	// the line is "Hello, wo!!!" and the cursor is before the first '!',
	// ("Hello, wo!!!", 9) is passed to the completer which may returns
	// ("Hello, ", {"world", "Word"}, "!!!") to have "Hello, world!!!".
	linerState.SetWordCompleter(func(line string, pos int) (head string, completions []string, tail string) {
		for _, n := range wordCompletions {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})
	*/
}

func openLinerHistory(line *liner.State) {
	if f, err := os.Open(historyFilename); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

func writeLinerHistory(line *liner.State) {
	if f, err := os.Create(historyFilename); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func main() {

	startupFileName := flag.String("l", "", "load a file at startup")
	showHelp := flag.Bool("help", false, "show the help")
	flag.Parse()
	exeFileName := flag.Arg(0)

	if showHelp != nil && *showHelp {
		fmt.Printf("Usage of piquant:\n")
		flag.PrintDefaults()
		return
	}

	// Setup liner

	line := liner.NewLiner()
	defer line.Close()
	openLinerHistory(line)
	configureLiner(line)

	// Initialize

	topLevelEnv := interpreter.NewTopLevelMapEnv()

	if len(exeFileName) != 0 {
		loadFile("prelude.v", topLevelEnv)
		loadFile(exeFileName, topLevelEnv)
		return
	}

	fmt.Printf("Piquant %s (%s)\n", version, versionDate)
	loadFile("prelude.v", topLevelEnv)
	fmt.Printf("(Press Ctrl+C or type :quit to exit)\n\n")

	// Loading of files

	if startupFileName != nil {
		loadFile(*startupFileName, topLevelEnv)
	}

	// REPL

	for {
		input, err := line.Prompt("> ")

		if err != nil {
			if err.Error() == "prompt aborted" {
				fmt.Printf("Quiting...\n")
			} else {
				fmt.Printf("Prompt error: %s\n", err)
			}
			return
		}

		line.AppendHistory(input)
		writeLinerHistory(line)

		switch {
		case input == ":quit":
			return
		case strings.HasPrefix(input, ":inspect "):
			withoutInspectPrefix := strings.Split(input, ":inspect ")[1]
			if result, err := interpreter.ParseEval(topLevelEnv, withoutInspectPrefix, "REPL"); err == nil {
				inspect(result)
			} else {
				fmt.Println(err.Error())
			}
		default:
			interpreter.ParseEvalPrint(topLevelEnv, input, "REPL", true)
		}
	}
}

func inspect(arg ast.Node) {
	switch val := arg.(type) {
	case *interpreter.EnvNode:
		fmt.Printf(
			"Environment\n  Name='%v'\n  Env=%v\n",
			val.Name(),
			val.Env.String())
	default:
		fmt.Printf("Don't know how to inspect: %v\n", val.String())
	}
}

func loadFile(fileName string, env interpreter.Env) {
	if len(fileName) > 0 {
		content, err := util.ReadFile(fileName)
		if err != nil {
			fmt.Printf("Error while loading file <%v>: %v\n", fileName, err.Error())
		} else {
			interpreter.ParseEvalPrint(env, content, fileName, false)
		}
	}
}
