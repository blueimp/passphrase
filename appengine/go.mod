module github.com/blueimp/passphrase/appengine

require github.com/blueimp/passphrase v1.0.0

replace github.com/blueimp/passphrase v1.0.0 => ../

// Use alternative replace pattern to deploy to App Engine:
//replace github.com/blueimp/passphrase v1.0.0 => ./passphrase
