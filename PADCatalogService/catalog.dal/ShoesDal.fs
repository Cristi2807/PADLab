namespace catalog.dal

open System
open Dapper
open FsToolkit.ErrorHandling
open catalog.ServiceUtils

module ShoesDal =

  let getShoes dalCtx =
    asyncResult {
      let! res =
        dalCtx.DbConnection.QueryAsync<Model.Shoes>(
          """SELECT "Id", "Color", "Size", "Price", "Brand", "Category", "Model"
	FROM catalog."Products";"""
        )
        |> ToAsyncResult

      return res |> Seq.map Shoes.toDomain
    }

  let getShoesById dalCtx (shoesId: Guid) =
    asyncResult {
      let! res =
        dalCtx.DbConnection.QueryFirstOrDefaultAsync<Model.Shoes>(
          """SELECT "Id", "Color", "Size", "Price", "Brand", "Category", "Model"
	FROM catalog."Products"
  WHERE "Id" = @shoesId;""",
          dict [ "shoesId" => shoesId ]
        )
        |> ToAsyncResult

      if res |> box |> isNull then
        return! "Object not found" |> NotFound |> Error
      else
        return res |> Shoes.toDomain
    }


  let createShoes dalCtx (dto: Model.Shoes) =
    asyncResult {
      let! res =
        dalCtx.DbConnection.ExecuteScalarAsync<Guid>(
          """INSERT INTO catalog."Products"(
	"Color", "Size", "Price", "Brand", "Category", "Model")
	VALUES (@color, @size, @price, @brand, @category, @model) RETURNING "Id";""",
          dict
            [ "color" => dto.Color
              "size" => dto.Size
              "price" => dto.Price
              "brand" => dto.Brand
              "category" => dto.Category
              "model" => dto.Model ]
        )
        |> ToAsyncResult

      return { dto with Id = res } |> Shoes.toDomain
    }

  let putShoes dalCtx (dto: Model.Shoes) =
    asyncResult {
      let! res =
        dalCtx.DbConnection.ExecuteAsync(
          """UPDATE catalog."Products"
	SET "Color"= @color, "Size"= @size, "Price"= @price, "Brand"= @brand, "Category"= @category, "Model"= @model
	WHERE "Id" = @id;""",
          dict
            [ "id" => dto.Id
              "color" => dto.Color
              "size" => dto.Size
              "price" => dto.Price
              "brand" => dto.Brand
              "category" => dto.Category
              "model" => dto.Model ]
        )
        |> ToAsyncResult
        |> Async.map (
          Result.bind (fun r ->
            if r = 0 then
              "Object to update not found or wrong identity" |> NotFound |> Error
            else
              r |> Ok // how many records have changed )
          )
        )

      return res
    }
