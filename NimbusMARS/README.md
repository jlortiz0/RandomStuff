# Nimbus MARS

It's [MARS](https://courses.missouristate.edu/KenVollmar/MARS/index.htm) but with some features to make it a little easier for me to use.

An incomplete list of features:
 - Dark mode
 - Smart home
 - Changed end key to be more consistent
 - Tabs to spaces
 - Memory address jump
 - Auto tab after label
 - Autoindent after label, reset after blank line
 - Faster capped execution speeds
 - "print string" syscall is now buffered
 - Removed rectangular selection mode (why did I do this?)
 - Partially fixed line number desync
 - Added errno button to view most recent syscall error (currently only works for I/O)
 - Warning when opening non-ASCII file
 - Tab to indent all highlighted lines
 - Tweaked defaults for Bitmap Viewer

This program has some known issues:
 - Not all menus have been modified to account for dark mode
 - Clicking will sometimes select instead of moving the cursor
 - Switching from "1 instr/sec" to "unlimited speed" while in the memory tab will hang
 - Probably some more issues