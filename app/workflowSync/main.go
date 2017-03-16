package main

import (
	clients "Retail/workflowSync/clients"
	"Retail/workflowSync/jobs"
	"github.com/jasonlvhit/gocron"
)

func main() {
	clients.CreateClientConnection()
	defer clients.CloseConnections();

	// running job for sending update price record for approval.
	gocron.Every(5).Seconds().Do(jobs.ApproveUpdatePriceJob)
	<-gocron.Start()
	defer gocron.Clear()

}
