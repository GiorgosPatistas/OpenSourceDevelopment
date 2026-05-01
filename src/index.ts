#!/usr/bin/env node

import { Command } from "commander";

// Declare the program 

const program = new Command();

// Add actions onto that CLI
program
    .argument("<string>", "strings to log")
    .action((message: string) => {
        console.log(`Hello ${message}`);
    })
    .description('Say hello');

program.command("add <numbers...>").action((numbers: number[]) => {
    const total = numbers.reduce((a, b) => a + b, 0);
    console.log(`The total is ${total}`);
}).description('Add numbers and log the total');

program
    .command("get-max-number <numbers...>").action((numbers: number[]) => {
    const max = Math.max(...numbers);
    console.log(`The maximum number is ${max}`);
}).description('Get the maximum number from a list of numbers');

// Execute the CLI with the given arguments

program.parse(process.argv);