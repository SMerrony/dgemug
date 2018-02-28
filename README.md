# dgemug
Go (Golang) packages of routines used in emulating historical Data General minicomputers

## cmd/dginstr
Dginstr generates Go instruction definitions from a CSV source.

## devices
Emulation of various DG peripherals
 * Bus
 * Disk6061 - Moving-head Disk - Type 6061
 * Magtape6026 - Magnetic Tape - Type 6026
 * TTO - console output

(No console input (TTI) as it is initmately tied to SCP)

## logging/debugLogs
DebugLogs is a fast memory-based circular logging subsystem with a facility to write out the logs to disk at the end of a run.

## memory
This package emulates the volatile memory of DG minis including the stacks and BMCDCH.

## util
This package contains mainly type conversions/extractions/tests for the datatypes used in DG mini emulation.

