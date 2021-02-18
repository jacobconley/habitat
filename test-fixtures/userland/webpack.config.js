const path = require('path'); 

module.exports = { 

    entry: [
        './src/index.js',   
        './src/index.css'
    ],

    output: { 
        path: path.resolve(__dirname, '.habitat/'),
    },

    mode: 'development',

}