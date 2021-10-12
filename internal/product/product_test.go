package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/egorovdmi/garagesale/internal/platform/tests"
	"github.com/egorovdmi/garagesale/internal/product"
	"github.com/egorovdmi/garagesale/internal/schema"
	"github.com/google/go-cmp/cmp"
)

func TestProducts(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	ctx := context.Background()
	newP := product.NewProduct{
		Name:     "Comic Book",
		Cost:     10,
		Quantity: 55,
	}
	now := time.Date(2021, time.August, 1, 0, 0, 0, 0, time.UTC)

	p0, err := product.Create(ctx, db, newP, now)
	if err != nil {
		t.Fatalf("creating product p0: %s", err)
	}

	p1, err := product.Single(ctx, db, p0.ID)
	if err != nil {
		t.Fatalf("getting product p0: %s", err)
	}

	if diff := cmp.Diff(p0, p1); diff != "" {
		t.Fatalf("fetched != created:\n %s", diff)
	}
}

func TestList(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	pl, err := product.List(ctx, db)
	if err != nil {
		t.Fatalf("listing products: %s", err)
	}

	if exp, got := 2, len(pl); got != exp {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
