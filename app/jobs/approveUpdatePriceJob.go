package jobs

import (
	"time"
	"log"
	"Retail/workflowSync/clients"
	workflow "github.com/RetailMarket/workFlowClient"
	priceManager "github.com/RetailMarket/priceManagerClient"
	"golang.org/x/net/context"
	"fmt"
)

func ApproveUpdatePriceJob() {
	processApproval();
	time.Sleep(time.Second * 100000)
}

func processApproval() {
	for range time.Tick(time.Second * 5) {
		log.Println("Fetching records with pending approval...")

		workflowResponse, err := clients.WorkflowClient.GetRecordsPendingForApproval(context.Background(), &workflow.GetProductsRequest{})

		if (err != nil) {
			log.Printf("Failed while fetching price update records\nError: %v", err)
			continue
		}
		log.Printf("Processing records : %v\n", workflowResponse.GetProducts())

		records := workflowResponse.GetProducts()
		notifyServices(records);
	}
}

func notifyServices(records []*workflow.Product) {
	if (len(records) != 0) {
		err := notifyPriceManagerService(records);
		if (err != nil) {
			log.Printf("Unable to change status to confirmed for entries in priceManager service %v\n Error: %v", records, err)
		} else {
			err := updateStatusInWorkflow(records)
			if (err != nil) {
				log.Printf("Unable to change status to comfirmed for entries in workflow service %v\n Error: %v", records, err)
			}
		}
	}
}

func createNotifyRequestForPriceManagerService(records []*workflow.Product) *priceManager.NotifyRequest {
	request := &priceManager.NotifyRequest{}
	for i := 0; i < len(records); i++ {
		priceObj := priceManager.Entry{
			ProductId: records[i].GetProductId(),
			Version: records[i].GetVersion()}
		request.Entries = append(request.Entries, &priceObj)
	}
	return request
}

func notifyPriceManagerService(records []*workflow.Product) error {
	log.Println("notifying price manager service")
	notifyRequest := createNotifyRequestForPriceManagerService(records);
	response, err := clients.PriceManagerClient.NotifySuccessfullyProcessed(context.Background(), notifyRequest)
	fmt.Print(response)
	//log.Printf("Price Service %s", response.Message)
	return err;
}

func updateStatusInWorkflow(records []*workflow.Product) error {
	log.Println("Updating status to completed in workflow service")
	request := &workflow.ProductsRequest{Products:records}
	response, err := clients.WorkflowClient.UpdateStatusToCompleted(context.Background(), request)
	log.Printf("Price Service %s", response.Message)
	return err;
}