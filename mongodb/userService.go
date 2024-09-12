package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/domain"
)

type MongoUserService struct {
	Collection *mongo.Collection
}

// CreateUser creates a new user in the MongoDB collection
func (m *MongoUserService) CreateUser(user *domain.User) error {
	_, err := m.Collection.InsertOne(context.TODO(), user)
	return err
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserName  string             `bson:"userName,required" json:"userName"`
	Password  string             `bson:"password,required" json:"password"`
	Role      primitive.ObjectID `bson:"role,omitempty" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// GetUser retrieves a user by ID from MongoDB
func (m *MongoUserService) GetUser(userFilter domain.UserFilter) (*domain.User, error) {
	filter := bson.D{}
	var err error
	if userFilter.ID != "" {
		objectID, err := primitive.ObjectIDFromHex(userFilter.ID)
		if err != nil {
			return nil, err
		}
		filter = append(filter, bson.E{Key: "_id", Value: objectID})

	}
	if userFilter.UserName != "" {
		filter = append(filter, bson.E{Key: "userName", Value: userFilter.UserName})
	}
	var user User
	err = m.Collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	apiUser, err := user.ToApi()

	if err != nil {
		return nil, err
	}
	return &apiUser, nil
}

// GetPaginatedUsers récupère les utilisateurs avec pagination.
func (s *MongoUserService) GetPaginatedUsers(ctx context.Context, params domain.PaginationOptions) (domain.PaginateResult[domain.User], error) {
	return Paginate[domain.User, *User](s.Collection, params)
}

func (u *User) GetRole() (Role, error) {
	var role Role
	coll := config.GetCollections().RoleModel

	filter := bson.D{{Key: "_id", Value: u.Role}}
	err := coll.FindOne(context.TODO(), filter).Decode(&role)
	if err != nil {
		return Role{}, err
	}
	return role, nil
}

func (u *User) ToApi() (domain.User, error) {
	role, err := u.GetRole()
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:        u.ID.Hex(),
		UserName:  u.UserName,
		Role:      role.ToApi(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Password:  u.Password,
	}, nil
}
