package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	arg "github.com/alexflint/go-arg"
	"github.com/isoextension/btgo/taylog"
)

var log *taylog.Logger = taylog.New("startheme", nil)
var args *Arguments = &Arguments{}

func main() {
	cfg, err := os.UserConfigDir()
	dir := fmt.Sprintf("%s%s", cfg, "/starship")
	file := fmt.Sprintf("%s%s", dir, ".toml")
	if err != nil {
		// we should have panicked here but that is not user friendly so
		// panic(err)
		log.BasicErrorf("an error occured and your config directory could not be fetched %v\n", err)
		log.BasicFatalf("cannot continue without a config directory, program has to exit")
		os.Exit(1)
	}
	p := arg.MustParse(args)
	// if args == nil {
	// 	panic("args is nil!")
	// }
	switch {
	case args.Get != nil:
		current, err := current(file)
		if err == nil {
			fmt.Printf("current: %v\n", current)
			if args.Get.ShowCode {
				fmt.Print(current.Code)
			}
			return
		}
		log.BasicErrorf("%v", err)
		fmt.Printf("current: %v\n", current)
		if args.Get.ShowPath {
			fmt.Print(current.path)
		}
		if args.Get.ShowCode {
			fmt.Print(current.Code)
		}
	case args.List != nil:
		var themes []*Theme
		var files []fs.DirEntry
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			files = append(files, d)
			return nil
		})
		if err != nil {
			log.BasicErrorf("%v\n", err)
		}
		for _, v := range files {
			theme, err := current(v.Name())
			if err != nil {
				log.BasicErrorf("%v\n", err)
			}
			themes = append(themes, theme)
		}
	case args.Switch != nil:
		name := args.Switch.Name
		theme, err := lookup(name, dir)
		if err != nil {
			log.BasicErrorf("%v\n", err)
			return
		}
		err = switchTheme(theme, file)
		if err != nil {
			log.BasicErrorf("%v\n", err)
			return
		}
		log.BasicInfo("successfully changed theme!")
	default:
		fmt.Println("no command provided")
		p.WriteHelp(os.Stderr)
	}
	fmt.Println("startheme", version, "build", Build)
}
