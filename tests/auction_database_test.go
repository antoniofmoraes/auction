package auction

import (
	"context"
	"errors"
	"fullcycle-auction_go/internal/infra/database/auction"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func mockEnv() {
	_ = os.Setenv("AUCTION_DURATION", "2s")
}

func TestSetExpiration(t *testing.T) {
	mockEnv()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	defer mt.Run("successful expiration update", func(mt *mtest.T) {
		ctx := context.Background()

		repo := auction.NewAuctionRepository(mt.DB)

		auctionId := "auction123"
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(),
			mtest.CreateSuccessResponse(),
		)

		repo.SetExpiration(ctx, auctionId)

		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		time.Sleep(3 * time.Second)

		expectedFilter := bson.D{{Key: "_id", Value: auctionId}}

		err := mt.Coll.FindOne(ctx, expectedFilter).Decode(nil)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			t.Errorf("SetExpiration did not update auction correctly: %v", err)
		}
	})

	mt.Run("failure due to incorrect duration format", func(mt *mtest.T) {
		_ = os.Setenv("AUCTION_DURATION", "invalid_duration")

		ctx := context.Background()
		repo := auction.NewAuctionRepository(mt.DB)

		auctionId := "auction456"

		mt.ClearMockResponses()
		repo.SetExpiration(ctx, auctionId)

		time.Sleep(1 * time.Second)
	})
}
