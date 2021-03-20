const path = require('path'); 
const Habitat = require('../../src-webpack/plugin');

module.exports = { 

    // entry: 'src/index.js',

    // output: { 
    //     path: path.resolve(__dirname, '.habitat/'),
    // },

    plugins: [ new Habitat() ],

    mode: 'development',

}