package mars.mips.instructions.syscalls;

import mars.ProcessingException;
import mars.ProgramStatement;
import mars.mips.hardware.RegisterFile;
import mars.util.Binary;
import mars.util.SystemIO;

/**
 * Service to display integer stored in $a0 on the console.
 *
 */
public class SyscallPrintIntBinary extends AbstractSyscall {

    /**
     * Build an instance of the Print Integer syscall. Default service number is
     * 1 and name is "PrintInt".
     */
    public SyscallPrintIntBinary() {
        super(35, "PrintIntBinary");
    }

    /**
     * Performs syscall function to print on the console the integer stored in
     * $a0, in hexadecimal format.
     */
    public void simulate(ProgramStatement statement) throws ProcessingException {
        SystemIO.printString(Binary.intToBinaryString(RegisterFile.getValue(4)));
    }
}
