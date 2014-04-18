var fs = require('fs'),
  filepath = process.argv[2],
  geojson = process.argv[3] ? true : false;

var exsist = fs.existsSync( __dirname + filepath );

if( !exsist ) throw('No file found @ ' + filepath);

var file = fs.readFileSync( __dirname + filepath );

if( !file ) throw('No content found @ ' + filepath);

try {
  file = JSON.parse( file.toString('utf8') );
} catch ( e ) {
  throw( e.message );
}


if ( geojson && Array.isArray( file ) ) {
  // this function is specific to the driver.json data
  function featurify ( raw ) {
    var feature = {
      type : "Feature",
      properties : {
        id : raw.id,
      },
      geometry : {
        type : 'LineString',
        coordinates : []
      }
    }

    
      raw.history.forEach(function( location ){
        console.log( location );
        feature.geometry.coordinates.push([ location.lng, location.lat ]);
      })
    

    return feature;
  }

  file = {
    type : "FeatureCollection",
    features : file.map(featurify)
  };

} 

fs.writeFileSync( __dirname + filepath + ( geojson ? '.geojson' : '' ), JSON.stringify(file, null, '\t'));
console.log('success file' + filepath + ' formatted');