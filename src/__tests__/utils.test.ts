import { describe, it, expect } from 'vitest'
import { getEnginePath } from '../utils'

describe('getEnginePath', () => {
    it('επιστρέφει το Windows binary στο win32', () => {
        const result = getEnginePath('win32', 'x64')
        expect(result).toContain('engine-windows.exe')
        expect(result).toContain('bin')
    })

    it('επιστρέφει το macOS arm64 binary σε darwin/arm64', () => {
        const result = getEnginePath('darwin', 'arm64')
        expect(result).toContain('engine-mac-arm64')
        expect(result).toContain('bin')
    })

    it('επιστρέφει το macOS Intel binary σε darwin/x64', () => {
        const result = getEnginePath('darwin', 'x64')
        expect(result).toContain('engine-mac')
        expect(result).not.toContain('arm64')
        expect(result).toContain('bin')
    })

    it('επιστρέφει το Linux binary στο linux', () => {
        const result = getEnginePath('linux', 'x64')
        expect(result).toContain('engine-linux')
        expect(result).toContain('bin')
    })

    it('κάθε platform επιστρέφει path με bin directory', () => {
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
