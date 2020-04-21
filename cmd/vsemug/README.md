# vsemug
User-level AOS/VS Emulator

Only if you are ***developing and changing instructions*** you will need to precede the build with: `go generate`

Build with: `go build -tags virtual`

Run with: `./vsemug -pr programs/LOOPS1.PR` 
then connect to port 10001 with a DASHER-compatible terminal emulator such as 
[DasherG](https://github.com/SMerrony/DasherG).

Current status is in [STATUS.md](./STATUS.md)
