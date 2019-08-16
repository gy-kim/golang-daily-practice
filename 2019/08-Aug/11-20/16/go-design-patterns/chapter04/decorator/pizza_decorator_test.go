package decorator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPizzaDecorator_AddIngredient(t *testing.T) {
	pizza := &PizzaDecorator{}
	pizzaResult, err := pizza.AddIngredient()
	if err != nil {
		t.Error(err)
	}
	expectedText := "Pizza with the following ingredients:"
	if !strings.Contains(pizzaResult, expectedText) {
		t.Errorf("When calling the add ingredient of the pizza decorator it "+
			"must return the text %s, not '%s'", expectedText, pizzaResult)
	}
}

func TestOnion_AddIngredient(t *testing.T) {
	onion := &Onion{}
	onionResult, err := onion.AddIngredient()
	assert.Error(t, err)

	onion = &Onion{&PizzaDecorator{}}
	onionResult, err = onion.AddIngredient()
	assert.NoError(t, err)
	assert.Contains(t, onionResult, "onion")
}

func TestMeat_AddIngredient(t *testing.T) {
	meat := &Meat{}
	meatResult, err := meat.AddIngredient()
	assert.Error(t, err)

	meat = &Meat{&PizzaDecorator{}}
	meatResult, err = meat.AddIngredient()
	assert.NoError(t, err)
	assert.Contains(t, meatResult, "meat")
}

func TestPizzaDecorator_FullStack(t *testing.T) {
	pizza := &Onion{&Meat{&PizzaDecorator{}}}
	pizzaResult, err := pizza.AddIngredient()
	assert.NoError(t, err)

	expectedText := "Pizza with the following ingredients: meat, onion"
	assert.Contains(t, pizzaResult, expectedText)

	t.Log(pizzaResult)
}
