//go:build test

package controller_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/app"
	"adeynack.net/lapiasse/pkg/controller"
	"adeynack.net/lapiasse/pkg/model"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.llib.dev/testcase"
	"gorm.io/gorm"
)

func TestBooksController(t *testing.T) {
	s := testcase.NewSpec(t)

	ctx := testcase.Let(s, func(t *testcase.T) context.Context {
		return app.CreateTestAppCtx(t)
	})

	ctrl := testcase.Let(s, func(t *testcase.T) *controller.BooksController {
		return &controller.BooksController{}
	})

	s.Describe("GetBooks", func(s *testcase.Spec) {
		request := testcase.Let(s, func(t *testcase.T) api.GetBooksRequestObject {
			return api.GetBooksRequestObject{}
		})

		response := testcase.Let(s, func(t *testcase.T) api.GetBooksResponseObject {
			resp, err := ctrl.Get(t).GetBooks(ctx.Get(t), request.Get(t))
			require.NoError(t, err)
			return resp
		})

		response200 := testcase.Let(s, func(t *testcase.T) api.GetBooks200JSONResponse {
			resp := response.Get(t)
			require.IsType(t, api.GetBooks200JSONResponse{}, resp)
			return resp.(api.GetBooks200JSONResponse)
		})

		responseBookNames := testcase.Let(s, func(t *testcase.T) []string {
			return lo.Map(response200.Get(t).Books, func(b api.Book, _ int) string { return b.Name })
		})

		s.When("there are no books", func(s *testcase.Spec) {
			s.Then("the response contains no book", func(t *testcase.T) {
				require.Empty(t, response200.Get(t).Books)
			})
		})

		s.When("some books exist", func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				db := ctxval.MustResolve[*gorm.DB](ctx.Get(t))
				for _, i := range rand.Perm(controller.DefaultPageSize + 10) {
					lo.Must0(gorm.G[model.Book](db).Create(ctx.Get(t), &model.Book{Name: fmt.Sprintf("Book %04d", i+1)}))
				}
			})

			s.Then("the response contains all books, up to the default limit", func(t *testcase.T) {
				expectedBookNames := lo.Map(lo.RangeFrom(1, controller.DefaultPageSize), func(i int, _ int) string {
					return fmt.Sprintf("Book %04d", i)
				})
				actual := responseBookNames.Get(t)
				require.Equal(t, expectedBookNames, actual)
			})

			s.When("a limit is given", func(s *testcase.Spec) {
				request.Let(s, func(t *testcase.T) api.GetBooksRequestObject {
					r := request.PreviousValue(t)
					r.Params.PageSize = lo.ToPtr(3)
					return r
				})

				s.Then("the response contains only the limited number of books", func(t *testcase.T) {
					require.Equal(t, []string{"Book 0001", "Book 0002", "Book 0003"}, responseBookNames.Get(t))
				})

				s.When("a page is given", func(s *testcase.Spec) {
					request.Let(s, func(t *testcase.T) api.GetBooksRequestObject {
						req := request.PreviousValue(t)
						req.Params.Page = lo.ToPtr(2)
						return req
					})

					s.Then("the requested page is returned", func(t *testcase.T) {
						require.Equal(t, []string{"Book 0004", "Book 0005", "Book 0006"}, responseBookNames.Get(t))
					})
				})
			})
		})
	})
}
