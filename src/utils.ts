import * as path from 'path';
import * as os from 'os';

/**
 * Returns the absolute path to the Go engine binary for the given platform/arch.
 *
 * Parameters have default values (os.platform() / os.arch()) so the function
 * works normally at runtime but is easily testable without mocking.
 *
 * @param platform - e.g. 'win32', 'darwin', 'linux'
 * @param arch     - e.g. 'x64', 'arm64'
 */
export function getEnginePath(
    platform: string = os.platform(),
    arch: string = os.arch()
): string {
    const binDir = path.resolve(__dirname, '../bin');

    if (platform === 'win32') {
        return path.join(binDir, 'engine-windows.exe');
    } else if (platform === 'darwin') {
        if (arch === 'arm64') {
            return path.join(binDir, 'engine-mac-arm64');
        }
        return path.join(binDir, 'engine-mac');
    } else {
        return path.join(binDir, 'engine-linux');
    }
}
