package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.2"

func PrintVersion() {
	fmt.Printf("Current nsfwdetector version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `
                 ____               __       __               __              
   ____   _____ / __/_      __ ____/ /___   / /_ ___   _____ / /_ ____   _____
  / __ \ / ___// /_ | | /| / // __  // _ \ / __// _ \ / ___// __// __ \ / ___/
 / / / /(__  )/ __/ | |/ |/ // /_/ //  __// /_ /  __// /__ / /_ / /_/ // /    
/_/ /_//____//_/    |__/|__/ \__,_/ \___/ \__/ \___/ \___/ \__/ \____//_/
`
	fmt.Printf("%s\n%75s\n\n", banner, "Current nsfwdetector version "+version)
}
