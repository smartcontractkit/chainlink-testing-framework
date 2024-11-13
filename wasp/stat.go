package wasp

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/pbnjay/memory"
	"github.com/rs/zerolog/log"
	"runtime"
	"sync"
	"time"
)

var (
	ResourcesThresholdCheckInterval = 5 * time.Second
	// CPUIdleThresholdPercentage is default CPU idle threshold
	CPUIdleThresholdPercentage = 20
	// MEMFreeThresholdPercentage is default MEM free threshold
	MEMFreeThresholdPercentage = 0
)

var once = &sync.Once{}

// CPUCheckLoop initializes a goroutine that continuously monitors the CPU and memory usage on a Linux system. 
// It checks the CPU idle percentage and the free memory percentage at regular intervals defined by 
// ResourcesThresholdCheckInterval. If the CPU idle percentage falls below the specified CPUIdleThresholdPercentage 
// or the free memory percentage falls below the MEMFreeThresholdPercentage, it logs a fatal error indicating 
// that the resource threshold has been triggered. This function ensures that resource usage is monitored 
// in a concurrent manner without blocking the main execution flow.
func CPUCheckLoop() {
	once.Do(func() {
		//nolint
		if runtime.GOOS == "linux" {
			go func() {
				for {
					time.Sleep(ResourcesThresholdCheckInterval)
					stat, err := linuxproc.ReadStat("/proc/stat")
					if err != nil {
						log.Fatal().Err(err).Send()
					}
					s := stat.CPUStatAll
					cpuPerc := float64((s.Idle * 100) / (s.User + s.Nice + s.System + s.Idle + s.IOWait + s.IRQ + s.SoftIRQ + s.Guest + s.GuestNice))
					log.Debug().Float64("CPUIdle", cpuPerc).Msg("Checking CPU load")
					freeMemPerc := float64(memory.FreeMemory()*100) / float64(memory.TotalMemory())
					log.Debug().Float64("FreeMEM", freeMemPerc).Msg("Free memory percentage")
					if cpuPerc <= float64(CPUIdleThresholdPercentage) || freeMemPerc <= float64(MEMFreeThresholdPercentage) {
						log.Fatal().Msgf("Resources threshold was triggered, CPUIdle: %.2f, FreeMEM: %.2f", cpuPerc, freeMemPerc)
					}
				}
			}()
		}
	})
}
