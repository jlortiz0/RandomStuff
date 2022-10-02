package main

import "time"
import "os"
import "syscall"
import "os/exec"
import "os/signal"
import "io"

const HEIGHT = "360"
const WIDTH = "480"

func main() {
    end := make(chan os.Signal, 1)
    end2 := make(chan bool, 0)
    signal.Notify(end, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
    reader := newFfmpegReader("/run/user/1000/apple.webm")
    go gnuplotRunner(end, end2)
    for len(end) == 0 {
        pipe, err := os.OpenFile("/run/user/1000/gnuplotPipe1", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
        if err != nil {
            panic(err)
        }
        io.CopyN(pipe, reader.bufReader, int64(reader.bytesToRead))
        pipe.Close()
        <-end2
        os.Rename("/run/user/1000/gnuplotPipe1", "/run/user/1000/gnuplotPipe")
    }
}

func gnuplotRunner(end chan os.Signal, end2 chan bool) {
    gnuplot := exec.Command("gnuplot")
    pipe, _ := gnuplot.StdinPipe()
    gnuplot.Stderr = os.Stderr
    gnuplot.Start()
    pipe.Write([]byte("set palette defined (0 \"black\", 255 \"white\")\nset cbrange [0:255]\n"))
    target := time.NewTicker(time.Second / 10)
    for len(end) == 0 {
        pipe.Write([]byte("plot [0:" + WIDTH + "] [0:" + HEIGHT + "] \"/run/user/1000/gnuplotPipe\" binary array=" + WIDTH + "x" + HEIGHT + " format='%uint8' flipy with image\n"))
        <-target.C
        end2 <- true
    }
    close(end2)
    target.Stop()
    pipe.Close()
    gnuplot.Process.Kill()
    gnuplot.Wait()
}
