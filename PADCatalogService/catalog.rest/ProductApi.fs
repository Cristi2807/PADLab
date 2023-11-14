namespace catalog.rest

open System
open Giraffe
open Microsoft.AspNetCore.Http
open catalog
open catalog.Domain
open catalog.dal
open catalog.ServiceUtils
open catalog.rest.ServiceUtils
open FsToolkit.ErrorHandling
open Thoth.Json.Net

module ProductApi =

  let mutable openTransactions: Map<Guid, DalContextFactory> = Map.empty

  let getShoes (_: HttpFunc) (ctx: HttpContext) =
    let dalFactory = new DalContextFactory()

    asyncResult {
      use! roDalCtx = dalFactory.GetReadContext()
      let! shoes = ShoesDal.getShoes roDalCtx

      return shoes |> Seq.map Encode.shoes |> Encode.seq |> Encode.toString 2
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx dalFactory

  let getShoesById (shoesId: Guid) (_: HttpFunc) (ctx: HttpContext) =
    let dalFactory = new DalContextFactory()

    asyncResult {
      use! roDalCtx = dalFactory.GetReadContext()
      let! shoes = ShoesDal.getShoesById roDalCtx shoesId

      return shoes |> Encode.shoes |> Encode.toString 2
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx dalFactory

  let postShoes (_: HttpFunc) (ctx: HttpContext) =
    let dalFactory = new DalContextFactory()

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
    |> toApiResponse ctx dalFactory

  let postShoesTwoPhase (_: HttpFunc) (ctx: HttpContext) =
    let dalFactory = new DalContextFactory()

    asyncResult {
      let! reqBody = ctx.ReadBodyFromRequestAsync()

      let! shoes =
        Decode.fromString (Decode.shoes ()) reqBody
        |> Result.mapError (fun ex -> ex |> BadRequest)

      let! wDalCtx = dalFactory.GetWriteContext()

      let! shoesRes = ShoesDal.createShoes wDalCtx (shoes |> Shoes.fromDomain)

      let (ShoesId sId) = shoesRes.Id
      openTransactions <- openTransactions.Add(sId, dalFactory)

      return shoesRes |> Encode.shoes |> Encode.toString 2
    }
    |> AsyncResult.teeError (fun _ ->
      dalFactory.Rollback()
      dalFactory.Dispose())
    |> toApiResponse ctx dalFactory

  let commit (shoesId: Guid) (_: HttpFunc) (ctx: HttpContext) =
    let dalFactory = openTransactions.Item shoesId
    openTransactions <- openTransactions.Remove shoesId

    asyncResult {
      dalFactory.CommitTransaction()
      dalFactory.Dispose()
      return "Transaction commit"
    }
    |> toApiResponse ctx dalFactory

  let rollback (shoesId: Guid) (_: HttpFunc) (ctx: HttpContext) =
    let dalFactory = openTransactions.Item shoesId
    openTransactions <- openTransactions.Remove shoesId

    asyncResult {
      dalFactory.Rollback()
      dalFactory.Dispose()
      return "Transaction rollback"
    }
    |> toApiResponse ctx dalFactory

  let putShoesById (shoesId: Guid) (_: HttpFunc) (ctx: HttpContext) =
    let dalFactory = new DalContextFactory()

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
    |> toApiResponse ctx dalFactory

  let getStatus (_: HttpFunc) (ctx: HttpContext) =
    task {
      ctx.SetStatusCode 200
      return! ctx.WriteJsonAsync("OK")
    }

  let getMetrics (_: HttpFunc) (ctx: HttpContext) =
    task {
      let mutable listToReturn =
        [ "# HELP http_requests_total The total number of HTTP requests."
          "# TYPE http_requests_total counter" ]

      requestCounts.Keys
      |> Seq.iter (fun key ->
        listToReturn <- listToReturn @ [ $"http_requests_total{{code=\"{key}\"}} {requestCounts[key]}" ])

      ctx.SetStatusCode 200
      return! ctx.WriteTextAsync(String.concat "\n" listToReturn)
    }

  let productRoutes: HttpHandler =
    choose
      [ GET >=> route "/status" >=> getStatus

        GET >=> route "/metrics" >=> getMetrics

        POST >=> routef "/commit/%O" commit

        POST >=> routef "/rollback/%O" rollback

        POST >=> route "/shoes/2phase" >=> postShoesTwoPhase

        route "/shoes" >=> choose [ GET >=> getShoes; POST >=> postShoes ]

        routef "/shoes/%O" (fun docId -> choose [ GET >=> getShoesById docId; PUT >=> putShoesById docId ]) ]
