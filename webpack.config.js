import path from 'path';

const tsLoaderConfig = {
    test: /\.ts$/i,
    use: {
        loader: 'ts-loader',
        options: {
            configFile: 'tsconfig.webpack.json'
        }
    },
    exclude: /node_modules/,
}

function baseConfig(rules = []) {
    return {
        resolve: {
            extensions: ['.ts', '...'],
        },
        mode: 'production',
        module: {
            rules: [
                tsLoaderConfig,
                ...rules
            ]
        }
    }
}

export default [
    // UMD build for browsers
    {
        entry: './static_src/telepath.ts',
        output: {
            path: path.resolve('telepath/static/telepath/'),
            filename: 'telepath.umd.js',
            library: {
                name: 'Telepath',
                type: 'umd',
                export: 'default',
            },
            globalObject: 'this',
        },
        ...baseConfig(),
    },
    // ESM build for Node
    {
        entry: './static_src/telepath.ts',
        output: {
            path: path.resolve('telepath/static/telepath/'),
            filename: 'telepath.js',
            library: {
                type: 'module',
            },
        },
        experiments: { outputModule: true },
        ...baseConfig(),
    },
];