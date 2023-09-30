namespace inventory.dal

open System
open Dapper
open FsToolkit.ErrorHandling
open inventory.ServiceUtils

module InventoryDal =

  let getTransactionsByShoesId dalCtx (shoesId: Guid) =
    asyncResult {
      let! res =
        dalCtx.DbConnection.QueryAsync<Transaction>(
          """SELECT "Id", "ShoesId", "CreationDate", "Quantity", "OperationType"
	FROM inventory."Transactions"
  WHERE "ShoesId" = @shoesId;""",
          dict [ "shoesId" => shoesId ]
        )
        |> ToAsyncResult

      return res |> Seq.map Transaction.toDomain
    }

  let getTurnaround dalCtx (shoesId: Guid) (opType: int) =
    asyncResult {
      let! res =
        dalCtx.DbConnection.QueryFirstOrDefaultAsync<int64>(
          """SELECT COALESCE(SUM("Quantity"),0) AS Turnaround
FROM inventory."Transactions"
WHERE "ShoesId" = @shoesId AND "OperationType" = @opType;""",
          dict [ "shoesId" => shoesId; "opType" => opType ]
        )
        |> ToAsyncResult

      return res
    }

  let getTurnaroundTimePeriod dalCtx (shoesId: Guid) (opType: int) (sinceWhen: DateTime) (untilWhen: DateTime) =
    asyncResult {
      let! res =
        dalCtx.DbConnection.QueryFirstOrDefaultAsync<int64>(
          """SELECT COALESCE(SUM("Quantity"),0) AS Turnaround
FROM inventory."Transactions"
WHERE "ShoesId" = @shoesId AND "OperationType" = @opType AND "CreationDate" BETWEEN @sinceWhen AND @untilWhen;""",
          dict
            [ "shoesId" => shoesId
              "opType" => opType
              "sinceWhen" => sinceWhen
              "untilWhen" => untilWhen ]
        )
        |> ToAsyncResult

      return res
    }

  let createTransaction dalCtx (dto: Model.Transaction) =
    asyncResult {
      let! res =
        dalCtx.DbConnection.ExecuteScalarAsync<Guid>(
          """INSERT INTO inventory."Transactions"(
	"ShoesId", "CreationDate", "Quantity", "OperationType")
	VALUES (@shoesId, @creationDate, @quantity, @operationType)  RETURNING "Id";""",
          dict
            [ "shoesId" => dto.ShoesId
              "creationDate" => dto.CreationDate
              "quantity" => dto.Quantity
              "operationType" => dto.OperationType ]
        )
        |> ToAsyncResult

      return { dto with Id = res } |> Transaction.toDomain
    }

  let getStock dalCtx (shoesId: Guid) =
    asyncResult {
      let! res =
        dalCtx.DbConnection.QueryFirstOrDefaultAsync<int64>(
          """SELECT COALESCE(SUM("Quantity" * "OperationType"),0) AS Turnaround
FROM inventory."Transactions"
WHERE "ShoesId" = @shoesId;""",
          dict [ "shoesId" => shoesId ]
        )
        |> ToAsyncResult

      return res
    }
