const repl = require('node:repl');


repl.start({
    prompt: '\0',
    input: process.stdin,
    output: process.stdout,
    ignoreUndefined: true,
    preview: false,
    terminal: false,
    writer: (obj) => {
        // catch errors and return the message
        // for some reason Uncaught doesn't pass obj instanceof Error
        if (obj instanceof Error || obj.constructor.name === 'Error') {
            return obj.message;
        }

        return '';
    }
});
