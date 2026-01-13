const path = require('path');
const nodeExternals = require('webpack-node-externals');

module.exports = [{
    target: 'electron-renderer',
    mode: 'development',
    entry: {
        app: './src/app/main.tsx',
    },
    output: {
        path: path.resolve(__dirname, 'build', 'debug'),
        filename: '[name].js',
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: {
                    loader: 'ts-loader',
                    options: {
                        configFile: "tsconfig.json"
                    },
                },
                exclude: /node_modules/
            },
            {
                test: /\.m?js$/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env', '@babel/preset-react']
                    }
                },
                exclude: /node_modules/
            },
            {
                test: /\.s?css$/,
                use: [{ loader: 'style-loader' }, { loader: 'css-loader' }],
                exclude: /node_modules/
            },
            {
                test: /\.(png|jp(e*)g|svg|gif)$/,
                use: ['file-loader'],
                exclude: /node_modules/
            }
        ]
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.jsx', '.js'],
        alias: {
          '@': path.resolve(__dirname, 'src/app')
        }
    },
    devtool: 'source-map'
}, {
    target: 'electron-main',
    mode: 'development',
    entry: {
        electron: './src/main.ts',
    },
    output: {
        path: path.resolve(__dirname, 'build', 'debug'),
        filename: '[name].js',
    },
    externals: [nodeExternals()],
    module: {
        rules: [
            {
                test: /\.ts$/,
                use: {
                    loader: 'ts-loader',
                    options: {
                        configFile: "tsconfig.json"
                    },
                },
                exclude: /node_modules/
            },
            {
                test: /\.m?js$/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env']
                    }
                },
                exclude: /node_modules/
            }
        ]
    },
    resolve: {
        extensions: ['.ts', '.js']
    },
    devtool: 'source-map'
}];
