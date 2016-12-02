package dondeestas

import (
	"fmt"
	"math/rand"
	"time"
)

func createRandomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 8)
	r.Read(b)

	return fmt.Sprintf("%x", b)
}

func createRandomStringSlice() []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	l := make([]string, r.Intn(20))
	for i, _ := range l {
		l[i] = createRandomString()
	}
	return l
}

func createRandomPerson() (*Person, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	p := new(Person)
	p.Id = createRandomString()
	p.Name = createRandomString()
	p.Position.Tov = time.Now()
	p.Position.Latitude = r.Float32()
	p.Position.Longitude = r.Float32()
	p.Position.Elevation = r.Float32()
	p.Visible = r.Int()%2 != 0
	p.Whitelist = createRandomStringSlice()
	p.Following = createRandomStringSlice()

	return p, nil
}

func arePersonEqual(p1, p2 *Person) bool {
	if p1.Id != p2.Id {
		fmt.Println("Ids don't match")
		return false
	}
	if p1.Name != p2.Name {
		fmt.Println("Names don't match")
		return false
	}
	if p1.Position.Tov != p2.Position.Tov {
		fmt.Println("TOVs don't match")
		return false
	}
	if p1.Position.Latitude != p2.Position.Latitude {
		fmt.Println("Lats don't match")
		return false
	}
	if p1.Position.Longitude != p2.Position.Longitude {
		fmt.Println("Lons don't match")
		return false
	}
	if p1.Position.Elevation != p2.Position.Elevation {
		fmt.Println("Elevations don't match")
		return false
	}
	if p1.Visible != p2.Visible {
		fmt.Println("Visibilities don't match")
		return false
	}
	/* TODO
	if p1.Whitelist != p2.Whitelist {
		return false
	}
	if p1.Following != p2.Following {
		return false
	}
	*/
	return true
}
