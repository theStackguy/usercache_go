//go:build linux
// +build linux

package src

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
)

func operatingSystemAvailableMemory(c chan<- error, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	close(c)
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		_ = osavailableMemory.Swap(ZERO)
		c <- errProcMem
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, freeMemoryConst) {
			fields := strings.Fields(line)
			if len(fields) >= TWO {
				value, err := strconv.ParseInt(fields[ONE], PARSE_START, PARSE_STOP)
				if err != nil {
					_ = osavailableMemory.Swap(ZERO)
					c <- errParseProcMem
					return
				}
				result := uint64(value * BINARY_SYS_NUMBER)
				_ = osavailableMemory.Swap(result)
				c <- nil
				return
			}
		}
	}
	if err := scanner.Err(); err != nil {
		_ = osavailableMemory.Swap(ZERO)
		c <- errReadProcMem
		return
	}
	_ = osavailableMemory.Swap(ZERO)
	c <- errKeyNotFoundInProcMem
	return
}
