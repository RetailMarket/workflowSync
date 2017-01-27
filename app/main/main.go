package main

import (
	clients "Retail/workflowSync/clients"
	"Retail/workflowSync/app/jobs"
)

func main() {
	clients.CreateClientConnection()
	defer clients.CloseConnections();

	// running job for sending update price record for approval.
	jobs.ApproveUpdatePriceJob();
}
