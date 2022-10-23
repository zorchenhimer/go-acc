package main

import (
	"io"
	"os"
	"os/signal"
	"fmt"
	"path/filepath"
	//"hash"
	"hash/crc32"
	"time"
	"regexp"
	"strings"
	//"net/url"
	"html"

	"golang.org/x/term"
	"github.com/alexflint/go-arg"
	"github.com/theckman/yacspin"
	"github.com/zorchenhimer/go-ed2k"
)

var re_hash = regexp.MustCompile(`([a-fA-F0-9]{8})`)
var re_hash_strict = regexp.MustCompile(`(\{|\[|\()([a-fA-F0-9]{8})(\}|\]|\))`)

type Arguments struct {
	InputFiles []string `arg:"positional,required" help:"Input file; accepts glob."`
	AddHash bool `arg:"-a,--add" help:"Add the calculated hash to the filename if none is found."`
	AddDelim string `arg:"-d,--add-delim" help:"Character to use before the added hash."`
	Ed2k bool `arg:"-e,--ed2k" help:"Print ED2K links."`
}

func handleInterrupt(spinner *yacspin.Spinner) {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return
	}

	ch := make(chan os.Signal, 5)
	signal.Notify(ch, os.Interrupt)
	<-ch

	if spinner != nil {
		spinner.Stop()
	}
	fmt.Print("\x1b[?25h\n")
	os.Exit(1)
}

func main() {
	args := &Arguments{}
	arg.MustParse(args)

	if args.AddDelim != "" {
		args.AddHash = true
	}

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

	isTerm := term.IsTerminal(int(os.Stdout.Fd()))
	var err error
	var spinner *yacspin.Spinner
	go handleInterrupt(spinner)

	for _, f := range args.InputFiles {
		if isTerm {
			cfg := yacspin.Config{
				Frequency: time.Millisecond * 500,
				CharSet: yacspin.CharSets[51],
				Suffix: "   ",
				Message: f,
				StopCharacter: "",
				StopColors: []string{"fgGreen"},
			}
			spinner, err = yacspin.New(cfg)
			if err != nil {
				return fmt.Errorf("unable to start spinner: %w", err)
			}
			spinner.Start()
		}

		if args.Ed2k {
			crc, err := ed2kFilename(f)
			if err != nil {
				return fmt.Errorf("error calculating ed2k for %q: %w", f, err)
			}

			if isTerm {
				spinner.Stop()
			}

			fmt.Println(crc)
			continue
		}

		crc, err := crcFilename(f)
		if err != nil {
			return fmt.Errorf("error calculating crc for %q: %w", f, err)
		}

		crc_filename := re_hash_strict.FindString(f)
		if crc_filename != "" {
			crc_filename = crc_filename[1:9]
		} else {
			crc_filename = re_hash.FindString(f)
		}

		color := "33"
		if crc_filename != "" {
			if crc_filename == crc {
				// found and is correct
				color = "32"
			} else {
				// found and is incorrect
				color = "31"
			}
		} else if args.AddHash {
			// not found

			ext := filepath.Ext(f)
			newName := f[:len(f)-len(ext)]+args.AddDelim+"["+crc+"]"+ext
			err = os.Rename(f, newName)
			if err != nil {
				return fmt.Errorf("unable to rename file: %w", err)
			}
		}

		if isTerm {
			spinner.Stop()

			f = strings.Replace(f, crc_filename, "\x1b["+color+"m"+crc_filename+"\x1b[0m", 1)
			fmt.Printf("\x1b[%sm%s\x1b[0m %s\n", color, crc, f)

		} else {
			fmt.Printf("%s %s\n", crc, f)
		}
	}

	return nil
}

// Generates a CRC32 hash of the given filename
func crcFilename(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hsh := crc32.New(crc32.IEEETable)
	_, err = io.Copy(hsh, file)
	if err != nil {
		return "", fmt.Errorf("error calculating hash: %w", err)
	}

	return fmt.Sprintf("%08X", hsh.Sum32()), nil
}

// Returns a full ed2k link.  These links are dumb AF.
func ed2kFilename(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hsh := ed2k.New()
	_, err = io.Copy(hsh, file)
	if err != nil {
		return "", fmt.Errorf("error calculating hash: %w", err)
	}

	hash, err := hsh.SumBlue()
	if err != nil {
		return "", err
	}

	st, err := os.Stat(filename)
	if err != nil {
		return "", err
	}

	// ed2k://|file|FILENAME|SIZE|HASH|/
	return fmt.Sprintf("ed2k://|file|%s|%d|%s|/", html.EscapeString(strings.ReplaceAll(filename, " ", "_")), st.Size(), hash), nil
}
