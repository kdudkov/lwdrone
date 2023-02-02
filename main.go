package main

import (
	"flag"
	"fmt"
	"github.com/kdudkov/lwdrone/msg"
	"os"
	"strings"
)

var gitRevision, gitBranch string

func getDrone() *msg.Lwdrone {
	drone := msg.NewDrone()
	if err := drone.SetTime(); err != nil {
		panic(err)
	}
	return drone
}

func main() {
	info := flag.Bool("info", false, "get drone info")
	photo := flag.Bool("photo", false, "take single photo")
	stream := flag.Bool("stream", false, "start video streaming")
	hq := flag.Bool("hq", false, "use high quality streaming")
	fname := flag.String("outfile", "h264out", "output H264 stream file name")

	flag.Parse()

	if *info {
		drone := getDrone()
		c, err := drone.GetConfig()
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		fmt.Printf("version: %s\n", c.Version)
		fmt.Printf("flash mounted: %d\n", c.SdcMounted)
		fmt.Printf("flash size: %d MiB\n", c.SdcSize/1024/1024)
		fmt.Printf("flash free: %d MiB (%.d%%)\n", c.SdcFree/1024/1024, 100.*c.SdcFree/c.SdcSize)
		fmt.Printf("time: %s\n", c.Time)
		return
	}

	if *stream {
		drone := getDrone()
		var f *os.File
		if *fname == "-" {
			f = os.Stdout
		} else {
			f, _ = os.Create(*fname)
		}
		if err := drone.StartStream(*hq, f); err != nil {
			fmt.Println(err)
		}
		return
	}

	if *photo {
		drone := getDrone()
		p, err := drone.TakePicture()
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		name := p.Path[strings.LastIndex(p.Path, "/")+1:]
		f, _ := os.Create(name)
		f.Write(p.Data)
		f.Close()
		fmt.Printf("writing file %s\n", name)
		return
	}

	fmt.Printf("version %s\n\n", gitRevision)
	flag.Usage()
}
