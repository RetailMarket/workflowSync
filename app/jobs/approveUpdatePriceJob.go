package jobs

import (
	"log"
	"Retail/workflowSync/clients"
	workflow "github.com/RetailMarket/workFlowClient"
	priceManager "github.com/RetailMarket/priceManagerClient"
	"golang.org/x/net/context"
)

func ApproveUpdatePriceJob() {
	log.Println("Fetching records with pending approval...")

	workflowResponse, err := clients.WorkflowClient.PendingRecords(context.Background(), &workflow.Request{})

	if (err != nil) {
		log.Printf("Failed while fetching price update records\nError: %v", err)
		return
	}
	entries := workflowResponse.GetEntries()
	log.Printf("Processing records : %v\n", entries)

	notifyServices(entries);
}

func notifyServices(records []*workflow.Entry) {
	if (len(records) != 0) {
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

func createNotifyRequestForPriceManagerService(records []*workflow.Entry) *priceManager.Records {
	request := &priceManager.Records{}
	for i := 0; i < len(records); i++ {
		priceObj := priceManager.Entry{
			ProductId: records[i].GetProductId(),
			Version: records[i].GetVersion()}
		request.Entries = append(request.Entries, &priceObj)
	}
	return request
}

func notifyPriceManagerService(records []*workflow.Entry) error {
	log.Println("notifying price manager service")
	notifyRequest := createNotifyRequestForPriceManagerService(records);
	response, err := clients.PriceManagerClient.NotifyRecordsProcessed(context.Background(), notifyRequest)
	log.Printf("priceManager response: %v", response.Message)
	return err;
}

func notifyWorkflowService(records []*workflow.Entry) error {
	log.Println("Updating status to completed in workflow service")
	request := &workflow.Records{Entries:records}
	response, err := clients.WorkflowClient.NotifyRecordsProcessed(context.Background(), request)
	log.Printf("Workflow Service %s", response.Message)
	return err;
}