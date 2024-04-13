const path = require('path');
const debugConfig = require('./webpack.config.debug');

const WebpackObfuscator = require('webpack-obfuscator');

module.exports = debugConfig.map(item => {
    item.mode = 'production';
    item.output.path = path.resolve(__dirname, 'build', 'release');
    if (!item.plugins) {
      item.plugins = [];
    }
    item.plugins.push(new WebpackObfuscator ({ rotateStringArray: true }));
    return item;
});
