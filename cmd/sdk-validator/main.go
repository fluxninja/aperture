package main

// var (
// 	logger   *log.Logger
// 	accepted int32
// )

// func init() {
// 	logger = log.NewLogger(log.GetPrettyConsoleWriter(), "debug")
// 	log.SetGlobalLogger(logger)

// 	accepted = 0
// }

// func main() {
// 	// setup flagset and flags
// 	fs := flag.NewFlagSet("sdk-validator", flag.ExitOnError)
// 	port := fs.String("port", "8080", "Port to start sdk-validator's grpc server on. Default is 8080.")
// 	accepts := fs.Int("accepts", 0, "Number of Check calls to accept. Default is 0 (accept all).")
// 	// parse flags
// 	err := fs.Parse(os.Args[1:])
// 	if err != nil {
// 		log.Error().Err(err).Msg("failed to parse flags")
// 		os.Exit(1)
// 	}

// 	// create listener for grpc server
// 	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	// instantiate flowcontrol
// 	f := &sdkvalidator.flowcontrolHandler{}

// 	// setup grpc server and register FlowControlServiceServer instance to it
// 	grpcServer := grpc.NewServer()
// 	reflection.Register(grpcServer)
// 	flowcontrolv1.RegisterFlowControlServiceServer(grpcServer, f)

// 	// start serving traffic on grpc server
// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }
