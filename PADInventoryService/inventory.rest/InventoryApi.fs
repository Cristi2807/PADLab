namespace inventory.rest

open System
open Giraffe
open Microsoft.AspNetCore.Http
open inventory
open inventory.Domain
open inventory.dal
open inventory.ServiceUtils
open FsToolkit.ErrorHandling
open Thoth.Json.Net

module InventoryApi =
  let getTransactionsByShoesId (shoesId: Guid) (_: HttpFunc) (ctx: HttpContext) =
    use dalFactory = new DalContextFactory()

    asyncResult {
      use! roDalCtx = dalFactory.GetReadContext()
      let! transactions = InventoryDal.getTransactionsByShoesId roDalCtx shoesId

      return transactions |> Seq.map Encode.transaction |> Encode.seq |> Encode.toString 2
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx

  let getTurnaround (shoesId: Guid) (operationType: int) (_: HttpFunc) (ctx: HttpContext) =
    use dalFactory = new DalContextFactory()

    asyncResult {
      use! roDalCtx = dalFactory.GetReadContext()

      let! turnAround =
        match operationType with
        | 1 -> InventoryDal.getTurnaround roDalCtx shoesId 1
        | -1 -> InventoryDal.getTurnaround roDalCtx shoesId -1
        | _ -> "Invalid operationType" |> BadRequest |> AsyncResult.returnError

      return
        [ "turnaround", turnAround |> Encode.int64 ]
        |> Encode.object
        |> Encode.toString 2
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx

  let getTurnaroundTimePeriod
    (shoesId: Guid)
    (operationType: int)
    (sinceWhen: string)
    (untilWhen: string)
    (_: HttpFunc)
    (ctx: HttpContext)
    =
    use dalFactory = new DalContextFactory()

    asyncResult {
      use! roDalCtx = dalFactory.GetReadContext()

      let! turnAround =
        match operationType, sinceWhen, untilWhen with
        | 1, ValidDateTime since, ValidDateTime until ->
          InventoryDal.getTurnaroundTimePeriod roDalCtx shoesId 1 since until
        | -1, ValidDateTime since, ValidDateTime until ->
          InventoryDal.getTurnaroundTimePeriod roDalCtx shoesId -1 since until
        | _, ValidDateTime _, ValidDateTime _ -> "Invalid operationType" |> BadRequest |> AsyncResult.returnError
        | _ -> "Invalid sinceWhen / untilWhen Time" |> BadRequest |> AsyncResult.returnError

      return
        [ "turnaround", turnAround |> Encode.int64 ]
        |> Encode.object
        |> Encode.toString 2
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx

  let getStock (shoesId: Guid) (_: HttpFunc) (ctx: HttpContext) =
    use dalFactory = new DalContextFactory()

    asyncResult {
      use! roDalCtx = dalFactory.GetReadContext()
      let! stock = InventoryDal.getStock roDalCtx shoesId

      return [ "stock", stock |> Encode.int64 ] |> Encode.object |> Encode.toString 2
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx

  let postTransaction (_: HttpFunc) (ctx: HttpContext) =
    use dalFactory = new DalContextFactory()

    asyncResult {
      let! reqBody = ctx.ReadBodyFromRequestAsync()

      let! transaction =
        Decode.fromString (Decode.transaction ()) reqBody
        |> Result.mapError (fun ex -> ex |> BadRequest)

      let! wDalCtx = dalFactory.GetWriteContext()

      let! transactionRes =
        match transaction.OperationType, transaction.Quantity with
        | 1s, x when x > 0 ->
          InventoryDal.createTransaction
            wDalCtx
            { (transaction |> Transaction.fromDomain) with
                CreationDate = DateTime.UtcNow }
        | -1s, x when x > 0 ->
          asyncResult {
            use! roDalCtx = dalFactory.GetReadContext()
            let (ShoesId sId) = transaction.ShoesId
            let! stock = InventoryDal.getStock roDalCtx sId

            match stock - transaction.Quantity with
            | x when x < 0 -> return! "No such quantity in stock to retrieve" |> BadRequest |> AsyncResult.returnError
            | _ ->
              return!
                InventoryDal.createTransaction
                  wDalCtx
                  { (transaction |> Transaction.fromDomain) with
                      CreationDate = DateTime.UtcNow }
          }
        | _ -> "Invalid operationType and / or quantity" |> BadRequest |> AsyncResult.returnError

      return transactionRes |> Encode.transaction |> Encode.toString 2
    }
    |> dalFactory.HandleTransaction
    |> toApiResponse ctx

  let inventoryRoutes: HttpHandler =
    choose
      [ GET >=> routef "/transaction/%O" getTransactionsByShoesId

        GET >=> routef "/stock/%O" getStock

        GET
        >=> routef "/turnaround/%O/%i" (fun (shoesId, operationType) -> getTurnaround shoesId operationType)

        GET
        >=> routef "/turnaround/%O/%i/%s/%s" (fun (shoesId, operationType, sinceWhen, untilWhen) ->
          getTurnaroundTimePeriod shoesId operationType sinceWhen untilWhen)

        POST >=> route "/transaction" >=> postTransaction ]
