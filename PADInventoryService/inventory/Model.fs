namespace inventory

open System

module Domain =

  type TransactionId = TransactionId of Guid

  type ShoesId = ShoesId of Guid

  type Transaction =
    { Id: TransactionId
      ShoesId: ShoesId
      CreationDate: DateTime
      Quantity: int64
      OperationType: int16 }


open Domain

module Decode =
  open Thoth.Json.Net

  let transaction () : Decoder<Transaction> =
    Decode.object (fun get ->
      { Transaction.Id =
          get.Optional.Field "id" Decode.guid
          |> Option.defaultValue Guid.Empty
          |> TransactionId
        Transaction.ShoesId = get.Required.Field "shoesId" Decode.guid |> ShoesId
        Transaction.CreationDate =
          get.Optional.Field "creationDate" Decode.datetimeUtc
          |> Option.defaultValue DateTime.UtcNow
        Transaction.Quantity = get.Required.Field "quantity" Decode.int64
        Transaction.OperationType = get.Required.Field "operationType" Decode.int16 })

module Encode =
  open Thoth.Json.Net

  let TransactionId (TransactionId tId) = Encode.guid tId
  let ShoesId (ShoesId shoesId) = Encode.guid shoesId

  let transaction (transaction: Transaction) =
    [ "id", TransactionId transaction.Id
      "shoesId", ShoesId transaction.ShoesId
      "creationDate", Encode.datetime transaction.CreationDate
      "quantity", Encode.int64 transaction.Quantity
      "operationType", Encode.int16 transaction.OperationType ]
    |> Encode.object
