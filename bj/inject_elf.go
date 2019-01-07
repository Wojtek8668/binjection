package bj

import (
	"debug/elf"
	"io/ioutil"
	"log"
)

// ElfBinject - Inject shellcode into an ELF binary
func ElfBinject(sourceFile string, destFile string, shellcode string, config *BinjectConfig) error {

	//
	// BEGIN CODE CAVE DETECTION SECTION
	//

	buf, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	type Cave struct {
		Start, End uint64
	}
	var caves []Cave

	MIN_CAVE_SIZE := 94

	count := 1
	caveStart := uint64(0)
	for i := uint64(0); i < uint64(len(buf)); i++ {
		switch buf[i] {
		case 0:
			if count == 1 {
				caveStart = i
			}
			count++
		default:
			if count >= MIN_CAVE_SIZE {
				caves = append(caves, Cave{Start: caveStart, End: i})
			}
			count = 1
		}
	}

	elfFile, err := elf.Open(sourceFile)
	if err != nil {
		return err
	}

	for _, cave := range caves {
		for _, section := range elfFile.Sections {
			if cave.Start >= section.Offset && cave.End <= (section.Size+section.Offset) &&
				cave.End-cave.Start >= uint64(MIN_CAVE_SIZE) {
				log.Printf("Cave found (start/end/size): %d / %d / %d \n", cave.Start, cave.End, cave.End-cave.Start)
			}
		}
	}

	//
	// END CODE CAVE DETECTION SECTION
	//

	return nil
}
