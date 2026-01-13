"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Log = exports.setConsoleMode = void 0;
var node_fs_1 = require("node:fs");
var node_path_1 = require("node:path");
var LOG_PATH = (0, node_path_1.resolve)(__dirname, './duck.log');
var consoleMode = false;
var setConsoleMode = function (flag) {
    if (flag && process.stdout._handle) {
        process.stdout._handle.setBlocking(true);
    }
    consoleMode = flag;
};
exports.setConsoleMode = setConsoleMode;
var Log = /** @class */ (function () {
    function Log() {
    }
    Log.print = function (message) {
        message = "".concat(new Date().toLocaleString(), " ").concat(message, " \n");
        consoleMode ? process.stdout.write(message) : (0, node_fs_1.appendFileSync)(LOG_PATH, message, 'utf8');
    };
    Log.debug = function (message) {
        this.print("[Debug] ".concat(message));
    };
    Log.info = function (message) {
        this.print("[Info] ".concat(message));
    };
    Log.warn = function (message) {
        if (message instanceof Error) {
            message = "".concat(message.name, ": ").concat(message.message, " \n ").concat(message.stack);
        }
        this.print("[Warn] ".concat(message));
    };
    Log.error = function (message) {
        if (message instanceof Error) {
            message = "".concat(message.name, ": ").concat(message.message, " \n ").concat(message.stack);
        }
        this.print("[Error] ".concat(message));
    };
    Log.kernelDebug = function (message) {
        this.print("[Debug] [kernel] ".concat(message));
    };
    Log.kernelInfo = function (message) {
        this.print("[Info] [kernel] ".concat(message));
    };
    Log.kernelWarn = function (message) {
        if (message instanceof Error) {
            message = "".concat(message.name, ": ").concat(message.message, " \n ").concat(message.stack);
        }
        this.print("[Warn] [kernel] ".concat(message));
    };
    Log.kernelError = function (message) {
        if (message instanceof Error) {
            message = "".concat(message.name, ": ").concat(message.message, " \n ").concat(message.stack);
        }
        this.print("[Error] [kernel] ".concat(message));
    };
    return Log;
}());
exports.Log = Log;
