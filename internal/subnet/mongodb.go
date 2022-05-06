package subnet

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubnetRepository struct {
	collection *mongo.Collection
}

func (r *SubnetRepository) Create(ctx context.Context, subnet *Subnet) (string, error) {
	result, err := r.collection.InsertOne(ctx, subnet)
	if err != nil {
		return "", err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

func (r *SubnetRepository) FindAll(ctx context.Context) ([]Subnet, error) {
	var d []Subnet

	cur, err := r.collection.Find(ctx, bson.D{{}})
	if err != nil {
		return d, err
	}

	if err = cur.All(ctx, &d); err != nil {
		return d, err
	}

	err = cur.Close(ctx)
	if err != nil {
		return d, err
	}

	if len(d) == 0 {
		d = make([]Subnet, 0)
	}

	return d, nil
}

func (r *SubnetRepository) FindOne(ctx context.Context, id string) (d Subnet, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return d, err
	}
	query := bson.M{"_id": oid}
	result := r.collection.FindOne(ctx, query)
	if result.Err() != nil {
		return d, result.Err()
	}
	if err = result.Decode(&d); err != nil {
		return d, err
	}
	return d, nil
}

func (r *SubnetRepository) Update(ctx context.Context, subnet *Subnet) error {
	oid, err := primitive.ObjectIDFromHex(subnet.Id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	deviceBytes, _ := bson.Marshal(subnet)
	var updateDevice bson.M
	err = bson.Unmarshal(deviceBytes, &updateDevice)
	if err != nil {
		return err
	}
	delete(updateDevice, "_id")

	update := bson.M{
		"$set": updateDevice,
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *SubnetRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func NewSubnetRepository(collection *mongo.Collection) *SubnetRepository {
	return &SubnetRepository{
		collection: collection,
	}
}
