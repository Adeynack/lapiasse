import "@typespec/http";

using Http;

@service(#{ title: "La Piasse" })
namespace lapiasse;

model Resource {
  id: string;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model List<T> {
  items: T[];
}

model PaginatedRequest {
  @minValue(1)
  page: int32 = 1;

  @minValue(1)
  @maxValue(1000)
  page_size: int32 = 50;
}

model Book {
  ...Resource;
  name: string;
  default_currency_iso_code: string;
}

model Created<T> {
  @statusCode statusCode: 201;
  item: T;
}

@route("/books")
interface Books {
  @get list(...PaginatedRequest): List<Book>;
  @post create(@body book: Book): Created<Book>;
}
