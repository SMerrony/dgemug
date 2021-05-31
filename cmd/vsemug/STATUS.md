# VSemuG Status
_Work has stopped on this project. It continues in https://github.com/SMerrony/dgemua_
* Last Updated: 31 May 2021
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
  | 21        | 32 | BOOTER.PR Unmapped write in XWSTA             | **** Unmapped: Check preceeding XLEF 2, 02,AC3 |
  | CB        | 32 | CB.PR - Wants to start in :PER                | ! |
  | CHESS     | 32 | Exits with no error                           | Shortly after ?IFPU |
  | DND       | 32 | Calling ?ERMSG after ?OPENs                   |  |
  | EMPIRE    | 32 | Runs for a bit!                               | Some screen corruption |
  | EMPIRE2   | 32 | Calling ?ERMSG after ?OPENs                   |  |
  | FERRET    | 32 | Decimal Type 5 nyi in WSTI                    | |
  | FISH      | 32 | ?GLIST nyi                                    | ?GLIST |
  | FOOBAR    | 32 | Error in Line 205                             |  |
  | MMM       | 32 | Unmapped read in WCMV                         | **** Unmapped |
  | MORTGAGE  | 32 | Syscall ?TASK nyi                             | **** ?TASK |
  | QUEST     | 32 | QUEST_SERVER.PR - Error 21 recreating shared_data_file | ?RECREATE bug |
  | QUEST     | 32 | QUEST.PR - Unmapped read in WCST              | **** Unmapped |
  | SCRABBLE  | 32 | Calling ?ERMSG after ?OPENing SEED file       | ?OPEN bug?  |
  | WUMPUS    | 32 | Hang/loop after displaying start screen       |  |
  | YAHTZEE   | 32 | Exits with Error Code: 71200                  | Almost immediately after ?MEMI |
  | ZORK      | 32 | Unmapped read in WCMV (underflow)             | **** Unmapped: Check last few instrs for math error |

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
  | STARTREK  | 16 | FSNER nyi                                     | FSNER |
  | SPACEWAR  | 16 | ?TASK (16-bit) nyi                            | **** ?TASK |
  | THISSALA  | 16 | Sys Call 0272 nyi                             | Call not listed in docs... |
 

## What's Next?

Do all the unmapped read/writes have a common cause?

Possbile issue in EMPIRE2T & DND is with ?OPEN of @CONSOLE
LCALL...