package faker

import (
	"math/rand"
	"src/user_proto"
)

// Helper functions to generate fake data:

// generateRandomName generates a random name.
func generateRandomName() string {

	names := []string{"John", "Jane", "David", "Mary", "Peter", "Susan", "Michael", "Emily", "William", "Olivia"}
	return names[rand.Intn(len(names))]
}

// generateRandomCity generates a random city name.
func generateRandomCity() string {
	cities := []string{"New York", "Los Angeles", "Chicago", "London", "Paris", "Tokyo", "Sydney", "Rome", "Berlin", "Madrid"}
	return cities[rand.Intn(len(cities))]
}

// generateRandomPhoneNumber generates a random phone number.
func generateRandomPhoneNumber() int64 {
	return int64(rand.Intn(9000000000) + 1000000000) // 10-digit phone number
}

// generateRandomHeight generates a random height between 4.5 and 6.5 feet.
func generateRandomHeight() float32 {
	height := float32(rand.Float64()*2.0 + 4.5)
	return float32(int(height*10)) / 10 // Round to 1 decimal place
}

// generateRandomMaritalStatus generates a random marital status (true/false).
func generateRandomMaritalStatus() bool {
	return rand.Intn(2) == 1 // 50/50 chance of being married
}

// GenerateFakeUsers generates a slice of fake user data.
func GenerateFakeUsers(numUsers int) []*user_proto.User {
	users := make([]*user_proto.User, numUsers)
	for i := 0; i < numUsers; i++ {
		users[i] = &user_proto.User{
			Id:      int32(i + 1), // Assign sequential IDs
			Fname:   generateRandomName(),
			City:    generateRandomCity(),
			Phone:   generateRandomPhoneNumber(),
			Height:  generateRandomHeight(),
			Married: generateRandomMaritalStatus(),
		}
	}
	return users
}

func GetSeed() []*user_proto.User {
	return []*user_proto.User{
		{Id: 1, Fname: "Jane", City: "Tokyo", Phone: 1764788495, Height: 4.9, Married: true},
		{Id: 2, Fname: "William", City: "Chicago", Phone: 3147005163, Height: 6.0, Married: true},
		{Id: 3, Fname: "Jane", City: "Sydney", Phone: 5660612778, Height: 5.2, Married: false},
		{Id: 4, Fname: "Olivia", City: "Paris", Phone: 7499978875, Height: 5.0, Married: false},
		{Id: 5, Fname: "Olivia", City: "New York", Phone: 4934669609, Height: 5.8, Married: true},
		{Id: 6, Fname: "Olivia", City: "Berlin", Phone: 3935422070, Height: 6.2, Married: false},
		{Id: 7, Fname: "Peter", City: "Chicago", Phone: 1389433165, Height: 5.7, Married: false},
		{Id: 8, Fname: "Susan", City: "Paris", Phone: 5307397290, Height: 6.1, Married: false},
		{Id: 9, Fname: "Olivia", City: "Madrid", Phone: 5199527895, Height: 4.6, Married: false},
		{Id: 10, Fname: "Olivia", City: "Madrid", Phone: 8354340417, Height: 5.5, Married: true},
	}
}
