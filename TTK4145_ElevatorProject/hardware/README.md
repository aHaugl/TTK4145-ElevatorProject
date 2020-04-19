# Hardware module

This module consists of driver code that's been handed out from [TTK4145](https://github.com/TTK4145). The code has been interfaced with GO
so that we can run the elevator with GO code. In addition, we've interfaced it with
elevator simulator also made by [TTK4145](https://github.com/TTK4145).

The harware module is neccesary to run the elevator and give feedback to the system.
Responsibilies of this module:
 - Setting/clearing lights (BUTTONS, FLOOR INDICATOR, STOP, DOOROPEN)
 - Setting motor direction (UP, DOWN, STOP)
 - Reading floor sensor signal
 - Reading button signals
