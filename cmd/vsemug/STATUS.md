# VSemuG Status
* Last Updated: 7 May 2020
* Last Significant Progress: 6 May 2020 (some byte-addressing issues resolved, some remain...)
  
## What Works? :+1:
The following 32-bit sample programs copied from a physical machine are working...
* HW.PR - Hello, World!
* HW2.PR - Hello, World! using CLI return message (18 Mar 2020)
* LOOPS1.PR - Basic looping constructs (15 Mar 2020)
* LOOPS2.PR - Further loops 
* LOOPS3.PR - Loops with LWDO and -ve values (18 Mar 2020)
* LOOPS4.PR - As LOOPS3 with external subroutines (18 Mar 2020)
* SPIGOT.PR - Calculate Pi to a thousand places using the spigot method (17 Mar 2020)
* STRINGTESTS.PR - Various string-handling routines (18 Mar 2020)
* TIMEOUT.PR - Uses ?GTMES and ?WDELAY to pause for n seconds (19 Mar 2020)

## What Doesn't Work? :-1:

The NADGUG library provides a good range of freely-available test targets...
  
* 32-bit NADGUG Games compiled for AOS/VS

  |    Game   | Bits |  Problem  |   Notes/Action   |
  |-----------|------|-----------|------------------|
  | 21        | 32 | Unmapped write in XWSTA                       | |
  | CB        | 32 | CB.PR - Wants to start in :PER                | ! |
  | CHESS     | 32 | Exits with no error                           | Shortly after ?IFPU |
  | DND       | 32 | Looping on output                             |  |
  | EMPIRE    | 32 | Instr 0xa079 nyi                              | |
  | EMPIRE2   | 32 | ?CON nyi                                      | ?CON |
  | FERRET    | 32 | Decimal Type 5 nyi in WSTI                    | |
  | FISH      | 32 | ?GLIST nyi                                    | ?GLIST |
  | FOOBAR    | 32 | LNADD nyi                                     | LNADD |
  | MMM       | 32 | LNADD nyi                                     | LNADD |
  | MORTGAGE  | 32 | Syscall ?TASK nyi                             | **** ?TASK |
  | QUEST     | 32 | QUEST_SERVER.PR - seems to be looping         | :-( |
  | QUEST     | 32 | QUEST.PR - QUEST Server is not up!            | :-)  |
  | SCRABBLE  | 32 | Keeps ?READing with no prompt                 | Is is an extended ?READ ? |
  | WUMPUS    | 32 | Hang/loop after displaying start screen       |  |
  | YAHTZEE   | 32 | Exits with Error Code: 71200                  | Almost immediately after ?MEMI |
  | ZORK      | 32 | Unmapped read in WCST                         |  |

  |  Folder  |  Program  | Bits |         Problem         |  Notes/Action  |
  |----------|-----------|------|-------------------------|----------------|
  | IMSLUTIL | HANGMAN   |  32  | Error in WDCMP          | WDCMP - After welcome shown |
  

* 16-bit NADGUG Games compiled for AOS/VS...  
  N.B. These might be handled quite differently by the OS - do not focus on them.  In particular, it seems
  that the initial memory setup may differ from ordinary 32-bit programs.

  |    Game   |  Bits  |  Problem  |   Notes/Action   |
  |-----------|--------|-----------|------------------|
  | ADVENTURE | 16 | Sys Call 0272 nyi                             | Call not listed in docs... |
  | ASTEROIDS | 16 | Reports error INSUFFICIENT MEMORY FOR PROGRAM | |
  | BRUTUS    | 16 | Tries to map already-mapped page in ?MEMI     | ?Not enough room between areas? check ?MEM | 
  | CONQUEST  | 16 | Reports error INSUFFICIENT MEMORY FOR PROGRAM | No room between unshared and shared areas, check ?MEM |
  | DICE      | 16 | ?TASK (16-bit) nyi                            | **** ?TASK |
  | HANGMAN   | 16 | Reports error INSUFFICIENT MEMORY FOR PROGRAM | No room between unshared and shared areas, check ?MEM |
  | OTHELLO   | 16 | Tries to map already-mapped page in ?MEMI     | |
  | PACMAN    | 16 | ?TASK (16-bit) nyi                            | **** ?TASK |
  | SERPENT   | 16 | JMPs to 0 after EJMP @0532                    |  |
  | SI        | 16 | ?TASK (16-bit) nyi                            | **** ?TASK |
  | STARTREK  | 16 | Not parsing Difficulty Level input            |  |
  | SPACEWAR  | 16 | ?TASK (16-bit) nyi                            | **** ?TASK |
  | THISSALA  | 16 | Sys Call 0272 nyi                             | Call not listed in docs... |
 

## What's Next?

Extended ?READ/?WRITE to @CONSOLE
