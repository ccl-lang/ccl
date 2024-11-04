package generator_test

import (
	"testing"
	"time"

	"github.com/ALiwoto/ccl/examples/go_gen/apiTypes_gen"
)

func TestGoGenerator1(t *testing.T) {
	currentTime := time.Now()
	var usersResult = &apiTypes_gen.GetUsersResult{
		Users: []*apiTypes_gen.UserInfo{
			{
				Id:        1,
				Username:  "user1",
				Email:     "email1@gmail.com",
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			{
				Id:        2,
				Username:  "user2",
				Email:     "email2@gmail.com",
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
		},
	}

	serializedData, err := usersResult.SerializeBinary()
	if err != nil {
		t.Fatalf("Failed to serialize data: %v", err)
	}

	var deserializedResult = &apiTypes_gen.GetUsersResult{}
	if err := deserializedResult.DeserializeBinary(serializedData); err != nil {
		t.Fatalf("Failed to deserialize data: %v", err)
	}

	// now compare the results
	if len(usersResult.Users) != len(deserializedResult.Users) {
		t.Fatalf("Expected %d users, but got %d", len(usersResult.Users), len(deserializedResult.Users))
	}

	for i := 0; i < len(usersResult.Users); i++ {
		expected := usersResult.Users[i]
		actual := deserializedResult.Users[i]

		if expected.Id != actual.Id {
			t.Fatalf("Expected ID %d, but got %d", expected.Id, actual.Id)
		}
		if expected.Username != actual.Username {
			t.Fatalf("Expected username %s, but got %s", expected.Username, actual.Username)
		}
		if expected.Email != actual.Email {
			t.Fatalf("Expected email %s, but got %s", expected.Email, actual.Email)
		}
		if !expected.CreatedAt.Equal(actual.CreatedAt) {
			t.Fatalf("Expected created at %v, but got %v", expected.CreatedAt, actual.CreatedAt)
		}
		if !expected.UpdatedAt.Equal(actual.UpdatedAt) {
			t.Fatalf("Expected updated at %v, but got %v", expected.UpdatedAt, actual.UpdatedAt)
		}
	}
}
