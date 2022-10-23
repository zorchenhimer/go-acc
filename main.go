package main

import (
	"io"
	"os"
	"fmt"
	//"path/filepath"
	//"hash"
	"hash/crc32"
	"time"
	"regexp"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/theckman/yacspin"
)

var re_hash = regexp.MustCompile(`([a-fA-F0-9]{8})`)
var re_hash_strict = regexp.MustCompile(`(\{|\[|\()([a-fA-F0-9]{8})(\}|\]|\))`)

type Arguments struct {
	InputFiles []string `arg:"positional,required" help:"Input file; accepts glob."`
}

func main() {
	args := &Arguments{}
	arg.MustParse(args)
	err := run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args *Arguments) error {
	if len(args.InputFiles) == 0 {
		return fmt.Errorf("no input file(s) given")
	}

	for _, f := range args.InputFiles {
		cfg := yacspin.Config{
			Frequency: time.Millisecond * 500,
			CharSet: yacspin.CharSets[51],
			Suffix: "   ",
			Message: f,
			StopCharacter: "",
			StopColors: []string{"fgGreen"},
		}
		spinner, err := yacspin.New(cfg)
		if err != nil {
			return fmt.Errorf("unable to start spinner: %w", err)
		}
		spinner.Start()

		crc, err := crcFilename(f)
		if err != nil {
			return fmt.Errorf("error calculating crc for %q: %w", f, err)
		}

		spinner.Stop()

		crc_filename := re_hash_strict.FindString(f)
		if crc_filename != "" {
			crc_filename = crc_filename[1:9]
		} else {
			crc_filename = re_hash.FindString(f)
		}
		color := "33"
		if crc_filename != "" {
			if crc_filename == crc {
				color = "32"
			} else {
				color = "31"
			}
		}
		f = strings.Replace(f, crc_filename, "\x1b["+color+"m"+crc_filename+"\x1b[0m", 1)
		fmt.Printf("\x1b[%sm%s\x1b[0m %s\n", color, crc, f)

	}

	return nil
}

// Generates a CRC32 hash of the given filename
func crcFilename(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	hsh := crc32.New(crc32.IEEETable)
	_, err = io.Copy(hsh, file)
	if err != nil {
		return "", fmt.Errorf("error calculating hash: %w", err)
	}

	return fmt.Sprintf("%08X", hsh.Sum32()), nil
}

