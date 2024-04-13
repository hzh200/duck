import { appendFileSync } from 'node:fs';
import { resolve } from 'node:path';

const LOG_PATH = resolve(__dirname, './duck.log');

let consoleMode = false;

const setConsoleMode = (flag: boolean) => {
    if (flag && (<any>process.stdout)._handle) {
        (<any>process.stdout)._handle.setBlocking(true);
    }
    consoleMode = flag;
};

class Log {
    static print(message: string) {
        message = `${new Date().toLocaleString()} ${message} \n`;
        consoleMode ? process.stdout.write(message) : appendFileSync(LOG_PATH, message, 'utf8');
    }

    static debug(message: string) {
        this.print(`[Debug] ${message}`);
    }

    static info(message: string) {
        this.print(`[Info] ${message}`);
    }

    static warn(message: string | Error) {
        if (message instanceof Error) {
            message = `${message.name}: ${message.message} \n ${message.stack}`;
        }
        this.print(`[Warn] ${message}`);
    }

    static error(message: string | Error) {
        if (message instanceof Error) {
            message = `${message.name}: ${message.message} \n ${message.stack}`;
        }
        this.print(`[Error] ${message}`);
    }

    static kernelDebug(message: string) {
        this.print(`[Debug] [kernel] ${message}`);
    }

    static kernelInfo(message: string) {
        this.print(`[Info] [kernel] ${message}`);
    }

    static kernelWarn(message: string | Error) {
        if (message instanceof Error) {
            message = `${message.name}: ${message.message} \n ${message.stack}`;
        }
        this.print(`[Warn] [kernel] ${message}`);
    }

    static kernelError(message: string | Error) {
        if (message instanceof Error) {
            message = `${message.name}: ${message.message} \n ${message.stack}`;
        }
        this.print(`[Error] [kernel] ${message}`);
    }
}

export { setConsoleMode, Log };
