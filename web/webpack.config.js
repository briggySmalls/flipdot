const HtmlWebpackPlugin = require('html-webpack-plugin')

module.exports = {
    entry: './src/client.js',
    output: {
        path: __dirname + '/dist',
        filename: 'main.js'
    },
    plugins: [
        new HtmlWebpackPlugin({
            title: 'Custom template',
            // Load a custom template (lodash by default see the FAQ for details)
            template: './src/index.html'
        })
    ]
}
