package main

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:    "[dir] [hold generation] [target extension...]",
		Args:   cobra.MinimumNArgs(3),
		Short:  "Rotating backup file with a given hold-generation number",
		RunE: func(cmd *cobra.Command, args []string) error {
			d := args[0]
			if _, err := os.Stat(d); err != nil {
				if os.IsNotExist(err) {
					return errors.New("The specified path does not exist")
				}
				return err
			}

			gen, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || gen == 0 {
				return errors.New("'generation' must be an integer and greater than zero")
			}

			ext := args[2:]

			files, err := os.ReadDir(d)
			if err != nil {
				return errors.New("Failed to get directory entries from the specified path")
			} else if len(files) < int(gen) {
				fmt.Println("The specified directory has fewer files than the specified generation limit")
				return nil
			}

			files = slices.DeleteFunc(
				files,
				func(e os.DirEntry) bool {
					isTargetExt := slices.Contains(ext, filepath.Ext(e.Name()))
					return e.IsDir() || !isTargetExt
				},
			)

			if len(files) <= int(gen) {
				fmt.Println("Filtered files fewer than the specified generation limit")
				return nil
			}

			timeSorted := make(map[int64][]os.DirEntry)
			for _, f := range files {
				var mod int64
				if info, err := f.Info(); err != nil {
					mod = time.Now().UnixMilli()
				} else {
					mod = info.ModTime().UnixMilli()
				}

				timeSorted[mod] = append(timeSorted[mod], f)
			}

			final := make([]os.DirEntry, 0, len(files))
			for _, t := range slices.Sorted(maps.Keys(timeSorted)) {
				v, ok := timeSorted[t]
				if !ok {
					continue
				}

				names := make([]string, 0, len(v))
				entries := make(map[string]os.DirEntry)
				for _, entry := range v {
					names = append(names, entry.Name())
					entries[entry.Name()] = entry
				}

				sorted := slices.Sorted(slices.Values(names))
				for _, s := range sorted {
					final = append(final, entries[s])
				}
			}

			deleteCount := len(final) - int(gen)
			for _, f := range final[:deleteCount] {
				if abs, err := filepath.Abs(filepath.Join(d, f.Name())); err == nil {
					if err := os.Remove(abs); err != nil {
						fmt.Printf("Failed to remove file (%s)\n", abs)
					} else {
						fmt.Printf("Success to remove file (%s)\n", abs)
					}
				}
			}

			return nil
		},
	}
}
