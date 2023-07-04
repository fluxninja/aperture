/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable no-undef */
const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')

module.exports = {
  entry: path.join(__dirname, 'src', 'index.tsx'),
  output: {
    path: path.resolve(__dirname, 'dist'),
  },
  module: {
    rules: [
      // `js` and `jsx` files are parsed using `babel`
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: ['babel-loader'],
      },
      // `ts` and `tsx` files are parsed using `ts-loader`
      {
        test: /\.(ts|tsx)$/,
        loader: 'ts-loader',
      },
      {
        test: /\.(png|jpe?g|gif|svg)$/i,
        use: [
          {
            loader: 'file-loader',
            options: {
              name: '[name].[ext]',
              outputPath: 'static',
              publicPath: 'static',
            },
          },
        ],
      },
    ],
  },
  devServer: {
    port: 3010,
    hot: true,
    open: true,
    proxy: {
      '/aperture': {
        target: 'http://localhost:3010',
        router: () => 'http://service1-demo-app.demoapp.svc.cluster.local',
        logLevel: 'debug',
        pathRewrite: { '^/aperture': '' },
      },
    },
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx'],
  },
  plugins: [
    new CopyWebpackPlugin({
      patterns: [{ from: './public/favicon.ico', to: 'favicon.ico' }],
    }),
    new HtmlWebpackPlugin({
      template: path.join(__dirname, 'public', 'index.html'),
    }),
  ],
}
