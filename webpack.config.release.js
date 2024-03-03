const path = require('path');
const debugConfig = require('./webpack.config.debug');

module.exports = debugConfig.map(item => {
    item.mode = 'production',
    item.output.path = path.resolve(__dirname, 'build', 'release')
    return item;
});