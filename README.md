# LWDrone

Simple program to communicate with a drone's lewei camera module.

based on [meekworth/pylwdrone](https://github.com/meekworth/pylwdrone)

get single photo:
`./lwdrone -photo`

stream to ffplay:
`./lwdrone -stream -hq -outfile - | ffplay -i -fflags nobuffer -flags low_delay -probesize 32 -sync ext -`