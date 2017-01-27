package jobs

import (
	"time"
	"log"
	"Retail/workflowSync/clients"
	workflow "github.com/RetailMarket/workFlowClient"
	priceManager "github.com/RetailMarket/priceManagerClient"
	"golang.org/x/net/context"
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
		changeStatus(records);
	}
}

func changeStatus(records []*workflow.Product) {
	if (len(records) != 0) {
		err := updateStatusInPriceManager(records);
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

func createStatusUpdateRequestForPriceManagerService(records []*workflow.Product) *priceManager.ChangeStatusRequest {
	request := &priceManager.ChangeStatusRequest{}
	for i := 0; i < len(records); i++ {
		priceObj := priceManager.ProductEntry{
			ProductId: records[i].GetProductId(),
			Version: records[i].GetVersion()}
		request.Products = append(request.Products, &priceObj)
	}
	return request
}

func updateStatusInPriceManager(records []*workflow.Product) error {
	log.Println("Updating status to completed in price manager service")
	priceRequest := createStatusUpdateRequestForPriceManagerService(records);
	_, err := clients.PriceManagerClient.ChangeStatusToCompleted(context.Background(), priceRequest)
	return err;
}

func updateStatusInWorkflow(records []*workflow.Product) error {
	log.Println("Updating status to completed in workflow service")
	request := &workflow.ProductsRequest{Products:records}
	_, err := clients.WorkflowClient.UpdateStatusToCompleted(context.Background(), request)
	return err;
}