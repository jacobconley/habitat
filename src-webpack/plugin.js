const { mkdir, mkdirSync } = require('fs');
const path = require('path'); 
const { Compilation } = require('webpack');

class Habitat { 

    apply(compiler) { 

        compiler.hooks.beforeRun.tap('Habitat', (compiler) => { 

            compiler.outputPath = '/home/jake/go/src/github.com/jacobconley/habitat/test-fixtures/userland/.habitat/out/packs'
            
            console.log(compiler.options.output.path)
            console.log(compiler.outputPath)


        })        


    }

}

module.exports = Habitat; 