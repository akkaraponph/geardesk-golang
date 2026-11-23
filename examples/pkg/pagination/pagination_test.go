package pagination

import (
	"net/url"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testTask struct {
	ID    uint
	Title string
}

func TestPaginateFiltersAndLinks(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&testTask{}); err != nil {
		t.Fatal(err)
	}
	for _, title := range []string{"alpha", "alphabet", "beta"} {
		if err := db.Create(&testTask{Title: title}).Error; err != nil {
			t.Fatal(err)
		}
	}
	p := NewPaginationParams[struct{}]("http://localhost/?title=alpha", "id", "asc", "1", "1", "")
	query := ApplyFilter(db.Model(&testTask{}), "title", "alpha", "contains").Order(NewOrderBy(SortParams{Sort: p.Sort, Order: p.Order}))
	result, err := Paginate[struct{}, []testTask](p, query)
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 2 || result.TotalPages != 2 || len(result.Rows) != 1 || result.Rows[0].Title != "alpha" {
		t.Fatalf("unexpected result: %+v", result)
	}
	next, err := url.Parse(result.Links.Next)
	if err != nil {
		t.Fatal(err)
	}
	if next.Query().Get("title") != "alpha" || next.Query().Get("page") != "2" {
		t.Fatalf("unexpected next link: %s", result.Links.Next)
	}
}

func TestInvalidPaginationAndSort(t *testing.T) {
	p := NewPaginationParams[struct{}]("", "", "", "-1", "invalid", "created_at desc")
	if p.Page != 1 || p.Limit != 10 {
		t.Fatalf("unexpected defaults: %+v", p)
	}
	if got := NewOrderBy(SortParams{Sort: "id; DROP TABLE tasks", Order: "asc", DefaultOrderBy: "updated_at DESC"}); got != "updated_at DESC" {
		t.Fatalf("unsafe sort: %q", got)
	}
}
