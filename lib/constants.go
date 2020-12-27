package lib

import "flag"

const DefaultPassword = "cloud_clipboard_password"
const ServerHelloWord = "Hey, brother, you are at the home of clipboard server."
const ClientHelloWord = "Hello, is my clipboard?"
const SplitKey = "]"
const DiscoverServicePort = 9266
const ServicePort = 5166

var ClientHelloFlag = flag.String("cw", ClientHelloWord, "Client Hello Word")
var ServerHelloFlag = flag.String("sw", ServerHelloWord, "Server Hello Word")
var ServerAuthFlag = flag.String("auth", DefaultPassword, "Server Auth Password. Cannot more than 32 char(256 bit)")
var DiscoveryServiceFlag = flag.Int("discoveryPort", DiscoverServicePort, "Discovery Service Port")
var ServerPortFlag = flag.Int("port", ServicePort, "Server Port")
