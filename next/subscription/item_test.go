package subscription_test

import (
	"testing"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"

	"github.com/stretchr/testify/assert"
)

type NullObserver struct {
	ID string
}

func (o NullObserver) GetID() string {
	return o.ID
}

func (o *NullObserver) Update(item subscription.Item) {

}

func TestRegisterObserver_addsObserverToList(t *testing.T) {
	assert := assert.New(t)

	item := subscription.Item{}
	assert.Equal(0, len(item.Observers()))

	ok, err := item.RegisterObserver(&NullObserver{ID: "observer 1"})
	assert.Equal(1, len(item.Observers()))
	assert.NoError(err)
	assert.True(ok)

	ok, err = item.RegisterObserver(&NullObserver{ID: "observer 2"})
	assert.Equal(2, len(item.Observers()))
	assert.NoError(err)
	assert.True(ok)

}

func TestRegisterObserver_doesNotAddObserverWithSameIDTwice(t *testing.T) {
	assert := assert.New(t)

	item := subscription.Item{}
	assert.Equal(0, len(item.Observers()))

	ok, err := item.RegisterObserver(&NullObserver{ID: "observer 1"})
	assert.Equal(1, len(item.Observers()))
	assert.NoError(err)
	assert.True(ok)

	ok, err = item.RegisterObserver(&NullObserver{ID: "observer 1"})
	assert.NoError(err)
	assert.False(ok)
	assert.Equal(1, len(item.Observers()))
}

func TestDeregisterObserver_removesObserver(t *testing.T) {
	assert := assert.New(t)

	item := subscription.Item{}
	assert.Equal(0, len(item.Observers()))

	observer := NullObserver{ID: "observer 1"}
	ok, err := item.RegisterObserver(&observer)
	assert.Equal(1, len(item.Observers()))
	assert.NoError(err)
	assert.True(ok)

	ok, err = item.DeregisterObserver(&observer)
	assert.NoError(err)
	assert.True(ok)
	assert.Equal(0, len(item.Observers()))
}

func TestDeregisterObserver_doesNotRemovesObserverTwice(t *testing.T) {
	assert := assert.New(t)

	item := subscription.Item{}
	assert.Equal(0, len(item.Observers()))

	observer := NullObserver{ID: "observer 1"}
	ok, err := item.RegisterObserver(&observer)
	assert.Equal(1, len(item.Observers()))
	assert.NoError(err)
	assert.True(ok)

	observer2 := NullObserver{ID: "observer 2"}
	ok, err = item.RegisterObserver(&observer2)
	assert.Equal(2, len(item.Observers()))
	assert.NoError(err)
	assert.True(ok)

	ok, err = item.DeregisterObserver(&observer)
	assert.NoError(err)
	assert.True(ok)
	assert.Equal(1, len(item.Observers()))

	ok, err = item.DeregisterObserver(&observer)
	assert.NoError(err)
	assert.False(ok)
	assert.Equal(1, len(item.Observers()))
}
