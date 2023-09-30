namespace catalog.rest

open System
open Giraffe
open Microsoft.AspNetCore.Http
open catalog
open catalog.Domain
open catalog.dal
open catalog.ServiceUtils
open FsToolkit.ErrorHandling
open Thoth.Json.Net

module ProductApi =

  let getShoes (_: HttpFunc) (ctx: HttpContext) =
    use dalFactory = new DalContextFactory()

    asyncResult {
      use! roDalCtx = dalFactory.GetReadContext()
      let! shoes = ShoesDal.getShoes roDalCtx

      return shoes |> Seq.map Encode.shoes |> Encode.seq |> Encode.toString 2
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx

  let getShoesById (shoesId: Guid) (_: HttpFunc) (ctx: HttpContext) =
    use dalFactory = new DalContextFactory()

    asyncResult {
      use! roDalCtx = dalFactory.GetReadContext()
      let! shoes = ShoesDal.getShoesById roDalCtx shoesId

      return shoes |> Encode.shoes |> Encode.toString 2
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx

  let postShoes (_: HttpFunc) (ctx: HttpContext) =
    use dalFactory = new DalContextFactory()

    asyncResult {
      let! reqBody = ctx.ReadBodyFromRequestAsync()

      let! shoes =
        Decode.fromString (Decode.shoes ()) reqBody
        |> Result.mapError (fun ex -> ex |> BadRequest)

      let! wDalCtx = dalFactory.GetWriteContext()

      let! shoesRes = ShoesDal.createShoes wDalCtx (shoes |> Shoes.fromDomain)

      return shoesRes |> Encode.shoes |> Encode.toString 2
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx

  let putShoesById (shoesId: Guid) (_: HttpFunc) (ctx: HttpContext) =
    use dalFactory = new DalContextFactory()

    asyncResult {
      let! reqBody = ctx.ReadBodyFromRequestAsync()

      let! shoes =
        Decode.fromString (Decode.shoes ()) reqBody
        |> Result.mapError (fun ex -> ex |> BadRequest)

      let! wDalCtx = dalFactory.GetWriteContext()

      let! rowsModified = ShoesDal.putShoes wDalCtx ({ shoes with Id = ShoesId shoesId } |> Shoes.fromDomain)

      return rowsModified |> string
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx

  let productRoutes: HttpHandler =
    choose
      [ route "/shoes" >=> choose [ GET >=> getShoes; POST >=> postShoes ]

        routef "/shoes/%O" (fun docId -> choose [ GET >=> getShoesById docId; PUT >=> putShoesById docId ]) ]
