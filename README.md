# dgemug
Go (Golang) packages of routines used in emulating historical Data General minicomputers

## aosvs
This package is used for partial emulation of an AOS/VS system at the user level. See cmd/vsemug below.

## cmd/dginstr
Dginstr generates Go DG CPU opcode definitions from a CSV source.

This command can be installed by performing a `go install` from its directory.  This is required prior to developing any of the related emulators.

## cmd/vsemug
VSemuG is an attempt at a user-level AOS/VS emulator.  
It is mainly intended to provide a testbed for the mvcpu package and is unlikely to be especially useful (or complete) in its own right.

It has its own [Readme](cmd/vsemug/README.md) and [Status](cmd/vsemug/STATUS.md) pages.

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

