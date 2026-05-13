#!/usr/bin/env node

import { Command } from 'commander';
import { spawn } from 'child_process';
import * as fs from 'fs/promises';
import * as path from 'path';
import * as os from 'os';

const program = new Command();

// --- 1. Η λογική για να βρούμε το εκτελέσιμο της Go ---
function getEnginePath(): string {
    const platform = os.platform();
    const arch = os.arch();
    // Τα εκτελέσιμα βρίσκονται στον φάκελο 'bin'
    // στο ίδιο επίπεδο με τον φάκελο 'dist'
    const binDir = path.resolve(__dirname, '../bin');

    if (platform === 'win32') {
        return path.join(binDir, 'engine-windows.exe');
    } else if (platform === 'darwin') {
        // Apple Silicon (M1/M2/M3) ή Intel
        if (arch === 'arm64') {
            return path.join(binDir, 'engine-mac-arm64');
        }
        return path.join(binDir, 'engine-mac');
    } else {
        return path.join(binDir, 'engine-linux');
    }
}

// --- 2. Η λογική που καλεί την Go ---
async function buildSite(targetDir: string) {
    console.log(`Ξεκινάει το χτίσιμο του site από τον φάκελο: ${targetDir}...`);

    try {
        // Έλεγχος αν ο φάκελος που έδωσε ο χρήστης υπάρχει
        const fullPath = path.resolve(process.cwd(), targetDir);
        await fs.access(fullPath);

        const enginePath = getEnginePath();

        console.log(`Καλείται η μηχανή Go: ${enginePath}`);

        // Εδώ "ξυπνάμε" το εκτελέσιμο της Go και του περνάμε τον φάκελο
        const goProcess = spawn(enginePath, [fullPath]);

        // "Ακούμε" τι μας στέλνει η Go στο stdout (π.χ. μηνύματα επιτυχίας)
        goProcess.stdout.on('data', (data) => {
            console.log(data.toString());
        });

        // "Ακούμε" για σφάλματα από την Go
        goProcess.stderr.on('data', (data) => {
            console.error(`Σφάλμα από τη μηχανή Go: ${data.toString()}`);
        });

        // Τι κάνουμε όταν τελειώσει η Go
        goProcess.on('close', (code) => {
            if (code === 0) {
                console.log('✅ Το site δημιουργήθηκε με επιτυχία!');
            } else {
                console.error(`❌ Η διαδικασία απέτυχε με κωδικό ${code}`);
            }
        });

    } catch (error) {
        console.error(`❌ Σφάλμα: Ο φάκελος '${targetDir}' δεν βρέθηκε ή δεν είναι προσβάσιμος.`);
    }
}

// --- 3. Το στήσιμο του Commander ---
program
  .name('my-ssg')
  .description('Ένας Static Site Generator με TypeScript και Go')
  .version('1.0.0');

program
  .command('build')
  .description('Μετατρέπει τα Markdown αρχεία ενός φακέλου σε HTML')
  .argument('<directory>', 'Ο φάκελος που περιέχει τα Markdown αρχεία')
  .action((directory) => {
      buildSite(directory);
  });

program.parse();