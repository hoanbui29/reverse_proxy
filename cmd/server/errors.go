package main

func (s *serverGateway) handleFatal(err error, msg string) {
	if err != nil {
		s.logger.Fatal(msg, map[string]string{"error": err.Error()})
	}
}
