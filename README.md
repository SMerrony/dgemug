# dgemug
Go (Golang) package of routines used in emulating historical Data General minicomputers

## cmd/dginstr
Dginstr generates Go instruction definitions from a CSV source.

## logging/debugLogs
DebugLogs is a fast memory-based circular logging subsystem with a facility to write out the logs to disk at the end of a run.

## memory/...
These files emulate the volatile memory of DG minis including the stacks and BMCDCH.

## util/...
These routines are mainly type conversions/extractions/tests for the datatypes used in DG mini emulation.

