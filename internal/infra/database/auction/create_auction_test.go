package auction_test

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSetExpiration(t *testing.T) {
	os.Setenv("AUCTION_DURATION", "2s") // Define duração para o teste

	mongot := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mongot.Run("Set expiration and update status", func(mt *mtest.T) {
		// Simula o contexto e a coleção
		ctx := context.Background()
		repo := auction.AuctionRepository{Collection: mt.Coll}

		auctionEntity := &auction_entity.Auction{
			Id:          "test_id",
			ProductName: "Test Product",
			Category:    "Test Category",
			Description: "Test Description",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		// Configura o comportamento esperado do MongoDB
		mt.AddMockResponses(mtest.CreateSuccessResponse()) // Para UpdateOne

		// Chama o método SetExpiration
		go repo.SetExpiration(ctx, auctionEntity)

		// Aguarda tempo suficiente para a execução do timer
		time.Sleep(3 * time.Second)

		// Verifica se a coleção recebeu a operação correta
		mt.Coll.FindOne(ctx, bson.D{{Key: "_id", Value: auctionEntity.Id}})
	})
}
