package mars.mips.instructions.syscalls;

import mars.Globals;
import mars.ProcessingException;
import mars.ProgramStatement;
import mars.mips.hardware.AddressErrorException;
import mars.mips.hardware.RegisterFile;
import mars.util.SystemIO;

/**
 * Service to display string stored starting at address in $a0 onto the console.
 */
public class SyscallPrintString extends AbstractSyscall {

    /**
     * Build an instance of the Print String syscall. Default service number is
     * 4 and name is "PrintString".
     */
    public SyscallPrintString() {
        super(4, "PrintString");
    }

    /**
     * Performs syscall function to print string stored starting at address in
     * $a0.
     */
    public void simulate(ProgramStatement statement) throws ProcessingException {
        int byteAddress = RegisterFile.getValue(4);
        char ch;
        char buf[] = new char[256];
        int bufCount = 0;
        try {
            ch = (char) Globals.memory.getByte(byteAddress);
            // won't stop until NULL byte reached!
            while (ch != 0) {
                //SystemIO.printString(Character.toString(ch));
                buf[bufCount] = ch;
                bufCount++;
                if (bufCount >= 256) {
                    SystemIO.printString(new String(buf));
                    bufCount = 0;
                }
                byteAddress++;
                ch = (char) Globals.memory.getByte(byteAddress);
            }
            if (bufCount != 0) {
                SystemIO.printString(new String(buf, 0, bufCount));
            }
        }
        catch (AddressErrorException e) {
            throw new ProcessingException(statement, e);
        }
    }
}
