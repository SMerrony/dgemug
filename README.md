# dgemug
Go (Golang) packages of routines used in emulating historical Data General minicomputers

## aosvs
This package is used to emulate an AOS/VS system at the user level.

## cmd/dginstr
Dginstr generates Go instruction definitions from a CSV source.

This command can be installed by performing a `go install` from its directory.  This is required prior to building any of the related emulators.

## devices
Emulation of various DG peripherals:
 * Bus
 * Disk4231
 * Disk6061 - Moving-head Disk, Type 6061 (AOS/VS - DPF)
 * Disk6239 - Moving-head Disk, Type 6239 (AOS/VS - DPJ)
 * Magtape6026 - Magnetic Tape, Type 6026
 * TTI - console input
 * TTO - console output

## logging/debugLogs
DebugLogs is a fast memory-based circular logging subsystem with a facility to write out the logs to disk at the end of a run.

## memory
This package emulates the volatile memory of DG minis including the stacks and BMCDCH.
N.B. There are physical (hardware) and logical (AOS/VS system) versions of some files here - 
choose which to use with an appropriate build tag.

## mvcpu
This package emulates an MV-class CPU at the machine instruction (opcode) level.

