const rawErrors = `
"ERROR_UNKNOWN":                         1,
"ERROR_DECODE":                          2,
"ERROR_NOT_IMPLEMENTED":                 3,
"ERROR_BUSY":                            4,
"ERROR_CONTINUOUS_COMMAND_INTERRUPTED":  14,
"ERROR_INVALID_PARAMETERS":              15,
"ERROR_STORAGE_NOT_READY":               5,
"ERROR_STORAGE_EXIST":                   6,
"ERROR_STORAGE_NOT_EXIST":               7,
"ERROR_STORAGE_INVALID_PARAMETER":       8,
"ERROR_STORAGE_DENIED":                  9,
"ERROR_STORAGE_INVALID_NAME":            10,
"ERROR_STORAGE_INTERNAL":                11,
"ERROR_STORAGE_NOT_IMPLEMENTED":         12,
"ERROR_STORAGE_ALREADY_OPEN":            13,
"ERROR_STORAGE_DIR_NOT_EMPTY":           18,
"ERROR_APP_CANT_START":                  16,
"ERROR_APP_SYSTEM_LOCKED":               17,
"ERROR_APP_NOT_RUNNING":                 21,
"ERROR_APP_CMD_ERROR":                   22,
"ERROR_VIRTUAL_DISPLAY_ALREADY_STARTED": 19,
"ERROR_VIRTUAL_DISPLAY_NOT_STARTED":     20,
"ERROR_GPIO_MODE_INCORRECT":             58,
"ERROR_GPIO_UNKNOWN_PIN_MODE":           59,
`.trim().split('\n');

const errorMessageCode = rawErrors.map(line => {
    const errorCode = line.substring(0, line.length - 1).split(' ');
    const rawError = errorCode.shift();
    const rawCode = errorCode.pop();

    const error = rawError.substring(7, rawError.length - 2).toLowerCase().split('');
    let variable = 'Err';
    let message = '';
    let upper = true;

    for (const char of error) {
        if (char === '_') {
            message += ' ';
            upper = true;
        } else if (upper) {
            variable += char.toUpperCase();
            message += char;

            upper = false;
        } else {
            variable += char;
            message += char;
        }
    }

    message += ' error';

    const code = parseInt(rawCode);

    return [variable, message, code];
});

const errors = errorMessageCode.map(([variable, message, _]) => {
    return `var ${variable} = errors.New("${message}");`;
}).join('\n');

const mapping = errorMessageCode.map(([variable, _, code]) => {
    return `${code}: ${variable},`
}).join('\n');

console.log(`
package rpc

${errors}

var (
    errorCodeMapping = map[flipper.CommandStatus]error{
${mapping}
    }
)
`.trim());
