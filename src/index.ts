#!/usr/bin/env node

import { Command } from 'commander';
import { spawn } from 'child_process';
import * as fs from 'fs/promises';
import * as path from 'path';
import { getEnginePath } from './utils';

const program = new Command();

// --- Η λογική που καλεί την Go ---
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

// --- Στήσιμο Commander ---
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
