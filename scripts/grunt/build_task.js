
module.exports = function(grunt) {
  "use strict";

  // Concat and Minify the public directory into dist
  grunt.registerTask('build', [
    'clean:release',
    'clean:build',
    'phantomjs',
    'exec:webpack',
  ]);

};
