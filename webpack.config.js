const path = require('path');

const tsLoaderConfig = {
    test: /\.ts$/i,
    use: 'ts-loader',
    exclude: '/node_modules/'
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

module.exports = [
    {
        entry: './static_src/telepath.ts',
        output: {
            'path': path.resolve(__dirname, 'telepath/static/'),
            'filename': 'telepath.js',
            'library': {
                name: 'Telepath',
                type: 'umd',
                export: 'default',
            },
            'globalObject': 'this',
        },
         ...baseConfig(),
    }
]