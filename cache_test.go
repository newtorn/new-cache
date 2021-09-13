package newcache

import (
	"context"
	"testing"
)

func TestCacheServiceSingleton(t *testing.T) {
	css := Singleton()
	k := "hello"
	v := "world"
	css.Set(k, v)
	cv, ok := css.Get(k)
	if !ok || v != cv {
		t.Errorf("excepted got value=%+v by key=%q, but got value=%+v\n", k, v, cv)
	}
	t.Log(v)
}

type User struct {
	Username string
	Password string
}
type userService struct {
	users []*User
	done  chan interface{}
}

func (us *userService) Done(ctx context.Context) (done <-chan interface{}) {
	return us.done
}
func (us *userService) LoadKeys(ctx context.Context, value interface{}) []string {
	return []string{
		value.(*User).Username,
	}
}
func (us *userService) LoadValues(ctx context.Context) []interface{} {
	values := make([]interface{}, len(us.users))
	for i := 0; i < len(us.users); i++ {
		values[i] = us.users[i]
	}
	return values
}
func TestRegisterFlushDaemon(t *testing.T) {
	us := &userService{
		done: make(chan interface{}),
		users: []*User{
			{"Jack", "123"},
			{"Jane", "123456"},
			{"Mark", "abd123"},
			{"Michale", "123aba"},
		},
	}
	ctx := context.Background()
	c := Singleton()
	c.Register(ctx, us)
	for i := 0; i < len(us.users); i++ {
		u, ok := c.Get(us.users[i].Username)
		if !ok {
			t.Error("could not get by key from get-started")
		}
		uv := *(u.(*User))
		uo := *us.users[i]
		if uv != uo {
			t.Error("could not get by key from get-started")
		} else {
			t.Log(uv)
		}
	}
	us.done <- "hello"
	t.Log(c.Get("Jack"))
}
