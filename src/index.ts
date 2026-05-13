#!/usr/bin/env node

import { Command } from 'commander';
import { spawn } from 'child_process';
import * as fs from 'fs/promises';
import * as path from 'path';
import { getEnginePath } from './utils';

const program = new Command();

// --- Logic that calls the Go engine ---
async function buildSite(targetDir: string) {
    console.log(`Starting site build from directory: ${targetDir}...`);

    try {
        // Check if the directory provided by the user exists
        const fullPath = path.resolve(process.cwd(), targetDir);
        await fs.access(fullPath);

        const enginePath = getEnginePath();

        console.log(`Calling Go engine: ${enginePath}`);

        // Spawn the Go executable and pass the directory path
        const goProcess = spawn(enginePath, [fullPath]);

        // Listen for output from Go's stdout (e.g. success messages)
        goProcess.stdout.on('data', (data) => {
            console.log(data.toString());
        });

        // Listen for errors from Go's stderr
        goProcess.stderr.on('data', (data) => {
            console.error(`Error from Go engine: ${data.toString()}`);
        });

        // Handle process exit
        goProcess.on('close', (code) => {
            if (code === 0) {
                console.log('✅ Site generated successfully!');
            } else {
                console.error(`❌ Process failed with exit code ${code}`);
            }
        });

    } catch (error) {
        console.error(`❌ Error: Directory '${targetDir}' not found or not accessible.`);
    }
}

// --- Commander setup ---
program
  .name('my-ssg')
  .description('A Static Site Generator built with TypeScript and Go')
  .version('1.0.0');

program
  .command('build')
  .description('Converts the Markdown files in a directory into HTML')
  .argument('<directory>', 'The directory containing the Markdown files')
  .action((directory) => {
      buildSite(directory);
  });

program.parse();
