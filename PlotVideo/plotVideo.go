package main

import "time"
import "os"
import "syscall"
import "os/exec"
import "os/signal"
import "io"

func main() {
    end := make(chan os.Signal, 1)
    end2 := make(chan bool, 0)
    signal.Notify(end, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
    reader := newFfmpegReader("/run/user/1000/apple.webm")
    go gnuplotRunner(end, end2)
    for len(end) == 0 {
        pipe, err := os.OpenFile("/run/user/1000/gnuplotPipe1", os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0600)
        if err != nil {
            panic(err)
        }
        io.CopyN(pipe, reader.bufReader, int64(reader.bytesToRead))
        pipe.Close()
        <-end2
        os.Rename("/run/user/1000/gnuplotPipe1", "/run/user/1000/gnuplotPipe")
    }
}

var target time.Time

func sleep() {
    time.Sleep(time.Until(target))
    target = time.Now().Add(time.Millisecond * 133)
}

func gnuplotRunner(end chan os.Signal, end2 chan bool) {
    gnuplot := exec.Command("gnuplot")
    pipe, _ := gnuplot.StdinPipe()
    gnuplot.Stderr = os.Stderr
    gnuplot.Start()
    pipe.Write([]byte("set palette defined (0 \"black\", 255 \"white\")\nset cbrange [0:255]\n"))
    for len(end) == 0 {
        pipe.Write([]byte("plot [0:960] [0:720] \"/run/user/1000/gnuplotPipe\" binary array=960x720 format='%uint8' flipy with image\n"))
        sleep()
        end2<-true
    }
    close(end2)
    pipe.Close()
    gnuplot.Process.Kill()
    gnuplot.Wait()
}

