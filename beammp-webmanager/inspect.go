package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func runInspect(m *Manager) {
	cfg, err := m.readConfig()
	if err != nil {
		fmt.Printf("cannot read %s in %s: %v\n", configName, m.dir, err)
		return
	}
	client, deactMods, _, _, customMaps := m.paths(cfg)
	fmt.Println("Directory:     ", m.dir)
	fmt.Println("ResourceFolder:", m.resourceFolder(cfg))
	fmt.Println("Current map:   ", strings.TrimSpace(tomlString(cfg, "General", "Map")))

	for _, d := range []struct{ label, dir string }{
		{"Resources/Client", client},
		{"custom_maps", customMaps},
		{"deactivated_mods", deactMods},
	} {
		files := zipsIn(d.dir)
		fmt.Printf("\n=== %s (%d zips) ===\n", d.label, len(files))
		for _, f := range files {
			name := filepath.Base(f)
			zi, err := inspectZip(f)
			if err != nil {
				fmt.Printf("\n  %s\n    ERROR: %v\n", name, err)
				continue
			}
			kind := "mod"
			if zi.mapPath != "" {
				kind = "map"
			} else if zi.hasLevels {
				kind = "has levels/ but no info.json below it -> shown nowhere"
			}
			fmt.Printf("\n  %s\n", name)
			fmt.Printf("    Type:     %s\n", kind)
			if zi.mapPath != "" {
				fmt.Printf("    Map path: %s\n", zi.mapPath)
				fmt.Printf("    Title:    %s\n", zi.mapTitle)
				if zi.mapImage != "" {
					fmt.Printf("    Image:    %s   (found via: %s)\n", zi.mapImage, zi.imageFrom)
				} else {
					fmt.Printf("    Image:    NONE\n")
				}
			} else {
				fmt.Printf("    Title:    %s\n", zi.modTitle)
				if zi.modImage != "" {
					fmt.Printf("    Image:    %s\n", zi.modImage)
				} else {
					fmt.Printf("    Image:    NONE (no icon.png/icon.jpg in the zip)\n")
				}
			}
			if zi.mapImage == "" && zi.modImage == "" {
				fmt.Printf("    Images in the zip (%d):\n", len(zi.images))
				for i, img := range zi.images {
					if i == 15 {
						fmt.Printf("      ... and %d more\n", len(zi.images)-15)
						break
					}
					mark := " "
					if !browserSafe(img) {
						mark = "x"
					}
					fmt.Printf("      %s %s\n", mark, img)
				}
			}
		}
	}
	fmt.Println("\n(x = format the browser cannot display, e.g. .dds)")
}
