# Data Channel and Burst Multiplexor Channel
_This document is paraphrased from the AOS/VS Internals Reference Manual for AOS/VS 5.0 and the later Principles of Operation (1983)_

The data channel (DCH) provides I/O communications for medium-speed devices (eg. tape drives) and synchronous communication.  The burst multiplexor channel (BMC) is a high speed communications pathway that transfers data directly between main memory and high-speed peripherals (eg. disk drives).  **I/O-to-memory transfers for both DCH and BMC always bypass the address translator.**

## DCH/BMC Maps
A map controls a DCH or BMC.  This map is a series of contiguous map slots, each of which contains a pair of map registers - an even-numbered register and its corresponding odd-numbered register.

An MV computer supports 16 DCH maps, each of which contains 32 map slots.  The DCH sends to the processor a logical address with each data transfer.  The processor translates the logical address into a physical address using the appropriate map slot for that address.

The device controller performing the data transfer controls the BMC.  No program control or CPU interaction is required, except when setting up the BMC's map table.  The BMC has two address modes and contains its own map.

### BMC Address Modes
The BMC operates in either the unmapped (physical) mode, or the mapped (logical) mode.

In the unmapped mode, the BMC receives 20-bit addresses from the device controllers and passes them directly to memory.  As the BMC transfers each data word to or from memory, it increments the destination address, causing successive words to move to or from consecutive locations in memory.

If the controller specifies the mapped mode for data transfer, the high-order 10 bits of the logical address form a logical page number, which the BMC map translates into a 10-bit physical page number.  This page number, combined with the 10 low-order bits from the logical address, forms a 20-bit physical address, which the BMC uses to access memory.

## BMC Map
The BMC uses its own map to translate logical page numbers into physical ones.  (On machines that implement it, the **SSPT** instruction defines the memory locations of the BMC map.)  The map contains 1024 map registers, the odd-numbered registers each containing a 10-bit physical page number.  The BMC uses the logical page number as an index into the map table, and the contents of the selected map register becomes the high-order 10 bits of the physical address.

Note that when the BMC performs a mapped transfer, it increments the destination address after it moves each data word.  If the increment causes an overflow out of the 10 low-order bits, this selects a new map register for subsequent address translation.  Depending upon the contents of the map table, this could mean that the BMC cannot transfer successive words to or from consecutive pages in memory.

## DCH/BMC Registers
An MV computer system contains 512 DCH registers and 1024 BMC registers.  The map registers are numbered from 0 through 07777.

| Registers | Description|
| ----------|------------|
|0000 - 3776 | Even-numbered regs, most significant half of BMC map posns 0 - 1777|
|0001 - 3777 | Odd-numbered regs, least significant half or BMC map posns 0 - 1777|
|4000 - 5776 | Even-numbered regs, most significant half of DCH map posns 0 - 777|
|4001 - 5777 | Odd-numbered regs, least significant half of DCH maps posns 0 - 777|
|001 - 7677 | (reserved)|
|7700        | I/O channel status register|
|7701        | I/O channel mask register|
|7702        | CPU dedication control|
|7703 - 7777 | (reserved)|

The register formats are the same for DCH and BMC registers.

### Even-Numbered Register Format
|V | D | Hardware Reserved|
|--|---|------------------|
|0 | 1 |2 -             15|

V - validity bit; if 1 then processor denies access
D - data bit; if 0 the channel transfers data, if 1 the channel transfers zeroes
Reserved should be written to with zeroes; reading these returns an undefined state.

### Odd-Numbered Register Format
|Res | Physical Page Number|
|----|---------------------|
|0   | 1 - 15 (N.B. was 2-15 early on)|

Res - hardware reserved
Physical Page Number - associated with logical page reference.

### I/O Channel Definition Register Format
_N.B. [SHM] There are differences between the earlier Internals doc and the later PoP..._

|ICE | Res | BVE | DVE | DCH | BMC | BAP | BDP | Res  | DME | 1 |
|----|-----|-----|-----|-----|-----|-----|-----|------|-----|---|
|0   | 1-2 | 3   | 4   | 5   | 6   | 7   | 8   | 9-13 | 14  | 15|

_N.B. Writing to bits 3,4,7,8 or 14 with a 1 complements these bits.  The **IORST** instruction clears these bits._

  * ICE - channel error flag
  * Res - reserved
  * BVE - BMC validity error flag, if 1 BMC protect error has occurred
  * DVE - DCH validity error flag, if 1 DCH protect errro has occurred
  * DCH - Data Channel transaction, if 1 a DCH transaction is in progress
  * BMC - BMC transfer flag, if 1 and BMC transfer is in progress
  * BAP - BMC address parity error
  * BDP - BMC data parity error
  * DIS - disable block transfer
  * DME - DCH mode, if 1 DCH mapping is enabled
  * 1 - always set to 1 (previous version was 0 - this may differentiate - SHM)

### I/O Channel Status Register Format

Old version...

|E | Res | XDCH | 1 | MSK | INT|
|--|-----|------|---|-----|----|
|0 | 1-11| 12   |13 | 14  | 15 |

  * E - error flag
  * Res - reserved
  * XDCH - DCH map slots and operations supported
  * 1 - always set to 1
  * MSK - prevents all devices connected to channel from interrupting the CPU
  * INT - Interrupt pending

New version...

|ERR | Res | DTO | MPE | 1 | 1 | CMB | INT|
|----|-----|-----|-----|---|---|-----|----|
|0   | 1-9 | 10  |  11 |12 |13 | 14  | 15 |

  * ERR - If 1 the I/O channel has detected an error
  * Res - Reserved
  * DTO - Data Time Out
  * MPE - Map Parity Error
  * 1 - Always set to 1 indicating extended DCH map slots and ops are supported
  * CMB - Channel Mask Bit, if 1 prevents all devices interrupting the CPU; however the **INTA** instruction returns the device code of any device with its DONE flag set
  * INT - Interrupt Pending, if 1 the channel is attempting to interrupt the CPU


### I/O Channel Mask Register Format
This write-only I/O channel mask register (07701) specified a mask flag for each channel.  When an I/O channel mask flag is set to 1, the CPU ignores all interrupt requests from devices on that channel.

The **INTA** instruction with a channel number returns on that channel the device code of the highest prioroty interrupting device which has its DONE flag set.  With channel 7, the **INTA** instruction returns the device code of the highest priority interrupting device on the highest prioroty channel, regardless of the state of the I/O channel mast register flags.

An I/O channnel Bus Reset **PRTRST** instruction sets the mask bit to 0 for one channel, or for all channels (7).

NOTE: A **CIO** read to the I/O/ channel mask register produces undefined results.

|Res | MK0 | MK1...MK6 | R |
|----|-----|-----------|---|
|0-7 | 8   |  9 - 14   | 15 |

  * Res - reserved
  * MK0 - prevents all devices on channel 0 from interrupting CPU
  * MK1 - prevents all devices on channel 1 from interrupting CPU
  * ...etc
  * R - reserved

NOTE: A system reset sets MK0 to 0 and MK1-6 to 1.

### CPU Dedication Control
Each IOC contains a 16-bit command I/O register (07702) which controls CPU dedication.  This read/write register is available only which the system is in dedicated mode.

|Res | CPU|
|----|----|
|1-14| 15 |

  * Res - reserved
  * CPU - CPU No. (0 or 1) to which all NOVA-type interrupts (except cross-interrupts) will be directed

On a system reset, CPU is set to the value of the initial CPU.  On execution of an **IORST** instruction, it is set to the value of the CPU which issued the instruction.

## DCH/BMC Map Instructions
The **CIO**, **CIOI**, and **WLMP** instructions initiate DCH/BMC map loads and reads when in mapped mode, with the **LPHY** instruction used for loads in unmapped mode.  The I/O channel sets its BUSY flag to 1 when a map load or read is in progress.  There is no DONE flag, and the channel never causes interrupts.
  * WLMP - Loads BMC/DCH map slots from memory
  * CIO, CIOI - Returns BMC/DCH status or loads map regs (1/2 slot) from accumulators
  * LPHY - Translates logical addresses to physical ones
  * IORST - Clears bits 3,4,7,8 & 14 of the I/O channel definition register, which disables data channel maps.
  
