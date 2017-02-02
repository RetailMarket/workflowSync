package jobs

import (
	"log"
	"Retail/workflowSync/clients"
	workflow "github.com/RetailMarket/workFlowClient"
	priceManager "github.com/RetailMarket/priceManagerClient"
	"golang.org/x/net/context"
	"encoding/json"
)

func ApproveUpdatePriceJob() {
	log.Println("Fetching records with pending approval...")

	records, err := clients.WorkflowClient.PendingRecords(context.Background(), &workflow.Request{})

	if (err != nil) {
		log.Printf("Failed while fetching price update records\nError: %v", err)
		return
	}
	entries := records.GetEntries()
	log.Printf("Processing records : %v\n", entries)

	notifyServices(records);
}

func notifyServices(records *workflow.Records) {
	if (len(records.GetEntries()) != 0) {
		err := notifyPriceManagerService(records);
		if (err != nil) {
			log.Printf("Unable to change status to confirmed for entries in priceManager service %v\n Error: %v", records, err)
		} else {
			err := notifyWorkflowService(records)
			if (err != nil) {
				log.Printf("Unable to change status to comfirmed for entries in workflow service %v\n Error: %v", records, err)
			}
		}
	}
}

func notifyPriceManagerService(records *workflow.Records) error {
	log.Println("notifying price manager service")
	recordsInBytes, err := json.Marshal(records)
	if (err != nil) {
		log.Printf("Unable to marshal records %v", records.Entries)
	}
	notifyRequest := &priceManager.Records{}
	json.Unmarshal(recordsInBytes, notifyRequest)
	response, err := clients.PriceManagerClient.NotifyRecordsProcessed(context.Background(), notifyRequest)
	log.Printf("priceManager response: %v", response)
	return err;
}

func notifyWorkflowService(records *workflow.Records) error {
	log.Println("Updating status to completed in workflow service")
	response, err := clients.WorkflowClient.NotifyRecordsProcessed(context.Background(), records)
	log.Printf("Workflow Service %s", response)
	return err;
}