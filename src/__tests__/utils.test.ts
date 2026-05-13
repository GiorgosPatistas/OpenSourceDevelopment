import { describe, it, expect } from 'vitest'
import { getEnginePath } from '../utils'

describe('getEnginePath', () => {
    it('returns the Windows binary on win32', () => {
        const result = getEnginePath('win32', 'x64')
        expect(result).toContain('engine-windows.exe')
        expect(result).toContain('bin')
    })

    it('returns the macOS arm64 binary on darwin/arm64', () => {
        const result = getEnginePath('darwin', 'arm64')
        expect(result).toContain('engine-mac-arm64')
        expect(result).toContain('bin')
    })

    it('returns the macOS Intel binary on darwin/x64', () => {
        const result = getEnginePath('darwin', 'x64')
        expect(result).toContain('engine-mac')
        expect(result).not.toContain('arm64')
        expect(result).toContain('bin')
    })

    it('returns the Linux binary on linux', () => {
        const result = getEnginePath('linux', 'x64')
        expect(result).toContain('engine-linux')
        expect(result).toContain('bin')
    })

    it('every platform returns a path containing the bin directory', () => {
        const cases = [
            ['win32', 'x64'],
            ['darwin', 'arm64'],
            ['darwin', 'x64'],
            ['linux', 'x64'],
        ] as const

        for (const [platform, arch] of cases) {
            expect(getEnginePath(platform, arch)).toContain('bin')
        }
    })
})
