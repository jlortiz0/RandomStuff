# PlotVideo

"Plots" a video in gnuplot-qt by streaming a sequence of uncompressed gray8 data to it. It's kind of cheating since it uses an image display mode.

I know that the sleep time of 133 ms seems strange, but for some reason if I use the expected sleep duration of 33 ms to get 30 FPS, it runs at 4x speed.

Oh yeah, it's hardcoded for 30 FPS 960x720 because that was the framerate and resolution of the version of Bad Apple that I had.