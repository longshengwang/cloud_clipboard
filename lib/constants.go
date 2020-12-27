package lib

import "flag"

var DefaultPassword = "cloud_clipboard_password"
var ServerHelloWord = "Hey, brother, you are at the home of clipboard server."
var ClientHelloWord = "Hello, is my clipboard?"
var SplitKey = "]"
var DiscoverServicePort = 9266
var ServicePort = 5166

var ClientHelloFlag = flag.String("cw", ClientHelloWord, "Client Hello Word")
var ServerHelloFlag = flag.String("sw", ServerHelloWord, "Server Hello Word")
var ServerAuthFlag = flag.String("auth", DefaultPassword, "Server Auth Word")
var DiscoveryServiceFlag = flag.Int("discoveryPort", DiscoverServicePort, "Discovery Service Port")
var ServerPortFlag = flag.Int("port", ServicePort, "Server Port")

