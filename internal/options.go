package internal

// Options contain settings from command line flags and env
type Options struct {
	LogLevel        string   `long:"log-level" env:"LOG_LEVEL" default:"info" description:"Logging level"`
	UploadLocation  string   `long:"location" env:"UPLOAD_LOCATION" required:"true" description:"Locations where to store media" default:"uploads"`
	Listen          string   `long:"listen" env:"LISTEN_HOST" default:"0.0.0.0" description:"where to listen for connections"`
	Port            int      `long:"port" env:"LISTEN_PORT" default:"8080" description:"Port"`
}
