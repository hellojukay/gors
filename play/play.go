package play

import (
	"fmt"
)

type Player struct {
	Filename string
}

func (p *Player) Execute() {
	fmt.Println("gors playing started.")

	// bytes, err := ioutil.ReadFile(p.Filename)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, frame := range destination.Frames {
	// 	d, _ := time.ParseDuration(fmt.Sprintf("%d%s", frame.Delay, "ms"))
	// 	time.Sleep(d)
	// 	fmt.Print(frame.Data)
	// }

	fmt.Println("gors recording finished.")
}
